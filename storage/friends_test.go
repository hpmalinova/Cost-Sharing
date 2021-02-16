package storage

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestFriends_AddFriend(t *testing.T) {
	t.Run("when both users exist", func(t *testing.T) {
		f := Friends{Friends: map[string]map[string]struct{}{peter: {}, george: {}}}
		err := f.Add(peter, george)
		assert.NoError(t, err)
		assert.Equal(t, map[string]struct{}{george: {}}, f.Friends[peter])
		assert.Equal(t, map[string]struct{}{peter: {}}, f.Friends[george])
	})
	t.Run("when only one user exists", func(t *testing.T) {
		f := Friends{Friends: map[string]map[string]struct{}{peter: {}}}
		err := f.Add(peter, george)
		assert.NoError(t, err)
		assert.Equal(t, map[string]struct{}{george: {}}, f.Friends[peter])
		assert.Equal(t, map[string]struct{}{peter: {}}, f.Friends[george])
	})
	t.Run("when already friends", func(t *testing.T) {
		f := Friends{Friends: map[string]map[string]struct{}{peter: {george: {}}, george: {peter: {}}}}
		err := f.Add(peter, george)
		assert.Error(t, err)
	})
}

func TestFriends_GetFriendsOf(t *testing.T) {
	t.Run("when has friends", func(t *testing.T) {
		f := Friends{Friends: map[string]map[string]struct{}{peter: {george: {}, lily: {}}}}
		expected := []string{george, lily}
		actual := f.GetFriendsOf(peter)
		sort.Strings(actual)
		assert.Equal(t, expected, actual)
	})
	t.Run("when has no friends", func(t *testing.T) {
		f := Friends{Friends: map[string]map[string]struct{}{peter: {}}}
		expected := []string{}
		actual := f.GetFriendsOf(peter)
		assert.Equal(t, expected, actual)
	})
}

func TestFriends_AreFriends(t *testing.T) {
	t.Run("when friends", func(t *testing.T) {
		f := Friends{Friends: map[string]map[string]struct{}{peter: {}, george: {}}}
		_ = f.Add(peter, george)

		assert.True(t, f.AreFriends(peter, george))
	})
	t.Run("when the first has no friends", func(t *testing.T) {
		f := Friends{Friends: map[string]map[string]struct{}{peter: {}, george: {}}}
		_ = f.Add(peter, george)

		assert.False(t, f.AreFriends(lily, peter))
	})
	t.Run("when the second has no friends", func(t *testing.T) {
		f := Friends{Friends: map[string]map[string]struct{}{peter: {}, george: {}}}

		assert.False(t, f.AreFriends(peter, lily))
	})
}