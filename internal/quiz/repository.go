package quiz

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"nilchan-hackaton/internal/storage"
)

type Repository struct {
	storage *storage.Storage
}

func NewRepository(storage *storage.Storage) *Repository {
	return &Repository{storage: storage}
}

func (r *Repository) Find(ctx context.Context, userID, resourceID int64) (*Quiz, error) {
	var quizID sql.NullInt64
	var title sql.NullString
	var questionsJSON sql.NullString

	err := r.storage.DB.QueryRowContext(ctx, `
		SELECT q.id, q.title, q.questions_json
		FROM resources r
		LEFT JOIN quizzes q ON q.resource_id = r.id
		WHERE r.id = ? AND r.user_id = ?
	`, resourceID, userID).Scan(&quizID, &title, &questionsJSON)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find quiz: %w", err)
	}
	if !quizID.Valid || !title.Valid || !questionsJSON.Valid {
		return nil, ErrUnavailable
	}

	found := Quiz{
		ID:         quizID.Int64,
		ResourceID: resourceID,
		Title:      title.String,
	}
	if err := json.Unmarshal([]byte(questionsJSON.String), &found.Questions); err != nil {
		return nil, fmt.Errorf("decode quiz questions: %w", err)
	}
	return &found, nil
}
