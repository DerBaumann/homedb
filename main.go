package main

import (
	"database/sql"
	"fmt"
	"homedb/config"
	"homedb/middleware"
	"homedb/repository"
	"net/http"
	"os"
)

func run() error {
	mux := http.NewServeMux()

	db, err := sql.Open("postgres", os.Getenv("DB_STRING"))
	if err != nil {
		return err
	}
	defer db.Close()

	repo := repository.New(db)

	config.SetupRoutes(mux, repo)

	stack := middleware.CreateStack(
		middleware.Logging,
	)

	server := http.Server{
		Addr:    ":8080",
		Handler: stack(mux),
	}

	fmt.Println("Server listening on port 8080")
	server.ListenAndServe()

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}
