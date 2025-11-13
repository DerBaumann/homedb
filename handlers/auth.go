package handlers

import (
	"homedb/repository"
	"homedb/services"
	"homedb/views/pages"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

func Logout(store *sessions.CookieStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, os.Getenv("SESSION_NAME"))
		delete(session.Values, "user_id")
		if err := session.Save(r, w); err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	})
}

func Login(repo *repository.Queries, store *sessions.CookieStore) http.Handler {
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
		session, _ := store.Get(r, os.Getenv("SESSION_NAME"))
		session.Values["user_id"] = user.ID.String()
		if err := session.Save(r, w); err != nil {
			pages.Login(err).Render(r.Context(), w)
			return
		}

		// redirect
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}

func Signup(repo *repository.Queries, store *sessions.CookieStore) http.Handler {
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
		session, _ := store.Get(r, os.Getenv("SESSION_NAME"))
		session.Values["user_id"] = user.ID.String()
		if err := session.Save(r, w); err != nil {
			pages.Signup([]error{err}).Render(r.Context(), w)
			return
		}

		// redirect
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}
