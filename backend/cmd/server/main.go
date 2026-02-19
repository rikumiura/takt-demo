package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	_ "modernc.org/sqlite"

	"todoapp/backend/internal/db"
	"todoapp/backend/internal/todo"
)

func main() {
	var (
		addr   = flag.String("addr", ":8080", "server listen address")
		dbPath = flag.String("db", "./todo.db", "sqlite database path")
	)
	flag.Parse()

	database, err := sql.Open("sqlite", *dbPath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer database.Close()

	if err := db.Migrate(database); err != nil {
		log.Fatalf("migrate db: %v", err)
	}
	if err := db.SeedIfEmpty(database); err != nil {
		log.Fatalf("seed db: %v", err)
	}

	repo := todo.NewRepository(database)
	handler := todo.NewHandler(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/todos", handler.ListTodos)
	mux.HandleFunc("POST /api/todos", handler.CreateTodo)
	mux.HandleFunc("PATCH /api/todos/{id}", handler.UpdateTodo)
	mux.HandleFunc("DELETE /api/todos/{id}", handler.DeleteTodo)

	server := &http.Server{
		Addr:    *addr,
		Handler: withCORS(mux),
	}

	log.Printf("server started on %s", *addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
