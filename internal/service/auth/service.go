package service_auth

import (
	"fmt"
	"nilchan-hackaton/internal/errors"
	pwd "nilchan-hackaton/internal/lib/utils/password"
	"nilchan-hackaton/internal/lib/utils/token"
	"nilchan-hackaton/internal/shared/models/users"
)

type Repository interface {
	Create(email, username, passwordHash string) (*users.User, error)
	FindOne(email string) (*users.User, error)
}

type AuthService struct {
	repo Repository
}

type AuthResult struct {
	user  *users.User
	token string
}

func NewAuthService(repo Repository) *AuthService {
	return &AuthService{repo: repo}
}

func (as *AuthService) Login(email string, password string) (*AuthResult, error) {
	res, err := as.repo.FindOne(email)
	if err != nil {
		return nil, err
	}

	if !pwd.CheckHash(password, res.PasswordHash) {
		return nil, errors.ErrInvalidCredentials
	}

	accessToken, err := token.Generate(res.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate accessToken: %w", err)
	}

	return &AuthResult{token: accessToken, user: res}, nil
}

func (as *AuthService) Register(email, username, password string) (*AuthResult, error) {
	pwdHash, err := pwd.Hash(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := as.repo.Create(email, username, pwdHash)
	if err != nil {
		return nil, err
	}

	accessToken, err := token.Generate(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate accessToken: %w", err)
	}

	return &AuthResult{token: accessToken, user: user}, nil
}
