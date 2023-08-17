package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Configurations exported
type Configurations struct {
	Secret                string  `mapstructure:"SECRET"`
	APIHost               string  `mapstructure:"API_HOST"`
	DebugHost             string  `mapstructure:"DEBUG_HOST"`
	ReadTimeout           int     `mapstructure:"READ_TIMEOUT"`
	WriteTimeout          int     `mapstructure:"WRITE_TIMEOUT"`
	ShutdownTimeout       int     `mapstructure:"SHUTDOWN_TIMEOUT"`
	DBUsername            string  `mapstructure:"DB_USERNAME"`
	DBPassword            string  `mapstructure:"DB_PASSWORD"`
	DBName                string  `mapstructure:"DB_NAME"`
	DBHost                string  `mapstructure:"DB_HOST"`
	MaxIdleConns          int     `mapstructure:"MAX_IDLE_CONN"`
	MaxOpenConns          int     `mapstructure:"MAX_OPEN_CONN"`
	DisableTLS            bool    `mapstructure:"DISABLE_TLS"`
	GOOGLE_SHEET_TOKEN    string  `mapstructure:"GOOGLE_SHEET_TOKEN"`
	NAME                  string  `mapstructure:"NAME"`
	TELEMETRY_URI         string  `mapstructure:"TELEMETRY_URI"`
	TELEMETRY_ENABLED     bool    `mapstructure:"telemetry_enabled"`
	PROBABILITY           float64 `mapstructure:"PROBABILITY"`
	DOMAIN_KAFKA_TOPIC    string  `mapstructure:"domain_kafka_topic"`
	DOMAIN_KAFKA_BROKER   string  `mapstructure:"domain_kafka_broker"`
}

type DatabaseConfigurations struct {
	DBUsername   string
	DBPassword   string
	DBName       string
	DBHost       string
	MaxIdleConns int
	MaxOpenConns int
	DisableTLS   bool
}

func LoadConfig(build, path string) (Configurations, error) {
	// =========================================================================
	// Configuration from file
	// Set the file name of the configurations file
	viper.AddConfigPath(path)
	viper.SetConfigName(build)
	viper.SetConfigType("env")
	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	cfg := Configurations{}
	if err := viper.ReadInConfig(); err != nil {
		return cfg, errors.Wrap(err, "Error reading config file")
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, errors.Wrap(err, "Unable to decode into struct")
	}

	return cfg, nil
}

func (cfg Configurations) GetDBConfig() DatabaseConfigurations {
	return DatabaseConfigurations{
		DBUsername:   cfg.DBUsername,
		DBPassword:   cfg.DBPassword,
		DBHost:       cfg.DBHost,
		DBName:       cfg.DBName,
		MaxIdleConns: cfg.MaxIdleConns,
		MaxOpenConns: cfg.MaxOpenConns,
		DisableTLS:   cfg.DisableTLS,
	}
}
