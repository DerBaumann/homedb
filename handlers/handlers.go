package handlers

import (
	"errors"
	"fmt"
	"homedb/repository"
	"homedb/utils"
	"homedb/views/pages"
	"net/http"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

func Home(repo *repository.Queries, store *sessions.CookieStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.WriteError(w, r, http.StatusMethodNotAllowed, errors.New("method not allowed"))
			return
		}

		q := r.URL.Query().Get("q")

		session, _ := store.Get(r, os.Getenv("SESSION_NAME"))
		userID, _ := uuid.Parse(session.Values["user_id"].(string))

		user, err := repo.GetUserByID(r.Context(), userID)
		if err != nil {
			utils.WriteError(w, r, 401, err)
			return
		}

		items, err := repo.FilterItemsByName(r.Context(), repository.FilterItemsByNameParams{
			UserID: user.ID,
			Name:   fmt.Sprintf("%%%s%%", q),
		})
		if err != nil {
			utils.WriteError(w, r, 401, err)
			return
		}

		pages.Home(items, nil).Render(r.Context(), w)
	})
}

func Add(repo *repository.Queries, store *sessions.CookieStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		session, _ := store.Get(r, os.Getenv("SESSION_NAME"))
		userID, _ := uuid.Parse(session.Values["user_id"].(string))

		_, err = repo.CreateItem(r.Context(), repository.CreateItemParams{
			Name:   name,
			Amount: int32(amount),
			Unit:   repository.ItemUnit(unit),
			UserID: userID,
		})
		if err != nil {
			pages.Add(err).Render(r.Context(), w)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	})
}
