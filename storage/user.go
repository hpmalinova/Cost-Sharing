package storage

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type Username string

type User struct {
	Username `json:"username"`
	Password string `json:"password"`
}

type Users struct {
	Users map[Username]User `json:"users"`
}

func (u *Users) GetPassword(username Username) string {
	return u.Users[username].Password
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (u *Users) Create(username Username, password string) error {
	// Check if username already exists:
	if _, ok := u.Users[username]; ok {
		return errors.New("this username is already taken")
	}

	password, err := HashPassword(password)
	if err != nil {
		log.Println(err)
		return errors.New("an error has occurred while creating user")
	}
	newUser := User{
		Username: username,
		Password: password,
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
