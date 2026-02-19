package todo

import (
	"encoding/json"
	"net/http"
)

type Lister interface {
	List() ([]Item, error)
}

type Handler struct {
	repo Lister
}

func NewHandler(repo Lister) *Handler {
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
