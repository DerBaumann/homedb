package utils

import (
	"net/http"

	"homedb/internal/views/pages"
)

func WriteError(w http.ResponseWriter, r *http.Request, status int, err error) {
	w.WriteHeader(status)
	pages.Error(status, err).Render(r.Context(), w)
}
