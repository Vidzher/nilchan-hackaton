package repository

import (
	"context"
	"errors"
	"fmt"

	"nilchan-hackaton/internal/resource"

	"github.com/mattn/go-sqlite3"
)

func (r *Repository) CheckCapacity(ctx context.Context, userID int64, purchaseOverflow bool) error {
	xp, ePoints, used, err := loadCapacity(ctx, r.storage.DB, userID)
	if err != nil {
		return err
	}
	return capacityError(xp, ePoints, used, purchaseOverflow)
}

func (r *Repository) CreateProcessing(ctx context.Context, draft resource.Resource, purchaseOverflow bool) (*resource.Resource, error) {
	tx, err := r.storage.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin resource creation: %w", err)
	}
	defer tx.Rollback()

	purchased, err := reserveCapacity(ctx, tx, draft.UserID, purchaseOverflow)
	if err != nil {
		return nil, err
	}

	draft.Status = resource.StatusProcessing
	draft.PurchasedOverflowSlot = purchased
	err = tx.QueryRowContext(ctx, `
		INSERT INTO resources(user_id, url, title, content, status, purchased_overflow_slot, created_at)
		VALUES(?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
		RETURNING id, created_at
	`, draft.UserID, draft.URL, draft.Title, draft.Content, draft.Status, purchased).Scan(&draft.ID, &draft.CreatedAt)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return nil, resource.ErrDuplicate
		}
		return nil, fmt.Errorf("create resource: %w", err)
	}

	for _, tag := range draft.Tags {
		if _, err := tx.ExecContext(ctx, "INSERT INTO resource_tags(resource_id, tag) VALUES(?, ?)", draft.ID, tag); err != nil {
			return nil, fmt.Errorf("create resource tag: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit resource creation: %w", err)
	}
	return &draft, nil
}

func (r *Repository) RetryFailed(ctx context.Context, resourceID, userID int64, purchaseOverflow bool) (bool, error) {
	tx, err := r.storage.DB.BeginTx(ctx, nil)
	if err != nil {
		return false, fmt.Errorf("begin resource retry: %w", err)
	}
	defer tx.Rollback()

	overflowSlotPurchased, err := reserveCapacity(ctx, tx, userID, purchaseOverflow)
	if err != nil {
		return false, err
	}
	result, err := tx.ExecContext(ctx, `
		UPDATE resources SET status = ?, purchased_overflow_slot = ?
		WHERE id = ? AND user_id = ? AND status = ?
	`, resource.StatusProcessing, overflowSlotPurchased, resourceID, userID, resource.StatusFailed)
	if err != nil {
		return false, fmt.Errorf("retry resource: %w", err)
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("inspect resource retry: %w", err)
	}
	if changed == 0 {
		return false, resource.ErrStateConflict
	}
	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("commit resource retry: %w", err)
	}
	return overflowSlotPurchased, nil
}
