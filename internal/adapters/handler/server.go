package handler

import (
	"celebrationhub/internal/core/ports"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewServer(service ports.Service, logger *log.Logger) http.Handler {
	mux := chi.NewMux()

	// Middleware
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	// Routes
	addRoutes(mux, service, logger)

	return mux
}
