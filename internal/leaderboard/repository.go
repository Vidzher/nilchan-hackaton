package leaderboard

import (
	"context"
	"database/sql"
	"fmt"

	"nilchan-hackaton/internal/storage"
)

const topLimit = 20

type Repository struct {
	db *sql.DB
}

func NewRepository(storage *storage.Storage) *Repository {
	return &Repository{db: storage.DB}
}

func (r *Repository) List(ctx context.Context, currentUserID int64) ([]Entry, error) {
	rows, err := r.db.QueryContext(ctx, `
		WITH ranked AS (
			SELECT
				ROW_NUMBER() OVER (ORDER BY p.xp DESC, u.id ASC) AS rank,
				u.id,
				u.username,
				p.xp,
				p.avatar_id,
				p.frame_id
			FROM users u
			JOIN user_progress p ON p.user_id = u.id
		)
		SELECT rank, id, username, xp, avatar_id, frame_id
		FROM ranked
		WHERE rank <= ? OR id = ?
		ORDER BY rank
	`, topLimit, currentUserID)
	if err != nil {
		return nil, fmt.Errorf("query leaderboard: %w", err)
	}
	defer rows.Close()

	entries := make([]Entry, 0, topLimit+1)
	for rows.Next() {
		var entry Entry
		if err := rows.Scan(
			&entry.Rank,
			&entry.UserID,
			&entry.Username,
			&entry.XP,
			&entry.AvatarID,
			&entry.FrameID,
		); err != nil {
			return nil, fmt.Errorf("scan leaderboard entry: %w", err)
		}
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate leaderboard: %w", err)
	}
	return entries, nil
}
