package shop

import (
	"context"

	"nilchan-hackaton/internal/cosmetics"
)

type purchaseRepository interface {
	Purchase(ctx context.Context, userID int64, itemID string, price int64) (int64, error)
}

type Service struct {
	repository purchaseRepository
}

type PurchaseResult struct {
	Item    cosmetics.Item
	EPoints int64
}

func NewService(repository purchaseRepository) *Service {
	return &Service{repository: repository}
}

func (s *Service) List() []cosmetics.Item {
	catalog := cosmetics.All()
	items := make([]cosmetics.Item, 0, len(catalog))
	for _, item := range catalog {
		if !item.Free {
			items = append(items, item)
		}
	}
	return items
}

func (s *Service) Purchase(ctx context.Context, userID int64, itemID string) (*PurchaseResult, error) {
	item, ok := cosmetics.GetByID(itemID)
	if !ok {
		return nil, ErrUnknownItem
	}
	if item.Free {
		return nil, ErrItemNotPurchasable
	}

	balance, err := s.repository.Purchase(ctx, userID, item.ID, item.Price)
	if err != nil {
		return nil, err
	}
	return &PurchaseResult{Item: item, EPoints: balance}, nil
}
