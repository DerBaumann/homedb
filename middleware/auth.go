package middleware

import (
	"errors"
	"homedb/utils"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

func Protected(store *sessions.CookieStore) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := store.Get(r, os.Getenv("SESSION_NAME"))

			_, ok := session.Values["user_id"]
			if !ok {
				utils.WriteError(w, r, http.StatusUnauthorized, errors.New("unauthorized"))
				return
			}

			// ctx := context.WithValue(r.Context(), sessions.ContextKey, session)

			// next.ServeHTTP(w, r.WithContext(ctx))
			next.ServeHTTP(w, r)
		})
	}
}
