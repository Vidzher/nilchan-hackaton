package llm

import "errors"

var (
	ErrInvalidModelName = errors.New("invalid model name")
	ErrInvalidAPIKey    = errors.New("invalid API key")
	ErrEmptyResponse    = errors.New("LLM returned an empty response")
)
