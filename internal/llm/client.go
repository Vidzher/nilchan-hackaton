package llm

import "context"

type CompletionRequest struct {
	SystemPrompt string
	UserPrompt   string
}

type Client interface {
	Generate(ctx context.Context, request CompletionRequest) (string, error)
}
