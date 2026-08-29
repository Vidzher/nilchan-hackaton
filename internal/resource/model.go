package resource

import "time"

type Status string

const (
	StatusProcessing   Status = "PROCESSING"
	StatusFailed       Status = "FAILED"
	StatusNotCompleted Status = "NOT_COMPLETED"
	StatusCompleted    Status = "COMPLETED"
)

type Resource struct {
	ID                    int64
	UserID                int64
	URL                   string
	Title                 string
	Tags                  []string
	Content               string
	Status                Status
	PurchasedOverflowSlot bool
	CreatedAt             time.Time
	CompletedAt           *time.Time
	XPEarned              *int
	EPointsEarned         *int
}

type Summary struct {
	ID            int64
	URL           string
	Title         string
	Tags          []string
	Status        Status
	CreatedAt     time.Time
	CompletedAt   *time.Time
	XPEarned      *int
	EPointsEarned *int
}

func (s Status) Valid() bool {
	switch s {
	case StatusProcessing, StatusFailed, StatusNotCompleted, StatusCompleted:
		return true
	default:
		return false
	}
}

func (s Status) CanTransitionTo(next Status) bool {
	switch s {
	case StatusProcessing:
		return next == StatusNotCompleted || next == StatusFailed
	case StatusFailed:
		return next == StatusProcessing
	case StatusNotCompleted:
		return next == StatusCompleted
	default:
		return false
	}
}
