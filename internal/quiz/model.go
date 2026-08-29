package quiz

import "time"

type Quiz struct {
	ID         int64
	ResourceID int64
	Title      string
	Questions  []Question
}

type Answer struct {
	QuestionIndex int
	SelectedIndex int
}

type CompletionRecord struct {
	CompletedAt   time.Time
	XPEarned      int
	EPointsEarned int
	XP            int64
	EPoints       int64
	CurrentStreak int
}

type CompletionResult struct {
	CompletedAt        time.Time
	TotalQuestions     int
	XPEarned           int
	EPointsEarned      int
	XP                 int64
	EPoints            int64
	CurrentStreak      int
	Level              int
	ActiveBacklogLimit int
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
