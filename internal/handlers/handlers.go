package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"homedb/internal/repository"
	"homedb/internal/utils"
	"homedb/internal/views/pages"

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

		userID, err := uuid.Parse(r.Context().Value("user_id").(string))
		if err != nil {
			utils.WriteError(w, r, http.StatusInternalServerError, err)
			return
		}

		items, err := repo.FilterItemsByName(r.Context(), repository.FilterItemsByNameParams{
			UserID: userID,
			Name:   fmt.Sprintf("%%%s%%", q),
		})
		if err != nil {
			utils.WriteError(w, r, http.StatusInternalServerError, err)
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

		userID, err := uuid.Parse(r.Context().Value("user_id").(string))
		if err != nil {
			utils.WriteError(w, r, 401, err)
			return
		}

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

func EditItemPage(repo *repository.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		itemID, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			utils.WriteError(w, r, 500, err)
			return
		}

		item, err := repo.GetItemByID(r.Context(), itemID)
		if err != nil {
			utils.WriteError(w, r, 500, err)
			return
		}

		pages.Edit(item, nil).Render(r.Context(), w)
	})
}

func EditItem(repo *repository.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		amountStr := r.FormValue("amount")
		unit := r.FormValue("unit")
		itemIDStr := r.PathValue("id")

		itemID, err := uuid.Parse(itemIDStr)
		if err != nil {
			utils.WriteError(w, r, http.StatusInternalServerError, err)
			return
		}

		if name == "" || amountStr == "" || unit == "" {
			item, err := repo.GetItemByID(r.Context(), itemID)
			if err != nil {
				utils.WriteError(w, r, 500, err)
			}

			pages.Edit(item, errors.New("all fields must be filled out"))
		}

		amount, err := strconv.Atoi(amountStr)
		if err != nil {
			item, err := repo.GetItemByID(r.Context(), itemID)
			if err != nil {
				utils.WriteError(w, r, 500, err)
			}

			pages.Edit(item, errors.New("malformed amount"))
		}

		if _, err := repo.UpdateItem(r.Context(), repository.UpdateItemParams{
			Name:   name,
			Amount: int32(amount),
			Unit:   repository.ItemUnit(unit),
			ID:     itemID,
		}); err != nil {
			utils.WriteError(w, r, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	})
}
