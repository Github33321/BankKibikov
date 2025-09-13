package main

import (
	"BankKibikov/internal/chat"
	"BankKibikov/internal/db"
	"BankKibikov/internal/handler"
	"BankKibikov/internal/logger"
	"BankKibikov/internal/repository"
	"BankKibikov/internal/security"
	"BankKibikov/internal/service"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`

	Database struct {
		DSN string `yaml:"dsn"`
	} `yaml:"database"`
}

func loadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func main() {
	cfg, err := loadConfig("configs/config.yaml")
	if err != nil {
		panic(err)
	}

	log, err := logger.New()
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	// БД
	pool, err := db.NewPool(cfg.Database.DSN)
	if err != nil {
		log.Fatal("failed to connect db", zap.Error(err))
	}
	defer pool.Close()

	// Репозитории
	userRepo := repository.NewUserRepository(pool)
	accountRepo := repository.NewAccountRepository(pool)
	newsRepo := repository.NewNewsRepository(pool)
	txRepo := repository.NewTransactionRepository(pool)
	loanRepo := repository.NewLoanRepository(pool)

	// Сервисы
	userService := service.NewUserService(userRepo, accountRepo)
	accountService := service.NewAccountService(accountRepo, txRepo)
	loanService := service.NewLoanService(loanRepo, accountRepo)

	// Security
	authService := security.NewAuthService(userRepo)

	// Chat hub
	chatHub := chat.NewHub()
	go chatHub.Run()

	// Хендлеры
	h := handler.NewHandler(
		log,
		userService,
		authService,
		newsRepo,
		accountService,
		loanService,
		chatHub, // ✅ прокидываем чат
		"",
	)

	// HTTP сервер
	r := gin.Default()
	h.InitRoutes(r)

	srv := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: r,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info("starting server", zap.String("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server error", zap.Error(err))
		}
	}()

	<-quit
	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown", zap.Error(err))
	}

	log.Info("server exited gracefully")
}
