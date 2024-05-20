package store

import "errors"

var ErrSessionNotFound = errors.New("session not found")

var SessionStore = newSessionStore()

type sessionStore struct {
	sessions map[string]*User
}

func newSessionStore() *sessionStore {
	return &sessionStore{
		sessions: make(map[string]*User),
	}
}

func (ss *sessionStore) AddSession(id string, user *User) {
	ss.sessions[id] = user
}

func (ss *sessionStore) GetSession(id string) (*User, error) {
	user, ok := ss.sessions[id]
	if !ok {
		return nil, ErrSessionNotFound
	}

	return user, nil
}
