package profile

import "time"

type GetProfileResponse struct {
	User      UserDTO      `json:"user"`
	Progress  ProgressDTO  `json:"progress"`
	Cosmetics CosmeticsDTO `json:"cosmetics"`
}

type UserDTO struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type ProgressDTO struct {
	XP                 int64      `json:"xp"`
	Level              int        `json:"level"`
	ActiveBacklogLimit int        `json:"activeBacklogLimit"`
	EPoints            int64      `json:"ePoints"`
	CurrentStreak      int        `json:"currentStreak"`
	LastCompletionAt   *time.Time `json:"lastCompletionAt,omitempty"`
}

type CosmeticsDTO struct {
	AvatarID         string   `json:"avatarId"`
	FrameID          string   `json:"frameId"`
	TitleID          *string  `json:"titleId,omitempty"`
	ShowcaseItemID   *string  `json:"showcaseItemId,omitempty"`
	OwnedCosmeticIDs []string `json:"ownedCosmeticIds"`
}

type UpdateCosmeticsRequest struct {
	AvatarID       *string `json:"avatarId,omitempty"`
	FrameID        *string `json:"frameId,omitempty"`
	TitleID        *string `json:"titleId,omitempty"`
	ShowcaseItemID *string `json:"showcaseItemId,omitempty"`
}
