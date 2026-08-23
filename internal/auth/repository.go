package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nilchan-hackaton/internal/storage"
	"nilchan-hackaton/internal/user"

	"github.com/mattn/go-sqlite3"
)

type repository struct {
	storage *storage.Storage
}

func NewRepository(storage *storage.Storage) *repository {
	return &repository{storage: storage}
}

func (ar *repository) checkRegistrationConflict(email, username string) error {
	var emailTaken, usernameTaken bool
	if err := ar.storage.DB.QueryRow(`
		SELECT
			EXISTS(SELECT 1 FROM users WHERE email = ?),
			EXISTS(SELECT 1 FROM users WHERE username = ?)
	`, email, username).Scan(&emailTaken, &usernameTaken); err != nil {
		return fmt.Errorf("check registration conflict: %w", err)
	}

	if emailTaken {
		return ErrEmailTaken
	}
	if usernameTaken {
		return ErrUsernameTaken
	}
	return nil
}

func (ar *repository) create(email, username, passwordHash string) (*user.User, error) {
	tx, err := ar.storage.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin user creation: %w", err)
	}
	defer tx.Rollback()

	var created user.User
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
		user.DefaultAvatarID,
		created.ID,
		user.DefaultFrameID,
	); err != nil {
		return nil, fmt.Errorf("grant default cosmetics: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit user creation: %w", err)
	}
	return &created, nil
}

func (ar *repository) findOne(email string) (*user.User, error) {
	var found user.User
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
