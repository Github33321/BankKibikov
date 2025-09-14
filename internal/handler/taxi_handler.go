package handler

import (
	"BankKibikov/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TaxiHandler struct {
	logger      *zap.Logger
	taxiService *service.TaxiService
}

func NewTaxiHandler(logger *zap.Logger, taxiService *service.TaxiService) *TaxiHandler {
	return &TaxiHandler{logger: logger, taxiService: taxiService}
}

func (h *TaxiHandler) OrderTaxi(c *gin.Context) {
	var req struct {
		UserID string  `json:"user_id"`
		From   string  `json:"from"`
		To     string  `json:"to"`
		Price  float64 `json:"price"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		handleClientError(c, h.logger, http.StatusBadRequest, "invalid request", err)
		return
	}

	order, err := h.taxiService.OrderTaxi(c.Request.Context(), req.UserID, req.From, req.To, req.Price)
	if err != nil {
		handleClientError(c, h.logger, http.StatusBadRequest, "cannot create taxi order", err)
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *TaxiHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	order, err := h.taxiService.GetOrder(c.Request.Context(), id)
	if err != nil {
		handleClientError(c, h.logger, http.StatusNotFound, "order not found", err)
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *TaxiHandler) GetUserOrders(c *gin.Context) {
	userID := c.Query("user_id")
	list, err := h.taxiService.GetUserOrders(c.Request.Context(), userID)
	if err != nil {
		handleClientError(c, h.logger, http.StatusInternalServerError, "cannot fetch orders", err)
		return
	}
	c.JSON(http.StatusOK, list)
}
