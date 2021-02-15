package storage

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
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
		return errors.New("an error has occurred while creating user")
	}

	newUser := User{
		Username: username,
		Password: password,
	}

	u.Users[username] = newUser
	return nil
}

func (u *Users) CheckCredentials(username Username, password string) error {
	if user, ok := u.Users[username]; ok {
		ok := CheckPasswordHash(password, user.Password)
		if !ok {
			return errors.New("invalid username or password")
		} else {
			return nil
		}
	} else {
		return errors.New("invalid username or password")
	}
}

func (u *Users) GetUsernames() []string {
	usernames := make([]string, 0, len(u.Users))
	for _, user := range u.Users {
		usernames = append(usernames, string(user.Username))
	}
	return usernames
}
