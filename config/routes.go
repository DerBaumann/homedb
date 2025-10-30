package config

import (
	"homedb/handlers"
	"homedb/middleware"
	"homedb/repository"
	"homedb/views/pages"
	"net/http"

	"github.com/a-h/templ"
)

func SetupRoutes(mux *http.ServeMux, repo *repository.Queries) {
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	mux.Handle("GET /", handlers.Home(repo))

	mux.Handle("GET /login", templ.Handler(pages.Login(nil)))
	mux.Handle("POST /login", handlers.Login(repo))

	mux.Handle("GET /signup", templ.Handler(pages.Signup(nil)))
	mux.Handle("POST /signup", handlers.Signup(repo))

	mux.Handle("GET /logout", handlers.Logout())
	mux.Handle("GET /protected", middleware.Protected(templ.Handler(pages.ProtectedPage())))
}
