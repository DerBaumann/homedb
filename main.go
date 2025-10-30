package main

import (
	"fmt"
	"homedb/config"
	"homedb/middleware"
	"net/http"
	"os"
)

func run() error {
	mux := http.NewServeMux()

	dependencies, err := config.SetupDependencies()
	if err != nil {
		return err
	}

	config.SetupRoutes(mux, dependencies)

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
