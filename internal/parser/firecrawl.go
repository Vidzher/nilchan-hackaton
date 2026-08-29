package parser

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

type Page struct {
	Title   string
	Content string
	Tags    []string
}

type scrapeRequest struct {
	URL             string   `json:"url"`
	Formats         []string `json:"formats"`
	OnlyMainContent bool     `json:"onlyMainContent"`
}

type scrapeResponse struct {
	Success bool       `json:"success"`
	Data    scrapeData `json:"data"`
	Error   string     `json:"error"`
}

type scrapeData struct {
	Markdown string         `json:"markdown"`
	Metadata scrapeMetadata `json:"metadata"`
}

type scrapeMetadata struct {
	Title    string `json:"title"`
	Keywords string `json:"keywords"`
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

func (fc *FirecrawlClient) ParsePage(ctx context.Context, pageURL string) (Page, error) {
	body, err := json.Marshal(scrapeRequest{
		URL:             pageURL,
		Formats:         []string{"markdown"},
		OnlyMainContent: true,
	})
	if err != nil {
		return Page{}, fmt.Errorf("marshal firecrawl request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fc.baseURL, strings.NewReader(string(body)))
	if err != nil {
		return Page{}, fmt.Errorf("create firecrawl request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+fc.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := fc.client.Do(req)
	if err != nil {
		return Page{}, fmt.Errorf("firecrawl request: %w", err)
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Page{}, fmt.Errorf("read firecrawl response: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return Page{}, fmt.Errorf("firecrawl returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(responseBody)))
	}

	var result scrapeResponse

	if err := json.Unmarshal(responseBody, &result); err != nil {
		return Page{}, fmt.Errorf("decode firecrawl response: %w", err)
	}
	if !result.Success {
		if result.Error == "" {
			result.Error = "unknown provider error"
		}
		return Page{}, fmt.Errorf("firecrawl scrape failed: %s", result.Error)
	}

	page := Page{
		Title:   strings.TrimSpace(result.Data.Metadata.Title),
		Content: strings.TrimSpace(result.Data.Markdown),
		Tags:    extractTags(pageURL, result.Data.Metadata.Keywords),
	}
	if page.Title == "" {
		return Page{}, fmt.Errorf("invalid extracted page: title is empty")
	}
	if page.Content == "" {
		return Page{}, fmt.Errorf("invalid extracted page: content is empty")
	}
	if len(page.Tags) == 0 {
		return Page{}, fmt.Errorf("invalid extracted page: tags are empty")
	}

	return page, nil
}
