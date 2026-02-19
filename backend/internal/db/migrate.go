package db

import (
	"database/sql"
	"fmt"
)

func Migrate(database *sql.DB) error {
	_, err := database.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			completed INTEGER NOT NULL DEFAULT 0
		);
	`)
	if err != nil {
		return err
	}

	hasCompleted, err := hasColumn(database, "todos", "completed")
	if err != nil {
		return err
	}
	if hasCompleted {
		return nil
	}

	_, err = database.Exec(`ALTER TABLE todos ADD COLUMN completed INTEGER NOT NULL DEFAULT 0;`)
	return err
}

func SeedIfEmpty(database *sql.DB) error {
	var count int
	if err := database.QueryRow(`SELECT COUNT(*) FROM todos`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	_, err := database.Exec(`
		INSERT INTO todos (title, completed)
		VALUES
			('Buy milk', 0),
			('Read Go docs', 0),
			('Build TODO list UI', 0);
	`)
	return err
}

func hasColumn(database *sql.DB, tableName string, columnName string) (bool, error) {
	rows, err := database.Query(fmt.Sprintf(`PRAGMA table_info(%s)`, tableName))
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			cid        int
			name       string
			columnType string
			notNull    int
			defaultVal sql.NullString
			primaryKey int
		)
		if err := rows.Scan(&cid, &name, &columnType, &notNull, &defaultVal, &primaryKey); err != nil {
			return false, err
		}
		if name == columnName {
			return true, nil
		}
	}

	if err := rows.Err(); err != nil {
		return false, err
	}
	return false, nil
}
