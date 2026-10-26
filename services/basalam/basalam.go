package basalam

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/goravel/framework/facades"
)

// Client is a thin Basalam OAuth (Authorization Code Flow) client used only to
// authenticate a user and read their identity. It does not persist tokens, since
// the application does not call Basalam APIs after login.
type Client struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Scope        string
	SsoURL       string
	TokenURL     string
	WhoamiURL    string
	HTTPClient   *http.Client
}

// BasalamUser is the identity returned by the /whoami endpoint.
type BasalamUser struct {
	ID     uint
	Name   string
	Mobile string
	HashID string
}

func NewClient() *Client {
	return &Client{
		ClientID:     facades.Config().GetString("basalam.client_id"),
		ClientSecret: facades.Config().GetString("basalam.client_secret"),
		RedirectURI:  facades.Config().GetString("basalam.redirect_uri"),
		Scope:        facades.Config().GetString("basalam.scope"),
		SsoURL:       facades.Config().GetString("basalam.sso_url"),
		TokenURL:     facades.Config().GetString("basalam.token_url"),
		WhoamiURL:    facades.Config().GetString("basalam.whoami_url"),
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// IsConfigured reports whether the minimum credentials are present.
func (c *Client) IsConfigured() bool {
	return c.ClientID != "" && c.ClientSecret != "" && c.RedirectURI != ""
}

// AuthorizationURL builds the Basalam SSO URL the user is redirected to.
func (c *Client) AuthorizationURL(state string) string {
	return fmt.Sprintf(
		"%s?client_id=%s&redirect_uri=%s&scope=%s&state=%s",
		c.SsoURL,
		url.QueryEscape(c.ClientID),
		c.RedirectURI,
		url.QueryEscape(c.Scope),
		url.QueryEscape(state),
	)
}

// ExchangeCode swaps the authorization code for an access token.
func (c *Client) ExchangeCode(code string) (string, error) {
	payload := map[string]any{
		"grant_type":    "authorization_code",
		"client_id":     c.ClientID,
		"client_secret": c.ClientSecret,
		"redirect_uri":  c.RedirectURI,
		"code":          code,
	}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", c.TokenURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token endpoint returned status: %s", resp.Status)
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("decode token response: %w", err)
	}
	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("empty access token in response")
	}

	return tokenResp.AccessToken, nil
}

// Whoami fetches the authenticated user's identity using the access token.
func (c *Client) Whoami(token string) (BasalamUser, error) {
	req, err := http.NewRequest("GET", c.WhoamiURL, nil)
	if err != nil {
		return BasalamUser{}, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return BasalamUser{}, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return BasalamUser{}, fmt.Errorf("whoami returned status: %s", resp.Status)
	}

	// id may be a JSON number or string; decode with UseNumber to stay safe.
	dec := json.NewDecoder(resp.Body)
	dec.UseNumber()
	var raw map[string]any
	if err := dec.Decode(&raw); err != nil {
		return BasalamUser{}, fmt.Errorf("decode whoami response: %w", err)
	}

	id, err := asUint(raw["id"])
	if err != nil {
		return BasalamUser{}, fmt.Errorf("parse user id: %w", err)
	}

	return BasalamUser{
		ID:     id,
		Name:   asString(raw["name"]),
		Mobile: asString(raw["mobile"]),
		HashID: asString(raw["hash_id"]),
	}, nil
}

func asString(v any) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

func asUint(v any) (uint, error) {
	s := asString(v)
	if s == "" {
		return 0, fmt.Errorf("empty value")
	}
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(n), nil
}
