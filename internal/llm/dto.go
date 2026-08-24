package llm

type DifficultyLevel string

const (
	DifficultyEasy   DifficultyLevel = "easy"
	DifficultyMedium DifficultyLevel = "medium"
	DifficultyHard   DifficultyLevel = "hard"
)

type RequestDTO struct {
	SourceTitle string
	SourceText  string
}

type ResponseDTO struct {
	Title      string           `json:"title"`
	Topic      string           `json:"topic"`
	Difficulty DifficultyLevel `json:"difficulty"`
	Questions  []QuestionDTO    `json:"questions"`
}

type QuestionDTO struct {
	Question     string   `json:"question"`
	Options      []string `json:"options"`
	CorrectIndex int      `json:"correctIndex"`
	Explanation  string   `json:"explanation"`
	Evidence     string   `json:"evidence"`
}
