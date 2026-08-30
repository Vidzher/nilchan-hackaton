package leaderboard

type Entry struct {
	Rank      int
	UserID    int64
	Username  string
	XP        int64
	AvatarID  string
	FrameID   string
	Level     int
	IsCurrent bool
}
