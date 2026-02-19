package todo

import (
	"database/sql"
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
			title TEXT NOT NULL
		);
		INSERT INTO todos (title) VALUES ('First'), ('Second');
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
	if items[0].Title != "First" || items[1].Title != "Second" {
		t.Fatalf("unexpected order or titles: %#v", items)
	}
}
