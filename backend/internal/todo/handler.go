package todo

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type ReaderWriter interface {
	List() ([]Item, error)
	Create(title string) (Item, error)
	Update(id int64, title *string, completed *bool) (Item, error)
	Delete(id int64) error
}

type Handler struct {
	repo ReaderWriter
}

func NewHandler(repo ReaderWriter) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) ListTodos(w http.ResponseWriter, _ *http.Request) {
	items, err := h.repo.List()
	if err != nil {
		http.Error(w, "failed to fetch todos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

type createTodoRequest struct {
	Title string `json:"title"`
}

func (h *Handler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req createTodoRequest
	if err := decodeJSON(r, &req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	item, err := h.repo.Create(title)
	if err != nil {
		http.Error(w, "failed to create todo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(item); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

type updateTodoRequest struct {
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
}

func (h *Handler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid todo id", http.StatusBadRequest)
		return
	}

	var req updateTodoRequest
	if err := decodeJSON(r, &req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Title == nil && req.Completed == nil {
		http.Error(w, "title or completed is required", http.StatusBadRequest)
		return
	}

	var title *string
	if req.Title != nil {
		trimmedTitle := strings.TrimSpace(*req.Title)
		if trimmedTitle == "" {
			http.Error(w, "title is required", http.StatusBadRequest)
			return
		}
		title = &trimmedTitle
	}

	item, err := h.repo.Update(id, title, req.Completed)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, "todo not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to update todo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(item); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid todo id", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(id); err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, "todo not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to delete todo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func decodeJSON(r *http.Request, target interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return errors.New("request body must contain only a single JSON object")
	}
	return nil
}

func parseID(raw string) (int64, error) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, err
	}
	if id <= 0 {
		return 0, errors.New("todo id must be positive")
	}
	return id, nil
}
