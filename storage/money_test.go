package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	username  = "Peter"
	username2 = "George"
	username3 = "Lily"
	reason    = "Food"
	amount    = 20
)

func TestMoneyExchange_AddUser(t *testing.T) {
	t.Run("when adding user", func(t *testing.T) {
		m := MoneyExchange{
			Owes:  map[string]To{},
			Lends: map[string]To{},
		}
		m.AddUser(username)

		expected := MoneyExchange{
			Owes:  map[string]To{username: {To: map[string]Debt{}}},
			Lends: map[string]To{username: {To: map[string]Debt{}}},
		}

		assert.Equal(t, expected.Lends, m.Lends)
		assert.Equal(t, expected.Owes, m.Owes)
	})
}

func TestMoneyExchange_GetLent(t *testing.T) {
	t.Run("when no one owes you money", func(t *testing.T) {
		m := MoneyExchange{
			Owes: map[string]To{},
			Lends: map[string]To{},
		}
		actual := m.GetLent(username)
		expected := []DebtC{}
		assert.Equal(t, expected, actual)
	})
	t.Run("when multiple people owe you money", func(t *testing.T) {
		m := MoneyExchange{Lends: map[string]To{
			username: {To: map[string]Debt{
				username2: {amount, reason},
				username3: {amount, reason}}}}}
		actual := m.GetLent(username)
		expected := []DebtC{{To: username2, Amount: amount, Reason: reason},
							{To: username3, Amount: amount, Reason: reason}}
		assert.True(t, containsAll(expected, actual))
	})
}

func TestMoneyExchange_GetOwed(t *testing.T) {
	t.Run("when you do not owe money to anyone", func(t *testing.T) {
		m := MoneyExchange{
			Owes: map[string]To{},
			Lends: map[string]To{},
		}
		actual := m.GetOwed(username)
		expected := []DebtC{}
		assert.Equal(t, expected, actual)
	})
	t.Run("when you owe money to multiple people", func(t *testing.T) {
		m := MoneyExchange{Owes: map[string]To{
			username: {To: map[string]Debt{
				username2: {amount, reason},
				username3: {amount, reason}}}}}
		actual := m.GetOwed(username)
		expected := []DebtC{{To: username2, Amount: amount, Reason: reason},
			{To: username3, Amount: amount, Reason: reason}}
		assert.True(t, containsAll(expected, actual))
	})
}

func containsAll(a []DebtC, b []DebtC) bool {
	if len(a) != len(b) {
		return false
	}

	for _, elem := range a {
		if !containsElem(elem, b) {
			return false
		}
	}

	return true
}

func containsElem(elem DebtC, elems []DebtC) bool {
	for _, e := range elems {
		if e == elem {
			return true
		}
	}
	return false
}
