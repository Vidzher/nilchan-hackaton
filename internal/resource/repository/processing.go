package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"nilchan-hackaton/internal/quiz"
	"nilchan-hackaton/internal/resource"
)

func (r *Repository) CompleteGeneration(ctx context.Context, resourceID int64, title string, questions []quiz.Question) error {
	tx, err := r.storage.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin generation completion: %w", err)
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, `
		UPDATE resources SET status = ? WHERE id = ? AND status = ?
	`, resource.StatusNotCompleted, resourceID, resource.StatusProcessing)
	if err != nil {
		return fmt.Errorf("complete resource generation: %w", err)
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("inspect resource completion: %w", err)
	}
	if changed != 1 {
		return resource.ErrStateConflict
	}

	result, err = tx.ExecContext(ctx, "INSERT INTO quizzes(resource_id, title, questions_json) VALUES(?, ?, '[]')", resourceID, title)
	if err != nil {
		return fmt.Errorf("create quiz: %w", err)
	}
	quizID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get quiz id: %w", err)
	}
	if err := quiz.AddVerificationData(quizID, questions); err != nil {
		return fmt.Errorf("add quiz verification data: %w", err)
	}

	payload, err := json.Marshal(questions)
	if err != nil {
		return fmt.Errorf("encode quiz questions: %w", err)
	}
	if _, err := tx.ExecContext(ctx, "UPDATE quizzes SET questions_json = ? WHERE id = ?", string(payload), quizID); err != nil {
		return fmt.Errorf("save quiz questions: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit generation completion: %w", err)
	}
	return nil
}

func (r *Repository) RecoverProcessing(ctx context.Context) (int64, error) {
	tx, err := r.storage.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin processing recovery: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `
		UPDATE user_progress
		SET e_points = e_points + ? * (
			SELECT COUNT(*) FROM resources
			WHERE user_id = user_progress.user_id
			  AND status = ?
			  AND purchased_overflow_slot = 1
		)
		WHERE EXISTS (
			SELECT 1 FROM resources
			WHERE user_id = user_progress.user_id
			  AND status = ?
			  AND purchased_overflow_slot = 1
		)
	`, overflowSlotPrice, resource.StatusProcessing, resource.StatusProcessing); err != nil {
		return 0, fmt.Errorf("refund interrupted overflow slots: %w", err)
	}

	result, err := tx.ExecContext(ctx, `
		UPDATE resources
		SET status = ?, purchased_overflow_slot = 0
		WHERE status = ?
	`, resource.StatusFailed, resource.StatusProcessing)
	if err != nil {
		return 0, fmt.Errorf("recover processing resources: %w", err)
	}
	recovered, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("inspect processing recovery: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit processing recovery: %w", err)
	}
	return recovered, nil
}

func (r *Repository) FailGeneration(ctx context.Context, resourceID int64) error {
	tx, err := r.storage.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin generation failure: %w", err)
	}
	defer tx.Rollback()

	var userID int64
	var purchased bool
	err = tx.QueryRowContext(ctx, `
		UPDATE resources SET status = ?
		WHERE id = ? AND status = ?
		RETURNING user_id, purchased_overflow_slot
	`, resource.StatusFailed, resourceID, resource.StatusProcessing).Scan(&userID, &purchased)
	if errors.Is(err, sql.ErrNoRows) {
		return resource.ErrStateConflict
	}
	if err != nil {
		return fmt.Errorf("fail resource generation: %w", err)
	}
	if purchased {
		if _, err := tx.ExecContext(ctx, "UPDATE user_progress SET e_points = e_points + ? WHERE user_id = ?", overflowSlotPrice, userID); err != nil {
			return fmt.Errorf("refund overflow slot: %w", err)
		}
		if _, err := tx.ExecContext(ctx, "UPDATE resources SET purchased_overflow_slot = 0 WHERE id = ?", resourceID); err != nil {
			return fmt.Errorf("clear overflow purchase: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit generation failure: %w", err)
	}
	return nil
}
