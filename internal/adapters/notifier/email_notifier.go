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

type EmailNotifier struct {
	apiKey string
	from   string
	client *http.Client
}

func NewEmailNotifier(apiKey, from string) *EmailNotifier {
	return &EmailNotifier{
		apiKey: apiKey,
		from:   from,
		client: &http.Client{},
	}
}

func (n *EmailNotifier) Send(ctx context.Context, event domain.Event) error {
	payload := map[string]any{
		"from":    n.from,
		"to":      []string{event.ContactDestination},
		"subject": fmt.Sprintf("Comemora: %s", event.Name),
		"text":    event.GetContent(),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("email: marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.resend.com/emails", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("email: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+n.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("email: send: %w", err)
	}
	defer func() {
		io.Copy(io.Discard, resp.Body) //nolint:errcheck
		resp.Body.Close()
	}()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("email: unexpected status %d", resp.StatusCode)
	}
	return nil
}
