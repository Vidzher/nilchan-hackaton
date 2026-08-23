package user

import "time"

const (
	DefaultAvatarID = "default_avatar"
	DefaultFrameID  = "default_frame"
)

type User struct {
	ID           int64
	Email        string
	Username     string
	PasswordHash string
	CreatedAt    time.Time
}
