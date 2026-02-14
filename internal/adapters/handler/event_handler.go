package handler

import (
	"comemora/internal/core/domain"
	"comemora/internal/core/ports"
	"io"
	"net/http"
)

func handleCreateEvent(service ports.Service) http.HandlerFunc {
	type request struct {
		Name               string                     `json:"name"`
		Day                int                        `json:"day"`
		Month              int                        `json:"month"`
		Year               int                        `json:"year"`
		Type               domain.EventType           `json:"type"`
		IsImportant        bool                       `json:"is_important"`
		PreferredChannel   domain.NotificationChannel `json:"preferred_channel"`
		ContactDestination string                     `json:"contact_destination"`
		CustomMessage      string                     `json:"custom_message"`
	}
	type response struct {
		ID uint `json:"id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req, err := decode[request](r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		event := &domain.Event{
			Name:               req.Name,
			Day:                req.Day,
			Month:              req.Month,
			Year:               req.Year,
			Type:               req.Type,
			IsImportant:        req.IsImportant,
			PreferredChannel:   req.PreferredChannel,
			ContactDestination: req.ContactDestination,
			CustomMessage:      req.CustomMessage,
		}

		if err := service.CreateEvent(r.Context(), event); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		encode(w, r, http.StatusCreated, response{ID: event.ID})
	}
}

func handleListEvents(service ports.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		events, err := service.ListEvents(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		encode(w, r, http.StatusOK, events)
	}
}

func handleExportEvents(service ports.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := service.ExportEvents(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=events.csv")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func handleImportEvents(service ports.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if err := service.ImportEvents(r.Context(), data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Imported successfully"))
	}
}

func handleTriggerCheck(service ports.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := service.CheckAndNotify(r.Context()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Check and Notify executed"))
	}
}
