package auth

import (
	"context"
	"fmt"

	pwd "nilchan-hackaton/internal/auth/password"
	"nilchan-hackaton/internal/auth/token"
	"nilchan-hackaton/internal/user"
)

type userRepository interface {
	create(ctx context.Context, email, username, passwordHash string) (*user.User, error)
	findByEmail(ctx context.Context, email string) (*user.User, error)
}

type Service struct {
	repo userRepository
}

type Result struct {
	User  *user.User
	Token string
}

func NewService(repo userRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Login(ctx context.Context, email string, password string) (*Result, error) {
	foundUser, err := s.repo.findByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !pwd.CheckHash(password, foundUser.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := token.Generate(foundUser.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate accessToken: %w", err)
	}

	return &Result{Token: accessToken, User: foundUser}, nil
}

func (s *Service) Register(ctx context.Context, email, username, password string) (*Result, error) {
	pwdHash, err := pwd.Hash(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := s.repo.create(ctx, email, username, pwdHash)
	if err != nil {
		return nil, err
	}

	accessToken, err := token.Generate(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate accessToken: %w", err)
	}

	return &Result{Token: accessToken, User: user}, nil
}
