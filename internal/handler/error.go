package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func handleClientError(c *gin.Context, logger *zap.Logger, httpStatus int, clientMsg string, err error) {
	logger.Error(clientMsg, zap.Error(err))
	c.JSON(httpStatus, gin.H{"error": clientMsg})
	c.Abort()
}
