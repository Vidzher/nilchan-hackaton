package quiz

type GetQuizResponse struct {
	ID         int64             `json:"id"`
	ResourceID int64             `json:"resourceId"`
	Title      string            `json:"title"`
	Questions  []QuizQuestionDTO `json:"questions"`
}

type QuizQuestionDTO struct {
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
