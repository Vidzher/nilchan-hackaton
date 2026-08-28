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
) (*GetProfileResponse, error) {
	data, err := s.repository.GetProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	levelInfo := progress.FromXP(data.Progress.XP)
	ownedIDs := make([]string, 0, len(data.Cosmetics))
	for _, cosmetic := range data.Cosmetics {
		ownedIDs = append(ownedIDs, cosmetic.ItemID)
	}

	return &GetProfileResponse{
		User: UserDTO{
			ID:       data.User.ID,
			Email:    data.User.Email,
			Username: data.User.Username,
		},
		Progress: ProgressDTO{
			XP:                 data.Progress.XP,
			Level:              levelInfo.Level,
			ActiveBacklogLimit: levelInfo.ActiveBacklogLimit,
			EPoints:            data.Progress.EPoints,
			CurrentStreak:      data.Progress.CurrentStreak,
			LastCompletionAt:   data.Progress.LastCompletionAt,
		},
		Cosmetics: CosmeticsDTO{
			AvatarID:         data.Progress.AvatarID,
			FrameID:          data.Progress.FrameID,
			TitleID:          data.Progress.TitleID,
			ShowcaseItemID:   data.Progress.ShowcaseItemID,
			OwnedCosmeticIDs: ownedIDs,
		},
	}, nil
}

func (s *Service) UpdateCosmetics(
	ctx context.Context,
	userID int64,
	req UpdateCosmeticsRequest,
) (*GetProfileResponse, error) {
	current, err := s.repository.GetProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	owned := make(map[string]struct{}, len(current.Cosmetics))
	for _, cosmetic := range current.Cosmetics {
		owned[cosmetic.ItemID] = struct{}{}
	}
	update := CosmeticsUpdate{
		AvatarID:       req.AvatarID,
		FrameID:        req.FrameID,
		TitleID:        req.TitleID,
		ShowcaseItemID: req.ShowcaseItemID,
	}
	if req.AvatarID != nil {
		if err := validateOwnedType(
			*req.AvatarID,
			cosmetics.ItemTypeAvatar,
			owned,
		); err != nil {
			return nil, err
		}
	}
	if req.FrameID != nil {
		if err := validateOwnedType(
			*req.FrameID,
			cosmetics.ItemTypeFrame,
			owned,
		); err != nil {
			return nil, err
		}
	}
	if req.TitleID.Set && req.TitleID.Value != nil {
		if err := validateOwnedType(
			*req.TitleID.Value,
			cosmetics.ItemTypeTitle,
			owned,
		); err != nil {
			return nil, err
		}
	}
	if req.ShowcaseItemID.Set && req.ShowcaseItemID.Value != nil {
		if err := validateOwnedType(
			*req.ShowcaseItemID.Value,
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
