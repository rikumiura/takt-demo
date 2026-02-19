package db

import "database/sql"

func Migrate(database *sql.DB) error {
	_, err := database.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL
		);
	`)
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
		INSERT INTO todos (title)
		VALUES
			('Buy milk'),
			('Read Go docs'),
			('Build TODO list UI');
	`)
	return err
}
