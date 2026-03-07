package services

import (
	"bytes"
	"comemora/internal/core/domain"
	"comemora/internal/core/ports"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
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
	if err := event.Validate(); err != nil {
		return err
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

	// Skip header row (line 1)
	if _, err := r.Read(); err != nil {
		return err
	}

	// line starts at 2 because line 1 is the header
	line := 2
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("line %d: %w", line, err)
		}
		line++

		if len(record) < 9 {
			return fmt.Errorf("line %d: expected 9 columns, got %d", line, len(record))
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

		if err := event.Validate(); err != nil {
			return fmt.Errorf("line %d: %w", line, err)
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
		log.Printf("Processing today event: %s", e.Name)
		if err := s.notifier.Send(ctx, e); err != nil {
			log.Printf("Failed to notify for event %d: %v", e.ID, err)
		}
	}

	for _, e := range upcomingEvents {
		if e.IsImportant {
			log.Printf("Processing upcoming important event: %s", e.Name)
			if err := s.notifier.Send(ctx, e); err != nil {
				log.Printf("Failed to notify for upcoming event %d: %v", e.ID, err)
			}
		}
	}

	return nil
}
