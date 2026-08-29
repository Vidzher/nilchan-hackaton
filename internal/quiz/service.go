package quiz

import (
	"context"
	"time"

	"nilchan-hackaton/internal/progress"
)

type quizRepository interface {
	Find(ctx context.Context, userID, resourceID int64) (*Quiz, error)
	Complete(
		ctx context.Context,
		userID, resourceID int64,
		xpEarned, ePointsEarned int,
		completedAt time.Time,
	) (*CompletionRecord, error)
}

type Service struct {
	repository quizRepository
}

func NewService(repository quizRepository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Get(ctx context.Context, userID, resourceID int64) (*Quiz, error) {
	found, err := s.repository.Find(ctx, userID, resourceID)
	if err != nil {
		return nil, err
	}
	if found == nil {
		return nil, ErrNotFound
	}
	return found, nil
}

func (s *Service) Complete(
	ctx context.Context,
	userID, resourceID int64,
	answers []Answer,
) (*CompletionResult, error) {
	found, err := s.Get(ctx, userID, resourceID)
	if err != nil {
		return nil, err
	}
	if !validAnswers(found.Questions, answers) {
		return nil, ErrInvalidAnswers
	}

	totalQuestions := len(found.Questions)
	record, err := s.repository.Complete(
		ctx,
		userID,
		resourceID,
		20+totalQuestions*5,
		totalQuestions,
		time.Now().UTC(),
	)
	if err != nil {
		return nil, err
	}

	level := progress.FromXP(record.XP)
	return &CompletionResult{
		CompletedAt:        record.CompletedAt,
		TotalQuestions:     totalQuestions,
		XPEarned:           record.XPEarned,
		EPointsEarned:      record.EPointsEarned,
		XP:                 record.XP,
		EPoints:            record.EPoints,
		CurrentStreak:      record.CurrentStreak,
		Level:              level.Level,
		ActiveBacklogLimit: level.ActiveBacklogLimit,
	}, nil
}

func validAnswers(questions []Question, answers []Answer) bool {
	if len(answers) != len(questions) {
		return false
	}

	seen := make([]bool, len(questions))
	for _, answer := range answers {
		if answer.QuestionIndex < 0 || answer.QuestionIndex >= len(questions) ||
			answer.SelectedIndex < 0 || answer.SelectedIndex >= len(questions[answer.QuestionIndex].Options) ||
			seen[answer.QuestionIndex] ||
			answer.SelectedIndex != questions[answer.QuestionIndex].CorrectIndex {
			return false
		}
		seen[answer.QuestionIndex] = true
	}
	return true
}
