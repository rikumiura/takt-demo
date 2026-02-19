package todo

import (
	"database/sql"
	"errors"
	"testing"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			completed INTEGER NOT NULL DEFAULT 0
		);
		INSERT INTO todos (title, completed) VALUES ('First', 0), ('Second', 1);
	`)
	if err != nil {
		t.Fatalf("setup schema: %v", err)
	}

	return db
}

func TestRepositoryList(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)
	items, err := repo.List()
	if err != nil {
		t.Fatalf("list todos: %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("expected 2 todos, got %d", len(items))
	}
	if items[0].Title != "First" || items[0].Completed {
		t.Fatalf("unexpected first todo: %#v", items[0])
	}
	if items[1].Title != "Second" || !items[1].Completed {
		t.Fatalf("unexpected second todo: %#v", items[1])
	}
}

func TestRepositoryCreate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)
	item, err := repo.Create("Created")
	if err != nil {
		t.Fatalf("create todo: %v", err)
	}

	if item.ID <= 0 {
		t.Fatalf("expected positive id, got %d", item.ID)
	}
	if item.Title != "Created" || item.Completed {
		t.Fatalf("unexpected created item: %#v", item)
	}
}

func TestRepositoryUpdateCompleted(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)
	item, err := repo.UpdateCompleted(1, true)
	if err != nil {
		t.Fatalf("update completed: %v", err)
	}

	if item.ID != 1 || item.Title != "First" || !item.Completed {
		t.Fatalf("unexpected updated item: %#v", item)
	}
}

func TestRepositoryUpdateCompleted_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)
	_, err := repo.UpdateCompleted(999, true)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestRepositoryDelete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)
	if err := repo.Delete(1); err != nil {
		t.Fatalf("delete todo: %v", err)
	}

	items, err := repo.List()
	if err != nil {
		t.Fatalf("list todos after delete: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 todo, got %d", len(items))
	}
	if items[0].ID != 2 {
		t.Fatalf("expected remaining id 2, got %d", items[0].ID)
	}
}

func TestRepositoryDelete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)
	err := repo.Delete(999)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
