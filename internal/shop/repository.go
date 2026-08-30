package shop

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"nilchan-hackaton/internal/storage"

	"github.com/mattn/go-sqlite3"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(storage *storage.Storage) *Repository {
	return &Repository{db: storage.DB}
}

func (r *Repository) Purchase(
	ctx context.Context,
	userID int64,
	itemID string,
	price int64,
) (int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin cosmetic purchase: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(
		ctx,
		"INSERT INTO user_cosmetics(user_id, item_id) VALUES(?, ?)",
		userID,
		itemID,
	); err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) &&
			(sqliteErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey ||
				sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique) {
			return 0, ErrAlreadyOwned
		}
		return 0, fmt.Errorf("grant purchased cosmetic: %w", err)
	}

	var balance int64
	err = tx.QueryRowContext(ctx, `
		UPDATE user_progress
		SET e_points = e_points - ?
		WHERE user_id = ? AND e_points >= ?
		RETURNING e_points
	`, price, userID, price).Scan(&balance)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrInsufficientEPoints
	}
	if err != nil {
		return 0, fmt.Errorf("deduct cosmetic price: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit cosmetic purchase: %w", err)
	}
	return balance, nil
}
