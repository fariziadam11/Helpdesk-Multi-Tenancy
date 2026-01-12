package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"werk-ticketing/internal/auth"
	"werk-ticketing/internal/config"
	"werk-ticketing/internal/constants"
	"werk-ticketing/internal/database"
	"werk-ticketing/internal/invgate"
	"werk-ticketing/internal/router"
	"werk-ticketing/internal/ticket"
	"werk-ticketing/internal/user"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	// Set Gin mode from configuration
	ginMode := cfg.GinMode
	if ginMode == "" {
		ginMode = gin.DebugMode
	}
	gin.SetMode(ginMode)

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("database error: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("database pooling error: %v", err)
	}
	defer sqlDB.Close()

	// Auto migrate all models
	// This ensures all tables are created/updated when the application starts
	if err := db.AutoMigrate(
		&user.User{}, // Users table
	); err != nil {
		log.Fatalf("auto migrate error: %v", err)
	}

	// Configure logger based on environment
	logger := configureLogger(cfg)

	// Initialize services
	invgateClient := invgate.NewService(cfg)
	userRepo := user.NewRepository(db)
	ticketService := ticket.NewService(invgateClient, userRepo, logger)
	ticketHandler := ticket.NewHandler(ticketService)

	authService := auth.NewService(
		userRepo,
		invgateClient,
		cfg.JWTSecret,
		logger,
		cfg.ArmMadaCompanyID,
		cfg.ArmMadaGroupID,
		cfg.ArmMadaLocationID,
	)
	authHandler := auth.NewHandler(authService)

	// Setup router
	appRouter := router.NewRouter(authHandler, ticketHandler, authService, logger)
	ginRouter := appRouter.SetupRoutes()

	// Create HTTP server with timeouts
	addr := ":" + cfg.ServerPort
	srv := &http.Server{
		Addr:         addr,
		Handler:      ginRouter,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.WithFields(logrus.Fields{
			"port":      cfg.ServerPort,
			"gin_mode":  gin.Mode(),
			"log_level": logger.GetLevel(),
			"env":       cfg.AppEnv,
		}).Info("backend server starting")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Error("server failed to start")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	// SIGINT (Ctrl+C) and SIGTERM (kill command)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down server...")

	// The context is used to inform the server it has 30 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), constants.GracefulShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.WithError(err).Error("server forced to shutdown")
	} else {
		logger.Info("server exited gracefully")
	}
}

// configureLogger sets up the logger based on configuration
func configureLogger(cfg *config.Config) *logrus.Logger {
	logger := logrus.New()

	// Set log level from configuration
	logLevel := cfg.LogLevel
	if logLevel == "" {
		logLevel = "info"
	}
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
		log.Printf("Invalid log level '%s', defaulting to 'info'", logLevel)
	}
	logger.SetLevel(level)

	// Set formatter based on configuration
	if cfg.LogFormat == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		// Use text formatter with colors only in non-production
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			ForceColors:     cfg.AppEnv != "production",
			TimestampFormat: time.RFC3339,
		})
	}

	return logger
}
