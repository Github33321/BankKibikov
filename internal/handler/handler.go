package handler

import (
	"BankKibikov/internal/chat"
	"BankKibikov/internal/repository"
	"BankKibikov/internal/security"
	"BankKibikov/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	Logger         *zap.Logger
	UserHandler    *UserHandler
	AuthHandler    *AuthHandler
	NewsHandler    *NewsHandler
	AccountHandler *AccountHandler
	StockHandler   *StockHandler
	LoanHandler    *LoanHandler
	ChatHandler    *ChatHandler // ✅ добавили чат
	jwtSecret      string
}

func NewHandler(
	logger *zap.Logger,
	userService *service.UserService,
	authService *security.AuthService,
	newsRepo *repository.NewsRepository,
	accountService *service.AccountService,
	loanService *service.LoanService,
	chatHub *chat.Hub, // ✅ чат
	secret string,
) *Handler {
	return &Handler{
		Logger:         logger,
		UserHandler:    NewUserHandler(logger, userService),
		AuthHandler:    NewAuthHandler(logger, authService),
		NewsHandler:    NewNewsHandler(logger, newsRepo),
		AccountHandler: NewAccountHandler(logger, accountService),
		StockHandler:   NewStockHandler(),
		LoanHandler:    NewLoanHandler(logger, loanService),
		ChatHandler:    NewChatHandler(chatHub), // ✅ инициализация
		jwtSecret:      secret,
	}
}

func (h *Handler) InitRoutes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Bank Kibikov API is running"})
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	secure := router.Group("/")
	{
		secure.GET("/balance", h.AccountHandler.GetBalance)
		secure.POST("/transfer", h.AccountHandler.Transfer)
		secure.GET("/transactions", h.AccountHandler.GetTransactions)

		// только для админов
		secure.POST("/deposit", security.AdminOnly(), h.AccountHandler.Deposit)
		secure.POST("/withdraw", security.AdminOnly(), h.AccountHandler.Withdraw)
	}

	// новости
	router.GET("/news", h.NewsHandler.GetNews)
	router.POST("/news", h.NewsHandler.CreateNews)

	// пользователи
	router.POST("/users", h.UserHandler.CreateUser)
	router.POST("/login", h.AuthHandler.Login)
	router.POST("/verify-otp", h.AuthHandler.VerifyOTP)
	router.GET("/users/:id", h.UserHandler.GetUserByID)
	router.GET("/users", h.UserHandler.GetUsers)

	// биржа
	router.GET("/moex-stocks", h.StockHandler.GetMoexStocks)

	// займы
	router.POST("/loan", h.LoanHandler.CreateLoan)
	router.GET("/loan/:id", h.LoanHandler.GetLoan)

	// чат
	router.StaticFile("/chat", "static/chat.html") // панель чата
	router.GET("/ws", h.ChatHandler.ServeWs)       // websocket
}
