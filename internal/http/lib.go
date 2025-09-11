package http

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
)

type ShardCloudClient struct {
	HttpClient *ShardCloudHttpClient
}

func (c *ShardCloudClient) Me(ctx context.Context) (*MeResponse, error) {
	return GetJSON[*MeResponse](ctx, c.HttpClient, "/user/me")
}

func (c *ShardCloudClient) GetLogs(ctx context.Context, appID string) (<-chan string, error) {
	url := fmt.Sprintf("%s/apps/%s/logs", c.HttpClient.baseurl, appID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Authorization", "Bearer "+c.HttpClient.token)

	resp, err := c.HttpClient.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	logChan := make(chan string)

	go func() {
		defer close(logChan)
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "data:") {
				data := strings.TrimPrefix(line, "data:")
				select {
				case logChan <- data:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return logChan, nil
}

func (c *ShardCloudClient) GetAppsQuiet(ctx context.Context) ([]*AppResponse, error) {
	return GetJSON[[]*AppResponse](ctx, c.HttpClient, "/apps?q=true")
}
func (c *ShardCloudClient) GetApps(ctx context.Context) ([]*AppResponse, error) {
	return GetJSON[[]*AppResponse](ctx, c.HttpClient, "/apps")
}

func (c *ShardCloudClient) GetApp(ctx context.Context, appID string) (*AppResponse, error) {
	return GetJSON[*AppResponse](ctx, c.HttpClient, fmt.Sprintf("/apps/%s?q=true", appID))
}

func (c *ShardCloudClient) DeleteApp(ctx context.Context, appID string) error {
	_, err := c.HttpClient.Delete(fmt.Sprintf("/apps/%s", appID))
	return err
}

func (c *ShardCloudClient) CreateApp(ctx context.Context, zip []byte) (*IdResponse, error) {
	return c.sendAppRequest(ctx, "POST", "/apps", zip)
}

func (c *ShardCloudClient) GetBackups(ctx context.Context, appID string) (*BackupsResponse, error) {
	return GetJSON[*BackupsResponse](ctx, c.HttpClient, fmt.Sprintf("/project/%s/backups", appID))
}

func (c *ShardCloudClient) RestoreBackup(ctx context.Context, appID string, backupID string) error {
	_, err := PostJSON[any](ctx, c.HttpClient, fmt.Sprintf("/project/%s/backups/%s/restore", appID, backupID), nil)
	return err
}

func (c *ShardCloudClient) CreateBackup(ctx context.Context, appID string) (*BackupResponse, error) {
	return PostJSON[BackupResponse](ctx, c.HttpClient, fmt.Sprintf("/project/%s/backups", appID), nil)
}

func (c *ShardCloudClient) UpdateAppStatus(ctx context.Context, appID string, status string) error {
	_, err := c.HttpClient.Post(fmt.Sprintf("/apps/%s/status", appID), "application/json", map[string]string{"status": status})
	return err
}

func (c *ShardCloudClient) UpdateApp(ctx context.Context, appID string, zip []byte) (*IdResponse, error) {
	endpoint := fmt.Sprintf("/apps/%s/file", appID)
	return c.sendAppRequest(ctx, "PUT", endpoint, zip)
}

func (c *ShardCloudClient) sendAppRequest(ctx context.Context, method, endpoint string, zip []byte) (*IdResponse, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="project"; filename="app.zip"`)
	h.Set("Content-Type", "application/zip")

	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, fmt.Errorf("failed to create form part: %w", err)
	}

	_, err = part.Write(zip)
	if err != nil {
		return nil, fmt.Errorf("failed to write zip data: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	contentType := writer.FormDataContentType()

	var response *http.Response
	switch method {
	case "POST":
		response, err = c.HttpClient.Post(endpoint, contentType, &body)
	case "PUT":
		response, err = c.HttpClient.Put(endpoint, contentType, &body)
	default:
		return nil, fmt.Errorf("unsupported HTTP method: %s", method)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to send %s request: %w", method, err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}
		return nil, fmt.Errorf("HTTP request failed with status %d: %s", response.StatusCode, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var idResponse IdResponse
	err = json.Unmarshal(bodyBytes, &idResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &idResponse, nil
}

func NewShardCloudClient(token, baseURL string) *ShardCloudClient {
	return &ShardCloudClient{
		HttpClient: NewShardCloudHttpClient(token, baseURL),
	}
}
