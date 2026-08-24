package parser

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestFirecrawlClientIntegration(t *testing.T) {
	if os.Getenv("FIRECRAWL_INTEGRATION") != "1" {
		t.Skip("set FIRECRAWL_INTEGRATION=1 to run")
	}

	apiKey := os.Getenv("FIRECRAWL_KEY")
	if apiKey == "" {
		t.Fatal("FIRECRAWL_KEY is required for the integration test")
	}

	pageURL := os.Getenv("FIRECRAWL_TEST_URL")
	if pageURL == "" {
		pageURL = "https://go.dev/blog/pipelines"
	}

	client, err := NewFirecrawlClient(
		apiKey,
		defaultFirecrawlURL,
		&http.Client{Timeout: 120 * time.Second},
	)
	if err != nil {
		t.Fatalf("create Firecrawl client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	startedAt := time.Now()
	page, err := client.ParsePage(ctx, pageURL)
	elapsed := time.Since(startedAt)
	if err != nil {
		t.Fatalf("parse page after %s: %v", elapsed, err)
	}

	t.Logf("elapsed=%s", elapsed)
	t.Logf("title=%q", page.Title)
	t.Logf("tags=%v", page.Tags)
	t.Logf("content_length=%d", len(page.Content))

	if page.Title == "" {
		t.Error("title is empty")
	}
	if page.Content == "" {
		t.Error("content is empty")
	}
	if len(page.Tags) < 1 || len(page.Tags) > 3 {
		t.Errorf("tags = %v, want 1 to 3", page.Tags)
	}
}
