package repository_auth

import "time"

type UserModel struct {
	ID           int64
	Email        string
	Username     string
	PasswordHash string
	CreatedAt    time.Time
}

type UserProgressModel struct {
	UserID        int64
	XP            int64
	EPoints       int64
	CurrentStreak int
	AvatarID      string
	FrameID       string
}

type UserCosmeticModel struct {
	UserID int64
	ItemID string
}
