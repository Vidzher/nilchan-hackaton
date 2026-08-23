package transport_http_resource

import "time"

type ResourceDTOResponse struct {
	ID                    int        `json:"id"`
	UserID                int        `json:"userId"`
	URL                   string     `json:"url"`
	Title                 string     `json:"title"`
	Tags                  []string   `json:"tags"`
	Status                string     `json:"status"`
	PurchasedOverflowSlot bool       `json:"purchasedOverflowSlot"`
	CreatedAt             time.Time  `json:"createdAt"`
	CompletedAt           *time.Time `json:"completedAt,omitempty"`
	XPEarned              int        `json:"xpEarned"`
	EPointsEarned         int        `json:"EPointEarned"`
}

type CreateResourceRequest struct {
	URL                  string `json:"url" validate:"required,url"`
	PurchaseOverflowSlot bool   `json:"purchaseOverflowSlot" validate:"omitempty"`
}

type CreateResourceResponse ResourceDTOResponse

type GetResourcesResponse []ResourceDTOResponse
