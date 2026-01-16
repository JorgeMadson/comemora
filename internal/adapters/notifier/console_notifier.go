package notifier

import (
	"celebrationhub/internal/core/domain"
	"context"
	"fmt"
	"log"
)

type ConsoleNotifier struct {
	logger *log.Logger
}

func NewConsoleNotifier(logger *log.Logger) *ConsoleNotifier {
	return &ConsoleNotifier{logger: logger}
}

func (n *ConsoleNotifier) Send(ctx context.Context, event domain.Event) error {
	msg := fmt.Sprintf("🎉 NOTIFICATION [%s] for %s: %s (Type: %s)", 
		event.PreferredChannel, 
		event.ContactDestination, 
		event.GetContent(), 
		event.Type,
	)
	
	n.logger.Println(msg)
	
	// Simulation of specific channel logic
	switch event.PreferredChannel {
	case domain.ChannelTeams:
		n.logger.Printf("[MOCK] Posting to Teams Webhook: %s", event.ContactDestination)
	case domain.ChannelWhatsApp:
		n.logger.Printf("[MOCK] Sending WhatsApp via API to: %s", event.ContactDestination)
	case domain.ChannelEmail:
		n.logger.Printf("[MOCK] Sending Email SMTP to: %s", event.ContactDestination)
	}
	
	return nil
}
