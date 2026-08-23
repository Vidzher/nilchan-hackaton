package repository_quizzes

type QuizModel struct {
	ID         int64
	ResourceID string
	Title      string
}
type QuestionModel struct {
	ID              int64
	QuizID          int64
	Title           string
	CorrectAnswerID int64
}
type AnswerModel struct {
	ID         int64
	QuestionID int64
	Title      string
}
