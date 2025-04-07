package pgclient

import (
	"context"
	"fmt"
	"l6/backend/internal/config"
	"l6/backend/pkg/logger"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

func InitDB(ctx context.Context, l *slog.Logger, cfg *config.PostgresConfig) (*sqlx.DB, func(ctx context.Context) error, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.Database,
		cfg.Password,
	)
	l.Debug("postgres", "dsn", dsn)

	var db *sqlx.DB

	err := DoWithAttempts(func() error {
		var connectErr error
		db, connectErr = sqlx.ConnectContext(ctx, "postgres", dsn)
		if connectErr != nil {
			l.Warn("Failed to connect to postgres", logger.ErrAttr(connectErr))

			return fmt.Errorf("connect to postgres: %w", connectErr)
		}

		return nil
	}, cfg.ConnectionAttempts, cfg.DelayBtwAttempts)

	if err != nil {
		return nil, nil, fmt.Errorf("connect to postgres: %w", err)
	}

	l.Debug("Successfully connected to postgresql")

	shutdown := func(_ context.Context) error {
		return db.Close()
	}

	return db, shutdown, nil
}

func DoWithAttempts(fn func() error, attempts int, delay time.Duration) error {
	var err error
	for attempts > 0 {
		if err = fn(); err == nil {
			return nil
		}
		time.Sleep(delay)
		attempts--
	}

	return err
}
