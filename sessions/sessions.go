package sessions

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID     uuid.UUID
	Expiry time.Time
}

func (s *Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}

type (
	contextKey   string
	SessionStore = map[uuid.UUID]Session
)

const (
	ContextKey contextKey = "session"
	fileName   string     = "session.json"
)

var (
	ErrSessionNotFound = errors.New("session does not exist")
	ErrUnauthorized    = errors.New("unauthorized")
)

var Sessions = map[uuid.UUID]Session{}

func Init() error {
	contents, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	var data SessionStore

	if err := json.Unmarshal(contents, &data); err != nil {
		return err
	}

	Sessions = data

	return nil
}

func save() error {
	data, err := json.MarshalIndent(Sessions, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(fileName, data, 0644); err != nil {
		return err
	}

	return nil
}

func Add(w http.ResponseWriter, userId uuid.UUID) (uuid.UUID, error) {
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

	if err := save(); err != nil {
		return uuid.Nil, err
	}

	return sessId, nil
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

	if err := save(); err != nil {
		return err
	}

	return nil
}
