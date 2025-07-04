package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	BankingHandler "github.com/rohanchauhan02/internal-transfer/domain/banking/delivery/https"
	BankingRepository "github.com/rohanchauhan02/internal-transfer/domain/banking/repository"
	BankingUsecase "github.com/rohanchauhan02/internal-transfer/domain/banking/usecase"
	HealthzHandler "github.com/rohanchauhan02/internal-transfer/domain/health/delivery/https"
	HealthzRepository "github.com/rohanchauhan02/internal-transfer/domain/health/repository"
	HealthzUsecase "github.com/rohanchauhan02/internal-transfer/domain/health/usecase"
	"github.com/rohanchauhan02/internal-transfer/utils"

	"github.com/rohanchauhan02/internal-transfer/models"
	"github.com/rohanchauhan02/internal-transfer/pkg/config"
	"github.com/rohanchauhan02/internal-transfer/pkg/ctx"
	"github.com/rohanchauhan02/internal-transfer/pkg/database"

	CustomMiddileware "github.com/rohanchauhan02/internal-transfer/pkg/middleware"
)

func main() {
	e := echo.New()

	// Load configuration
	cnf := config.NewImmutableConfigs()

	// Initialize PostgreSQL client
	postgresClient := database.NewPostgres(cnf)
	db, err := postgresClient.InitClient(context.Background())
	if err != nil {
		log.Panicf("Failed to initialize database: %s ", err.Error())
	}
	// Auto migrate models
	if err := db.AutoMigrate(
		&models.Account{},
		&models.Transaction{},
	); err != nil {
		log.Panicf("Failed to auto migrate models: %s ", err.Error())
	}

	// add request ID middleware
	e.Use(CustomMiddileware.MiddlewareRequestID())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Gzip())
	e.Use(middleware.CORS())

	// Middleware to inject dependencies into the request context
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			customCtx := &ctx.CustomApplicationContext{
				Context:    c,
				PostgresDB: db,
			}
			return next(customCtx)
		}
	})

	// Set validator globally
	validator := utils.DefaultValidator()
	e.Validator = validator

	// Set up repositories for subdomains
	healthzRepo := HealthzRepository.NewHealthRepository(db)
	bankingRepo := BankingRepository.NewBankingRepository(db)

	// Set up use cases for subdomains
	healthzUsecase := HealthzUsecase.NewHealthUsecase(healthzRepo)
	bankingUsecase := BankingUsecase.NewBankingUsecase(bankingRepo)

	// Set up handlers for subdomains
	HealthzHandler.NewHealthHandler(e, healthzUsecase)
	BankingHandler.NewBankingHandler(e, bankingUsecase)

	// Start server in a separate goroutine
	serverAddr := fmt.Sprintf(":%d", cnf.GetPort())
	go func() {
		if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server shutdown unexpectedly: %v", err)
		}
	}()

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Errorf("Server forced to shutdown: %v", err)
	}
	log.Info("Server exited properly.")
}
