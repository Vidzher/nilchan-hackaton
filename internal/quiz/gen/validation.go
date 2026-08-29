package gen

import (
	"fmt"
	"strings"
	"unicode"
)

func validateGeneratedQuiz(response GeneratedQuiz, sourceText string) error {
	if len(response.Questions) < 5 || len(response.Questions) > 10 {
		return invalidQuizError(fmt.Sprintf("expected 5 to 10 questions, got %d", len(response.Questions)))
	}

	normalizedSource := normalizeEvidence(sourceText)
	seenQuestions := make(map[string]struct{}, len(response.Questions))

	for questionIndex, question := range response.Questions {
		questionText := strings.TrimSpace(question.Text)
		if questionText == "" {
			return invalidQuizError(fmt.Sprintf("question %d is empty", questionIndex+1))
		}

		normalizedQuestion := strings.ToLower(normalizeText(questionText))
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
			normalizedOption := strings.ToLower(normalizeText(option))
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

		evidence := normalizeEvidence(question.Evidence)
		if evidence == "" {
			return invalidQuizError(fmt.Sprintf("question %d evidence is empty", questionIndex+1))
		}
		if !strings.Contains(" "+normalizedSource+" ", " "+evidence+" ") {
			return invalidQuizError(fmt.Sprintf("question %d evidence was not found in the source", questionIndex+1))
		}
	}

	return nil
}

func normalizeText(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

func normalizeEvidence(value string) string {
	var normalized strings.Builder
	for _, current := range value {
		if unicode.IsLetter(current) || unicode.IsNumber(current) {
			normalized.WriteRune(unicode.ToLower(current))
			continue
		}
		normalized.WriteByte(' ')
	}
	return strings.Join(strings.Fields(normalized.String()), " ")
}
