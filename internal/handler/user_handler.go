package handler

import (
	"BankKibikov/internal/models"
	"BankKibikov/internal/service"
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	logger      *zap.Logger
	userService *service.UserService
}

func NewUserHandler(logger *zap.Logger, userService *service.UserService) *UserHandler {
	return &UserHandler{
		logger:      logger,
		userService: userService,
	}
}

// POST /users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var u models.User
	if err := c.ShouldBindJSON(&u); err != nil {
		handleClientError(c, h.logger, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if err := h.userService.CreateUser(context.Background(), &u); err != nil {
		handleClientError(c, h.logger, http.StatusInternalServerError, "cannot create user", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "user created",
		"id":     u.ID,
	})
}

// GET /users/:id
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	_, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		handleClientError(c, h.logger, http.StatusBadRequest, "invalid id", err)
		return
	}

	user, err := h.userService.GetUser(context.Background(), idStr)
	if err != nil {
		handleClientError(c, h.logger, http.StatusNotFound, "user not found", err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// GET /users
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.userService.GetUsers(context.Background())
	if err != nil {
		handleClientError(c, h.logger, http.StatusInternalServerError, "cannot fetch users", err)
		return
	}

	c.JSON(http.StatusOK, users)
}
