package llm

import (
	"fmt"
	"strings"
)

func validateQuiz(response ResponseDTO, sourceText string) error {
	if strings.TrimSpace(response.Title) == "" {
		return invalidQuizError("title is empty")
	}
	if strings.TrimSpace(response.Topic) == "" {
		return invalidQuizError("topic is empty")
	}
	if !isValidDifficulty(response.Difficulty) {
		return invalidQuizError("difficulty must be easy, medium, or hard")
	}
	if len(response.Questions) != 5 {
		return invalidQuizError(fmt.Sprintf("expected 5 questions, got %d", len(response.Questions)))
	}

	normalizedSource := normalizeText(sourceText)
	seenQuestions := make(map[string]struct{}, len(response.Questions))

	for questionIndex, question := range response.Questions {
		questionText := strings.TrimSpace(question.Question)
		if questionText == "" {
			return invalidQuizError(fmt.Sprintf("question %d is empty", questionIndex+1))
		}

		normalizedQuestion := strings.ToLower(questionText)
		if _, exists := seenQuestions[normalizedQuestion]; exists {
			return invalidQuizError(fmt.Sprintf("question %d is duplicated", questionIndex+1))
		}
		seenQuestions[normalizedQuestion] = struct{}{}

		if len(question.Options) != 4 {
			return invalidQuizError(fmt.Sprintf("question %d must have 4 options", questionIndex+1))
		}
		if question.CorrectIndex < 0 || question.CorrectIndex >= len(question.Options) {
			return invalidQuizError(fmt.Sprintf("question %d has an invalid correctIndex", questionIndex+1))
		}

		seenOptions := make(map[string]struct{}, len(question.Options))
		for optionIndex, option := range question.Options {
			normalizedOption := strings.ToLower(strings.TrimSpace(option))
			if normalizedOption == "" {
				return invalidQuizError(fmt.Sprintf("question %d option %d is empty", questionIndex+1, optionIndex+1))
			}
			if _, exists := seenOptions[normalizedOption]; exists {
				return invalidQuizError(fmt.Sprintf("question %d contains duplicate options", questionIndex+1))
			}
			seenOptions[normalizedOption] = struct{}{}
		}

		if strings.TrimSpace(question.Explanation) == "" {
			return invalidQuizError(fmt.Sprintf("question %d explanation is empty", questionIndex+1))
		}

		evidence := normalizeText(question.Evidence)
		if evidence == "" {
			return invalidQuizError(fmt.Sprintf("question %d evidence is empty", questionIndex+1))
		}
		if !strings.Contains(normalizedSource, evidence) {
			return invalidQuizError(fmt.Sprintf("question %d evidence was not found in the source", questionIndex+1))
		}
	}

	return nil
}

func isValidDifficulty(difficulty DifficultyLevel) bool {
	switch difficulty {
	case DifficultyEasy, DifficultyMedium, DifficultyHard:
		return true
	default:
		return false
	}
}

func normalizeText(value string) string {
	return strings.Join(strings.Fields(value), " ")
}
