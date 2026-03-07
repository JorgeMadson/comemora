package notifier

import (
	"comemora/internal/core/domain"
	"comemora/internal/core/ports"
	"context"
)

type MultiNotifier struct {
	channels map[domain.NotificationChannel]ports.Notifier
	fallback ports.Notifier
}

func NewMultiNotifier(fallback ports.Notifier) *MultiNotifier {
	return &MultiNotifier{
		channels: make(map[domain.NotificationChannel]ports.Notifier),
		fallback: fallback,
	}
}

func (m *MultiNotifier) Register(channel domain.NotificationChannel, n ports.Notifier) {
	m.channels[channel] = n
}

func (m *MultiNotifier) Send(ctx context.Context, event domain.Event) error {
	n, ok := m.channels[event.PreferredChannel]
	if !ok {
		return m.fallback.Send(ctx, event)
	}
	return n.Send(ctx, event)
}
