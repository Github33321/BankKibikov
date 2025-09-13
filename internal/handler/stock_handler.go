package handler

import (
	"BankKibikov/internal/market"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StockHandler struct{}

func NewStockHandler() *StockHandler {
	return &StockHandler{}
}

// Получение котировок с Московской Биржи (без токена)
func (h *StockHandler) GetMoexStocks(c *gin.Context) {
	quotes, err := market.GetMoexPrices(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot fetch moex prices"})
		return
	}
	c.JSON(http.StatusOK, quotes)
}
