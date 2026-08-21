package pparser

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const defaultFirecrawlURL = "https://api.firecrawl.dev/v2/scrape"

type FirecrawlClient struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

type scrapeRequest struct {
	URL string `json:"url"`
}

type scrapeResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Markdown string `json:"markdown"`
	} `json:"data"`
	Error string `json:"error"`
}

func NewFirecrawlClient(apiKey, baseURL string, client *http.Client) (*FirecrawlClient, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, fmt.Errorf("firecrawl api key is required")
	}
	if strings.TrimSpace(baseURL) == "" {
		baseURL = defaultFirecrawlURL
	}
	if client == nil {
		client = http.DefaultClient
	}
	return &FirecrawlClient{apiKey: apiKey, baseURL: baseURL, client: client}, nil
}

func (fc *FirecrawlClient) ParsePage(ctx context.Context, pageURL string) (string, error) {
	body, err := json.Marshal(scrapeRequest{URL: pageURL})
	if err != nil {
		return "", fmt.Errorf("marshal firecrawl request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fc.baseURL, strings.NewReader(string(body)))
	if err != nil {
		return "", fmt.Errorf("create firecrawl request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+fc.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := fc.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("firecrawl request: %w", err)
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read firecrawl response: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("firecrawl returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(responseBody)))
	}

	var result scrapeResponse

	if err := json.Unmarshal(responseBody, &result); err != nil {
		return "", fmt.Errorf("decode firecrawl response: %w", err)
	}
	if !result.Success {
		if result.Error == "" {
			result.Error = "unknown provider error"
		}
		return "", fmt.Errorf("firecrawl scrape failed: %s", result.Error)
	}

	return result.Data.Markdown, nil
}
