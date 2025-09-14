package handler

import (
	"BankKibikov/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AccountHandler struct {
	logger         *zap.Logger
	accountService *service.AccountService
}

func NewAccountHandler(logger *zap.Logger, accountService *service.AccountService) *AccountHandler {
	return &AccountHandler{logger: logger, accountService: accountService}
}

func (h *AccountHandler) GetBalance(c *gin.Context) {
	userID := c.Query("user_id") // пока через query (можно привязать к JWT позже)

	acc, err := h.accountService.GetBalance(c.Request.Context(), userID)
	if err != nil {
		handleClientError(c, h.logger, http.StatusInternalServerError, "cannot fetch balance", err)
		return
	}
	c.JSON(http.StatusOK, acc)
}

func (h *AccountHandler) Transfer(c *gin.Context) {
	from := c.PostForm("from_user")
	to := c.PostForm("to_user")
	amountStr := c.PostForm("amount")

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		handleClientError(c, h.logger, http.StatusBadRequest, "invalid amount", err)
		return
	}

	if err := h.accountService.Transfer(c.Request.Context(), from, to, amount); err != nil {
		handleClientError(c, h.logger, http.StatusBadRequest, "transfer failed", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "transfer successful"})
}

func (h *AccountHandler) GetTransactions(c *gin.Context) {
	userID := c.Query("user_id")

	list, err := h.accountService.GetTransactions(c.Request.Context(), userID)
	if err != nil {
		handleClientError(c, h.logger, http.StatusInternalServerError, "cannot fetch transactions", err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *AccountHandler) Deposit(c *gin.Context) {
	userID := c.PostForm("user_id")
	amountStr := c.PostForm("amount")

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		handleClientError(c, h.logger, http.StatusBadRequest, "invalid amount", err)
		return
	}

	if err := h.accountService.Deposit(c.Request.Context(), userID, amount); err != nil {
		handleClientError(c, h.logger, http.StatusBadRequest, "deposit failed", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deposit successful"})
}

func (h *AccountHandler) Withdraw(c *gin.Context) {
	userID := c.PostForm("user_id")
	amountStr := c.PostForm("amount")

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		handleClientError(c, h.logger, http.StatusBadRequest, "invalid amount", err)
		return
	}

	if err := h.accountService.Withdraw(c.Request.Context(), userID, amount); err != nil {
		handleClientError(c, h.logger, http.StatusBadRequest, "withdraw failed", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "withdraw successful"})
}
