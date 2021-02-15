package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFriends_AddFriend(t *testing.T) {
	t.Run("when both users exist", func(t *testing.T) {
		f := Friends{Friends: map[string]map[string]struct{}{"Pesho": {}, "Gosho": {}}}
		err := f.Add("Pesho", "Gosho")
		assert.NoError(t, err)
		assert.Equal(t, map[string]struct{}{"Gosho": {}}, f.Friends["Pesho"])
		assert.Equal(t, map[string]struct{}{"Pesho": {}}, f.Friends["Gosho"])
	})
	t.Run("when only one user exists", func(t *testing.T) {
		f := Friends{Friends: map[string]map[string]struct{}{"Pesho": {}}}
		err := f.Add("Pesho", "Gosho")
		assert.NoError(t, err)
		assert.Equal(t, map[string]struct{}{"Gosho": {}}, f.Friends["Pesho"])
		assert.Equal(t, map[string]struct{}{"Pesho": {}}, f.Friends["Gosho"])
	})
	t.Run("when already friends", func(t *testing.T) {
		f := Friends{Friends: map[string]map[string]struct{}{"Pesho": {"Gosho": {}}, "Gosho": {"Pesho": {}}}}
		err := f.Add("Pesho", "Gosho")
		assert.Error(t, err)
	})
}

func TestFriends_GetFriendsOf(t *testing.T) {
	t.Run("when has friends", func(t *testing.T) {
		f := Friends{Friends: map[string]map[string]struct{}{"Pesho": {"Gosho": {}, "Maria": {}}}}
		expected := []string{"Gosho", "Maria"}
		actual := f.GetFriendsOf("Pesho")
		assert.Equal(t, expected, actual)
	})
	t.Run("when has no friends", func(t *testing.T) {
		f := Friends{Friends: map[string]map[string]struct{}{"Pesho": {}}}
		expected := []string{}
		actual := f.GetFriendsOf("Pesho")
		assert.Equal(t, expected, actual)
	})
}
