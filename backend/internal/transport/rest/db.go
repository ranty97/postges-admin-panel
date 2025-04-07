package rest

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Service interface {
	Tables(ctx context.Context) ([]string, error)
	ExecuteQuery(ctx context.Context, query string) (string, error)
}

type Handler struct {
	logger  *slog.Logger
	service Service
}

func NewHandler(service Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) InitRoutes(router *gin.Engine) {
	router.GET("/tables", h.Tables)
	router.POST("/execute", h.Execute)
}

type TableResponse struct {
	Tables []string `json:"tables"`
}

func (h *Handler) Tables(c *gin.Context) {
	h.logger.Info("Tables request received")
	tables, err := h.service.Tables(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, TableResponse{Tables: tables})
}

type ExecuteRequest struct {
	Query string `json:"query"`
}

func (h *Handler) Execute(c *gin.Context) {
	h.logger.Info("Execute request received")
	var request ExecuteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		h.logger.Error("Failed to bind request", "error", err)
		return
	}
	result, err := h.service.ExecuteQuery(c, request.Query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logger.Error("Failed to execute query", "error", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": result})
	h.logger.Info("Query executed successfully", "result", result)
}
