package rest

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Database API
// @version 1.0
// @description API для работы с базой данных
// @host localhost:8080
// @BasePath /

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
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// TableResponse represents the response for the tables endpoint
type TableResponse struct {
	Tables []string `json:"tables"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// @Summary Get list of tables
// @Description Returns a list of all tables in the database
// @Tags tables
// @Accept json
// @Produce json
// @Success 200 {object} TableResponse
// @Failure 500 {object} ErrorResponse
// @Router /tables [get]
func (h *Handler) Tables(c *gin.Context) {
	h.logger.Info("Tables request received")
	tables, err := h.service.Tables(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, TableResponse{Tables: tables})
}

type ExecuteRequest struct {
	Query string `json:"query"`
}

// @Summary Execute SQL query
// @Description Executes an arbitrary SQL query and returns the result
// @Tags execute
// @Accept json
// @Produce json
// @Param request body ExecuteRequest true "SQL query"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /execute [post]
func (h *Handler) Execute(c *gin.Context) {
	h.logger.Info("Execute request received")
	var request ExecuteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		h.logger.Error("Failed to bind request", "error", err)
		return
	}
	result, err := h.service.ExecuteQuery(c, request.Query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		h.logger.Error("Failed to execute query", "error", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": result})
	h.logger.Info("Query executed successfully", "result", result)
}
