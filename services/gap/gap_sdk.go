package gap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/goravel/framework/facades"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	AuthToken  string
	UserID     string
	HTTPClient *http.Client
}

func NewClient() *Client {
	baseURL := facades.Config().GetString("GAP_BASE_URL")
	authToken := facades.Config().GetString("GAP_AUTH_TOKEN")
	userID := facades.Config().GetString("GAP_USER_ID")

	return &Client{
		BaseURL:   baseURL,
		AuthToken: authToken,
		UserID:    userID,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) SendMessage(text string, roomID string) error {
	if roomID == "" {
		roomID = facades.Config().GetString("GAP_ROOM_ID")
	}

	payload := map[string]interface{}{
		"roomId": roomID,
		"text":   text,
	}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	url := c.BaseURL + "/api/v1/chat.postMessage"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", c.AuthToken)
	req.Header.Set("X-User-Id", c.UserID)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	return nil
}

func (c *Client) GetUserInfo(username string) (map[string]interface{}, error) {
	var query string
	query = "username=" + username
	url := fmt.Sprintf("%s/api/v1/users.info?%s", c.BaseURL, query)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", c.AuthToken)
	req.Header.Set("X-User-Id", c.UserID)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return result, nil
}

func (c *Client) GetUserId(username string) (string, error) {
	info, err := c.GetUserInfo(username)
	if err != nil {
		return "", err
	}

	userData, ok := info["user"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("user data not found or in wrong format")
	}

	userID, ok := userData["_id"].(string)
	if !ok {
		return "", fmt.Errorf("user _id not found or not a string")
	}

	return userID, nil
}

func (c *Client) SendDirectMessage(username, text string) error {
	payload := map[string]interface{}{
		"username": username,
	}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	url := c.BaseURL + "/api/v1/im.create"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", c.AuthToken)
	req.Header.Set("X-User-Id", c.UserID)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("im.create failed: %s", resp.Status)
	}

	var createResp struct {
		Room struct {
			ID string `json:"_id"`
		} `json:"room"`
		Success bool `json:"success"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		return fmt.Errorf("decode createResp: %w", err)
	}
	if !createResp.Success {
		return fmt.Errorf("im.create API returned success=false")
	}

	roomID := createResp.Room.ID
	if roomID == "" {
		return fmt.Errorf("roomId is empty from im.create response")
	}

	return c.SendMessage(text, roomID)
}
