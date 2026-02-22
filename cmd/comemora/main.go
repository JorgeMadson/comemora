package main

import (
	"comemora/internal/adapters/handler"
	"comemora/internal/adapters/notifier"
	"comemora/internal/adapters/repository"
	"comemora/internal/core/services"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, w io.Writer, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// Loads .env in dev; silently ignored in production where real env vars are set
	_ = godotenv.Load()

	logger := log.New(w, "[Comemora] ", log.LstdFlags)

	// 1. Init Dependencies
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "user"),
		getEnv("DB_PASSWORD", "password"),
		getEnv("DB_NAME", "comemora"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_SSLMODE", "disable"),
	)
	repo, err := repository.NewPostgresRepository(dsn)
	if err != nil {
		return fmt.Errorf("failed to init db: %w", err)
	}

	notif := notifier.NewConsoleNotifier(logger)
	service := services.NewEventService(repo, notif)

	// 2. Init Server
	serverPort := getEnv("SERVER_PORT", "8080")
	srvHandler := handler.NewServer(service, logger)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort("0.0.0.0", serverPort),
		Handler: srvHandler,
	}

	// 3. Start Server
	go func() {
		logger.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()

	// 4. Graceful Shutdown
	<-ctx.Done()
	logger.Println("shutting down gracefully, press Ctrl+C again to force")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("error shutting down http server: %w", err)
	}

	logger.Println("server stopped")
	return nil
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
