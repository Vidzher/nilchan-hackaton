package shop

import "nilchan-hackaton/internal/cosmetics"

type Service struct{}

func NewService() *Service {
	return &Service{}
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
