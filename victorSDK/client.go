package victorSDK

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
	"victorgo/daemon/cmd/http_daemon"
	"victorgo/daemon/pkg/routes"
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

	baseURL := fmt.Sprintf("http://%s:%s", options.Host, options.Port)
	isLocal := options.Host == "localhost" || options.Host == "127.0.0.1"

	client := &Client{
		HttpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		BaseURL: baseURL,
		IsLocal: isLocal,
	}

	if isLocal && options.AutoStartDaemon {

		if !isPortInUse(options.Port) {
			port, _ := strconv.Atoi(options.Port)
			daemonServer := http_daemon.NewServer(port)
			go daemonServer.Start()
			client.Daemon = daemonServer

			time.Sleep(100 * time.Millisecond)
		}
	}

	return client
}

func (c *Client) Close() error {
	if c.IsLocal && c.Daemon != nil {
		return c.Daemon.Stop()
	}
	return nil
}

func isPortInUse(port string) bool {
	conn, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return true // El puerto está en uso
	}
	conn.Close()
	return false
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

	// bodyBytes, _ := io.ReadAll(resp.Body)
	// fmt.Println("Respuesta JSON sin procesar:", string(bodyBytes))
	// // Restaurar el body para la decodificación
	// resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

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
