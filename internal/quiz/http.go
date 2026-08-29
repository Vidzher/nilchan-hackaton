package quiz

import "time"

type GetQuizResponse struct {
	ID         int64                  `json:"id"`
	ResourceID int64                  `json:"resourceId"`
	Title      string                 `json:"title"`
	Questions  []QuizQuestionResponse `json:"questions"`
}

type QuizQuestionResponse struct {
	Text              string    `json:"text"`
	Options           [4]string `json:"options"`
	Explanation       string    `json:"explanation"`
	Evidence          string    `json:"evidence"`
	VerificationSalt  string    `json:"verificationSalt"`
	CorrectAnswerHash string    `json:"correctAnswerHash"`
}

type CompleteQuizRequest struct {
	Answers []SubmittedAnswer `json:"answers"`
}

type SubmittedAnswer struct {
	QuestionIndex int `json:"questionIndex"`
	SelectedIndex int `json:"selectedIndex"`
}

type CompleteQuizResponse struct {
	Completion CompletionDetails `json:"completion"`
	Progress   ProgressSnapshot  `json:"progress"`
}

type CompletionDetails struct {
	CompletedAt    time.Time `json:"completedAt"`
	TotalQuestions int       `json:"totalQuestions"`
	XPEarned       int       `json:"xpEarned"`
	EPointsEarned  int       `json:"ePointsEarned"`
}

type ProgressSnapshot struct {
	XP                 int64 `json:"xp"`
	EPoints            int64 `json:"ePoints"`
	CurrentStreak      int   `json:"currentStreak"`
	Level              int   `json:"level"`
	ActiveBacklogLimit int   `json:"activeBacklogLimit"`
}
