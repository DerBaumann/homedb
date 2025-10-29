package utils

import (
	"homedb/views/pages"
	"net/http"
)

func WriteError(w http.ResponseWriter, r *http.Request, status int, err error) {
	w.WriteHeader(status)
	pages.Error(status, err).Render(r.Context(), w)
}
