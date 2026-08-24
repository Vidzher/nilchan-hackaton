package gen

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"

	"nilchan-hackaton/internal/llm"
)

const (
	maxSourceChars     = 50_000
	generationAttempts = 2
)

type completer interface {
	Complete(ctx context.Context, request llm.Request) (string, error)
}

type Generator struct {
	client completer
	schema map[string]any
}

func NewGenerator(client completer) (*Generator, error) {
	if client == nil {
		return nil, ErrClientNotConfigured
	}

	var schema map[string]any
	if err := json.Unmarshal([]byte(generationSchemaJSON), &schema); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidSchema, err)
	}

	return &Generator{client: client, schema: schema}, nil
}

func (g *Generator) Generate(ctx context.Context, request GenerationRequest) (*GeneratedQuiz, error) {
	request.SourceText = strings.TrimSpace(request.SourceText)
	request.SourceTitle = strings.TrimSpace(request.SourceTitle)

	if request.SourceText == "" {
		return nil, ErrEmptySource
	}
	if utf8.RuneCountInString(request.SourceText) > maxSourceChars {
		return nil, ErrSourceTooLong
	}

	completionRequest := llm.Request{
		SystemPrompt:       systemPrompt,
		UserPrompt:         buildUserPrompt(request),
		ResponseSchemaName: generationSchemaName,
		ResponseSchema:     g.schema,
		Temperature:        0.2,
	}
	var lastErr error

	for range generationAttempts {
		rawResponse, err := g.client.Complete(ctx, completionRequest)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrGenerationFailed, err)
		}

		var response GeneratedQuiz
		if err := json.Unmarshal([]byte(rawResponse), &response); err != nil {
			lastErr = fmt.Errorf("decode response: %w", err)
			continue
		}

		if err := validateGeneratedQuiz(response, request.SourceText); err != nil {
			lastErr = err
			continue
		}

		return &response, nil
	}

	return nil, fmt.Errorf("%w after %d attempts: %w", ErrInvalidResponse, generationAttempts, lastErr)
}
