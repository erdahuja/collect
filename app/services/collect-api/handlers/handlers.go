// Package handlers contains the full set of handler functions and routes
// supported by the web api.
package handlers

import (
	"net/http"
	"os"

	v1 "collect/app/services/collect-api/handlers/v1"
	"collect/foundation/config"
	"collect/foundation/web"

	"collect/business/auth"
	"collect/business/mid"

	"go.opentelemetry.io/otel/trace"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Shutdown     chan os.Signal
	Log          *zap.SugaredLogger
	DB           *sqlx.DB
	Auth         *auth.Auth
	Tracer       trace.Tracer
	ServerErrors chan error
	Cfg          config.Configurations
}

// APIMux constructs a http.Handler with all application routes defined.
func APIMux(cfg APIMuxConfig) http.Handler {
	var app *web.App

	if app == nil {
		app = web.NewApp(
			cfg.Shutdown,
			cfg.Tracer,
			mid.Logger(cfg.Log),
			mid.Errors(cfg.Log),
			mid.Metrics(),
			mid.Panics(),
		)
	}

	v1.Register(app, v1.Config{
		Log:          cfg.Log,
		DB:           cfg.DB,
		ServerErrors: cfg.ServerErrors,
		Cfg:          &cfg.Cfg,
	})

	return app
}
