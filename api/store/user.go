package store

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

var UserStore = newUserStore()

type User struct {
	Id             string `json:"id"`
	Email          string `json:"email"`
	HashedPassword string `json:"-"`
}

func newUserStore() *userStore {
	return &userStore{
		users: make(map[string]*User),
	}
}

type userStore struct {
	users map[string]*User
}

func (us *userStore) SaveUser(email, hashedPassword string) (*User, error) {
	id := uuid.NewString()

	u := &User{
		Id:             id,
		Email:          email,
		HashedPassword: hashedPassword,
	}

	us.users[u.Id] = u

	return u, nil
}

func (us *userStore) GetUserById(id string) (*User, error) {
	for _, u := range us.users {
		if u.Id == id {
			return u, nil
		}
	}

	return nil, ErrUserNotFound
}

func (us *userStore) GetUserByEmail(email string) (*User, error) {
	for _, u := range us.users {
		if u.Email == email {
			return u, nil
		}
	}

	return nil, ErrUserNotFound
}
