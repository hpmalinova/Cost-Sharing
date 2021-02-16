package storage

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestGroups_CreateGroup(t *testing.T) {
	t.Run("when adding new group", func(t *testing.T) {
		g := getGroups()
		participants := []string{username, username2, username3}
		sort.Strings(participants)

		groupID := g.CreateGroup(groupName1, participants)
		assert.Equal(t, len(g.Groups), 1)
		_, ok := g.Groups[groupID]
		assert.True(t, ok)
	})
}

func TestGroups_GetGroupNames(t *testing.T) {
	t.Run("when having groups", func(t *testing.T) {
		g := Groups{map[uuid.UUID]Group{}}
		g.Groups[groupID1] = Group{Name: groupName1, MoneyExchange: MoneyExchange{}}
		g.Groups[groupID2] = Group{Name: groupName2, MoneyExchange: MoneyExchange{}}

		groupIDs := []uuid.UUID{groupID1, groupID2}
		expected := []string{groupName1, groupName2}
		sort.Strings(expected)

		actual := g.GetGroupNames(groupIDs)
		sort.Strings(actual)

		assert.Equal(t, expected, actual)
	})
	t.Run("when having no groups", func(t *testing.T) {
		g := getGroups()

		groupIDs := []uuid.UUID{groupID1, groupID2}
		expected := []string{}
		actual := g.GetGroupNames(groupIDs)

		assert.Equal(t, expected, actual)
	})
	t.Run("when having some groups", func(t *testing.T) {
		g := Groups{map[uuid.UUID]Group{}}
		g.Groups[groupID1] = Group{Name: groupName1, MoneyExchange: MoneyExchange{}}
		g.Groups[groupID2] = Group{Name: groupName2, MoneyExchange: MoneyExchange{}}

		groupIDs := []uuid.UUID{groupID1, groupID2, uuid.New()}
		expected := []string{groupName1, groupName2}
		sort.Strings(expected)

		actual := g.GetGroupNames(groupIDs)
		sort.Strings(actual)

		assert.Equal(t, expected, actual)
	})
}

func TestGroups_AddDebt(t *testing.T) {
	t.Run("when no debts prior to addDebt", func(t *testing.T) {
		g := getGroups()
		participants := []string{username, username2, username3}
		id1 := g.CreateGroup(groupName1, participants)

		g.AddDebt(username, id1, participants, amount, reason)

		expected := map[string][]DebtC{groupName1: {
			{username2, amount, reason},
			{username3, amount, reason}}}
		actual := g.GetLent(username, []uuid.UUID{id1})
		fmt.Println("ACTUAL", actual)

		assert.True(t, equal(expected, actual))
	})
}

func TestGroups_ReturnDebt(t *testing.T) {
	t.Run("return debt", func(t *testing.T) {
		g := getGroups()
		participants := []string{username, username2, username3}
		id1 := g.CreateGroup(groupName1, participants)

		participants2 := []string{username, username2}
		id2 := g.CreateGroup(groupName1, participants2)

		// Peter lends money to George and Lily
		g.AddDebt(username, id1, participants, amount, reason)
		// George returns the money
		g.ReturnDebt(username2, id1, username, amount)

		expected := map[string][]DebtC{}
		actual := g.GetOwed(username, []uuid.UUID{id1, id2})

		assert.Equal(t, expected, actual)
	})
	t.Run("when returning less", func(t *testing.T) {
		g := getGroups()
		participants := []string{username, username2, username3}
		id1 := g.CreateGroup(groupName1, participants)

		// Peter lends money to George and Lily
		g.AddDebt(username, id1, participants, amount, reason)
		// George returns half of the money
		g.ReturnDebt(username2, id1, username, amount/2)

		actualOwed := g.GetOwed(username2, []uuid.UUID{id1})
		expectedOwed := map[string][]DebtC{groupName1: {{username, amount / 2, reason}}}

		actualLent := g.GetLent(username, []uuid.UUID{id1})
		expectedLent := map[string][]DebtC{groupName1: {{username2, amount / 2, reason},
			{username3, amount, reason}}}

		assert.True(t, equal(expectedOwed, actualOwed))
		assert.True(t, equal(expectedLent, actualLent))
	})
	t.Run("when returning more", func(t *testing.T) {
		g := getGroups()
		participants := []string{username, username2, username3}
		id1 := g.CreateGroup(groupName1, participants)

		// Peter lends money to George and Lily
		g.AddDebt(username, id1, participants, amount, reason)
		// George returns double the money
		g.ReturnDebt(username2, id1, username, amount*2)

		// Now Peter owes money to George
		newReason := "Repay"
		actualOwed := g.GetOwed(username, []uuid.UUID{id1})
		expectedOwed := map[string][]DebtC{groupName1: {{username2, amount, newReason}}}
		fmt.Println("ActualOwed", actualOwed, expectedOwed)

		actualLent := g.GetLent(username, []uuid.UUID{id1})
		expectedLent := map[string][]DebtC{groupName1: {{username3, amount, reason}}}
		fmt.Println("ActualLent", actualLent, expectedLent)

		assert.True(t, equal(expectedOwed, actualOwed))
		assert.True(t, equal(expectedLent, actualLent))
	})
}

func TestGroups_FindGroupID(t *testing.T) {
	t.Run("when user participates in group", func(t *testing.T) {
		g := Groups{map[uuid.UUID]Group{}}
		g.Groups[groupID1] = Group{Name: groupName1, MoneyExchange: MoneyExchange{}}
		g.Groups[groupID2] = Group{Name: groupName2, MoneyExchange: MoneyExchange{}}

		participatesIn := []uuid.UUID{groupID1, groupID2}
		groupID, err := g.FindGroupID(groupName1, participatesIn)

		assert.Equal(t, groupID, groupID1)
		assert.NoError(t, err)
	})
	t.Run("when user does not participate in group", func(t *testing.T) {
		g := Groups{map[uuid.UUID]Group{}}
		g.Groups[groupID1] = Group{Name: groupName1, MoneyExchange: MoneyExchange{}}
		g.Groups[groupID2] = Group{Name: groupName2, MoneyExchange: MoneyExchange{}}

		participatesIn := []uuid.UUID{groupID2}
		groupID, err := g.FindGroupID(groupName1, participatesIn)

		assert.Equal(t, groupID, uuid.Nil)
		assert.Error(t, err)
	})
}

func TestGroups_GetOwed(t *testing.T) {
	t.Run("when you owe money to multiple groups", func(t *testing.T) {
		groups := getGroups()

		id1 := groups.CreateGroup(groupName1, []string{username, username2, username3})
		exchange1 := groups.Groups[id1].MoneyExchange
		// Peter owes George and Lily 20lv
		exchange1.AddDebt(username, username2, amount, reason)
		exchange1.AddDebt(username, username3, amount, reason)

		id2 := groups.CreateGroup(groupName2, []string{username, username3})
		exchange2 := groups.Groups[id2].MoneyExchange
		// Peter owes Lily 100lv ()
		exchange2.AddDebt(username, username3, amount2, reason2)

		actual := groups.GetOwed(username, []uuid.UUID{id1, id2})
		expected := map[string][]DebtC{groupName1: {{username2, amount, reason},
			{username3, amount, reason}},
			groupName2: {{username3, amount2, reason2}}}

		assert.True(t, equal(expected, actual))
	})
	t.Run("when you do not owe money to any group", func(t *testing.T) {
		groups := getGroups()

		id1 := groups.CreateGroup(groupName1, []string{username, username2, username3})
		id2 := groups.CreateGroup(groupName2, []string{username, username3})

		actual := groups.GetOwed(username, []uuid.UUID{id1, id2})
		expected := map[string][]DebtC{}
		assert.True(t, equal(expected, actual))
	})
}

func TestGroups_GetLent(t *testing.T) {
	t.Run("when you have lent money to multiple groups", func(t *testing.T) {
		groups := getGroups()

		id1 := groups.CreateGroup(groupName1, []string{username, username2, username3})
		exchange1 := groups.Groups[id1].MoneyExchange
		// Peter lends George and Lily 20lv
		exchange1.AddDebt(username2, username, amount, reason)
		exchange1.AddDebt(username3, username, amount, reason)

		id2 := groups.CreateGroup(groupName2, []string{username, username3})
		exchange2 := groups.Groups[id2].MoneyExchange
		// Peter lends Lily 100lv ()
		exchange2.AddDebt(username3, username, amount2, reason2)

		actual := groups.GetLent(username, []uuid.UUID{id1, id2})
		expected := map[string][]DebtC{groupName1: {{username2, amount, reason},
			{username3, amount, reason}},
			groupName2: {{username3, amount2, reason2}}}

		assert.True(t, equal(expected, actual))
	})
	t.Run("when you have not lent money to any group", func(t *testing.T) {
		groups := getGroups()

		id1 := groups.CreateGroup(groupName1, []string{username, username2, username3})
		id2 := groups.CreateGroup(groupName2, []string{username, username3})

		actual := groups.GetOwed(username, []uuid.UUID{id1, id2})
		expected := map[string][]DebtC{}
		assert.True(t, equal(expected, actual))
	})
}
