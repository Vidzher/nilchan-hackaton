package gen

type GenerationRequest struct {
	SourceTitle string
	SourceText  string
}

type GeneratedQuiz struct {
	Questions []GeneratedQuestion `json:"questions"`
}

type GeneratedQuestion struct {
	Text         string    `json:"text"`
	Options      [4]string `json:"options"`
	CorrectIndex int       `json:"correctIndex"`
	Explanation  string    `json:"explanation"`
	Evidence     string    `json:"evidence"`
}
