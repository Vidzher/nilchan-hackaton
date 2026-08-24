package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	maxSourceChars     = 50_000
	generationAttempts = 2
)

type Service struct {
	client Client
}

func NewService(client Client) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) GenerateQuiz(ctx context.Context, dto RequestDTO) (*ResponseDTO, error) {
	if s == nil || s.client == nil {
		return nil, ErrClientNotConfigured
	}

	dto.SourceText = strings.TrimSpace(dto.SourceText)
	dto.SourceTitle = strings.TrimSpace(dto.SourceTitle)

	if dto.SourceText == "" {
		return nil, ErrEmptySource
	}
	if utf8.RuneCountInString(dto.SourceText) > maxSourceChars {
		return nil, ErrSourceTooLong
	}

	completionRequest := buildCompletionRequest(dto)
	var lastErr error

	for range generationAttempts {
		rawResponse, err := s.client.Generate(ctx, completionRequest)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrGenerationFailed, err)
		}

		var response ResponseDTO
		if err := json.Unmarshal([]byte(rawResponse), &response); err != nil {
			lastErr = fmt.Errorf("decode response: %w", err)
			continue
		}

		if err := validateQuiz(response, dto.SourceText); err != nil {
			lastErr = err
			continue
		}

		return &response, nil
	}

	return nil, fmt.Errorf("%w after %d attempts: %v", ErrInvalidResponse, generationAttempts, lastErr)
}
