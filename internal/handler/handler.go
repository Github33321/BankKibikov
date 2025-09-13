package handler

import (
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
	jwtSecret      string
}

func NewHandler(
	logger *zap.Logger,
	userService *service.UserService,
	authService *security.AuthService,
	newsRepo *repository.NewsRepository,
	accountService *service.AccountService,
	secret string,
) *Handler {
	return &Handler{
		Logger:         logger,
		UserHandler:    NewUserHandler(logger, userService),
		AuthHandler:    NewAuthHandler(logger, authService),
		NewsHandler:    NewNewsHandler(logger, newsRepo),
		AccountHandler: NewAccountHandler(logger, accountService),
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
		// для всех авторизованных
		secure.GET("/balance", h.AccountHandler.GetBalance)
		secure.POST("/transfer", h.AccountHandler.Transfer)
		secure.GET("/transactions", h.AccountHandler.GetTransactions)

		// только для админов
		secure.POST("/deposit", security.AdminOnly(), h.AccountHandler.Deposit)
		secure.POST("/withdraw", security.AdminOnly(), h.AccountHandler.Withdraw)
	}
	// новости (public)
	router.GET("/news", h.NewsHandler.GetNews)
	router.POST("/news", h.NewsHandler.CreateNews)

	//балансы
	//router.GET("/balance", h.AccountHandler.GetBalance)
	//router.POST("/transfer", h.AccountHandler.Transfer)
	//router.GET("/transactions", h.AccountHandler.GetTransactions)
	//router.POST("/deposit", h.AccountHandler.Deposit)
	//router.POST("/withdraw", h.AccountHandler.Withdraw)

	// пользователи
	router.POST("/users", h.UserHandler.CreateUser)
	router.POST("/login", h.AuthHandler.Login)
	router.POST("/verify-otp", h.AuthHandler.VerifyOTP)
	router.GET("/users/:id", h.UserHandler.GetUserByID)
	router.GET("/users", h.UserHandler.GetUsers)
}
