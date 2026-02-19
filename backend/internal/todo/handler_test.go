package todo

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeRepo struct {
	items []Item
	err   error
}

func (f fakeRepo) List() ([]Item, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.items, nil
}

func TestListTodos_Success(t *testing.T) {
	h := NewHandler(fakeRepo{items: []Item{{ID: 1, Title: "test"}}})

	req := httptest.NewRequest(http.MethodGet, "/api/todos", nil)
	rr := httptest.NewRecorder()

	h.ListTodos(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if got := rr.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected application/json, got %q", got)
	}
	if rr.Body.String() != "[{\"id\":1,\"title\":\"test\"}]\n" {
		t.Fatalf("unexpected body: %q", rr.Body.String())
	}
}

func TestListTodos_Error(t *testing.T) {
	h := NewHandler(fakeRepo{err: errors.New("boom")})

	req := httptest.NewRequest(http.MethodGet, "/api/todos", nil)
	rr := httptest.NewRecorder()

	h.ListTodos(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rr.Code)
	}
}
