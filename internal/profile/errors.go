package profile

import "errors"

var (
	ErrProfileNotFound      = errors.New("profile not found")
	ErrUnknownCosmetic      = errors.New("unknown cosmetic")
	ErrCosmeticNotOwned     = errors.New("cosmetic is not owned")
	ErrCosmeticTypeMismatch = errors.New("cosmetic type does not match slot")
)
