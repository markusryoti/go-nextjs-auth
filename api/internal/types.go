package internal

import "time"

type User struct {
	Id             string `json:"id"`
	Email          string `json:"email"`
	HashedPassword string `json:"-"`
	Salt           string `json:"-"`
}

type SessionData struct {
	User      *User     `json:"user"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func (sd SessionData) Expired() bool {
	return time.Now().After(sd.ExpiresAt)
}
