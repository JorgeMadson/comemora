package services

import (
	"bytes"
	"celebrationhub/internal/core/domain"
	"celebrationhub/internal/core/ports"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"
)

type EventService struct {
	repo     ports.EventRepository
	notifier ports.Notifier
}

func NewEventService(repo ports.EventRepository, notifier ports.Notifier) *EventService {
	return &EventService{
		repo:     repo,
		notifier: notifier,
	}
}

func (s *EventService) CreateEvent(ctx context.Context, event *domain.Event) error {
	// Basic validation could go here
	if event.Name == "" {
		return fmt.Errorf("event name is required")
	}
	return s.repo.Save(ctx, event)
}

func (s *EventService) ListEvents(ctx context.Context) ([]domain.Event, error) {
	return s.repo.List(ctx)
}

func (s *EventService) ExportEvents(ctx context.Context) ([]byte, error) {
	events, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	b := &bytes.Buffer{}
	w := csv.NewWriter(b)

	// Header
	if err := w.Write([]string{"ID", "Name", "Day", "Month", "Year", "Type", "IsImportant", "Channel", "Contact"}); err != nil {
		return nil, err
	}

	for _, e := range events {
		record := []string{
			fmt.Sprintf("%d", e.ID),
			e.Name,
			fmt.Sprintf("%d", e.Day),
			fmt.Sprintf("%d", e.Month),
			fmt.Sprintf("%d", e.Year),
			string(e.Type),
			fmt.Sprintf("%t", e.IsImportant),
			string(e.PreferredChannel),
			e.ContactDestination,
		}
		if err := w.Write(record); err != nil {
			return nil, err
		}
	}
	w.Flush()
	return b.Bytes(), w.Error()
}

func (s *EventService) ImportEvents(ctx context.Context, data []byte) error {
	r := csv.NewReader(bytes.NewReader(data))

	// Skip header
	if _, err := r.Read(); err != nil { // Read header
		return err
	}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Simple mapping - assumes strict CSV format matching Export
		if len(record) < 9 {
			continue // skip invalid lines
		}

		day, _ := strconv.Atoi(record[2])
		month, _ := strconv.Atoi(record[3])
		year, _ := strconv.Atoi(record[4])
		isImp, _ := strconv.ParseBool(record[6])

		event := &domain.Event{
			Name:               record[1],
			Day:                day,
			Month:              month,
			Year:               year,
			Type:               domain.EventType(record[5]),
			IsImportant:        isImp,
			PreferredChannel:   domain.NotificationChannel(record[7]),
			ContactDestination: record[8],
		}

		if err := s.repo.Save(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

func (s *EventService) CheckAndNotify(ctx context.Context) error {
	now := time.Now()
	// 1. Check for events TODAY
	todayEvents, err := s.repo.FindByDate(ctx, now.Day(), int(now.Month()))
	if err != nil {
		return err
	}

	// 2. Check for IMPORTANT events coming up (e.g., in 3 days)
	future := now.AddDate(0, 0, 3)
	upcomingEvents, err := s.repo.FindByDate(ctx, future.Day(), int(future.Month()))
	if err != nil {
		return err
	}

	// Combine and notify
	// Note: In a real system you'd distinct or have better logic,
	// here we just send notifications.

	for _, e := range todayEvents {
		// Log or wrap context
		fmt.Printf("Processing Today Event: %s\n", e.Name)
		if err := s.notifier.Send(ctx, e); err != nil {
			fmt.Printf("Failed to notify for event %d: %v\n", e.ID, err)
			// don't break, try others
		}
	}

	for _, e := range upcomingEvents {
		if e.IsImportant {
			fmt.Printf("Processing Upcoming Important Event: %s\n", e.Name)
			// Maybe modify message to say "Upcoming in 3 days"
			if err := s.notifier.Send(ctx, e); err != nil {
				fmt.Printf("Failed to notify for upcoming event %d: %v\n", e.ID, err)
			}
		}
	}

	return nil
}
