package middleware

import (
	"errors"
	"homedb/sessions"
	"homedb/utils"
	"net/http"
)

func Protected(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := sessions.Get(r)
		if err != nil {
			if err == http.ErrNoCookie {
				utils.WriteError(w, r, http.StatusUnauthorized, errors.New("you do not posess sufficient permissions to view this page"))
				return
			}
			utils.WriteError(w, r, http.StatusInternalServerError, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}
