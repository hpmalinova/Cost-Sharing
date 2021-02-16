package storage

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestParticipants_Add(t *testing.T) {
	t.Run("when adding a new group", func(t *testing.T) {
		participants := []string{username, username2}

		p := Participants{Participants: map[uuid.UUID]map[string]struct{}{}}
		id := groupID1
		p.Add(id, participants)
		expected := map[uuid.UUID]map[string]struct{}{groupID1: {username: {}, username2: {}}}
		actual := p.Participants
		assert.Equal(t, expected, actual)
	})
}

func TestParticipants_HasParticipant(t *testing.T) {
	t.Run("when having the participant", func(t *testing.T) {
		p := Participants{Participants: map[uuid.UUID]map[string]struct{}{groupID1: {username: {}, username2: {}}}}
		ok := p.HasParticipant(groupID1, username)

		assert.True(t, ok)
	})
	t.Run("when not having the participant", func(t *testing.T) {
		p := Participants{Participants: map[uuid.UUID]map[string]struct{}{groupID1: {username: {}, username2: {}}}}
		ok := p.HasParticipant(groupID1, username3)

		assert.False(t, ok)
	})
}

func TestParticipants_GetParticipants(t *testing.T) {
	t.Run("when having participants", func(t *testing.T) {
		p := Participants{Participants: map[uuid.UUID]map[string]struct{}{groupID1: {username: {}, username2: {}}}}
		actual := p.GetParticipants(groupID1)
		sort.Strings(actual)
		expected := []string{username2, username}
		assert.Equal(t, expected, actual)
	})
	t.Run("when not having participants", func(t *testing.T) {
		p := Participants{Participants: map[uuid.UUID]map[string]struct{}{groupID1: {}}}
		actual := p.GetParticipants(groupID1)
		sort.Strings(actual)
		expected := []string{}
		assert.Equal(t, expected, actual)
	})
	t.Run("when not having this groupID", func(t *testing.T) {
		p := Participants{Participants: map[uuid.UUID]map[string]struct{}{}}
		actual := p.GetParticipants(groupID1)
		sort.Strings(actual)
		expected := []string{}
		assert.Equal(t, expected, actual)
	})
}

func TestParticipants_GetNumberOfParticipants(t *testing.T) {
	t.Run("when having multiple participants", func(t *testing.T) {
		p := Participants{Participants: map[uuid.UUID]map[string]struct{}{groupID1: {username: {}, username2: {}, username3: {}}}}
		actual := p.GetNumberOfParticipants(groupID1)
		expected := 3
		assert.Equal(t, expected, actual)
	})
	t.Run("when having no participants", func(t *testing.T) {
		p := Participants{Participants: map[uuid.UUID]map[string]struct{}{groupID1: {}}}
		actual := p.GetNumberOfParticipants(groupID1)
		expected := 0
		assert.Equal(t, expected, actual)
	})
}
