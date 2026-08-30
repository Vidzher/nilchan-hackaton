package leaderboard

type Entry struct {
	Rank           int
	UserID         int64
	Username       string
	XP             int64
	AvatarID       string
	FrameID        string
	TitleID        *string
	ShowcaseItemID *string
	Level          int
	IsCurrent      bool
}
