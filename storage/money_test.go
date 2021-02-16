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
	reason2   = "Travel"
	amount    = 20
	amount2   = 100
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
			Owes:  map[string]To{},
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
			Owes:  map[string]To{},
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

func TestMoneyExchange_AddDebt(t *testing.T) {
	t.Run("test owes when the struct is empty", func(t *testing.T) {
		m := getMoneyExchange()
		m.AddDebt(username, username2, amount, reason)

		expected := []DebtC{{username2, amount, reason}}

		actual := m.GetOwed(username)

		assert.Equal(t, expected, actual)
		//expected := MoneyExchange{Owes: map[string]To{username: {To: map[string]Debt{username2: {amount, reason} }}}}
		//debt := DebtC{To: username2, Amount: amount, Reason: reason}
		//assert.Equal(t, expected, m)
	})
	t.Run("test lends when the struct is empty", func(t *testing.T) {
		m := getMoneyExchange()
		m.AddDebt(username, username2, amount, reason)

		expected := []DebtC{{username, amount, reason}}

		actual := m.GetLent(username2)

		assert.Equal(t, expected, actual)
	})
	t.Run("when the creditor owes more to the debtor", func(t *testing.T) {
		m := getMoneyExchange()
		// Peter takes 20lv from George(in the past)
		m.AddDebt(username, username2, amount, reason)
		// George takes 100lv from Peter (now)
		m.AddDebt(username2, username, amount2, reason2)
		// => George owes Peter 80lv
		expectedOwed := []DebtC{{username, amount2 - amount, reason2}}
		actualOwed := m.GetOwed(username2)

		// => Peter lends 80lv to George
		expectedLent := []DebtC{{username2, amount2 - amount, reason2}}
		actualLent := m.GetLent(username)

		// => Peter does not owe any money to George
		expectedEmpty := []DebtC{}
		actualEmptyOwed := m.GetOwed(username)
		actualEmptyLent := m.GetLent(username2)

		assert.Equal(t, expectedOwed, actualOwed)
		assert.Equal(t, expectedLent, actualLent)
		assert.Equal(t, expectedEmpty, actualEmptyOwed)
		assert.Equal(t, expectedEmpty, actualEmptyLent)

	})
	t.Run("when the creditor owes the same to the debtor", func(t *testing.T) {
		m := getMoneyExchange()
		// Peter takes 20lv from George (in the past)
		m.AddDebt(username, username2, amount, reason)
		// George takes 20lv from Peter (now)
		m.AddDebt(username2, username, amount, reason2)
		// => Friendship is restored
		expectedEmpty := []DebtC{}
		actualOwed1 := m.GetOwed(username)
		actualOwed2 := m.GetOwed(username2)
		actualLent1 := m.GetLent(username)
		actualLent2 := m.GetLent(username2)

		assert.Equal(t, expectedEmpty, actualOwed1)
		assert.Equal(t, expectedEmpty, actualOwed2)
		assert.Equal(t, expectedEmpty, actualLent1)
		assert.Equal(t, expectedEmpty, actualLent2)
	})
	t.Run("when the creditor owes less to the debtor", func(t *testing.T) {
		m := getMoneyExchange()
		// Peter takes 100lv from George(in the past)
		m.AddDebt(username, username2, amount2, reason)
		// George takes 20lv from Peter (now)
		m.AddDebt(username2, username, amount, reason2)
		// => Peter owes George 80lv
		expectedOwed := []DebtC{{username2, amount2 - amount, reason}}
		actualOwed := m.GetOwed(username)
		// => George lends 80lv to Peter
		expectedLent := []DebtC{{username, amount2 - amount, reason}}
		actualLent := m.GetLent(username2)

		// => George does not owe any money to Peter
		expectedEmpty := []DebtC{}
		actualEmptyOwed := m.GetOwed(username2)
		actualEmptyLent := m.GetLent(username)

		assert.Equal(t, expectedOwed, actualOwed)
		assert.Equal(t, expectedLent, actualLent)
		assert.Equal(t, expectedEmpty, actualEmptyOwed)
		assert.Equal(t, expectedEmpty, actualEmptyLent)

	})
	t.Run("when the debtor owes more money to the creditor", func(t *testing.T) {
		m := getMoneyExchange()
		// Peter takes 20lv from George(in the past)
		m.AddDebt(username, username2, amount, reason)
		// Peter takes 100lv from George (now)
		m.AddDebt(username, username2, amount2, reason2)
		// => Peter owes George 120lv
		newReason := "" + reason + ", " + reason2
		expectedOwed := []DebtC{{username2, amount + amount2, newReason}}
		actualOwed := m.GetOwed(username)
		// => George lends 80lv to Peter
		expectedLent := []DebtC{{username, amount2 + amount, newReason}}
		actualLent := m.GetLent(username2)

		// => George does not owe any money to Peter
		expectedEmpty := []DebtC{}
		actualEmptyOwed := m.GetOwed(username2)
		actualEmptyLent := m.GetLent(username)

		assert.Equal(t, expectedOwed, actualOwed)
		assert.Equal(t, expectedLent, actualLent)
		assert.Equal(t, expectedEmpty, actualEmptyOwed)
		assert.Equal(t, expectedEmpty, actualEmptyLent)
	})
}

// Help functions
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

func getMoneyExchange() MoneyExchange {
	m := MoneyExchange{Owes: map[string]To{}, Lends: map[string]To{}}
	m.AddUser(username)
	m.AddUser(username2)
	return m
}
