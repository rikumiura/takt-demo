package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithCORS_SetsHeaders(t *testing.T) {
	nextCalled := false
	handler := withCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/todos", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if !nextCalled {
		t.Fatalf("expected next handler to be called")
	}
	if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatalf("unexpected allow origin: %q", rr.Header().Get("Access-Control-Allow-Origin"))
	}
	if rr.Header().Get("Access-Control-Allow-Methods") != "GET, POST, PATCH, DELETE, OPTIONS" {
		t.Fatalf("unexpected allow methods: %q", rr.Header().Get("Access-Control-Allow-Methods"))
	}
	if rr.Header().Get("Access-Control-Allow-Headers") != "Content-Type" {
		t.Fatalf("unexpected allow headers: %q", rr.Header().Get("Access-Control-Allow-Headers"))
	}
}

func TestWithCORS_HandlesPreflight(t *testing.T) {
	nextCalled := false
	handler := withCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodOptions, "/api/todos/1", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if nextCalled {
		t.Fatalf("expected next handler not to be called for preflight")
	}
	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rr.Code)
	}
}
