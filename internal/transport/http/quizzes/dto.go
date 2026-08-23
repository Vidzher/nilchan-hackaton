package transport_http_quizzes

type GetQuizResponse struct {
	ID         int               `json:"id"`
	ResourceID int               `json:"resourceId"`
	Title      string            `json:"title"`
	Questions  []QuizQuestionDTO `json:"questions"`
}

type QuizQuestionDTO struct {
	ID      int             `json:"id"`
	Title   string          `json:"title"`
	Answers []QuizAnswerDTO `json:"answers"`
}

type QuizAnswerDTO struct {
	ID         int    `json:"id"`
	QuestionID int    `json:"questionId"`
	Title      string `json:"title"`
}

// {
//   "answers": [
//     {
//       "questionId": 101,
//       "answerId": 1001
//     },
//     {
//       "questionId": 102,
//       "answerId": 1008
//     },
//     {
//       "questionId": 103,
//       "answerId": 1012
//     }
//   ]
// }

type CompleteQuizRequest struct {
	Answers []SubmittedAnswer `json:"answers"`
}

type SubmittedAnswer struct {
	QuestionID int `json:"questionId"`
	AnswerID   int `json:"answerId"`
}
