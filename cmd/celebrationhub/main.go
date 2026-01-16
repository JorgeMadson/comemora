package main

import (
	"celebrationhub/internal/adapters/handler"
	"celebrationhub/internal/adapters/notifier"
	"celebrationhub/internal/adapters/repository"
	"celebrationhub/internal/core/services"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
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

	logger := log.New(w, "[CelebrationHub] ", log.LstdFlags)

	// 1. Init Dependencies
	repo, err := repository.NewSQLiteRepository("celebration.db")
	if err != nil {
		return fmt.Errorf("failed to init db: %w", err)
	}

	notif := notifier.NewConsoleNotifier(logger)
	service := services.NewEventService(repo, notif)

	// 2. Init Server
	srvHandler := handler.NewServer(service, logger)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort("0.0.0.0", "8080"),
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
