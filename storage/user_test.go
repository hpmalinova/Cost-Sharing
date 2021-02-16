package storage

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestUsers_Create(t *testing.T) {
	t.Run("when creating user with free username", func(t *testing.T) {
		u := Users{Users: map[string]User{}}
		err := u.Create(peter, password)

		expected := []string{peter}
		assert.Equal(t, expected, u.GetUsernames())
		assert.NoError(t, err)
	})
	t.Run("when creating user with taken username", func(t *testing.T) {
		u := Users{Users: map[string]User{peter: {peter, password}}}
		err := u.Create(peter, "123456")

		expected := []string{peter}
		assert.Equal(t, expected, u.GetUsernames())
		assert.Error(t, err)
	})
}

func TestUsers_GetUsernames(t *testing.T) {
	t.Run("when there are multiple users", func(t *testing.T) {
		u1 := User{Username: peter, Password: password}
		u2 := User{Username: george, Password: password2}

		expected := []string{peter, george}
		sort.Strings(expected)

		u := Users{Users: map[string]User{peter: u1, george: u2}}
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
		u1 := User{Username: peter, Password: password}
		u := Users{Users: map[string]User{peter: u1}}

		assert.True(t, u.DoesExist(peter))
	})
	t.Run("when user does not exists", func(t *testing.T) {
		u := Users{Users: map[string]User{}}

		assert.False(t, u.DoesExist(peter))
	})
}

func TestUsers_CheckCredentials(t *testing.T) {
	u := Users{Users: map[string]User{}}
	_ = u.Create(peter, password)

	t.Run("when legitimate user", func(t *testing.T) {
		err := u.CheckCredentials(peter, password)
		assert.NoError(t, err)
	})
	t.Run("when non-existing username", func(t *testing.T) {
		err := u.CheckCredentials(george, password2)
		assert.Error(t, err)
	})
	t.Run("when wrong password", func(t *testing.T) {
		err := u.CheckCredentials(peter, password2)
		assert.Error(t, err)
	})
}
