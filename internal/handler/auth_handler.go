package handler

import (
	"BankKibikov/internal/security"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type otpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type verifyRequest struct {
	Username string `json:"username"`
	OTP      string `json:"otp"`
}

type AuthHandler struct {
	logger      *zap.Logger
	authService *security.AuthService
}

func NewAuthHandler(logger *zap.Logger, authService *security.AuthService) *AuthHandler {
	return &AuthHandler{
		logger:      logger,
		authService: authService,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req otpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleClientError(c, h.logger, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if err := h.authService.RequestOTP(c.Request.Context(), req.Username, req.Password); err != nil {
		handleClientError(c, h.logger, http.StatusUnauthorized, "invalid username or password", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "OTP sent to email"})
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req verifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleClientError(c, h.logger, http.StatusBadRequest, "invalid request body", err)
		return
	}

	status, err := h.authService.VerifyOTP(c.Request.Context(), req.Username, req.OTP)
	if err != nil {
		handleClientError(c, h.logger, http.StatusUnauthorized, "invalid or expired OTP", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": status})
}
