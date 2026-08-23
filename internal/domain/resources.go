package domain

import (
	"time"
)

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
	ErrorCode             *string
	PurchasedOverflowSlot bool
	CreatedAt             time.Time
	CompletedAt           *time.Time
	XPEarned              *int
	EPointsEarned         *int
	OldBacklogBonus       *int
	FullBacklogBonus      *int
}

type CreateResourceInput struct {
	UserID               int64
	URL                  string
	PurchaseOverflowSlot bool
}

type ListResourcesInput struct {
	UserID int64
	Status *Status
	Tag    *string
}
