package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"restobill/internal/config"
	"restobill/internal/db"
	"restobill/internal/logger"
	"restobill/internal/router"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.FatalF("failed to load config: %v", err)
	}

	// Connect to database
	dbClient, err := db.NewClient(cfg.DatabaseURL())
	if err != nil {
		logger.FatalF("failed to connect to database: %v", err)
	}
	defer dbClient.Close()

	// Setup router
	r := router.Setup(dbClient, cfg)

	// Start server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Run server in a goroutine
	go func() {
		logger.InfoF("RestoBill server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.FatalF("server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.FatalF("server forced to shutdown: %v", err)
	}

	logger.Info("server stopped")
}
