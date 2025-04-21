package victorSDK

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"victorgo/victorSDK/pkg/routes"
)

func NewClient(options *ClientOptions) *Client {
	if options == nil {
		options = &ClientOptions{
			Host:            "localhost",
			Port:            "7007",
			AutoStartDaemon: true,
		}
	}

	if options.Host == "" {
		options.Host = "localhost"
	}

	if options.Port == "" {
		options.Port = "7007"
	}

	baseURL := ""
	isLocal := options.Host == "localhost" || options.Host == "127.0.0.1"
	if isLocal {
		baseURL = fmt.Sprintf("http://%s:%s", options.Host, options.Port)
	} else {
		baseURL = fmt.Sprintf("https://%s:%s", options.Host, options.Port)
	}

	client := &Client{
		HttpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		BaseURL: baseURL,
		IsLocal: isLocal,
	}

	return client
}

func (c *Client) CreateIndex(input *CreateIndexCommandInput) (*CreateIndexCommandOutput, error) {
	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create index request: %w", err)
	}

	fmt.Println("URL "+c.BaseURL+fmt.Sprintf(routes.CreateIndex, input.IndexName), input.IndexName)

	req, err := http.NewRequest(http.MethodPost, c.BaseURL+fmt.Sprintf(routes.CreateIndex, input.IndexName), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	fmt.Println("RESP" + resp.Status)

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			return nil, fmt.Errorf("API error: %s", errorResp["message"])
		}
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode)
	}

	var output CreateIndexCommandOutput
	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &output, nil
}

func (c *Client) InsertVector(input *InsertVectorCommandInput) (*InsertVectorCommandOutput, error) {
	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal insert vector request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.BaseURL+fmt.Sprintf(routes.InsertVector, input.IndexName), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errorResp map[string]string

		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil && errorResp["message"] != "" {
			return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, errorResp["message"])
		}
	}

	var output InsertVectorCommandOutput
	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &output, nil
}

func (c *Client) SearchVector(input *SearchVectorCommandInput) (*SearchCommandOutput, error) {
	vectorValues := make([]string, len(input.Vector))
	for i, v := range input.Vector {
		vectorValues[i] = fmt.Sprintf("%f", v)
	}
	vectorStr := strings.Join(vectorValues, ",")

	req, err := http.NewRequest(http.MethodGet, c.BaseURL+fmt.Sprintf(routes.SearchVector, input.IndexName)+fmt.Sprintf("?vector=%v&k=%v", vectorStr, input.TopK), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			return nil, fmt.Errorf("API error: %s", errorResp["message"])
		}
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode)
	}

	var output SearchCommandOutput
	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &output, nil
}

func (c *Client) DeleteVector(input *DeleteVectorCommandInput) (*DeleteVectorCommandOutput, error) {
	req, err := http.NewRequest(http.MethodDelete, c.BaseURL+fmt.Sprintf(routes.DeleteVector, input.IndexName, input.VectorID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			return nil, fmt.Errorf("API error: %s", errorResp["message"])
		}
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode)
	}

	var output DeleteVectorCommandOutput
	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &output, nil
}
