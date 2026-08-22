package users

import "time"

const (
	DefaultAvatarID = "default_avatar"
	DefaultFrameID  = "default_frame"
)

type User struct {
	ID           int       `json:"userId"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
}

type UserProgress struct {
	UserID           int        `json:"userId"`
	XP               int        `json:"xp"`
	EPoints          int        `json:"ePoints"`
	CurrentStreak    int        `json:"currentStreak"`
	LastCompletionAt *time.Time `json:"lastCompletionAt,omitempty"`
	AvatarID         string     `json:"avatarId"`
	FrameID          string     `json:"frameId"`
	TitleID          *string    `json:"titleId,omitempty"`
	ShowcaseItemID   *string    `json:"showcaseItemId,omitempty"`
	OwnedCosmeticIDs []string   `json:"ownedCosmeticIds"`
}
