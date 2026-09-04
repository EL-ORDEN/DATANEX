package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Request holds the data needed for an outbound HTTP request.
type Request struct {
	Method      string
	URL         string
	Headers     map[string]string
	Body        any
	BearerToken string
}

// Response stores the result of an HTTP request, including timing and status.
type Response struct {
	StatusCode int
	Status     string
	Time       time.Duration
	Headers    http.Header
	Body       []byte
	Error      string
}

// Client is a small HTTP client wrapper for CLI usage.
type Client struct {
	HTTPClient *http.Client
}

// NewClient creates a client with defaults.
func NewClient() *Client {
	return &Client{HTTPClient: &http.Client{Timeout: 15 * time.Second}}
}

// Do executes an HTTP request.
func (c *Client) Do(req Request) (*Response, error) {
	if c == nil || c.HTTPClient == nil {
		return nil, fmt.Errorf("http client is not initialized")
	}
	if strings.TrimSpace(req.URL) == "" {
		return nil, fmt.Errorf("request URL is required")
	}
	if strings.TrimSpace(req.Method) == "" {
		req.Method = http.MethodGet
	}

	var body io.Reader
	if req.Body != nil {
		payload, err := json.Marshal(req.Body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		body = bytes.NewReader(payload)
	}

	httpReq, err := http.NewRequest(req.Method, req.URL, body)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	if req.Body != nil {
		httpReq.Header.Set("Content-Type", "application/json")
	}
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}
	if req.BearerToken != "" {
		httpReq.Header.Set("Authorization", "Bearer "+req.BearerToken)
	}

	start := time.Now()
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Time:       time.Since(start),
		Headers:    resp.Header,
		Body:       payload,
	}, nil
}
