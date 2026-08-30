package leaderboard

type EntryResponse struct {
	Rank           int     `json:"rank"`
	UserID         int64   `json:"userId"`
	Username       string  `json:"username"`
	XP             int64   `json:"xp"`
	Level          int     `json:"level"`
	AvatarID       string  `json:"avatarId"`
	FrameID        string  `json:"frameId"`
	TitleID        *string `json:"titleId,omitempty"`
	ShowcaseItemID *string `json:"showcaseItemId,omitempty"`
	IsCurrent      bool    `json:"isCurrent"`
}
