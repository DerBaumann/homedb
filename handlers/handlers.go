package handlers

import (
	"homedb/repository"
	"homedb/views/pages"
	"net/http"
)

func Home(repo *repository.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		users, err := repo.ListUsers(r.Context())
		if err != nil {
			pages.Home(nil, err).Render(r.Context(), w)
			return
		}

		pages.Home(users, nil).Render(r.Context(), w)
	})
}
