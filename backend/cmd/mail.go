package main

import (
	"context"
	"fmt"
	"l6/backend/internal/config"
	"l6/backend/pkg/logger"
	"l6/backend/pkg/pgclient"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	notifyCtx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGABRT, syscall.SIGQUIT, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	shutdowns := make([]func(context.Context) error, 0)

	l := slog.Default()

	l.Info("Initializing config")
	cfg, err := config.LoadConfig(l)
	if err != nil {
		l.Error("Failed to load config", logger.ErrAttr(err))
		return
	}

	defer func() {
		closeConnections(shutdowns, l, cfg.App.ShutdownTimeout)
	}()

	l = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     cfg.Level(),
		AddSource: true,
	}))

	l.Info("PostgreSQL initializing")
	postgresDB, shutdown, err := pgclient.InitDB(notifyCtx, l, &cfg.Postgres)
	if err != nil {
		l.Error("Failed to initialize postgresSQL", logger.ErrAttr(err))

		return
	}

	fmt.Println(postgresDB)

	shutdowns = append(shutdowns, shutdown)

	<-notifyCtx.Done()

	l.Info("Termination signal received, shutting down...")
}

func closeConnections(shutdowns []func(context.Context) error, l *slog.Logger, timeout time.Duration) {
	l.Warn("Closing connections")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(len(shutdowns))

	for _, shutdownFn := range shutdowns {
		go func(fn func(context.Context) error) {
			defer wg.Done()
			err := fn(ctx)
			if err != nil {
				l.Error("Closing connection", logger.ErrAttr(err))
			}
		}(shutdownFn)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		l.Error("Context deadline exceeded before shutdown completed")
	case <-done:
		l.Info("All connections closed successfully")
	}

	l.Info("App terminated")
}
