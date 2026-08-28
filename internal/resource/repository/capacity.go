package repository

import (
	"context"
	"database/sql"
	"fmt"

	"nilchan-hackaton/internal/resource"
)

const overflowSlotPrice = 25

type capacityQuerier interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

func loadCapacity(ctx context.Context, querier capacityQuerier, userID int64) (int, int, int, error) {
	var xp, ePoints, used int
	err := querier.QueryRowContext(ctx, `
		SELECT p.xp, p.e_points,
		       (SELECT COUNT(*) FROM resources r WHERE r.user_id = p.user_id AND r.status IN (?, ?))
		FROM user_progress p WHERE p.user_id = ?
	`, resource.StatusProcessing, resource.StatusNotCompleted, userID).Scan(&xp, &ePoints, &used)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("load backlog capacity: %w", err)
	}
	return xp, ePoints, used, nil
}

func capacityError(xp, ePoints, used int, purchase bool) error {
	limit := backlogLimit(xp)
	if used < limit {
		return nil
	}
	if used >= limit+1 || !purchase {
		return resource.ErrBacklogFull
	}
	if ePoints < overflowSlotPrice {
		return resource.ErrInsufficientEPoints
	}
	return nil
}

func backlogLimit(xp int) int {
	switch {
	case xp >= 1000:
		return 10
	case xp >= 600:
		return 8
	case xp >= 300:
		return 7
	case xp >= 120:
		return 6
	default:
		return 5
	}
}

func reserveCapacity(ctx context.Context, tx *sql.Tx, userID int64, purchaseOverflow bool) (bool, error) {
	xp, ePoints, used, err := loadCapacity(ctx, tx, userID)
	if err != nil {
		return false, err
	}
	if err := capacityError(xp, ePoints, used, purchaseOverflow); err != nil {
		return false, err
	}
	if used < backlogLimit(xp) {
		return false, nil
	}

	result, err := tx.ExecContext(ctx, "UPDATE user_progress SET e_points = e_points - ? WHERE user_id = ? AND e_points >= ?", overflowSlotPrice, userID, overflowSlotPrice)
	if err != nil {
		return false, fmt.Errorf("purchase overflow slot: %w", err)
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("inspect overflow purchase: %w", err)
	}
	if changed != 1 {
		return false, resource.ErrInsufficientEPoints
	}
	return true, nil
}
