package config

import (
	"homedb/internal/handlers"
	"homedb/internal/middleware"
	"homedb/internal/repository"
	"homedb/internal/views/pages"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gorilla/sessions"
)

func SetupRoutes(mux *http.ServeMux, repo *repository.Queries, store *sessions.CookieStore) {
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	mux.Handle("GET /{$}", middleware.Protected(store)(handlers.Home(repo, store)))

	mux.Handle("GET /login", templ.Handler(pages.Login(nil)))
	mux.Handle("POST /login", handlers.Login(repo, store))

	mux.Handle("GET /signup", templ.Handler(pages.Signup(nil)))
	mux.Handle("POST /signup", handlers.Signup(repo, store))

	mux.Handle("GET /logout", handlers.Logout(store))

	itemRouter := http.NewServeMux()
	itemRouter.Handle("GET /new", templ.Handler(pages.Add(nil)))
	itemRouter.Handle("POST /new", handlers.Add(repo, store))

	itemRouter.Handle("GET /{id}/edit", handlers.EditItemPage(repo))
	itemRouter.Handle("POST /{id}/edit", handlers.EditItem(repo))

	itemRouter.Handle("GET /{id}/delete", handlers.DeleteItem(repo))

	mux.Handle("/items/", http.StripPrefix("/items", middleware.Protected(store)(itemRouter)))
}
