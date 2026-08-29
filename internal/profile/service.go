package profile

import (
	"context"

	"nilchan-hackaton/internal/cosmetics"
	"nilchan-hackaton/internal/progress"
)

type profileRepository interface {
	GetProfile(
		ctx context.Context,
		userID int64,
	) (*Profile, error)
	UpdateCosmetics(
		ctx context.Context,
		userID int64,
		update CosmeticsUpdate,
	) error
}

type Service struct {
	repository profileRepository
}

func NewService(repository profileRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetProfile(
	ctx context.Context,
	userID int64,
) (*ProfileResult, error) {
	data, err := s.repository.GetProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	levelInfo := progress.FromXP(data.Progress.XP)
	return &ProfileResult{
		Profile:            data,
		Level:              levelInfo.Level,
		ActiveBacklogLimit: levelInfo.ActiveBacklogLimit,
	}, nil
}

func (s *Service) UpdateCosmetics(
	ctx context.Context,
	userID int64,
	update CosmeticsUpdate,
) (*ProfileResult, error) {
	current, err := s.repository.GetProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	owned := make(map[string]struct{}, len(current.Cosmetics))
	for _, cosmetic := range current.Cosmetics {
		owned[cosmetic.ItemID] = struct{}{}
	}
	if update.AvatarID != nil {
		if err := validateOwnedType(
			*update.AvatarID,
			cosmetics.ItemTypeAvatar,
			owned,
		); err != nil {
			return nil, err
		}
	}
	if update.FrameID != nil {
		if err := validateOwnedType(
			*update.FrameID,
			cosmetics.ItemTypeFrame,
			owned,
		); err != nil {
			return nil, err
		}
	}
	if update.TitleID.Set && update.TitleID.Value != nil {
		if err := validateOwnedType(
			*update.TitleID.Value,
			cosmetics.ItemTypeTitle,
			owned,
		); err != nil {
			return nil, err
		}
	}
	if update.ShowcaseItemID.Set && update.ShowcaseItemID.Value != nil {
		if err := validateOwnedType(
			*update.ShowcaseItemID.Value,
			cosmetics.ItemTypeShowcase,
			owned,
		); err != nil {
			return nil, err
		}
	}
	if err := s.repository.UpdateCosmetics(ctx, userID, update); err != nil {
		return nil, err
	}
	return s.GetProfile(ctx, userID)
}

func validateOwnedType(
	itemID string,
	expectedType cosmetics.ItemType,
	owned map[string]struct{},
) error {
	item, ok := cosmetics.GetByID(itemID)
	if !ok {
		return ErrUnknownCosmetic
	}
	if item.Type != expectedType {
		return ErrCosmeticTypeMismatch
	}
	if _, ok := owned[itemID]; !ok {
		return ErrCosmeticNotOwned
	}
	return nil
}
