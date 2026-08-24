package quiz

type Quiz struct {
	ID         int64
	ResourceID int64
	Title      string
	Questions  []Question
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
