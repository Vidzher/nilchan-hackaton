package auth

import (
	"fmt"
	pwd "nilchan-hackaton/internal/lib/utils/password"
	"nilchan-hackaton/internal/lib/utils/token"
	"nilchan-hackaton/internal/shared/models/users"
)

type Repository interface {
	Create(email string, password string) error
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

	if !pwd.CheckHash(password, res.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	accessToken, err := token.Generate(res.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate accessToken: %w", err)
	}

	return &AuthResult{token: accessToken, user: res}, nil
}

func (as *AuthService) Register(email string, password string) error {
	pwdHas, err := pwd.Hash(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err.Error())
	}
	err = as.repo.Create(email, pwdHas)

	return err
}
