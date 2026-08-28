package resource

import "time"

type ResourceResponse struct {
	ID            int64      `json:"id"`
	URL           string     `json:"url"`
	Title         string     `json:"title"`
	Tags          []string   `json:"tags"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"createdAt"`
	CompletedAt   *time.Time `json:"completedAt,omitempty"`
	XPEarned      *int       `json:"xpEarned,omitempty"`
	EPointsEarned *int       `json:"ePointsEarned,omitempty"`
}

type CreateResourceRequest struct {
	URL                  string `json:"url" validate:"required"`
	PurchaseOverflowSlot bool   `json:"purchaseOverflowSlot"`
}

type CreateResourceResponse ResourceResponse

type ListResourcesResponse []ResourceResponse
