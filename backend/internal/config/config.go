package config

import (
	"fmt"
	"l6/backend/pkg/logger"
	"log/slog"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App       AppConfig       `yaml:"app"`
	Postgres  PostgresConfig  `yaml:"postgres"`
	AppServer AppServerConfig `yaml:"appServer"`
}

type AppConfig struct {
	LogLevel        string        `env:"APP_LOG_LEVEL"          yaml:"logLevel"`
	ShutdownTimeout time.Duration `env:"APP_SHUTDOWN_TIMEOUT"   yaml:"shutdownTimeout"`
}

type PostgresConfig struct {
	Host               string        `env:"POSTGRES_HOST"                yaml:"host"`
	Port               string        `env:"POSTGRES_PORT"                yaml:"port"`
	Database           string        `env:"POSTGRES_DATABASE"            yaml:"database"`
	Username           string        `env:"POSTGRES_USERNAME"            yaml:"username"`
	Password           string        `env:"POSTGRES_PASSWORD"            yaml:"password"`
	ConnectionAttempts int           `env:"POSTGRES_CONNECTION_ATTEMPTS" yaml:"connectionAttempts"`
	DelayBtwAttempts   time.Duration `env:"POSTGRES_DELAY_BTW_ATTEMPTS"  yaml:"delayBtwAttempts"`
	QueryTimeout       time.Duration `env:"POSTGRES_QUERY_TIMEOUT"       yaml:"queryTimeout"`
}

type AppServerConfig struct {
	Host               string        `env:"APP_HTTP_HOST"                 yaml:"host"`
	Port               string        `env:"APP_HTTP_PORT"                 yaml:"port"`
	ReadTimeout        time.Duration `env:"APP_HTTP_READ_TIMEOUT"         yaml:"readTimeout"`
	WriteTimeout       time.Duration `env:"APP_HTTP_WRITE_TIMEOUT"        yaml:"writeTimeout"`
	LongPollingTimeout time.Duration `env:"APP_HTTP_LONG_POLLING_TIMEOUT" yaml:"longPollingTimeout"`
}

func LoadConfig(l *slog.Logger) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("local.yaml", cfg)
	if err != nil {
		l.Warn("Error loading config from file", logger.ErrAttr(err))
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		l.Error("Error reading environment variables", logger.ErrAttr(err))

		return nil, fmt.Errorf("reading environment variables: %w", err)
	}

	return cfg, nil
}

func (c Config) Level() slog.Level {
	switch c.App.LogLevel {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
