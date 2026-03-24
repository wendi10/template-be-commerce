package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql" // goose MySQL dialect driver
	"github.com/pressly/goose/v3"
	"github.com/template-be-commerce/config"
	"github.com/template-be-commerce/internal/domain"
	handler "github.com/template-be-commerce/internal/handler/http"
	"github.com/template-be-commerce/internal/payment"
	"github.com/template-be-commerce/internal/payment/doku"
	mysqlrepo "github.com/template-be-commerce/internal/repository/mysql"
	"github.com/template-be-commerce/internal/usecase"
	jwtpkg "github.com/template-be-commerce/pkg/jwt"
	"github.com/template-be-commerce/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// ─── Config ──────────────────────────────────────────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// ─── Logger ──────────────────────────────────────────────────────────────────
	logger.Init(cfg.App.LogLevel, cfg.App.Debug)
	defer logger.Sync()

	logger.Info("starting ecommerce api",
		zap.String("env", cfg.App.Env),
		zap.String("port", cfg.App.Port),
	)

	// ─── Database (GORM + MySQL) ───────────────────────────────────────────────
	gormDB, err := mysqlrepo.NewDB(cfg.Database)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}

	// Retrieve underlying *sql.DB for goose migrations and graceful close
	sqlDB, err := gormDB.DB()
	if err != nil {
		logger.Fatal("failed to get underlying sql.DB", zap.Error(err))
	}
	defer sqlDB.Close()

	logger.Info("database connection established")

	// ─── Migrations (goose) ────────────────────────────────────────────────────
	if err := runMigrations(sqlDB); err != nil {
		logger.Fatal("failed to run migrations", zap.Error(err))
	}

	// ─── Repositories ──────────────────────────────────────────────────────────
	userRepo := mysqlrepo.NewUserRepository(gormDB)
	customerRepo := mysqlrepo.NewCustomerRepository(gormDB)
	addressRepo := mysqlrepo.NewAddressRepository(gormDB)
	productRepo := mysqlrepo.NewProductRepository(gormDB)
	productImageRepo := mysqlrepo.NewProductImageRepository(gormDB)
	categoryRepo := mysqlrepo.NewCategoryRepository(gormDB)
	bannerRepo := mysqlrepo.NewBannerRepository(gormDB)
	promoRepo := mysqlrepo.NewPromoRepository(gormDB)
	cartRepo := mysqlrepo.NewCartRepository(gormDB)
	orderRepo := mysqlrepo.NewOrderRepository(gormDB)
	paymentRepo := mysqlrepo.NewPaymentRepository(gormDB)

	// ─── JWT Manager ───────────────────────────────────────────────────────────
	jwtManager := jwtpkg.NewManager(
		cfg.JWT.AccessTokenSecret,
		cfg.JWT.RefreshTokenSecret,
		cfg.JWT.AccessTokenTTL,
		cfg.JWT.RefreshTokenTTL,
	)

	// ─── Payment Gateways ──────────────────────────────────────────────────────
	gateways := map[domain.PaymentProvider]payment.Gateway{
		domain.PaymentProviderDoku: doku.NewProvider(
			cfg.Payment.DokuClientID,
			cfg.Payment.DokuSecretKey,
			cfg.Payment.DokuBaseURL,
		),
	}

	// ─── Use Cases ─────────────────────────────────────────────────────────────
	promoUC := usecase.NewPromoUseCase(promoRepo)
	authUC := usecase.NewAuthUseCase(userRepo, customerRepo, jwtManager)
	customerUC := usecase.NewCustomerUseCase(customerRepo, addressRepo)
	productUC := usecase.NewProductUseCase(productRepo, productImageRepo, categoryRepo)
	bannerUC := usecase.NewBannerUseCase(bannerRepo)
	cartUC := usecase.NewCartUseCase(cartRepo, customerRepo, productRepo, promoUC)
	orderUC := usecase.NewOrderUseCase(orderRepo, cartRepo, customerRepo, productRepo, promoRepo, addressRepo, promoUC)
	paymentUC := usecase.NewPaymentUseCase(paymentRepo, orderRepo, customerRepo, addressRepo, gateways, cfg.Payment.CallbackURL)

	// ─── Handlers ──────────────────────────────────────────────────────────────
	handlers := handler.Handlers{
		Auth:     handler.NewAuthHandler(authUC),
		Customer: handler.NewCustomerHandler(customerUC),
		Product:  handler.NewProductHandler(productUC),
		Banner:   handler.NewBannerHandler(bannerUC),
		Promo:    handler.NewPromoHandler(promoUC),
		Cart:     handler.NewCartHandler(cartUC),
		Order:    handler.NewOrderHandler(orderUC),
		Payment:  handler.NewPaymentHandler(paymentUC),
		Admin:    handler.NewAdminHandler(customerUC),
	}

	// ─── Router ────────────────────────────────────────────────────────────────
	router := handler.NewRouter(handlers, jwtManager)

	// ─── HTTP Server ───────────────────────────────────────────────────────────
	srv := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("http server listening", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("server error", zap.Error(err))
		}
	}()

	// ─── Graceful Shutdown ─────────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server forced shutdown", zap.Error(err))
	}
	logger.Info("server stopped")
}

// runMigrations runs all pending goose migrations against the given *sql.DB.
func runMigrations(db *sql.DB) error {
	if err := goose.SetDialect("mysql"); err != nil {
		return fmt.Errorf("goose set dialect: %w", err)
	}
	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}
	return nil
}
