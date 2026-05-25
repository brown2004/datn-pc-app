package ipc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func NewDefaultClient(timeout time.Duration) *Client {
	return NewClient(DefaultBaseURL, timeout)
}

func (c *Client) GetStatus(ctx context.Context) (Status, error) {
	var status Status
	if err := c.doJSON(ctx, http.MethodGet, "/status", nil, &status); err != nil {
		return Status{}, err
	}

	return status, nil
}

func (c *Client) GetIdentity(ctx context.Context) (Identity, error) {
	var identity Identity
	if err := c.doJSON(ctx, http.MethodGet, "/identity", nil, &identity); err != nil {
		return Identity{}, err
	}

	return identity, nil
}

func (c *Client) GetCredentialsStatus(ctx context.Context) (CredentialsStatus, error) {
	var status CredentialsStatus
	if err := c.doJSON(ctx, http.MethodGet, "/credentials/status", nil, &status); err != nil {
		return CredentialsStatus{}, err
	}

	return status, nil
}

func (c *Client) SetProtection(ctx context.Context, enabled bool) (Status, error) {
	var status Status
	if err := c.doJSON(ctx, http.MethodPatch, "/protection", ProtectionRequest{Enabled: enabled}, &status); err != nil {
		return Status{}, err
	}

	return status, nil
}

func (c *Client) SetCredentials(ctx context.Context, pcAgentID string, agentSecret string) (Status, error) {
	var status Status
	if err := c.doJSON(ctx, http.MethodPost, "/credentials", CredentialsRequest{
		PCAgentID:   pcAgentID,
		AgentSecret: agentSecret,
	}, &status); err != nil {
		return Status{}, err
	}

	return status, nil
}

func (c *Client) doJSON(ctx context.Context, method string, path string, body any, target any) error {
	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reader)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr ErrorResponse
		_ = json.NewDecoder(resp.Body).Decode(&apiErr)
		if apiErr.Error != "" {
			return fmt.Errorf("ipc %d: %s", resp.StatusCode, apiErr.Error)
		}

		return fmt.Errorf("ipc returned %d", resp.StatusCode)
	}

	if target == nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(target)
}
