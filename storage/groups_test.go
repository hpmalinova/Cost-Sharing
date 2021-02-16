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
		participants := []string{peter, george, lily}
		sort.Strings(participants)

		groupID := g.CreateGroup(christmasParty, participants)
		assert.Equal(t, len(g.Groups), 1)
		_, ok := g.Groups[groupID]
		assert.True(t, ok)
	})
}

func TestGroups_GetGroupNames(t *testing.T) {
	t.Run("when having groups", func(t *testing.T) {
		g := Groups{map[uuid.UUID]Group{}}
		g.Groups[groupID1] = Group{Name: christmasParty, MoneyExchange: MoneyExchange{}}
		g.Groups[groupID2] = Group{Name: travelToJapan, MoneyExchange: MoneyExchange{}}

		groupIDs := []uuid.UUID{groupID1, groupID2}
		expected := []string{christmasParty, travelToJapan}
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
		g.Groups[groupID1] = Group{Name: christmasParty, MoneyExchange: MoneyExchange{}}
		g.Groups[groupID2] = Group{Name: travelToJapan, MoneyExchange: MoneyExchange{}}

		groupIDs := []uuid.UUID{groupID1, groupID2, uuid.New()}
		expected := []string{christmasParty, travelToJapan}
		sort.Strings(expected)

		actual := g.GetGroupNames(groupIDs)
		sort.Strings(actual)

		assert.Equal(t, expected, actual)
	})
}

func TestGroups_AddDebt(t *testing.T) {
	t.Run("when no debts prior to addDebt", func(t *testing.T) {
		g := getGroups()
		participants := []string{peter, george, lily}
		id1 := g.CreateGroup(christmasParty, participants)

		g.AddDebt(peter, id1, participants, amount20, food)

		expected := map[string][]DebtC{christmasParty: {
			{george, amount20, food},
			{lily, amount20, food}}}
		actual := g.GetLent(peter, []uuid.UUID{id1})
		fmt.Println("ACTUAL", actual)

		assert.True(t, equal(expected, actual))
	})
}

func TestGroups_ReturnDebt(t *testing.T) {
	t.Run("return debt", func(t *testing.T) {
		g := getGroups()
		participants := []string{peter, george, lily}
		id1 := g.CreateGroup(christmasParty, participants)

		participants2 := []string{peter, george}
		id2 := g.CreateGroup(christmasParty, participants2)

		// Peter lends money to George and Lily
		g.AddDebt(peter, id1, participants, amount20, food)
		// George returns the money
		g.ReturnDebt(george, id1, peter, amount20)

		expected := map[string][]DebtC{}
		actual := g.GetOwed(peter, []uuid.UUID{id1, id2})

		assert.Equal(t, expected, actual)
	})
	t.Run("when returning less", func(t *testing.T) {
		g := getGroups()
		participants := []string{peter, george, lily}
		id1 := g.CreateGroup(christmasParty, participants)

		// Peter lends money to George and Lily
		g.AddDebt(peter, id1, participants, amount20, food)
		// George returns half of the money
		g.ReturnDebt(george, id1, peter, amount20/2)

		actualOwed := g.GetOwed(george, []uuid.UUID{id1})
		expectedOwed := map[string][]DebtC{christmasParty: {{peter, amount20 / 2, food}}}

		actualLent := g.GetLent(peter, []uuid.UUID{id1})
		expectedLent := map[string][]DebtC{christmasParty: {{george, amount20 / 2, food},
			{lily, amount20, food}}}

		assert.True(t, equal(expectedOwed, actualOwed))
		assert.True(t, equal(expectedLent, actualLent))
	})
	t.Run("when returning more", func(t *testing.T) {
		g := getGroups()
		participants := []string{peter, george, lily}
		id1 := g.CreateGroup(christmasParty, participants)

		// Peter lends money to George and Lily
		g.AddDebt(peter, id1, participants, amount20, food)
		// George returns double the money
		g.ReturnDebt(george, id1, peter, amount20*2)

		// Now Peter owes money to George
		newReason := "Repay"
		actualOwed := g.GetOwed(peter, []uuid.UUID{id1})
		expectedOwed := map[string][]DebtC{christmasParty: {{george, amount20, newReason}}}
		fmt.Println("ActualOwed", actualOwed, expectedOwed)

		actualLent := g.GetLent(peter, []uuid.UUID{id1})
		expectedLent := map[string][]DebtC{christmasParty: {{lily, amount20, food}}}
		fmt.Println("ActualLent", actualLent, expectedLent)

		assert.True(t, equal(expectedOwed, actualOwed))
		assert.True(t, equal(expectedLent, actualLent))
	})
}

func TestGroups_FindGroupID(t *testing.T) {
	t.Run("when user participates in group", func(t *testing.T) {
		g := Groups{map[uuid.UUID]Group{}}
		g.Groups[groupID1] = Group{Name: christmasParty, MoneyExchange: MoneyExchange{}}
		g.Groups[groupID2] = Group{Name: travelToJapan, MoneyExchange: MoneyExchange{}}

		participatesIn := []uuid.UUID{groupID1, groupID2}
		groupID, err := g.FindGroupID(christmasParty, participatesIn)

		assert.Equal(t, groupID, groupID1)
		assert.NoError(t, err)
	})
	t.Run("when user does not participate in group", func(t *testing.T) {
		g := Groups{map[uuid.UUID]Group{}}
		g.Groups[groupID1] = Group{Name: christmasParty, MoneyExchange: MoneyExchange{}}
		g.Groups[groupID2] = Group{Name: travelToJapan, MoneyExchange: MoneyExchange{}}

		participatesIn := []uuid.UUID{groupID2}
		groupID, err := g.FindGroupID(christmasParty, participatesIn)

		assert.Equal(t, groupID, uuid.Nil)
		assert.Error(t, err)
	})
}

func TestGroups_GetOwed(t *testing.T) {
	t.Run("when you owe money to multiple groups", func(t *testing.T) {
		groups := getGroups()

		id1 := groups.CreateGroup(christmasParty, []string{peter, george, lily})
		exchange1 := groups.Groups[id1].MoneyExchange
		// Peter owes George and Lily 20lv
		exchange1.AddDebt(peter, george, amount20, food)
		exchange1.AddDebt(peter, lily, amount20, food)

		id2 := groups.CreateGroup(travelToJapan, []string{peter, lily})
		exchange2 := groups.Groups[id2].MoneyExchange
		// Peter owes Lily 100lv ()
		exchange2.AddDebt(peter, lily, amount100, travel)

		actual := groups.GetOwed(peter, []uuid.UUID{id1, id2})
		expected := map[string][]DebtC{christmasParty: {{george, amount20, food},
			{lily, amount20, food}},
			travelToJapan: {{lily, amount100, travel}}}

		assert.True(t, equal(expected, actual))
	})
	t.Run("when you do not owe money to any group", func(t *testing.T) {
		groups := getGroups()

		id1 := groups.CreateGroup(christmasParty, []string{peter, george, lily})
		id2 := groups.CreateGroup(travelToJapan, []string{peter, lily})

		actual := groups.GetOwed(peter, []uuid.UUID{id1, id2})
		expected := map[string][]DebtC{}
		assert.True(t, equal(expected, actual))
	})
}

func TestGroups_GetLent(t *testing.T) {
	t.Run("when you have lent money to multiple groups", func(t *testing.T) {
		groups := getGroups()

		id1 := groups.CreateGroup(christmasParty, []string{peter, george, lily})
		exchange1 := groups.Groups[id1].MoneyExchange
		// Peter lends George and Lily 20lv
		exchange1.AddDebt(george, peter, amount20, food)
		exchange1.AddDebt(lily, peter, amount20, food)

		id2 := groups.CreateGroup(travelToJapan, []string{peter, lily})
		exchange2 := groups.Groups[id2].MoneyExchange
		// Peter lends Lily 100lv ()
		exchange2.AddDebt(lily, peter, amount100, travel)

		actual := groups.GetLent(peter, []uuid.UUID{id1, id2})
		expected := map[string][]DebtC{christmasParty: {{george, amount20, food},
			{lily, amount20, food}},
			travelToJapan: {{lily, amount100, travel}}}

		assert.True(t, equal(expected, actual))
	})
	t.Run("when you have not lent money to any group", func(t *testing.T) {
		groups := getGroups()

		id1 := groups.CreateGroup(christmasParty, []string{peter, george, lily})
		id2 := groups.CreateGroup(travelToJapan, []string{peter, lily})

		actual := groups.GetOwed(peter, []uuid.UUID{id1, id2})
		expected := map[string][]DebtC{}
		assert.True(t, equal(expected, actual))
	})
}
