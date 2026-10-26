package notify

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

// Send posts a simple {"text": ...} JSON payload to an incoming webhook URL
// (e.g. a GAP/Rocket.Chat-style hook). Best-effort; errors are ignored so a
// failed notification never blocks the user action.
func Send(webhookURL, text string) {
	if webhookURL == "" {
		return
	}

	payload, err := json.Marshal(map[string]string{"text": text})
	if err != nil {
		return
	}

	client := &http.Client{Timeout: 8 * time.Second}
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(payload))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	_ = resp.Body.Close()
}

// SendAsync sends the notification in the background.
func SendAsync(webhookURL, text string) {
	if webhookURL == "" {
		return
	}
	go Send(webhookURL, text)
}
