package storage

import (
	"errors"
	"fmt"
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

func (u *Users) String() { // ?
	for _, user := range u.Users {
		fmt.Println(user.Username, " ")
	}
	fmt.Println()
}

// TODO showAllUsers (only usernames)