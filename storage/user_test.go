package storage

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

const (
	password  = "123456"
	password2 = "qwerty"
)

func TestUsers_Create(t *testing.T) {
	t.Run("when creating user with free username", func(t *testing.T) {
		u := Users{Users: map[string]User{}}
		err := u.Create(username, password)

		expected := []string{username}
		assert.Equal(t, expected, u.GetUsernames())
		assert.NoError(t, err)
	})
	t.Run("when creating user with taken username", func(t *testing.T) {
		u := Users{Users: map[string]User{username: {username, password}}}
		err := u.Create(username, "123456")

		expected := []string{username}
		assert.Equal(t, expected, u.GetUsernames())
		assert.Error(t, err)
	})
}

func TestUsers_GetUsernames(t *testing.T) {
	t.Run("when there are multiple users", func(t *testing.T) {
		u1 := User{Username: username, Password: password}
		u2 := User{Username: username2, Password: password2}

		expected := []string{username, username2}
		sort.Strings(expected)

		u := Users{Users: map[string]User{username: u1, username2: u2}}
		actual := u.GetUsernames()
		sort.Strings(actual)

		assert.Equal(t, expected, actual)
	})
	t.Run("when there are no users", func(t *testing.T) {
		expected := []string{}

		u := Users{Users: map[string]User{}}
		actual := u.GetUsernames()

		assert.Equal(t, expected, actual)
	})
}

func TestUsers_DoesExist(t *testing.T) {
	t.Run("when user exists", func(t *testing.T) {
		u1 := User{Username: username, Password: password}
		u := Users{Users: map[string]User{username: u1}}

		assert.True(t, u.DoesExist(username))
	})
	t.Run("when user does not exists", func(t *testing.T) {
		u := Users{Users: map[string]User{}}

		assert.False(t, u.DoesExist(username))
	})
}

func TestUsers_CheckCredentials(t *testing.T) {
	u := Users{Users: map[string]User{}}
	_ = u.Create(username, password)

	t.Run("when legitimate user", func(t *testing.T) {
		err := u.CheckCredentials(username, password)
		assert.NoError(t, err)
	})
	t.Run("when non-existing username", func(t *testing.T) {
		err := u.CheckCredentials(username2, password2)
		assert.Error(t, err)
	})
	t.Run("when wrong password", func(t *testing.T) {
		err := u.CheckCredentials(username, password2)
		assert.Error(t, err)
	})
}
