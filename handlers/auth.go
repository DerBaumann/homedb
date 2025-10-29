package handlers

import (
	"homedb/repository"
	"homedb/services"
	"homedb/sessions"
	"homedb/views/pages"
	"net/http"

	_ "github.com/lib/pq"
)

func Logout() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := sessions.Delete(w, r); err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	})
}

func Login(repo *repository.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// check if user exists
		user, err := services.Login(r.Context(), repo, username, password)
		if err != nil {
			pages.Login(err).Render(r.Context(), w)
			return
		}

		// session
		sessions.Add(w, user.ID)

		// redirect
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}

func Signup(repo *repository.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		passwordRepeat := r.FormValue("password-repeat")

		user, errs := services.Signup(r.Context(), repo, username, email, password, passwordRepeat)
		if errs != nil {
			pages.Signup(errs).Render(r.Context(), w)
			return
		}

		// session
		sessions.Add(w, user.ID)

		// redirect
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}
