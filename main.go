package main

import (
	"database/sql"
	"fmt"
	"homedb/handlers"
	"homedb/middleware"
	"homedb/repository"
	"homedb/views/pages"
	"log"
	"net/http"
	"os"

	"github.com/a-h/templ"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("DB_STRING"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repository.New(db)

	router := http.NewServeMux()

	fs := http.FileServer(http.Dir("./static"))
	router.Handle("GET /static/", http.StripPrefix("/static/", fs))

	router.Handle("GET /", handlers.Home(repo))

	router.Handle("GET /login", templ.Handler(pages.Login(nil)))
	router.Handle("POST /login", handlers.Login(repo))

	router.Handle("GET /signup", templ.Handler(pages.Signup(nil)))
	router.Handle("POST /signup", handlers.Signup(repo))

	router.Handle("GET /logout", handlers.Logout())

	protectedRoutes := http.NewServeMux()
	protectedRoutes.HandleFunc("GET /protected", func(w http.ResponseWriter, r *http.Request) {})

	stack := middleware.CreateStack(
		middleware.Logging,
	)

	server := http.Server{
		Addr:    ":8080",
		Handler: stack(router),
	}

	fmt.Println("Server listening on port 8080")
	server.ListenAndServe()
}
