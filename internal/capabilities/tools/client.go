package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	baseURL = "https://airquality.googleapis.com/v1"
)

// Client represents an Air Quality API client
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new Air Quality API client
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetCurrentConditions retrieves current air quality conditions
func (c *Client) GetCurrentConditions(req CurrentConditionsRequest) (*CurrentConditionsResponse, error) {
	url := fmt.Sprintf("%s/currentConditions:lookup?key=%s", baseURL, c.apiKey)

	var resp CurrentConditionsResponse
	if err := c.doPostRequest(url, req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetForecast retrieves air quality forecast
func (c *Client) GetForecast(req ForecastRequest) (*ForecastResponse, error) {
	url := fmt.Sprintf("%s/forecast:lookup?key=%s", baseURL, c.apiKey)

	var resp ForecastResponse
	if err := c.doPostRequest(url, req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetHistory retrieves historical air quality data
func (c *Client) GetHistory(req HistoryRequest) (*HistoryResponse, error) {
	url := fmt.Sprintf("%s/history:lookup?key=%s", baseURL, c.apiKey)

	var resp HistoryResponse
	if err := c.doPostRequest(url, req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetHeatmapTile retrieves a heatmap tile image
func (c *Client) GetHeatmapTile(mapType MapType, zoom, x, y int) ([]byte, error) {
	url := fmt.Sprintf("%s/mapTypes/%s/heatmapTiles/%d/%d/%d?key=%s",
		baseURL, mapType, zoom, x, y, c.apiKey)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get heatmap tile: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return data, nil
}

// doPostRequest performs a POST request with JSON payload
func (c *Client) doPostRequest(url string, reqBody interface{}, respBody interface{}) error {
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, respBody); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}
