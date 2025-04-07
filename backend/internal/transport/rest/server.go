package rest

import (
	"context"
	"errors"
	"fmt"
	"l6/backend/internal/config"
	"l6/backend/pkg/logger"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Server struct {
	cfg     *config.AppServerConfig
	server  *http.Server
	handler *Handler
}

func NewServer( //nolint:gocritic // no named params
	l *slog.Logger, cfg *config.AppServerConfig, handler *Handler,
) (*Server, func(ctx context.Context) error) {
	s := &Server{
		cfg:     cfg,
		handler: handler,
	}

	l.Info("AppServer server initializing...")
	l.Info("bind AppServer to host: %s and port: %s", s.cfg.Host, s.cfg.Port)

	router := gin.New()

	router.Use(requestLoggerMiddleware(l))

	s.handler.InitRoutes(router)

	s.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port),
		Handler:      router,
		WriteTimeout: s.cfg.WriteTimeout,
		ReadTimeout:  s.cfg.ReadTimeout,
	}

	shutdown := func(ctx context.Context) error {
		return s.server.Shutdown(ctx)
	}

	return s, shutdown
}

func (s *Server) StartHTTP(l *slog.Logger) error {
	l.Info("AppServer server started")
	if err := s.server.ListenAndServe(); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			l.Info("AppServer server shutdown")
		default:
			l.Error("failed to start server", logger.ErrAttr(err))

			return fmt.Errorf("start server: %w", err)
		}
	}

	return nil
}

func requestLoggerMiddleware(l *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := uuid.New().String()

		newLogger := l.With("req_id", reqID)

		newLogger.Debug("Incoming request", "url", c.Request.URL.String())

		ctx := logger.ContextWithLogger(c.Request.Context(), newLogger)

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
