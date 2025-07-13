package http

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

// We'll use a package-level random source for jitter calculations
// As of Go 1.20, there's no need to seed the global rand source
var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

// HTTPClient defines the interface for HTTP clients
type HTTPClient interface {
	Get(url string) ([]byte, error)
	GetWithHeaders(url string, headers map[string]string) ([]byte, error)
	GetJSON(url string, v any) error
}

// Client is a wrapper around http.Client
type Client struct {
	client *http.Client
}

// NewClient creates a new Client
func NewClient() HTTPClient {
	return &Client{
		client: &http.Client{
			Timeout: 15 * time.Second, // Increased timeout
		},
	}
}

// Get performs a GET request to the specified URL
func (c *Client) Get(url string) ([]byte, error) {
	return c.GetWithHeaders(url, nil)
}

// GetWithHeaders performs a GET request to the specified URL with custom headers
func (c *Client) GetWithHeaders(url string, headers map[string]string) ([]byte, error) {
	// Create request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set default headers
	req.Header.Set("User-Agent", "GossipBot/1.0 (https://github.com/ducminhgd/gossip-bot; contact@example.com)")
	req.Header.Set("Accept", "application/json")

	// Set custom headers (will override defaults if same key)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Perform request with retry logic
	var resp *http.Response
	var lastErr error
	maxRetries := 3

	for i := range maxRetries {
		// Add a small delay between retries with some jitter
		if i > 0 {
			delay := time.Duration(500+rnd.Intn(500)) * time.Millisecond
			time.Sleep(delay)
		}

		resp, err = c.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("failed to perform request (attempt %d/%d): %w", i+1, maxRetries, err)
			continue
		}

		// If we got a response, break out of the retry loop
		break
	}

	// If all retries failed, return the last error
	if resp == nil {
		return nil, lastErr
	}

	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

// GetJSON performs a GET request to the specified URL and unmarshals the response into v
func (c *Client) GetJSON(url string, v any) error {
	body, err := c.Get(url)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}
