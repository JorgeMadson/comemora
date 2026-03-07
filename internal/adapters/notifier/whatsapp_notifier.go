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

type WhatsAppNotifier struct {
	apiKey  string
	baseURL string
	from    string
	client  *http.Client
}

func NewWhatsAppNotifier(apiKey, baseURL, from string) *WhatsAppNotifier {
	return &WhatsAppNotifier{
		apiKey:  apiKey,
		baseURL: baseURL,
		from:    from,
		client:  &http.Client{},
	}
}

func (n *WhatsAppNotifier) Send(ctx context.Context, event domain.Event) error {
	payload := map[string]any{
		"from": n.from,
		"to":   event.ContactDestination,
		"content": map[string]string{
			"text": event.GetContent(),
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("whatsapp: marshal payload: %w", err)
	}

	url := fmt.Sprintf("https://%s.api.infobip.com/whatsapp/1/message/text", n.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("whatsapp: create request: %w", err)
	}
	req.Header.Set("Authorization", "App "+n.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("whatsapp: send: %w", err)
	}
	defer func() {
		io.Copy(io.Discard, resp.Body) //nolint:errcheck
		resp.Body.Close()
	}()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("whatsapp: unexpected status %d", resp.StatusCode)
	}
	return nil
}
