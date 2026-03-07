package main

import (
	"comemora/internal/adapters/handler"
	"comemora/internal/adapters/notifier"
	"comemora/internal/adapters/repository"
	"comemora/internal/core/domain"
	"comemora/internal/core/services"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
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
	// Prefer DATABASE_URL (Railway standard) over individual DB_* vars
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			getEnv("DB_HOST", "localhost"),
			getEnv("DB_USER", "user"),
			getEnv("DB_PASSWORD", "password"),
			getEnv("DB_NAME", "comemora"),
			getEnv("DB_PORT", "5432"),
			getEnv("DB_SSLMODE", "disable"),
		)
	} else if strings.HasPrefix(dsn, "postgresql://") || strings.HasPrefix(dsn, "postgres://") {
		u, err := url.Parse(dsn)
		if err != nil {
			return fmt.Errorf("invalid DATABASE_URL: %w", err)
		}
		password, _ := u.User.Password()
		port := u.Port()
		if port == "" {
			port = "5432"
		}
		sslmode := u.Query().Get("sslmode")
		if sslmode == "" {
			sslmode = "disable"
		}
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			u.Hostname(), u.User.Username(), password,
			strings.TrimPrefix(u.Path, "/"), port, sslmode,
		)
	}
	repo, err := repository.NewPostgresRepository(dsn)
	if err != nil {
		return fmt.Errorf("failed to init db: %w", err)
	}

	console := notifier.NewConsoleNotifier(logger)
	multi := notifier.NewMultiNotifier(console)

	if key := os.Getenv("RESEND_API_KEY"); key != "" {
		multi.Register(domain.ChannelEmail, notifier.NewEmailNotifier(key, getEnv("RESEND_FROM", "Comemora <noreply@comemora.app>")))
		logger.Println("notifier: Email (Resend) enabled")
	}
	if k, b, f := os.Getenv("INFOBIP_API_KEY"), os.Getenv("INFOBIP_BASE_URL"), os.Getenv("INFOBIP_FROM_NUMBER"); k != "" && b != "" && f != "" {
		multi.Register(domain.ChannelWhatsApp, notifier.NewWhatsAppNotifier(k, b, f))
		logger.Println("notifier: WhatsApp (Infobip) enabled")
	}
	if url := os.Getenv("TEAMS_WEBHOOK_URL"); url != "" {
		multi.Register(domain.ChannelTeams, notifier.NewTeamsNotifier(url))
		logger.Println("notifier: Teams enabled")
	}
	if token := os.Getenv("TELEGRAM_BOT_TOKEN"); token != "" {
		multi.Register(domain.ChannelTelegram, notifier.NewTelegramNotifier(token))
		logger.Println("notifier: Telegram enabled")
	}
	if url := os.Getenv("DISCORD_WEBHOOK_URL"); url != "" {
		multi.Register(domain.ChannelDiscord, notifier.NewDiscordNotifier(url))
		logger.Println("notifier: Discord enabled")
	}

	service := services.NewEventService(repo, multi)

	// 2. Init Server
	// Railway sets PORT; SERVER_PORT as fallback for local dev
	serverPort := getEnv("PORT", getEnv("SERVER_PORT", "8080"))
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
