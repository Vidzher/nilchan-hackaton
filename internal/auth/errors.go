package auth

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailTaken         = errors.New("email already exists")
	ErrUsernameTaken      = errors.New("username already exists")
)
