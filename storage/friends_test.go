package storage

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestFriends_AddFriend(t *testing.T) {
	t.Run("when both users exist", func(t *testing.T) {
		f := Friends{Friends: map[string]map[string]struct{}{username: {}, username2: {}}}
		err := f.Add(username, username2)
		assert.NoError(t, err)
		assert.Equal(t, map[string]struct{}{username2: {}}, f.Friends[username])
		assert.Equal(t, map[string]struct{}{username: {}}, f.Friends[username2])
	})
	t.Run("when only one user exists", func(t *testing.T) {
		f := Friends{Friends: map[string]map[string]struct{}{username: {}}}
		err := f.Add(username, username2)
		assert.NoError(t, err)
		assert.Equal(t, map[string]struct{}{username2: {}}, f.Friends[username])
		assert.Equal(t, map[string]struct{}{username: {}}, f.Friends[username2])
	})
	t.Run("when already friends", func(t *testing.T) {
		f := Friends{Friends: map[string]map[string]struct{}{username: {username2: {}}, username2: {username: {}}}}
		err := f.Add(username, username2)
		assert.Error(t, err)
	})
}

func TestFriends_GetFriendsOf(t *testing.T) {
	t.Run("when has friends", func(t *testing.T) {
		f := Friends{Friends: map[string]map[string]struct{}{username: {username2: {}, username3: {}}}}
		expected := []string{username2, username3}
		actual := f.GetFriendsOf(username)
		sort.Strings(actual)
		assert.Equal(t, expected, actual)
	})
	t.Run("when has no friends", func(t *testing.T) {
		f := Friends{Friends: map[string]map[string]struct{}{username: {}}}
		expected := []string{}
		actual := f.GetFriendsOf(username)
		assert.Equal(t, expected, actual)
	})
}
