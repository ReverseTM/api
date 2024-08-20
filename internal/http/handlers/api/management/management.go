package management

import (
	"api/internal/errors"
	"api/internal/services/management"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type ReadRequest struct {
	Keys []string `json:"keys" binding:"required"`
}

type ReadResponse struct {
	Error string         `json:"error,omitempty"`
	Data  map[string]any `json:"data"`
}

type WriteRequest struct {
	Data map[string]any `json:"data" binding:"required"`
}

type WriteResponse struct {
	Error  string `json:"error,omitempty"`
	Status string `json:"status,omitempty"`
}

const Success = "success"

type ManagementHandler struct {
	log               *slog.Logger
	managementService *management.ManagementService
}

func NewManagementHandler(
	log *slog.Logger,
	managementService *management.ManagementService,
) *ManagementHandler {
	return &ManagementHandler{
		log:               log,
		managementService: managementService,
	}
}

func (h *ManagementHandler) Read(c *gin.Context) {
	const op = "handlers.api.management.Read"

	log := h.log.With(
		slog.String("op", op),
	)

	var request ReadRequest
	if err := c.BindJSON(&request); err != nil {
		log.Error("failed at parsing request")
		c.JSON(http.StatusBadRequest, ReadResponse{Error: errors.ErrInvalidRequest})
		return
	}

	result, err := h.managementService.Read(request.Keys)
	if err != nil {
		log.Error("failed to read data")
		c.JSON(http.StatusInternalServerError, ReadResponse{Error: errors.ErrInternalServer})
		return
	}

	c.JSON(http.StatusOK, ReadResponse{Data: result})
	return
}

func (h *ManagementHandler) Write(c *gin.Context) {
	const op = "handlers.api.management.Write"

	log := h.log.With(
		slog.String("op", op),
	)

	var request WriteRequest
	if err := c.BindJSON(&request); err != nil {
		log.Error("failed at parsing request")
		c.JSON(http.StatusBadRequest, WriteResponse{Error: errors.ErrInvalidRequest})
		return
	}

	err := h.managementService.Write(request.Data)
	if err != nil {
		log.Error("failed to write data")
		c.JSON(http.StatusInternalServerError, WriteResponse{Error: errors.ErrInternalServer})
		return
	}

	c.JSON(http.StatusOK, WriteResponse{Status: Success})
	return
}
