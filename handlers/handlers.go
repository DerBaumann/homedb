package handlers

import (
	"homedb/repository"
	"homedb/sessions"
	"homedb/utils"
	"homedb/views/pages"
	"net/http"
)

func Home(repo *repository.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := sessions.Get(r)
		if err != nil {
			utils.WriteError(w, r, 401, err)
			return
		}

		user, err := repo.GetUserByID(r.Context(), session.ID)
		if err != nil {
			utils.WriteError(w, r, 401, err)
			return
		}

		items, err := repo.ListItems(r.Context(), user.ID)
		if err != nil {
			utils.WriteError(w, r, 401, err)
			return
		}

		pages.Home(items, nil).Render(r.Context(), w)
	})
}
