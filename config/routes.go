package config

import (
	"homedb/handlers"
	"homedb/middleware"
	"homedb/repository"
	"homedb/views/pages"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gorilla/sessions"
)

func SetupRoutes(mux *http.ServeMux, repo *repository.Queries, store *sessions.CookieStore) {
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	indexStack := middleware.CreateStack(
		middleware.NotFound,
		middleware.Protected(store),
	)
	mux.Handle("/", indexStack(handlers.Home(repo, store)))

	mux.Handle("GET /login", templ.Handler(pages.Login(nil)))
	mux.Handle("POST /login", handlers.Login(repo, store))

	mux.Handle("GET /signup", templ.Handler(pages.Signup(nil)))
	mux.Handle("POST /signup", handlers.Signup(repo, store))

	mux.Handle("GET /logout", handlers.Logout(store))

	mux.Handle("GET /add", middleware.Protected(store)(templ.Handler(pages.Add(nil))))
	mux.Handle("POST /add", middleware.Protected(store)(handlers.Add(repo, store)))
}
