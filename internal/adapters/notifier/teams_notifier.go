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

type TeamsNotifier struct {
	webhookURL string
	client     *http.Client
}

func NewTeamsNotifier(webhookURL string) *TeamsNotifier {
	return &TeamsNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

func (n *TeamsNotifier) Send(ctx context.Context, event domain.Event) error {
	payload := map[string]string{
		"text": event.GetContent(),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("teams: marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.webhookURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("teams: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("teams: send: %w", err)
	}
	defer func() {
		io.Copy(io.Discard, resp.Body) //nolint:errcheck
		resp.Body.Close()
	}()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("teams: unexpected status %d", resp.StatusCode)
	}
	return nil
}
