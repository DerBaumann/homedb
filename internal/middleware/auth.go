package middleware

import (
	"context"
	"errors"
	"net/http"
	"os"

	"homedb/internal/utils"

	"github.com/gorilla/sessions"
)

func Protected(store *sessions.CookieStore) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := store.Get(r, os.Getenv("SESSION_NAME"))

			userID, ok := session.Values["user_id"].(string)
			if !ok {
				utils.WriteError(w, r, http.StatusUnauthorized, errors.New("unauthorized"))
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
