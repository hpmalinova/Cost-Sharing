package storage

import (
	"errors"
)

type Username string

type User struct {
	Username `json:"username"`
	Password string `json:"password"`
}

type Users struct {
	Users map[Username]User `json:"users"`
}

func (u *Users) Create(username Username, password string) error {
	// Check if username already exists:
	if _, ok := u.Users[username]; ok {
		return errors.New("this username is already taken")
	}

	newUser := User{
		Username: username,
		Password: password, // TODO hash
	}
	u.Users[username] = newUser
	return nil
}

func (u *Users) GetUsernames() []Username {
	usernames := make([]Username, 0, len(u.Users))
	for _, user := range u.Users {
		usernames = append(usernames, user.Username)
	}
	return usernames
}
