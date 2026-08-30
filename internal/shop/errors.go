package shop

import "errors"

var (
	ErrUnknownItem         = errors.New("unknown shop item")
	ErrItemNotPurchasable  = errors.New("shop item is not purchasable")
	ErrAlreadyOwned        = errors.New("cosmetic is already owned")
	ErrInsufficientEPoints = errors.New("insufficient e-points")
)
