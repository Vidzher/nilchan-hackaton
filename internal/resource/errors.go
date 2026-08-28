package resource

import "errors"

var (
	ErrInvalidURL          = errors.New("invalid resource URL")
	ErrBlockedURL          = errors.New("resource URL is blocked")
	ErrDuplicate           = errors.New("resource already exists")
	ErrBacklogFull         = errors.New("active backlog is full")
	ErrInsufficientEPoints = errors.New("insufficient e-points")
	ErrContentTooShort     = errors.New("resource content is too short")
	ErrContentTooLong      = errors.New("resource content is too long")
	ErrFirecrawlFailed     = errors.New("content extraction failed")
	ErrFirecrawlTimeout    = errors.New("content extraction timed out")
	ErrStateConflict       = errors.New("resource state conflict")
)
