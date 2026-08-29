package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"nilchan-hackaton/internal/cosmetics"
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

func (r *repository) create(ctx context.Context, email, username, passwordHash string) (*user.User, error) {
	tx, err := r.storage.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin user creation: %w", err)
	}
	defer tx.Rollback()

	var created user.User
	err = tx.QueryRowContext(ctx, `
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
			return nil, ErrUserAlreadyExists
		}
		return nil, fmt.Errorf("create user: %w", err)
	}

	if _, err := tx.ExecContext(ctx, "INSERT INTO user_progress(user_id) VALUES(?)", created.ID); err != nil {
		return nil, fmt.Errorf("create user progress: %w", err)
	}

	if _, err := tx.ExecContext(
		ctx,
		"INSERT INTO user_cosmetics(user_id, item_id) VALUES(?, ?), (?, ?)",
		created.ID,
		cosmetics.DefaultAvatarID,
		created.ID,
		cosmetics.DefaultFrameID,
	); err != nil {
		return nil, fmt.Errorf("grant default cosmetics: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit user creation: %w", err)
	}
	return &created, nil
}

func (r *repository) findByEmail(ctx context.Context, email string) (*user.User, error) {
	var found user.User
	err := r.storage.DB.QueryRowContext(ctx, `
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
