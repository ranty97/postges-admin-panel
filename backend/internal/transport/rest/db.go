package rest

import (
	"context"
	"fmt"
	"l6/internal/domain"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @title Database API
// @version 1.0
// @description API для работы с базой данных
// @host localhost:8080
// @BasePath /

type Service interface {
	Tables(ctx context.Context) ([]string, error)
	ExecuteQuery(ctx context.Context, query string) (string, error)
	ListBackups(ctx context.Context) ([]domain.Backup, error)
	CreateBackup(ctx context.Context) (domain.BackupCreated, error)
	DownloadBackup(ctx context.Context, filename string) ([]byte, error)
	DeleteBackup(ctx context.Context, filename string) error
	RestoreBackup(ctx context.Context, filename string) error
	DeleteAllTables(ctx context.Context) error
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
	router.GET("/backup/list", h.ListBackups)
	router.POST("/backup/create", h.CreateBackup)
	router.GET("/backup/download/:filename", h.DownloadBackup)
	router.DELETE("/backup/delete/:filename", h.DeleteBackup)
	router.POST("/backup/restore/:filename", h.RestoreBackup)
	router.DELETE("/tables/delete/all", h.DeleteAllTables)
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
	h.logger.Info("Executing query", "query", request.Query)
	result, err := h.service.ExecuteQuery(c, request.Query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		h.logger.Error("Failed to execute query", "error", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": result})
	h.logger.Info("Query executed successfully", "result", result)
}

// @Summary Get list of backups
// @Description Returns a list of all available backups
// @Tags backup
// @Accept json
// @Produce json
// @Success 200 {object} map[string][]domain.Backup
// @Failure 500 {object} ErrorResponse
// @Router /backup/list [get]
func (h *Handler) ListBackups(c *gin.Context) {
	h.logger.Info("ListBackups request received")
	backups, err := h.service.ListBackups(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		h.logger.Error("Failed to list backups", "error", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"backups": backups})
	h.logger.Info("Backups listed successfully", "backups", backups)
}

// @Summary Create new backup
// @Description Creates a new backup of the database
// @Tags backup
// @Accept json
// @Produce json
// @Success 200 {object} map[string]domain.BackupCreated
// @Failure 500 {object} ErrorResponse
// @Router /backup/create [post]
func (h *Handler) CreateBackup(c *gin.Context) {
	h.logger.Info("CreateBackup request received")
	backup, err := h.service.CreateBackup(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		h.logger.Error("Failed to create backup", "error", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"backup": backup})
	h.logger.Info("Backup created successfully")
}

// @Summary Download backup
// @Description Downloads a specific backup file
// @Tags backup
// @Accept json
// @Produce application/sql
// @Param filename path string true "Backup filename"
// @Success 200 {file} application/sql
// @Failure 500 {object} ErrorResponse
// @Router /backup/download/{filename} [get]
func (h *Handler) DownloadBackup(c *gin.Context) {
	h.logger.Info("DownloadBackup request received")
	filename := c.Param("filename")
	backup, err := h.service.DownloadBackup(c, filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		h.logger.Error("Failed to download backup", "error", err)
		return
	}
	c.Header("Content-Type", "application/sql")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(http.StatusOK, "application/sql", backup)
	h.logger.Info("Backup downloaded successfully", "filename", filename)
}

// @Summary Delete backup
// @Description Deletes a specific backup file
// @Tags backup
// @Accept json
// @Produce json
// @Param filename path string true "Backup filename"
// @Success 200 {object} domain.BackupDeleted
// @Failure 500 {object} ErrorResponse
// @Router /backup/delete/{filename} [delete]
func (h *Handler) DeleteBackup(c *gin.Context) {
	h.logger.Info("DeleteBackup request received")
	filename := c.Param("filename")
	err := h.service.DeleteBackup(c, filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		h.logger.Error("Failed to delete backup", "error", err)
		return
	}

	c.JSON(http.StatusOK, domain.BackupDeleted{Success: true, Message: "Backup deleted successfully"})
	h.logger.Info("Backup deleted successfully", "filename", filename)
}

// @Summary Restore backup
// @Description Restores the database from a specific backup
// @Tags backup
// @Accept json
// @Produce json
// @Param filename path string true "Backup filename"
// @Success 200 {object} domain.BackupCreated
// @Failure 500 {object} ErrorResponse
// @Router /backup/restore/{filename} [post]
func (h *Handler) RestoreBackup(c *gin.Context) {
	h.logger.Info("RestoreBackup request received")
	filename := c.Param("filename")
	err := h.service.RestoreBackup(c, filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		h.logger.Error("Failed to restore backup", "error", err)
		return
	}
	c.JSON(http.StatusOK, domain.BackupCreated{Success: true, Message: "Backup restored successfully"})
	h.logger.Info("Backup restored successfully", "filename", filename)
}

// @Summary Delete all tables
// @Description Deletes all tables from the database
// @Tags tables
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} ErrorResponse
// @Router /tables/delete/all [delete]
func (h *Handler) DeleteAllTables(c *gin.Context) {
	h.logger.Info("DeleteAllTables request received")
	err := h.service.DeleteAllTables(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		h.logger.Error("Failed to delete all tables", "error", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "All tables deleted successfully"})
	h.logger.Info("All tables deleted successfully")
}
