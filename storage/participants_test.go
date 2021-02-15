package storage

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestParticipants_Add(t *testing.T) {
	t.Run("when adding a new group", func(t *testing.T) {
		uuid1 := uuid.New()
		participants := []string{"Pesho", "Gosho"}

		p := Participants{Participants: map[uuid.UUID]map[string]struct{}{}}
		p.Add(uuid1, participants)
		expected := map[uuid.UUID]map[string]struct{}{uuid1: {"Pesho": {}, "Gosho": {}}}
		actual := p.Participants
		assert.Equal(t, expected, actual)
	})
}

func TestParticipants_HasParticipant(t *testing.T) {
	t.Run("when having the participant", func(t *testing.T) {
		uuid1 := uuid.New()

		p := Participants{Participants: map[uuid.UUID]map[string]struct{}{uuid1: {"Pesho": {}, "Maria": {}}}}
		ok := p.HasParticipant(uuid1, "Pesho")

		assert.True(t, ok)
	})
	t.Run("when not having the participant", func(t *testing.T) {
		uuid1 := uuid.New()

		p := Participants{Participants: map[uuid.UUID]map[string]struct{}{uuid1: {"Pesho": {}, "Maria": {}}}}
		ok := p.HasParticipant(uuid1, "lily")

		assert.False(t, ok)
	})
}

func TestParticipants_GetParticipants(t *testing.T) {
	t.Run("when having participants", func(t *testing.T) {
		uuid1 := uuid.New()

		p := Participants{Participants: map[uuid.UUID]map[string]struct{}{uuid1: {"Pesho": {}, "Maria": {}}}}
		actual := p.GetParticipants(uuid1)
		sort.Strings(actual)
		expected := []string{"Maria", "Pesho"}
		assert.Equal(t, expected, actual)
	})
	t.Run("when not having participants", func(t *testing.T) {
		uuid1 := uuid.New()

		p := Participants{Participants: map[uuid.UUID]map[string]struct{}{uuid1: {}}}
		actual := p.GetParticipants(uuid1)
		sort.Strings(actual)
		expected := []string{}
		assert.Equal(t, expected, actual)
	})
	t.Run("when not having this groupID", func(t *testing.T) {
		uuid1 := uuid.New()

		p := Participants{Participants: map[uuid.UUID]map[string]struct{}{}}
		actual := p.GetParticipants(uuid1)
		sort.Strings(actual)
		expected := []string{}
		assert.Equal(t, expected, actual)
	})
}

func TestParticipants_GetNumberOfParticipants(t *testing.T) {
	t.Run("when having multiple participants", func(t *testing.T) {
		uuid1 := uuid.New()

		p := Participants{Participants: map[uuid.UUID]map[string]struct{}{uuid1: {"Pesho": {}, "Maria": {}, "Silviya": {}}}}
		actual := p.GetNumberOfParticipants(uuid1)
		expected := 3
		assert.Equal(t, expected, actual)
	})
	t.Run("when having no participants", func(t *testing.T) {
		uuid1 := uuid.New()

		p := Participants{Participants: map[uuid.UUID]map[string]struct{}{uuid1: {}}}
		actual := p.GetNumberOfParticipants(uuid1)
		expected := 0
		assert.Equal(t, expected, actual)
	})
}
