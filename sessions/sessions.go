package sessions

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
)

var ErrSessionNotFound = errors.New("session does not exist")

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

func Get(r *http.Request) (*Session, error) {
	sessionCookie, err := r.Cookie("session-id")
	if err != nil {
		return nil, err
	}

	sessId, err := uuid.Parse(sessionCookie.Value)
	if err != nil {
		return nil, err
	}

	session, ok := Sessions[sessId]
	if !ok {
		return nil, ErrSessionNotFound
	}

	return &session, nil
}

func Delete(w http.ResponseWriter, r *http.Request) error {
	sessionCookie, err := r.Cookie("session-id")
	if err != nil {
		return err
	}

	sessId, err := uuid.Parse(sessionCookie.Value)
	if err != nil {
		return err
	}

	delete(Sessions, sessId)

	http.SetCookie(w, &http.Cookie{
		Name:     "session-id",
		Value:    "",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now(),
		MaxAge:   -1,
		Secure:   false, // TODO: Change to true for prod
	})

	return nil
}
