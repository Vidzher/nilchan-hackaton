package domain

type Quiz struct {
	ID         int
	ResourceID string
	Title      string
	Questions  []Question
}

type Question struct {
	ID              int
	QuizID          int
	Title           string
	Answers         []Answer
	CorrectAnswerID int
}

type Answer struct {
	ID         int
	QuestionID int
	Title      string
}
