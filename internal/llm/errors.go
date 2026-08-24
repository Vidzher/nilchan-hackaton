package llm

import (
	"errors"
	"fmt"
)

var (
	ErrEmptySource         = errors.New("source text is empty")
	ErrSourceTooLong       = errors.New("source text is too long")
	ErrClientNotConfigured = errors.New("LLM client is not configured")
	ErrGenerationFailed    = errors.New("LLM generation failed")
	ErrInvalidResponse     = errors.New("LLM returned an invalid quiz")
	ErrInvalidModelName    = errors.New("invalid model name")
	ErrInvalidSchema       = errors.New("invalid schema")
	ErrInvalidApiKey       = errors.New("invalid API key")
	ErrEmptyLLMResponse    = errors.New("LLM returned an empty response")
)

func invalidQuizError(message string) error {
	return fmt.Errorf("invalid quiz: %s", message)
}
