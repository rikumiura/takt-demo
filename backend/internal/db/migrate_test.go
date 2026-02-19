package db

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()

	database, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	return database
}

func TestMigrateCreatesTodosWithCompleted(t *testing.T) {
	database := openTestDB(t)
	defer database.Close()

	if err := Migrate(database); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	if !columnExists(t, database, "todos", "completed") {
		t.Fatalf("expected completed column to exist")
	}
}

func TestMigrateAddsCompletedToExistingTable(t *testing.T) {
	database := openTestDB(t)
	defer database.Close()

	_, err := database.Exec(`
		CREATE TABLE todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL
		);
	`)
	if err != nil {
		t.Fatalf("create old schema: %v", err)
	}

	if err := Migrate(database); err != nil {
		t.Fatalf("migrate old schema: %v", err)
	}

	if !columnExists(t, database, "todos", "completed") {
		t.Fatalf("expected completed column to exist")
	}
}

func TestSeedIfEmptyInsertsInitialTodos(t *testing.T) {
	database := openTestDB(t)
	defer database.Close()

	if err := Migrate(database); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if err := SeedIfEmpty(database); err != nil {
		t.Fatalf("seed: %v", err)
	}

	var count int
	if err := database.QueryRow(`SELECT COUNT(*) FROM todos`).Scan(&count); err != nil {
		t.Fatalf("count rows: %v", err)
	}
	if count != 3 {
		t.Fatalf("expected 3 rows, got %d", count)
	}

	var completedCount int
	if err := database.QueryRow(`SELECT COUNT(*) FROM todos WHERE completed = 1`).Scan(&completedCount); err != nil {
		t.Fatalf("count completed rows: %v", err)
	}
	if completedCount != 0 {
		t.Fatalf("expected no completed seeds, got %d", completedCount)
	}
}

func columnExists(t *testing.T, database *sql.DB, tableName string, columnName string) bool {
	t.Helper()

	ok, err := hasColumn(database, tableName, columnName)
	if err != nil {
		t.Fatalf("check column existence: %v", err)
	}
	return ok
}
