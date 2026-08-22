package storage

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestNewStorageConfiguresSQLiteAndAppliesSchema(t *testing.T) {
	path := filepath.Join(t.TempDir(), "learning-backlog.db")

	store, err := NewStorage(path)
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}

	var foreignKeys int
	if err := store.DB.QueryRow("PRAGMA foreign_keys").Scan(&foreignKeys); err != nil {
		t.Fatalf("read foreign_keys: %v", err)
	}
	if foreignKeys != 1 {
		t.Fatalf("foreign_keys = %d, want 1", foreignKeys)
	}

	var journalMode string
	if err := store.DB.QueryRow("PRAGMA journal_mode").Scan(&journalMode); err != nil {
		t.Fatalf("read journal_mode: %v", err)
	}
	if strings.ToLower(journalMode) != "wal" {
		t.Fatalf("journal_mode = %q, want wal", journalMode)
	}

	var busyTimeout int
	if err := store.DB.QueryRow("PRAGMA busy_timeout").Scan(&busyTimeout); err != nil {
		t.Fatalf("read busy_timeout: %v", err)
	}
	if busyTimeout != 5000 {
		t.Fatalf("busy_timeout = %d, want 5000", busyTimeout)
	}

	_, err = store.DB.Exec(`
		INSERT INTO resources (
			id, user_id, url, title, content, status, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`, "resource-1", 999, "https://example.com", "Title", "content", "PROCESSING", "2026-01-01T00:00:00Z")
	if err == nil {
		t.Fatal("inserting a resource for a missing user succeeded")
	}

	if err := store.Close(); err != nil {
		t.Fatalf("close storage: %v", err)
	}

	store, err = NewStorage(path)
	if err != nil {
		t.Fatalf("reopen storage: %v", err)
	}
	defer store.Close()
}
