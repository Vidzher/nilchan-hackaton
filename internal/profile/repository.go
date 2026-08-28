package profile

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nilchan-hackaton/internal/storage"
)

type repository struct {
	db *sql.DB
}

func NewRepository(storage *storage.Storage) *repository {
	return &repository{db: storage.DB}
}

func (r *repository) GetProfile(
	ctx context.Context,
	userID int64,
) (*Profile, error) {
	result := &Profile{}

	err := r.db.QueryRowContext(ctx, `
		SELECT
			u.id,
			u.email,
			u.username,
			u.password_hash,
			u.created_at,
			p.user_id,
			p.xp,
			p.e_points,
			p.current_streak,
			p.last_completion_at,
			p.avatar_id,
			p.frame_id,
			p.title_id,
			p.showcase_item_id
		FROM users u
		JOIN user_progress p
			ON p.user_id = u.id
		WHERE u.id = ?
	`,
		userID,
	).Scan(
		&result.User.ID,
		&result.User.Email,
		&result.User.Username,
		&result.User.PasswordHash,
		&result.User.CreatedAt,
		&result.Progress.UserID,
		&result.Progress.XP,
		&result.Progress.EPoints,
		&result.Progress.CurrentStreak,
		&result.Progress.LastCompletionAt,
		&result.Progress.AvatarID,
		&result.Progress.FrameID,
		&result.Progress.TitleID,
		&result.Progress.ShowcaseItemID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProfileNotFound
		}
		return nil, fmt.Errorf("get profile user and progress: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT user_id, item_id
		FROM user_cosmetics
		WHERE user_id = ?
		ORDER BY item_id
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("get owned cosmetics: %w", err)
	}
	defer rows.Close()

	result.Cosmetics = make([]UserCosmetic, 0)
	for rows.Next() {
		var cosmetic UserCosmetic
		if err := rows.Scan(
			&cosmetic.UserID,
			&cosmetic.ItemID,
		); err != nil {
			return nil, fmt.Errorf("scan owned cosmetic: %w", err)
		}
		result.Cosmetics = append(result.Cosmetics, cosmetic)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate owned cosmetics: %w", err)
	}

	return result, nil
}

func (r *repository) UpdateCosmetics(
	ctx context.Context,
	userID int64,
	update CosmeticsUpdate,
) error {
	setClauses := make([]string, 0, 4)
	args := make([]any, 0, 5)
	if update.AvatarID != nil {
		setClauses = append(setClauses, "avatar_id = ?")
		args = append(args, *update.AvatarID)
	}
	if update.FrameID != nil {
		setClauses = append(setClauses, "frame_id = ?")
		args = append(args, *update.FrameID)
	}
	if update.TitleID.Set {
		setClauses = append(setClauses, "title_id = ?")
		args = append(args, update.TitleID.Value)
	}
	if update.ShowcaseItemID.Set {
		setClauses = append(setClauses, "showcase_item_id = ?")
		args = append(args, update.ShowcaseItemID.Value)
	}
	if len(setClauses) == 0 {
		return nil
	}

	args = append(args, userID)
	result, err := r.db.ExecContext(
		ctx,
		"UPDATE user_progress SET "+strings.Join(setClauses, ", ")+" WHERE user_id = ?",
		args...,
	)
	if err != nil {
		return fmt.Errorf("update profile cosmetics: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get updated profile rows: %w", err)
	}
	if rowsAffected == 0 {
		return ErrProfileNotFound
	}

	return nil
}
