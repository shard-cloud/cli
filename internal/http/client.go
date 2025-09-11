package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/http2"
)

type ShardCloudHttpClient struct {
	http.Client
	token   string
	baseurl string
}

func (c *ShardCloudHttpClient) SetToken(token string) {
	c.token = token
}

func (c *ShardCloudHttpClient) BuildURL(endpoint string) string {
	if c.baseurl == "" {
		return endpoint
	}
	if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") {
		return endpoint
	}

	endpoint = strings.TrimPrefix(endpoint, "/")

	return c.baseurl + "/" + endpoint
}

func (c *ShardCloudHttpClient) BuildHeaders() http.Header {
	return http.Header{
		"Authorization": []string{"Bearer " + c.token},
	}
}

func (c *ShardCloudHttpClient) BuildRequest(method, endpoint string, body interface{}, contentType string) (*http.Request, error) {
	fullURL := c.BuildURL(endpoint)
	var reqBody io.Reader
	if contentType == "application/json" {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewReader(jsonBody)
	} else {
		rawBody := body.(*bytes.Buffer)
		reqBody = rawBody
	}
	req, err := http.NewRequest(method, fullURL, reqBody)
	req.Header = c.BuildHeaders()
	req.Header.Set("Content-Type", contentType)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (c *ShardCloudHttpClient) Get(endpoint string) (*http.Response, error) {
	fullURL := c.BuildURL(endpoint)
	request, err := c.BuildRequest(http.MethodGet, fullURL, nil, "application/json")
	if err != nil {
		return nil, err
	}
	return c.Do(request)
}

func PostJSON[T any](ctx context.Context, c *ShardCloudHttpClient, endpoint string, body interface{}) (*T, error) {
	var result T
	response, err := c.Post(endpoint, "application/json", body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP error: %d %s", response.StatusCode, response.Status)
	}
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return &result, nil
}

func GetJSON[T any](ctx context.Context, c *ShardCloudHttpClient, endpoint string) (T, error) {
	var result T
	response, err := c.Get(endpoint)
	if err != nil {
		return result, err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return result, fmt.Errorf("HTTP error: %d %s", response.StatusCode, response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return result, fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return result, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return result, nil
}

func (c *ShardCloudHttpClient) Delete(endpoint string) (*http.Response, error) {
	fullURL := c.BuildURL(endpoint)
	request, err := c.BuildRequest(http.MethodDelete, fullURL, nil, "application/json")
	if err != nil {
		return nil, err
	}
	return c.Do(request)
}

func (c *ShardCloudHttpClient) Put(endpoint string, contentType string, body interface{}) (*http.Response, error) {
	fullURL := c.BuildURL(endpoint)
	request, err := c.BuildRequest(http.MethodPut, fullURL, body, contentType)
	if err != nil {
		return nil, err
	}

	return c.Do(request)
}

func (c *ShardCloudHttpClient) Do(request *http.Request) (*http.Response, error) {
	response, err := c.Client.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}
		return nil, errors.New(string(body))
	}
	return response, nil
}

func (c *ShardCloudHttpClient) Post(endpoint string, contentType string, body interface{}) (*http.Response, error) {
	fullURL := c.BuildURL(endpoint)
	request, err := c.BuildRequest(http.MethodPost, fullURL, body, contentType)

	if err != nil {
		return nil, err
	}
	return c.Do(request)
}

func NewShardCloudHttpClient(token, baseurl string) *ShardCloudHttpClient {
	transport := &http2.Transport{}
	return &ShardCloudHttpClient{
		Client: http.Client{
			//	Timeout:   10 * time.Second,
			Transport: transport,
		},
		token:   token,
		baseurl: baseurl,
	}
}
