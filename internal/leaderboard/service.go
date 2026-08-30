package leaderboard

import (
	"context"

	"nilchan-hackaton/internal/progress"
)

type leaderboardRepository interface {
	List(ctx context.Context, currentUserID int64) ([]Entry, error)
}

type Service struct {
	repository leaderboardRepository
}

func NewService(repository leaderboardRepository) *Service {
	return &Service{repository: repository}
}

func (s *Service) List(ctx context.Context, currentUserID int64) ([]Entry, error) {
	entries, err := s.repository.List(ctx, currentUserID)
	if err != nil {
		return nil, err
	}
	for index := range entries {
		entries[index].Level = progress.FromXP(entries[index].XP).Level
		entries[index].IsCurrent = entries[index].UserID == currentUserID
	}
	return entries, nil
}
