package quiz

import "time"

type Quiz struct {
	ID         string
	ResourceID string
	Title      string
	Questions  []Question
	CreatedAt  time.Time
}

type Question struct {
	Text              string
	Options           [4]string
	CorrectIndex      int
	Explanation       string
	Evidence          string
	VerificationSalt  string
	CorrectAnswerHash string
}
