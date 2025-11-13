package main

import (
	"database/sql"
	"fmt"
	"homedb/config"
	"homedb/middleware"
	"homedb/repository"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

func NotFound(mux *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h, pattern := mux.Handler(r)
		if pattern == "" {
			return
		}
		h.ServeHTTP(w, r)
	})
}

func run() error {
	mux := http.NewServeMux()

	store := sessions.NewCookieStore(
		[]byte(os.Getenv("SESSION_AUTH_KEY")),
		[]byte(os.Getenv("SESSION_ENCRYPTION_KEY")),
	)

	db, err := sql.Open("postgres", os.Getenv("DB_STRING"))
	if err != nil {
		return err
	}
	defer db.Close()

	repo := repository.New(db)

	config.SetupRoutes(mux, repo, store)

	stack := middleware.CreateStack(
		middleware.Logging,
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: stack(mux),
	}

	fmt.Printf("Server listening on port %s\n", port)
	server.ListenAndServe()

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}
