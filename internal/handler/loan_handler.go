package handler

import (
	"BankKibikov/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LoanHandler struct {
	logger      *zap.Logger
	loanService *service.LoanService
}

func NewLoanHandler(logger *zap.Logger, loanService *service.LoanService) *LoanHandler {
	return &LoanHandler{logger: logger, loanService: loanService}
}

// POST /loan
func (h *LoanHandler) CreateLoan(c *gin.Context) {
	var req struct {
		UserID string  `json:"user_id"`
		Amount float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	loan, err := h.loanService.CreateLoan(c.Request.Context(), req.UserID, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create loan"})
		return
	}

	c.JSON(http.StatusCreated, loan)
}

// GET /loan/:id
func (h *LoanHandler) GetLoan(c *gin.Context) {
	id := c.Param("id")
	current, err := h.loanService.GetLoanWithInterest(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "loan not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"loan_id": id, "current_debt": current})
}
