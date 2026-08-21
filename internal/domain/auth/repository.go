package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"nilchan-hackaton/internal/shared/models/users"
	"nilchan-hackaton/internal/storage"

	"github.com/mattn/go-sqlite3"
)

type AuthRepository struct {
	storage *storage.Storage
}

func NewAuthRepository(storage *storage.Storage) *AuthRepository {
	return &AuthRepository{storage: storage}
}

func (ar *AuthRepository) Create(email string, password string) error {
	stmt, err := ar.storage.DB.Prepare("INSERT INTO users(email, password) VALUES(?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(email, password)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("user exists: %v", email)
		}

		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}

func (ar *AuthRepository) FindOne(email string) (*users.User, error) {
	stmt, err := ar.storage.DB.Prepare("SELECT * FROM users WHERE email=?")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}

	defer stmt.Close()

	var found users.User
	err = stmt.QueryRow(email).Scan(&found.ID, &found.Email, &found.Password, &found.Balance)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("invalid credentials")
	}
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}

	return &found, nil
}
