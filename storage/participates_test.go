package storage

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestParticipates_Add(t *testing.T) {
	t.Run("when having no user", func(t *testing.T) {
		p := Participates{map[string]map[uuid.UUID]struct{}{}}
		expected := map[string]map[uuid.UUID]struct{}{peter: {groupID1: {}}}
		p.Add(peter, groupID1)
		actual := p.Participates

		assert.Equal(t, expected, actual)
	})
	t.Run("when user already has another group", func(t *testing.T) {
		p := Participates{map[string]map[uuid.UUID]struct{}{peter: {groupID1: {}}}}
		expected := map[string]map[uuid.UUID]struct{}{peter: {groupID1: {}, groupID2: {}}}
		p.Add(peter, groupID2)
		actual := p.Participates

		assert.Equal(t, expected, actual)
	})
}

func TestParticipates_DoesParticipate(t *testing.T) {
	t.Run("when user participates in the group", func(t *testing.T) {
		p := Participates{Participates: map[string]map[uuid.UUID]struct{}{peter: {groupID1: {}, groupID2: {}}}}
		ok := p.DoesParticipate(peter, groupID1)

		assert.True(t, ok)
	})
	t.Run("when user does not participate in the group", func(t *testing.T) {
		p := Participates{Participates: map[string]map[uuid.UUID]struct{}{peter: {groupID1: {}}}}
		ok := p.DoesParticipate(peter, groupID2)

		assert.False(t, ok)
	})
}

func TestParticipates_GetGroups(t *testing.T) {
	t.Run("when user participates in multiple groups", func(t *testing.T) {
		p := Participates{Participates: map[string]map[uuid.UUID]struct{}{peter: {groupID1: {}, groupID2: {}}}}
		actual := p.GetGroups(peter)
		sort.Slice(actual, func(i, j int) bool {
			return actual[i].String() < actual[j].String()
		})
		expected := []uuid.UUID{groupID1, groupID2}
		sort.Slice(expected, func(i, j int) bool {
			return expected[i].String() < expected[j].String()
		})
		assert.Equal(t, expected, actual)

	})
	t.Run("when user does not participate in any group", func(t *testing.T) {
		p := Participates{Participates: map[string]map[uuid.UUID]struct{}{peter: {}}}
		actual := p.GetGroups(peter)
		expected := []uuid.UUID{}

		assert.Equal(t, expected, actual)
	})
	t.Run("when not having this user", func(t *testing.T) {
		p := Participates{Participates: map[string]map[uuid.UUID]struct{}{}}
		actual := p.GetGroups(peter)
		expected := []uuid.UUID{}

		assert.Equal(t, expected, actual)
	})
}
