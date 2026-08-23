package repository_resources

import "time"

type ResourceModel struct {
	ID                    int64
	UserID                int64
	URL                   string
	Title                 string
	Content               string
	Status                string
	PurchasedOverflowSlot bool
	CreatedAt             time.Time
	CompletedAt           *time.Time
	XPEarned              *int
	EPointsEarned         *int
}

type ResourceTagModel struct {
	ResourceID int64
	Tag        string
}
