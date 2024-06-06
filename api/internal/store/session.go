package store

import (
	"errors"
	"time"

	"github.com/markusryoti/next-js-go-auth/internal"
)

var ErrSessionNotFound = errors.New("session not found")

var SessionStore = newSessionStore()

type sessionStore struct {
	sessions map[string]*internal.SessionData
}

func newSessionStore() *sessionStore {
	return &sessionStore{
		sessions: make(map[string]*internal.SessionData),
	}
}

func (ss *sessionStore) AddSession(id string, user *internal.User, expiresAt time.Time) {
	ss.sessions[id] = &internal.SessionData{
		User:      user,
		ExpiresAt: expiresAt,
	}
}

func (ss *sessionStore) GetSession(id string) (*internal.SessionData, error) {
	session, ok := ss.sessions[id]
	if !ok {
		return nil, ErrSessionNotFound
	}

	return session, nil
}

func (ss *sessionStore) UpdateExpiration(id string, expiresAt time.Time) error {
	session, err := ss.GetSession(id)
	if err != nil {
		return err
	}

	session.ExpiresAt = expiresAt

	return nil
}

func (ss *sessionStore) RemoveSession(id string) {
	delete(ss.sessions, id)
}
