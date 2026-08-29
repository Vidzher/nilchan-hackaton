package parser

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFirecrawlClientParsePageReturnsValidatedPage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Errorf("Authorization = %q, want %q", got, "Bearer test-key")
		}

		var request scrapeRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if request.URL != "https://example.com/article" {
			t.Errorf("URL = %q, want article URL", request.URL)
		}
		if len(request.Formats) != 1 || request.Formats[0] != "markdown" {
			t.Errorf("Formats = %v, want [markdown]", request.Formats)
		}
		if !request.OnlyMainContent {
			t.Error("OnlyMainContent = false, want true")
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"success": true,
			"data": {
				"markdown": "article content",
				"metadata": {
					"title": "Article title",
					"keywords": "Go, Concurrency"
				}
			}
		}`))
	}))
	defer server.Close()

	client, err := NewFirecrawlClient("test-key", server.URL, server.Client())
	if err != nil {
		t.Fatalf("NewFirecrawlClient() error = %v", err)
	}

	page, err := client.ParsePage(context.Background(), "https://example.com/article")
	if err != nil {
		t.Fatalf("ParsePage() error = %v", err)
	}
	if page.Content != "article content" {
		t.Errorf("Content = %q, want article content", page.Content)
	}
	if page.Title != "Article title" {
		t.Errorf("Title = %q, want Article title", page.Title)
	}
	if len(page.Tags) != 2 || page.Tags[0] != "go" || page.Tags[1] != "concurrency" {
		t.Errorf("Tags = %v, want [go concurrency]", page.Tags)
	}
}
