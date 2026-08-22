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
	ID                    string
	UserID                int
	URL                   string
	Title                 string
	Tags                  []string
	Content               string
	Status                Status
	ErrorCode             *string
	PurchasedOverflowSlot bool
	CreatedAt             time.Time
	CompletedAt           *time.Time
	XPEarned              *int
	EPointsEarned         *int
	OldBacklogBonus       *int
	FullBacklogBonus      *int
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
