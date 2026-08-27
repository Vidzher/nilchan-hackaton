package profile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

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

type OptionalString struct {
	Set   bool
	Value *string
}

func (o *OptionalString) UnmarshalJSON(data []byte) error {
	o.Set = true
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		o.Value = nil
		return nil
	}
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return fmt.Errorf("expected string or null: %w", err)
	}
	if value == "" {
		return fmt.Errorf("value must not be empty")
	}
	o.Value = &value
	return nil
}

type UpdateCosmeticsRequest struct {
	AvatarID       *string        `json:"avatarId,omitempty"`
	FrameID        *string        `json:"frameId,omitempty"`
	TitleID        OptionalString `json:"titleId"`
	ShowcaseItemID OptionalString `json:"showcaseItemId"`
}
