package sessions

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID     uuid.UUID
	Expiry time.Time
}

var Sessions = map[uuid.UUID]Session{}

func (s *Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}

func Add(w http.ResponseWriter, userId uuid.UUID) uuid.UUID {
	sessId := uuid.New()
	session := Session{ID: userId, Expiry: time.Now().Add(24 * time.Hour)}

	Sessions[sessId] = session

	http.SetCookie(w, &http.Cookie{
		Name:     "session-id",
		Value:    sessId.String(),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   false, // TODO: Change to true for prod
	})

	return sessId
}
