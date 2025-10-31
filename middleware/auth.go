package middleware

import (
	"context"
	"homedb/sessions"
	"homedb/utils"
	"net/http"
)

func Protected(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := sessions.Get(r)
		if err != nil {
			if err == http.ErrNoCookie || err == sessions.ErrSessionNotFound {
				utils.WriteError(w, r, http.StatusUnauthorized, err)
				return
			}
			utils.WriteError(w, r, http.StatusInternalServerError, err)
			return
		}

		ctx := context.WithValue(r.Context(), sessions.ContextKey, session)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
