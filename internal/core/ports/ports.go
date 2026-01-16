package ports

import (
	"celebrationhub/internal/core/domain"
	"context"
)

type EventRepository interface {
	Save(ctx context.Context, event *domain.Event) error
	FindByID(ctx context.Context, id uint) (*domain.Event, error)
	List(ctx context.Context) ([]domain.Event, error)
	Delete(ctx context.Context, id uint) error
	FindByDate(ctx context.Context, day, month int) ([]domain.Event, error)
}

type Notifier interface {
	Send(ctx context.Context, event domain.Event) error
}

type Service interface {
	CreateEvent(ctx context.Context, event *domain.Event) error
	ListEvents(ctx context.Context) ([]domain.Event, error)
	ExportEvents(ctx context.Context) ([]byte, error)     // Returns CSV bytes
	ImportEvents(ctx context.Context, data []byte) error // Parses CSV bytes
	CheckAndNotify(ctx context.Context) error
}
