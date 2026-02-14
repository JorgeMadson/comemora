package handler

import (
	"comemora/internal/core/ports"
	"log"

	"github.com/go-chi/chi/v5"
)

func addRoutes(
	mux *chi.Mux,
	service ports.Service,
	logger *log.Logger,
) {
	mux.Route("/events", func(r chi.Router) {
		r.Post("/", handleCreateEvent(service))
		r.Get("/", handleListEvents(service))
		r.Get("/export", handleExportEvents(service))
		r.Post("/import", handleImportEvents(service))
	})

	mux.Get("/trigger-check", handleTriggerCheck(service))
}
