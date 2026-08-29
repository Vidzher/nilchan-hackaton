package quiz

import "errors"

var (
	ErrNotFound       = errors.New("quiz not found")
	ErrUnavailable    = errors.New("quiz is not available")
	ErrInvalidAnswers = errors.New("quiz answers are invalid")
)
