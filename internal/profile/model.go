package profile

import (
	"time"

	"nilchan-hackaton/internal/user"
)

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
	User      user.User
	Progress  UserProgress
	Cosmetics []UserCosmetic
}

type ProfileResult struct {
	Profile            *Profile
	Level              int
	ActiveBacklogLimit int
}

type OptionalCosmeticID struct {
	Set   bool
	Value *string
}

type CosmeticsUpdate struct {
	AvatarID       *string
	FrameID        *string
	TitleID        OptionalCosmeticID
	ShowcaseItemID OptionalCosmeticID
}
