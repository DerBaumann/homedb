package main

import (
	"homedb/handlers"
	"homedb/views/pages"
	"log"
	"net/http"

	"github.com/a-h/templ"
)

func main() {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	mux.Handle("GET /", templ.Handler(pages.Home()))

	mux.Handle("GET /login", templ.Handler(pages.Login(nil)))
	mux.HandleFunc("POST /login", handlers.Login)

	mux.Handle("GET /signup", templ.Handler(pages.Signup(nil)))
	mux.HandleFunc("POST /signup", handlers.Signup)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
