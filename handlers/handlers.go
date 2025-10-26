package handlers

import (
	"database/sql"
	"homedb/repository"
	"homedb/views/pages"
	"net/http"
	"os"
)

func Home(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", os.Getenv("DB_STRING"))
	if err != nil {
		pages.Home(nil, err).Render(r.Context(), w)
		return
	}
	defer db.Close()

	repo := repository.New(db)

	users, err := repo.ListUsers(r.Context())
	if err != nil {
		pages.Home(nil, err).Render(r.Context(), w)
		return
	}

	pages.Home(users, nil).Render(r.Context(), w)
}
