package todo

import "errors"

var ErrNotFound = errors.New("todo not found")

type Item struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}
