package gen

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"nilchan-hackaton/internal/llm"
)

const (
	maxSourceChars     = 50_000
	generationAttempts = 2
	generationTimeout  = 120 * time.Second
)

type completer interface {
	Complete(ctx context.Context, request llm.Request) (string, error)
}

type Generator struct {
	client         completer
	responseSchema ResponseSchema
}

func NewGenerator(client completer) (*Generator, error) {
	if client == nil {
		return nil, ErrClientNotConfigured
	}

	var schema map[string]any
	if err := json.Unmarshal([]byte(generationSchemaJSON), &schema); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidSchema, err)
	}

	return &Generator{client: client, responseSchema: ResponseSchema{
		Name:   generationSchemaName,
		Schema: schema,
	}}, nil
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
		ResponseSchemaName: g.responseSchema.Name,
		ResponseSchema:     g.responseSchema.Schema,
		Temperature:        0.2,
	}
	var lastErr error

	for range generationAttempts {
		generationCtx, cancel := context.WithTimeout(ctx, generationTimeout)
		rawResponse, err := g.client.Complete(generationCtx, completionRequest)
		cancel()

		if err != nil {
			lastErr = fmt.Errorf("%w: %w", ErrGenerationFailed, err)
			continue
		}

		var response GeneratedQuiz
		if err := json.Unmarshal([]byte(rawResponse), &response); err != nil {
			lastErr = fmt.Errorf("decode response: %w", err)
			continue
		}

		if err := validateGeneratedQuiz(response); err != nil {
			lastErr = err
			continue
		}

		return &response, nil
	}

	return nil, fmt.Errorf("%w after %d attempts: %w", ErrInvalidResponse, generationAttempts, lastErr)
}
