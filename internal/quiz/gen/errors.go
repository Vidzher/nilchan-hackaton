package gen

import (
	"errors"
	"fmt"
)

var (
	ErrEmptySource         = errors.New("source text is empty")
	ErrSourceTooLong       = errors.New("source text is too long")
	ErrClientNotConfigured = errors.New("LLM client is not configured")
	ErrGenerationFailed    = errors.New("quiz generation failed")
	ErrInvalidResponse     = errors.New("LLM returned an invalid quiz")
	ErrInvalidSchema       = errors.New("invalid quiz schema")
)

func invalidQuizError(message string) error {
	return fmt.Errorf("invalid quiz: %s", message)
}
