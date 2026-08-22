package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"nilchan-hackaton/internal/shared/models/users"
	"nilchan-hackaton/internal/storage"
	"strings"

	"github.com/mattn/go-sqlite3"
)

type AuthRepository struct {
	storage *storage.Storage
}

func NewAuthRepository(storage *storage.Storage) *AuthRepository {
	return &AuthRepository{storage: storage}
}

func (ar *AuthRepository) Create(email, username, passwordHash string) (*users.User, error) {
	tx, err := ar.storage.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin user creation: %w", err)
	}
	defer tx.Rollback()

	var created users.User
	err = tx.QueryRow(`
		INSERT INTO users(email, username, password_hash)
		VALUES(?, ?, ?)
		RETURNING id, email, username, password_hash, created_at
	`, email, username, passwordHash).Scan(
		&created.ID,
		&created.Email,
		&created.Username,
		&created.PasswordHash,
		&created.CreatedAt,
	)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			switch {
			case strings.Contains(err.Error(), "users.email"):
				return nil, ErrEmailTaken
			case strings.Contains(err.Error(), "users.username"):
				return nil, ErrUsernameTaken
			}
		}
		return nil, fmt.Errorf("create user: %w", err)
	}

	if _, err := tx.Exec("INSERT INTO user_progress(user_id) VALUES(?)", created.ID); err != nil {
		return nil, fmt.Errorf("create user progress: %w", err)
	}

	if _, err := tx.Exec(
		"INSERT INTO user_cosmetics(user_id, item_id) VALUES(?, ?), (?, ?)",
		created.ID,
		users.DefaultAvatarID,
		created.ID,
		users.DefaultFrameID,
	); err != nil {
		return nil, fmt.Errorf("grant default cosmetics: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit user creation: %w", err)
	}
	return &created, nil
}

func (ar *AuthRepository) FindOne(email string) (*users.User, error) {
	var found users.User
	err := ar.storage.DB.QueryRow(`
		SELECT id, email, username, password_hash, created_at
		FROM users
		WHERE email = ?
	`, email).Scan(
		&found.ID,
		&found.Email,
		&found.Username,
		&found.PasswordHash,
		&found.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	return &found, nil
}
