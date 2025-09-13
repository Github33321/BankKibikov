package handler

import (
	"BankKibikov/internal/models"
	"BankKibikov/internal/repository"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type NewsHandler struct {
	logger *zap.Logger
	repo   *repository.NewsRepository
}

func NewNewsHandler(logger *zap.Logger, repo *repository.NewsRepository) *NewsHandler {
	return &NewsHandler{logger: logger, repo: repo}
}

// GET /news
func (h *NewsHandler) GetNews(c *gin.Context) {
	news, err := h.repo.GetAll(context.Background())
	if err != nil {
		handleClientError(c, h.logger, http.StatusInternalServerError, "cannot fetch news", err)
		return
	}
	c.JSON(http.StatusOK, news)
}

// POST /news (для админа или тестов)
func (h *NewsHandler) CreateNews(c *gin.Context) {
	var n models.News
	if err := c.ShouldBindJSON(&n); err != nil {
		handleClientError(c, h.logger, http.StatusBadRequest, "invalid request body", err)
		return
	}
	if err := h.repo.Create(context.Background(), &n); err != nil {
		handleClientError(c, h.logger, http.StatusInternalServerError, "cannot create news", err)
		return
	}
	c.JSON(http.StatusCreated, n)
}
