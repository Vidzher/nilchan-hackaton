package auth

import (
	"fmt"

	pwd "nilchan-hackaton/internal/auth/password"
	"nilchan-hackaton/internal/auth/token"
	"nilchan-hackaton/internal/user"
)

type userRepository interface {
	create(email, username, passwordHash string) (*user.User, error)
	findOne(email string) (*user.User, error)
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

func (as *Service) Login(email string, password string) (*Result, error) {
	res, err := as.repo.findOne(email)
	if err != nil {
		return nil, err
	}

	if !pwd.CheckHash(password, res.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := token.Generate(res.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate accessToken: %w", err)
	}

	return &Result{Token: accessToken, User: res}, nil
}

func (as *Service) Register(email, username, password string) (*Result, error) {
	pwdHash, err := pwd.Hash(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := as.repo.create(email, username, pwdHash)
	if err != nil {
		return nil, err
	}

	accessToken, err := token.Generate(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate accessToken: %w", err)
	}

	return &Result{Token: accessToken, User: user}, nil
}
