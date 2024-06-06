package store

import (
	"errors"

	"github.com/google/uuid"
	"github.com/markusryoti/next-js-go-auth/internal"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

var UserStore = newUserStore()

func newUserStore() *userStore {
	return &userStore{
		users: make(map[string]*internal.User),
	}
}

type userStore struct {
	users map[string]*internal.User
}

func (us *userStore) SaveUser(email, hashedPassword, salt string) (*internal.User, error) {
	id := uuid.NewString()

	u := &internal.User{
		Id:             id,
		Email:          email,
		HashedPassword: hashedPassword,
		Salt:           salt,
	}

	us.users[u.Id] = u

	return u, nil
}

func (us *userStore) GetUserById(id string) (*internal.User, error) {
	for _, u := range us.users {
		if u.Id == id {
			return u, nil
		}
	}

	return nil, ErrUserNotFound
}

func (us *userStore) GetUserByEmail(email string) (*internal.User, error) {
	for _, u := range us.users {
		if u.Email == email {
			return u, nil
		}
	}

	return nil, ErrUserNotFound
}
