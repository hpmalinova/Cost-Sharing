package storage

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Users Creates a connection between usernames and users
type Users struct {
	Users map[string]User `json:"users"`
}

// Create Creates a new user, if the given <username> is not taken
// It Hashes the password and stores it in the Users struct
func (u *Users) Create(username string, password string) error {
	// Check if username already exists:
	if u.DoesExist(username) {
		return errors.New("this username is already taken")
	}

	password, err := hashPassword(password)
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

func (u *Users) GetPassword(username string) string {
	return u.Users[username].Password
}

func (u *Users) CheckCredentials(username string, password string) error {
	if user, ok := u.Users[username]; ok {
		ok := checkPasswordHash(password, user.Password)
		if !ok {
			return errors.New("invalid username or password")
		} else {
			return nil
		}
	} else {
		return errors.New("invalid username or password")
	}
}

// GetUsernames Returns a collection of all registered usernames
func (u *Users) GetUsernames() []string {
	usernames := make([]string, 0, len(u.Users))
	for _, user := range u.Users {
		usernames = append(usernames, user.Username)
	}
	return usernames
}

// DoesExist Checks if a user with <username> exists
func (u *Users) DoesExist(username string) bool {
	if _, ok := u.Users[username]; ok {
		return true
	}
	return false
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
