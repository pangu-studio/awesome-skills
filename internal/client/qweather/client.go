package qweather

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

// ClientOption is a functional option for configuring a Client.
type ClientOption func(*Client)

// WithLogger sets the logger used for diagnostic output.
// All diagnostic messages are written at Debug level to avoid polluting normal output.
func WithLogger(logger *slog.Logger) ClientOption {
	return func(c *Client) {
		c.logger = logger
	}
}

// Client is the QWeather API client
type Client struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
	logger     *slog.Logger
}

// NewClient creates a new QWeather API client.
// Optional ClientOptions can be passed to customise the client (e.g. WithLogger).
func NewClient(apiKey, baseURL string, opts ...ClientOption) *Client {
	if baseURL == "" {
		baseURL = "https://devapi.qweather.com"
	} else if baseURL[0:4] != "http" {
		baseURL = "https://" + baseURL
	}

	c := &Client{
		APIKey:  apiKey,
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: slog.Default(),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// doRequest performs an HTTP request to the QWeather API
func (c *Client) doRequest(ctx context.Context, endpoint string, params url.Values, result interface{}) error {
	// Build URL
	u, err := url.Parse(c.BaseURL + endpoint)
	if err != nil {
		return fmt.Errorf("parse URL: %w", err)
	}

	// Add query parameters
	if params == nil {
		params = url.Values{}
	}
	u.RawQuery = params.Encode()

	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	c.logger.Debug("sending request", "url", req.URL.String())

	// Set API key in header for better security
	req.Header.Set("X-QW-Api-Key", c.APIKey)

	// Perform request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("perform request: %w", err)
	}
	defer resp.Body.Close()

	// Handle gzip-compressed responses
	var bodyReader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return fmt.Errorf("read gzip response: %w", err)
		}
		defer gzipReader.Close()
		bodyReader = gzipReader
	}

	// Read response body
	body, err := io.ReadAll(bodyReader)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		bodyStr := string(body)
		if len(bodyStr) > 500 {
			bodyStr = bodyStr[:500] + "... (truncated)"
		}
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, bodyStr)
	}

	// Parse JSON response
	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}

	return nil
}
