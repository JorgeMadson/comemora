package notifier

import (
	"bytes"
	"comemora/internal/core/domain"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type TelegramNotifier struct {
	botToken string
	client   *http.Client
}

func NewTelegramNotifier(botToken string) *TelegramNotifier {
	return &TelegramNotifier{
		botToken: botToken,
		client:   &http.Client{},
	}
}

func (n *TelegramNotifier) Send(ctx context.Context, event domain.Event) error {
	payload := map[string]any{
		"chat_id": event.ContactDestination,
		"text":    event.GetContent(),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("telegram: marshal payload: %w", err)
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", n.botToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("telegram: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("telegram: send: %w", err)
	}
	defer func() {
		io.Copy(io.Discard, resp.Body) //nolint:errcheck
		resp.Body.Close()
	}()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("telegram: unexpected status %d", resp.StatusCode)
	}
	return nil
}
