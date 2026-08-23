package repository_profile

import "time"

type UserModel struct {
	ID           int64
	Email        string
	Username     string
	PasswordHash string
	CreatedAt    time.Time
}

type UserProgressModel struct {
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

type UserCosmeticModel struct {
	UserID int64
	ItemID string
}

type ProfileModel struct {
	User      UserModel
	Progress  UserProgressModel
	Cosmetics []UserCosmeticModel
}
