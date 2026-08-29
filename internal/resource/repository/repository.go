package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"nilchan-hackaton/internal/resource"
	"nilchan-hackaton/internal/storage"
)

type Repository struct {
	storage *storage.Storage
}

func New(storage *storage.Storage) *Repository {
	return &Repository{storage: storage}
}

func (r *Repository) FindByURL(ctx context.Context, userID int64, resourceURL string) (*resource.Resource, error) {
	row := r.storage.DB.QueryRowContext(ctx, `
		SELECT id, user_id, url, title, content, status, purchased_overflow_slot,
		       created_at, completed_at, xp_earned, e_points_earned
		FROM resources WHERE user_id = ? AND url = ?
	`, userID, resourceURL)
	found, err := scanResource(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find resource: %w", err)
	}
	found.Tags, err = r.loadTags(ctx, found.ID)
	if err != nil {
		return nil, err
	}
	return found, nil
}

func (r *Repository) FindByID(ctx context.Context, userID, resourceID int64) (*resource.Summary, error) {
	row := r.storage.DB.QueryRowContext(ctx, `
		SELECT id, url, title, status, created_at, completed_at, xp_earned, e_points_earned
		FROM resources WHERE id = ? AND user_id = ?
	`, resourceID, userID)
	found, err := scanSummary(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find resource by ID: %w", err)
	}
	found.Tags, err = r.loadTags(ctx, found.ID)
	if err != nil {
		return nil, err
	}
	return found, nil
}

func (r *Repository) List(ctx context.Context, userID int64, status, tag string) ([]resource.Summary, error) {
	rows, err := r.storage.DB.QueryContext(ctx, `
		SELECT r.id, r.url, r.title, r.status, r.created_at, r.completed_at,
		       r.xp_earned, r.e_points_earned
		FROM resources r
		WHERE r.user_id = ?
		  AND (? = '' OR r.status = ?)
		  AND (? = '' OR EXISTS (
		      SELECT 1 FROM resource_tags rt
		      WHERE rt.resource_id = r.id AND rt.tag = ?
		  ))
		ORDER BY r.created_at DESC, r.id DESC
	`, userID, status, status, tag, tag)
	if err != nil {
		return nil, fmt.Errorf("list resources: %w", err)
	}
	result := make([]resource.Summary, 0)
	for rows.Next() {
		found, err := scanSummary(rows)
		if err != nil {
			rows.Close()
			return nil, fmt.Errorf("scan resource: %w", err)
		}
		result = append(result, *found)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, fmt.Errorf("iterate resources: %w", err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("close resource rows: %w", err)
	}

	for index := range result {
		tags, err := r.loadTags(ctx, result[index].ID)
		if err != nil {
			return nil, err
		}
		result[index].Tags = tags
	}
	return result, nil
}

func (r *Repository) loadTags(ctx context.Context, resourceID int64) ([]string, error) {
	rows, err := r.storage.DB.QueryContext(ctx, "SELECT tag FROM resource_tags WHERE resource_id = ? ORDER BY tag", resourceID)
	if err != nil {
		return nil, fmt.Errorf("load resource tags: %w", err)
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, fmt.Errorf("scan resource tag: %w", err)
		}
		tags = append(tags, tag)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate resource tags: %w", err)
	}
	return tags, nil
}

type rowScanner interface {
	Scan(...any) error
}

func scanSummary(row rowScanner) (*resource.Summary, error) {
	var value resource.Summary
	err := row.Scan(
		&value.ID,
		&value.URL,
		&value.Title,
		&value.Status,
		&value.CreatedAt,
		&value.CompletedAt,
		&value.XPEarned,
		&value.EPointsEarned,
	)
	return &value, err
}

func scanResource(row rowScanner) (*resource.Resource, error) {
	var value resource.Resource
	err := row.Scan(
		&value.ID,
		&value.UserID,
		&value.URL,
		&value.Title,
		&value.Content,
		&value.Status,
		&value.PurchasedOverflowSlot,
		&value.CreatedAt,
		&value.CompletedAt,
		&value.XPEarned,
		&value.EPointsEarned,
	)
	return &value, err
}
