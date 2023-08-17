package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	tr "go.opentelemetry.io/otel/trace"

	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"

	"collect/app/services/collect-api/handlers"
	"collect/app/services/debug-api"
	"collect/foundation/logger"

	"collect/business/sys/database"
	"collect/foundation/config"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	_ "github.com/spf13/viper/remote"
)

const build = "dev"

func main() {
	log, err := logger.New("COLLECT-API")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()
	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		log.Sync()
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {
	opt := maxprocs.Logger(log.Infof)
	if _, err := maxprocs.Set(opt); err != nil {
		return fmt.Errorf("maxprocs: %w", err)
	}

	fmt.Printf("Go version: %s\n", runtime.Version())

	cfg, err := config.LoadConfig(build, ".")
	if err != nil {
		return err
	}

	// =========================================================================
	// Start Database

	log.Infow("startup", "status", "initializing database support", "host", cfg.DBHost)
	db, err := setupDB(cfg)
	if err != nil {
		log.Errorw("main: unable to connect to database")
		return err
	}

	defer func() {
		log.Infow("shutdown", "status", "stopping database support", "host", cfg.DBHost)
		db.Close()
	}()

	log.Infow("main:", "config", cfg)

	// =========================================================================
	// Start Tracing Support

	log.Infow("startup", "status", "initializing OT/Zipkin tracing support")

	traceProvider, err := startTracing(
		cfg.NAME,
		cfg.TELEMETRY_URI,
		cfg.PROBABILITY,
		log.Desugar(),
	)
	if err != nil {
		return fmt.Errorf("starting tracing: %w", err)
	}
	defer traceProvider.Shutdown(context.Background())

	tracer := traceProvider.Tracer("collect-service")

	// =========================================================================
	// Start Debug Service

	go func() {
		if err := http.ListenAndServe(cfg.DebugHost, debug.Mux(build, log, db)); err != nil {
			log.Fatal("shutdown", "status", "debug v1 router closed", cfg.DebugHost, "ERROR", err)
		}
		log.Infow("startup", "status", "debug v1 router started", "host", cfg.DebugHost)
	}()

	// =========================================================================
	// Start API Service

	log.Infow("startup", "status", "initializing V1 API support", "host", cfg.APIHost)

	return runAPI(cfg, db, log, tracer)
}

func setupDB(cfg config.Configurations) (*sqlx.DB, error) {
	db, err := database.Open(cfg.GetDBConfig())
	if err != nil {
		return db, errors.Wrap(err, "connecting to db")
	}

	return db, nil
}

func runAPI(cfg config.Configurations, db *sqlx.DB, log *zap.SugaredLogger, tracer tr.Tracer) error {

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	serverErrors := make(chan error, 1)

	apiMux := handlers.APIMux(handlers.APIMuxConfig{
		Shutdown:     shutdown,
		Log:          log,
		DB:           db,
		Tracer:       tracer,
		ServerErrors: serverErrors,
		Cfg:          cfg,
	})

	api := http.Server{
		Addr:         cfg.APIHost,
		Handler:      apiMux,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		ErrorLog:     zap.NewStdLog(log.Desugar()),
	}

	go func() {
		log.Infow("startup", "status", "api router started", "host", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	return waitForCompletion(serverErrors, shutdown, &api, time.Duration(cfg.ShutdownTimeout)*time.Second)
}

// startTracing configure open telemetry to be used with zipkin.
func startTracing(name string, telemetryURI string, prob float64, log *zap.Logger) (*trace.TracerProvider, error) {

	exporter, err := zipkin.New(
		telemetryURI,
		zipkin.WithLogger(zap.NewStdLog(log.With(zap.Any("service", name)))),
	)
	if err != nil {
		return nil, fmt.Errorf("creating new exporter: %w", err)
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.TraceIDRatioBased(prob)),
		trace.WithBatcher(exporter,
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
			trace.WithBatchTimeout(trace.DefaultScheduleDelay*time.Millisecond),
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
		),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(name),
				attribute.String("exporter", "zipkin"),
			),
		),
	)

	otel.SetTracerProvider(traceProvider)

	return traceProvider, nil
}

func waitForCompletion(serverErrors chan error, shutdown chan os.Signal, api *http.Server, timeout time.Duration) error {
	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "server error")

	case sig := <-shutdown:
		log.Printf("main: %+v : Start shutdown", sig)
		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Asking listener to shutdown and shed load.
		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return errors.Wrap(err, "could not stop server gracefully")
		}
	}
	return nil
}
