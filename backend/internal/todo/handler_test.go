package todo

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeRepo struct {
	listItems []Item
	listErr   error

	createItem  Item
	createErr   error
	createTitle string

	updateItem      Item
	updateErr       error
	updateID        int64
	updateCompleted bool

	deleteErr error
	deleteID  int64
}

func (f *fakeRepo) List() ([]Item, error) {
	if f.listErr != nil {
		return nil, f.listErr
	}
	return f.listItems, nil
}

func (f *fakeRepo) Create(title string) (Item, error) {
	f.createTitle = title
	if f.createErr != nil {
		return Item{}, f.createErr
	}
	return f.createItem, nil
}

func (f *fakeRepo) UpdateCompleted(id int64, completed bool) (Item, error) {
	f.updateID = id
	f.updateCompleted = completed
	if f.updateErr != nil {
		return Item{}, f.updateErr
	}
	return f.updateItem, nil
}

func (f *fakeRepo) Delete(id int64) error {
	f.deleteID = id
	return f.deleteErr
}

func TestListTodos_Success(t *testing.T) {
	repo := &fakeRepo{listItems: []Item{{ID: 1, Title: "test", Completed: true}}}
	h := NewHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/todos", nil)
	rr := httptest.NewRecorder()

	h.ListTodos(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if got := rr.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected application/json, got %q", got)
	}
	if rr.Body.String() != "[{\"id\":1,\"title\":\"test\",\"completed\":true}]\n" {
		t.Fatalf("unexpected body: %q", rr.Body.String())
	}
}

func TestListTodos_Error(t *testing.T) {
	repo := &fakeRepo{listErr: errors.New("boom")}
	h := NewHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/todos", nil)
	rr := httptest.NewRecorder()

	h.ListTodos(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rr.Code)
	}
}

func TestCreateTodo_Success(t *testing.T) {
	repo := &fakeRepo{
		createItem: Item{ID: 3, Title: "created", Completed: false},
	}
	h := NewHandler(repo)

	req := httptest.NewRequest(http.MethodPost, "/api/todos", strings.NewReader(`{"title":"created"}`))
	rr := httptest.NewRecorder()

	h.CreateTodo(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rr.Code)
	}
	if repo.createTitle != "created" {
		t.Fatalf("expected title to be passed to repo, got %q", repo.createTitle)
	}

	var body Item
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.ID != 3 || body.Title != "created" || body.Completed {
		t.Fatalf("unexpected response body: %#v", body)
	}
}

func TestCreateTodo_InvalidJSON(t *testing.T) {
	repo := &fakeRepo{}
	h := NewHandler(repo)

	req := httptest.NewRequest(http.MethodPost, "/api/todos", strings.NewReader(`{"title":`))
	rr := httptest.NewRecorder()

	h.CreateTodo(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
}

func TestCreateTodo_EmptyTitle(t *testing.T) {
	repo := &fakeRepo{}
	h := NewHandler(repo)

	req := httptest.NewRequest(http.MethodPost, "/api/todos", strings.NewReader(`{"title":"   "}`))
	rr := httptest.NewRecorder()

	h.CreateTodo(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
}

func TestCreateTodo_Error(t *testing.T) {
	repo := &fakeRepo{createErr: errors.New("boom")}
	h := NewHandler(repo)

	req := httptest.NewRequest(http.MethodPost, "/api/todos", strings.NewReader(`{"title":"created"}`))
	rr := httptest.NewRecorder()

	h.CreateTodo(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rr.Code)
	}
}

func TestUpdateTodo_Success(t *testing.T) {
	repo := &fakeRepo{
		updateItem: Item{ID: 2, Title: "updated", Completed: true},
	}
	h := NewHandler(repo)

	req := httptest.NewRequest(http.MethodPatch, "/api/todos/2", strings.NewReader(`{"completed":true}`))
	req.SetPathValue("id", "2")
	rr := httptest.NewRecorder()

	h.UpdateTodo(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if repo.updateID != 2 || !repo.updateCompleted {
		t.Fatalf("expected update args (2, true), got (%d, %v)", repo.updateID, repo.updateCompleted)
	}

	var body Item
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.ID != 2 || body.Title != "updated" || !body.Completed {
		t.Fatalf("unexpected response body: %#v", body)
	}
}

func TestUpdateTodo_InvalidID(t *testing.T) {
	repo := &fakeRepo{}
	h := NewHandler(repo)

	req := httptest.NewRequest(http.MethodPatch, "/api/todos/abc", strings.NewReader(`{"completed":true}`))
	req.SetPathValue("id", "abc")
	rr := httptest.NewRecorder()

	h.UpdateTodo(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
}

func TestUpdateTodo_InvalidJSON(t *testing.T) {
	repo := &fakeRepo{}
	h := NewHandler(repo)

	req := httptest.NewRequest(http.MethodPatch, "/api/todos/1", strings.NewReader(`{"completed":`))
	req.SetPathValue("id", "1")
	rr := httptest.NewRecorder()

	h.UpdateTodo(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
}

func TestUpdateTodo_MissingCompleted(t *testing.T) {
	repo := &fakeRepo{}
	h := NewHandler(repo)

	req := httptest.NewRequest(http.MethodPatch, "/api/todos/1", strings.NewReader(`{"title":"nope"}`))
	req.SetPathValue("id", "1")
	rr := httptest.NewRecorder()

	h.UpdateTodo(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
}

func TestUpdateTodo_NotFound(t *testing.T) {
	repo := &fakeRepo{updateErr: ErrNotFound}
	h := NewHandler(repo)

	req := httptest.NewRequest(http.MethodPatch, "/api/todos/10", strings.NewReader(`{"completed":true}`))
	req.SetPathValue("id", "10")
	rr := httptest.NewRecorder()

	h.UpdateTodo(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rr.Code)
	}
}

func TestUpdateTodo_Error(t *testing.T) {
	repo := &fakeRepo{updateErr: errors.New("boom")}
	h := NewHandler(repo)

	req := httptest.NewRequest(http.MethodPatch, "/api/todos/10", strings.NewReader(`{"completed":true}`))
	req.SetPathValue("id", "10")
	rr := httptest.NewRecorder()

	h.UpdateTodo(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rr.Code)
	}
}

func TestDeleteTodo_Success(t *testing.T) {
	repo := &fakeRepo{}
	h := NewHandler(repo)

	req := httptest.NewRequest(http.MethodDelete, "/api/todos/5", nil)
	req.SetPathValue("id", "5")
	rr := httptest.NewRecorder()

	h.DeleteTodo(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rr.Code)
	}
	if repo.deleteID != 5 {
		t.Fatalf("expected delete id 5, got %d", repo.deleteID)
	}
}

func TestDeleteTodo_InvalidID(t *testing.T) {
	repo := &fakeRepo{}
	h := NewHandler(repo)

	req := httptest.NewRequest(http.MethodDelete, "/api/todos/abc", nil)
	req.SetPathValue("id", "abc")
	rr := httptest.NewRecorder()

	h.DeleteTodo(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
}

func TestDeleteTodo_NotFound(t *testing.T) {
	repo := &fakeRepo{deleteErr: ErrNotFound}
	h := NewHandler(repo)

	req := httptest.NewRequest(http.MethodDelete, "/api/todos/42", nil)
	req.SetPathValue("id", "42")
	rr := httptest.NewRecorder()

	h.DeleteTodo(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rr.Code)
	}
}

func TestDeleteTodo_Error(t *testing.T) {
	repo := &fakeRepo{deleteErr: errors.New("boom")}
	h := NewHandler(repo)

	req := httptest.NewRequest(http.MethodDelete, "/api/todos/42", nil)
	req.SetPathValue("id", "42")
	rr := httptest.NewRecorder()

	h.DeleteTodo(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rr.Code)
	}
}
