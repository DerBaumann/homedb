package sessions

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID     uuid.UUID
	Expiry time.Time
}

var Sessions map[uuid.UUID]Session

func (s *Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}

func New(id uuid.UUID) Session {
	return Session{ID: id, Expiry: time.Now().Add(24 * time.Hour)}
}
