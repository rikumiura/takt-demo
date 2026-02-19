package todo

import (
	"database/sql"
	"errors"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List() ([]Item, error) {
	rows, err := r.db.Query(`SELECT id, title, completed FROM todos ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Item, 0)
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Title, &item.Completed); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Repository) Create(title string) (Item, error) {
	result, err := r.db.Exec(`INSERT INTO todos (title, completed) VALUES (?, ?)`, title, false)
	if err != nil {
		return Item{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Item{}, err
	}

	return Item{
		ID:        id,
		Title:     title,
		Completed: false,
	}, nil
}

func (r *Repository) UpdateCompleted(id int64, completed bool) (Item, error) {
	result, err := r.db.Exec(`UPDATE todos SET completed = ? WHERE id = ?`, completed, id)
	if err != nil {
		return Item{}, err
	}

	updatedRows, err := result.RowsAffected()
	if err != nil {
		return Item{}, err
	}
	if updatedRows == 0 {
		return Item{}, ErrNotFound
	}

	var item Item
	if err := r.db.QueryRow(`SELECT id, title, completed FROM todos WHERE id = ?`, id).Scan(&item.ID, &item.Title, &item.Completed); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Item{}, ErrNotFound
		}
		return Item{}, err
	}
	return item, nil
}

func (r *Repository) Delete(id int64) error {
	result, err := r.db.Exec(`DELETE FROM todos WHERE id = ?`, id)
	if err != nil {
		return err
	}

	deletedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if deletedRows == 0 {
		return ErrNotFound
	}

	return nil
}
