package auth

import (
	"api/internal/errors"
	"api/internal/services/auth"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type Request struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Response struct {
	Error string `json:"error,omitempty"`
	Token string `json:"token,omitempty"`
}

type LoginHandler struct {
	log         *slog.Logger
	authService *auth.AuthService
}

func NewAuthHandler(
	log *slog.Logger,
	authService *auth.AuthService,
) *LoginHandler {
	return &LoginHandler{
		log:         log,
		authService: authService,
	}
}

func (h *LoginHandler) Login(c *gin.Context) {
	var request Request

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: errors.ErrInvalidRequest})
		return
	}

	token, err := h.authService.Login(request.Username, request.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, Response{Error: errors.ErrInvalidCredentials})
		return
	}

	c.JSON(http.StatusOK, Response{Token: token})
	return
}
