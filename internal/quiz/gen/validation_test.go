package gen

import (
	"encoding/json"
	"strconv"
	"strings"
	"testing"
)

const (
	validEvidence = `Use [context.Context](https://pkg.go.dev/context) to cancel work — even when status != OK & retries >= 2.`
	validSource   = `# Go & "Concurrency" 🚀

> **Rule:** source text is messy; validation should not shit the bed.

Use [context.Context](https://pkg.go.dev/context) to cancel work — even when status != OK & retries >= 2.

- Escaped-ish symbols: <tag attr="value">, {json: true}, [brackets], (parentheses), #hash, @user.
- Inline code: ` + "`defer cancel()`" + ` and ` + "`if err != nil`" + `.
- URL garbage: https://example.com/docs?q=go%20context&utm_source=test#section.

---

Unicode survives too: café, naïve, 中文, русский, Ελληνικά, 🦆.`
)

func validGeneratedQuiz(questionCount int) GeneratedQuiz {
	questions := make([]GeneratedQuestion, questionCount)
	for i := range questions {
		questions[i] = GeneratedQuestion{
			Text:         "Question " + string(rune('A'+i)),
			Options:      []string{"Option A", "Option B", "Option C", "Option D"},
			CorrectIndex: 0,
			Explanation:  "Explanation",
			Evidence:     validEvidence,
		}
	}
	return GeneratedQuiz{Questions: questions}
}

func TestValidateGeneratedQuizAcceptsValidQuiz(t *testing.T) {
	for _, questionCount := range []int{5, 10} {
		t.Run(strconv.Itoa(questionCount)+" questions", func(t *testing.T) {
			quiz := validGeneratedQuiz(questionCount)
			quiz.Questions[0].Evidence = "Use   [context.Context](https://pkg.go.dev/context)\n\tto cancel work — even when status != OK & retries >= 2."

			if err := validateGeneratedQuiz(quiz); err != nil {
				t.Fatalf("validate quiz: %v", err)
			}
		})
	}
}

func TestValidateGeneratedQuizRejectsInvalidQuiz(t *testing.T) {
	tests := []struct {
		name       string
		mutate     func(*GeneratedQuiz)
		wantErrExp string
	}{
		{
			name: "too few questions",
			mutate: func(quiz *GeneratedQuiz) {
				quiz.Questions = quiz.Questions[:4]
			},
			wantErrExp: `expected 5 to 10 questions, got 4`,
		},
		{
			name: "too many questions",
			mutate: func(quiz *GeneratedQuiz) {
				*quiz = validGeneratedQuiz(11)
			},
			wantErrExp: `expected 5 to 10 questions, got 11`,
		},
		{
			name: "empty question",
			mutate: func(quiz *GeneratedQuiz) {
				quiz.Questions[0].Text = " \t\n"
			},
			wantErrExp: `question 1 is empty`,
		},
		{
			name: "duplicate question",
			mutate: func(quiz *GeneratedQuiz) {
				quiz.Questions[1].Text = "  QUESTION\nA  "
			},
			wantErrExp: `question 2 is duplicated`,
		},
		{
			name: "too few options",
			mutate: func(quiz *GeneratedQuiz) {
				quiz.Questions[0].Options = quiz.Questions[0].Options[:3]
			},
			wantErrExp: `question 1 must have 4 options`,
		},
		{
			name: "too many options",
			mutate: func(quiz *GeneratedQuiz) {
				quiz.Questions[0].Options = append(quiz.Questions[0].Options, "Option E")
			},
			wantErrExp: `question 1 must have 4 options`,
		},
		{
			name: "negative correct index",
			mutate: func(quiz *GeneratedQuiz) {
				quiz.Questions[0].CorrectIndex = -1
			},
			wantErrExp: `question 1 has an invalid correctIndex`,
		},
		{
			name: "correct index out of range",
			mutate: func(quiz *GeneratedQuiz) {
				quiz.Questions[0].CorrectIndex = 4
			},
			wantErrExp: `question 1 has an invalid correctIndex`,
		},
		{
			name: "empty option",
			mutate: func(quiz *GeneratedQuiz) {
				quiz.Questions[0].Options[1] = " \t\n"
			},
			wantErrExp: `question 1 option 2 is empty`,
		},
		{
			name: "duplicate option",
			mutate: func(quiz *GeneratedQuiz) {
				quiz.Questions[0].Options[1] = "  OPTION\nA  "
			},
			wantErrExp: `question 1 contains duplicate options`,
		},
		{
			name: "empty explanation",
			mutate: func(quiz *GeneratedQuiz) {
				quiz.Questions[0].Explanation = " \t\n"
			},
			wantErrExp: `question 1 explanation is empty`,
		},
		{
			name: "empty evidence",
			mutate: func(quiz *GeneratedQuiz) {
				quiz.Questions[0].Evidence = " \t\n"
			},
			wantErrExp: `question 1 evidence is empty`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quiz := validGeneratedQuiz(5)
			tt.mutate(&quiz)

			err := validateGeneratedQuiz(quiz)
			if err == nil {
				t.Fatal("expected validation error")
			}
			if !strings.Contains(err.Error(), tt.wantErrExp) {
				t.Fatalf("expected error containing %q, got %q", tt.wantErrExp, err)
			}
		})
	}
}

func TestGeneratedQuestionPreservesExcessOptionsForValidation(t *testing.T) {
	question := `{
		"text":"Question",
		"options":["one","two","three","four","five"],
		"correctIndex":0,
		"explanation":"Explanation",
		"evidence":"Use [context.Context](https://pkg.go.dev/context) to cancel work — even when status != OK & retries >= 2."
	}`
	raw := `{"questions":[` + strings.TrimSuffix(strings.Repeat(question+`,`, 5), `,`) + `]}`

	var quiz GeneratedQuiz
	if err := json.Unmarshal([]byte(raw), &quiz); err != nil {
		t.Fatalf("decode quiz: %v", err)
	}
	if got := len(quiz.Questions[0].Options); got != 5 {
		t.Fatalf("expected 5 decoded options, got %d", got)
	}

	err := validateGeneratedQuiz(quiz)
	if err == nil || !strings.Contains(err.Error(), "must have 4 options") {
		t.Fatalf("expected option count error, got %v", err)
	}
}
