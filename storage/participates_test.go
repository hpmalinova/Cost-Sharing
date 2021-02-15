package storage

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestParticipates_Add(t *testing.T) {
	t.Run("when having no user", func(t *testing.T) {
		uuid1 := uuid.New()
		username := "Pesho"

		p := Participates{map[string]map[uuid.UUID]struct{}{}}
		expected := map[string]map[uuid.UUID]struct{}{username: {uuid1: {}}}
		p.Add(username, uuid1)
		actual := p.Participates

		assert.Equal(t, expected, actual)
	})
	t.Run("when user already has another group", func(t *testing.T) {
		uuid1 := uuid.New()
		uuid2 := uuid.New()
		username := "Pesho"

		p := Participates{map[string]map[uuid.UUID]struct{}{username: {uuid1: {}}}}
		expected := map[string]map[uuid.UUID]struct{}{username: {uuid1: {}, uuid2: {}}}
		p.Add(username, uuid2)
		actual := p.Participates

		assert.Equal(t, expected, actual)
	})
}

func TestParticipates_DoesParticipate(t *testing.T) {
	t.Run("when user participates in the group", func(t *testing.T) {
		uuid1 := uuid.New()
		uuid2 := uuid.New()
		username := "Pesho"

		p := Participates{Participates: map[string]map[uuid.UUID]struct{}{username: {uuid1: {}, uuid2: {}}}}
		ok := p.DoesParticipate(username, uuid1)

		assert.True(t, ok)
	})
	t.Run("when user does not participate in the group", func(t *testing.T) {
		uuid1 := uuid.New()
		uuid2 := uuid.New()
		username := "Pesho"

		p := Participates{Participates: map[string]map[uuid.UUID]struct{}{username: {uuid1: {}}}}
		ok := p.DoesParticipate(username, uuid2)

		assert.False(t, ok)
	})
}

func TestParticipates_GetGroups(t *testing.T) {
	t.Run("when user participates in multiple groups", func(t *testing.T) {
		uuid1 := uuid.New()
		uuid2 := uuid.New()
		username := "Pesho"

		p := Participates{Participates: map[string]map[uuid.UUID]struct{}{username: {uuid1: {}, uuid2: {}}}}
		actual := p.GetGroups(username)
		sort.Slice(actual, func(i, j int) bool {
			return actual[i].String() < actual[j].String()
		})
		expected := []uuid.UUID{uuid1, uuid2}
		sort.Slice(expected, func(i, j int) bool {
			return expected[i].String() < expected[j].String()
		})
		assert.Equal(t, expected, actual)

	})
	t.Run("when user does not participate in any group", func(t *testing.T) {
		username := "Pesho"
		p := Participates{Participates: map[string]map[uuid.UUID]struct{}{username: {}}}
		actual := p.GetGroups(username)
		expected := []uuid.UUID{}

		assert.Equal(t, expected, actual)
	})
	t.Run("when not having this user", func(t *testing.T) {
		username := "Pesho"
		p := Participates{Participates: map[string]map[uuid.UUID]struct{}{}}
		actual := p.GetGroups(username)
		expected := []uuid.UUID{}

		assert.Equal(t, expected, actual)
	})
}
