package utils

import (
	"homedb/internal/views/pages"
	"net/http"
)

func WriteError(w http.ResponseWriter, r *http.Request, status int, err error) {
	_, isLoggedIn := r.Context().Value("user_id").(string)

	w.WriteHeader(status)
	pages.Error(status, err, isLoggedIn).Render(r.Context(), w)
}
