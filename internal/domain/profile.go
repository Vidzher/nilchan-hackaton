package domain

import "time"

type User struct {
	ID           int64
	Email        string
	Username     string
	PasswordHash string
	CreatedAt    time.Time
}

type UserProgress struct {
	UserID           int64
	XP               int64
	EPoints          int64
	CurrentStreak    int
	LastCompletionAt *time.Time
	AvatarID         string
	FrameID          string
	TitleID          *string
	ShowcaseItemID   *string
}

type UserCosmetic struct {
	UserID int64
	ItemID string
}

type Profile struct {
	User      User
	Progress  UserProgress
	Cosmetics []UserCosmetic
}
