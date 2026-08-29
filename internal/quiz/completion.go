package quiz

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"nilchan-hackaton/internal/progress"
)

func (r *Repository) Complete(
	ctx context.Context,
	userID, resourceID int64,
	xpEarned, ePointsEarned int,
	completedAt time.Time,
) (*CompletionRecord, error) {
	completedAt = completedAt.UTC()

	tx, err := r.storage.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin quiz completion: %w", err)
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, `
		UPDATE resources
		SET status = 'COMPLETED', completed_at = ?, xp_earned = ?, e_points_earned = ?
		WHERE id = ? AND user_id = ? AND status = 'NOT_COMPLETED'
	`, completedAt, xpEarned, ePointsEarned, resourceID, userID)
	if err != nil {
		return nil, fmt.Errorf("complete resource: %w", err)
	}

	changed, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("inspect quiz completion: %w", err)
	}
	if changed == 0 {
		completion, err := loadExistingCompletion(ctx, tx, userID, resourceID)
		if err != nil {
			return nil, err
		}
		if err := tx.Commit(); err != nil {
			return nil, fmt.Errorf("commit existing quiz completion: %w", err)
		}
		return completion, nil
	}

	var xp int64
	var ePoints int64
	var currentStreak int
	var lastCompletionAt sql.NullTime
	if err := tx.QueryRowContext(ctx, `
		SELECT xp, e_points, current_streak, last_completion_at
		FROM user_progress
		WHERE user_id = ?
	`, userID).Scan(&xp, &ePoints, &currentStreak, &lastCompletionAt); err != nil {
		return nil, fmt.Errorf("load user progress for quiz completion: %w", err)
	}

	var previousCompletionAt *time.Time
	if lastCompletionAt.Valid {
		previousCompletionAt = &lastCompletionAt.Time
	}
	currentStreak = progress.NextStreak(currentStreak, previousCompletionAt, completedAt)
	xp += int64(xpEarned)
	ePoints += int64(ePointsEarned)

	result, err = tx.ExecContext(ctx, `
		UPDATE user_progress
		SET xp = ?, e_points = ?, current_streak = ?, last_completion_at = ?
		WHERE user_id = ?
	`, xp, ePoints, currentStreak, completedAt, userID)
	if err != nil {
		return nil, fmt.Errorf("update progress for quiz completion: %w", err)
	}
	changed, err = result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("inspect progress update: %w", err)
	}
	if changed != 1 {
		return nil, fmt.Errorf("update progress for quiz completion: user progress not found")
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit quiz completion: %w", err)
	}

	return &CompletionRecord{
		CompletedAt:   completedAt,
		XPEarned:      xpEarned,
		EPointsEarned: ePointsEarned,
		XP:            xp,
		EPoints:       ePoints,
		CurrentStreak: currentStreak,
	}, nil
}

func loadExistingCompletion(
	ctx context.Context,
	tx *sql.Tx,
	userID, resourceID int64,
) (*CompletionRecord, error) {
	var completedAt sql.NullTime
	var xpEarned sql.NullInt64
	var ePointsEarned sql.NullInt64
	err := tx.QueryRowContext(ctx, `
		SELECT completed_at, xp_earned, e_points_earned
		FROM resources
		WHERE id = ? AND user_id = ?
	`, resourceID, userID).Scan(&completedAt, &xpEarned, &ePointsEarned)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("load existing quiz completion: %w", err)
	}
	if !completedAt.Valid || !xpEarned.Valid || !ePointsEarned.Valid {
		return nil, ErrUnavailable
	}

	var xp int64
	var ePoints int64
	var currentStreak int
	if err := tx.QueryRowContext(ctx, `
		SELECT xp, e_points, current_streak
		FROM user_progress
		WHERE user_id = ?
	`, userID).Scan(&xp, &ePoints, &currentStreak); err != nil {
		return nil, fmt.Errorf("load progress for existing quiz completion: %w", err)
	}

	return &CompletionRecord{
		CompletedAt:   completedAt.Time,
		XPEarned:      int(xpEarned.Int64),
		EPointsEarned: int(ePointsEarned.Int64),
		XP:            xp,
		EPoints:       ePoints,
		CurrentStreak: currentStreak,
	}, nil
}
