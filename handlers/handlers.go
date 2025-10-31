package handlers

import (
	"errors"
	"homedb/repository"
	"homedb/sessions"
	"homedb/utils"
	"homedb/views/pages"
	"net/http"
	"strconv"
)

func Home(repo *repository.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.WriteError(w, r, http.StatusMethodNotAllowed, errors.New("method not allowed"))
			return
		}

		sess, ok := r.Context().Value(sessions.ContextKey).(*sessions.Session)
		if !ok {
			utils.WriteError(w, r, 401, sessions.ErrUnauthorized)
			return
		}

		user, err := repo.GetUserByID(r.Context(), sess.ID)
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

func ShowAdd() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(sessions.ContextKey).(*sessions.Session)
		if !ok {
			utils.WriteError(w, r, 401, sessions.ErrUnauthorized)
			return
		}

		pages.Add(nil).Render(r.Context(), w)
	})
}

func Add(repo *repository.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, ok := r.Context().Value(sessions.ContextKey).(*sessions.Session)
		if !ok {
			utils.WriteError(w, r, 401, sessions.ErrUnauthorized)
			return
		}

		name := r.FormValue("name")
		rawAmount := r.FormValue("amount")
		unit := r.FormValue("unit")

		if name == "" || rawAmount == "" || unit == "" {
			pages.Add(errors.New("all fields must be filled out")).Render(r.Context(), w)
			return
		}

		amount, err := strconv.Atoi(rawAmount)
		if err != nil {
			pages.Add(err).Render(r.Context(), w)
			return
		}

		_, err = repo.CreateItem(r.Context(), repository.CreateItemParams{
			Name:   name,
			Amount: int32(amount),
			Unit:   repository.ItemUnit(unit),
			UserID: session.ID,
		})
		if err != nil {
			pages.Add(err).Render(r.Context(), w)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	})
}
