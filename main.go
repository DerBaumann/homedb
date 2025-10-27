package main

import (
	"fmt"
	"homedb/handlers"
	"homedb/middleware"
	"homedb/views/pages"
	"net/http"

	"github.com/a-h/templ"
)

func main() {
	router := http.NewServeMux()

	fs := http.FileServer(http.Dir("./static"))
	router.Handle("GET /static/", http.StripPrefix("/static/", fs))

	router.HandleFunc("GET /", handlers.Home)

	router.Handle("GET /login", templ.Handler(pages.Login(nil)))
	router.HandleFunc("POST /login", handlers.Login)

	router.Handle("GET /signup", templ.Handler(pages.Signup(nil)))
	router.HandleFunc("POST /signup", handlers.Signup)

	router.HandleFunc("GET /logout", handlers.Logout)

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
