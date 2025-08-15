package handlers

import (
	"deni-be-crm/internal/common"
	"deni-be-crm/internal/models"
	"deni-be-crm/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service services.IAuthService
}

func NewAuthHandler(s services.IAuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		if httpErr, ok := err.(*common.HTTPError); ok {
			c.JSON(httpErr.Code, common.Error(httpErr.Message))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Error("Failed to login"))
		return
	}

	c.JSON(http.StatusOK, common.Success(gin.H{"token_access": token}))
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Role == "" {
		req.Role = "sales"
	}

	if err := h.service.Register(&req); err != nil {
		if httpErr, ok := err.(*common.HTTPError); ok {
			c.JSON(httpErr.Code, common.Error(httpErr.Message))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Error("Failed to register"))
		return
	}

	c.JSON(http.StatusCreated, common.Success("User registered successfully"))
}
