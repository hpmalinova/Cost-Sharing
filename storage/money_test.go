package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMoneyExchange_AddUser(t *testing.T) {
	t.Run("when adding user", func(t *testing.T) {
		m := MoneyExchange{
			Owes:  map[string]To{},
			Lends: map[string]To{},
		}
		m.AddUser(peter)

		expected := MoneyExchange{
			Owes:  map[string]To{peter: {To: map[string]Debt{}}},
			Lends: map[string]To{peter: {To: map[string]Debt{}}},
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
		actual := m.GetLent(peter)
		expected := []DebtC{}
		assert.Equal(t, expected, actual)
	})
	t.Run("when multiple people owe you money", func(t *testing.T) {
		m := MoneyExchange{Lends: map[string]To{
			peter: {To: map[string]Debt{
				george: {amount20, food},
				lily:   {amount20, food}}}}}
		actual := m.GetLent(peter)
		expected := []DebtC{{To: george, Amount: amount20, Reason: food},
			{To: lily, Amount: amount20, Reason: food}}
		assert.True(t, ContainsAll(expected, actual))
	})
}

func TestMoneyExchange_GetOwed(t *testing.T) {
	t.Run("when you do not owe money to anyone", func(t *testing.T) {
		m := MoneyExchange{
			Owes:  map[string]To{},
			Lends: map[string]To{},
		}
		actual := m.GetOwed(peter)
		expected := []DebtC{}
		assert.Equal(t, expected, actual)
	})
	t.Run("when you owe money to multiple people", func(t *testing.T) {
		m := MoneyExchange{Owes: map[string]To{
			peter: {To: map[string]Debt{
				george: {amount20, food},
				lily:   {amount20, food}}}}}
		actual := m.GetOwed(peter)
		expected := []DebtC{{To: george, Amount: amount20, Reason: food},
			{To: lily, Amount: amount20, Reason: food}}
		assert.True(t, ContainsAll(expected, actual))
	})
}

func TestMoneyExchange_AddDebt(t *testing.T) {
	t.Run("test owes when the struct is empty", func(t *testing.T) {
		m := getMoneyExchange(peter, george)
		m.AddDebt(peter, george, amount20, food)

		expected := []DebtC{{george, amount20, food}}

		actual := m.GetOwed(peter)

		assert.Equal(t, expected, actual)
	})
	t.Run("test lends when the struct is empty", func(t *testing.T) {
		m := getMoneyExchange(peter, george)
		m.AddDebt(peter, george, amount20, food)

		expected := []DebtC{{peter, amount20, food}}

		actual := m.GetLent(george)

		assert.Equal(t, expected, actual)
	})
	t.Run("when the creditor owes more to the debtor", func(t *testing.T) {
		m := getMoneyExchange(peter, george)
		// Peter takes 20lv from George(in the past)
		m.AddDebt(peter, george, amount20, food)
		// George takes 100lv from Peter (now)
		m.AddDebt(george, peter, amount100, travel)
		// => George owes Peter 80lv
		expectedOwed := []DebtC{{peter, amount100 - amount20, travel}}
		actualOwed := m.GetOwed(george)

		// => Peter lends 80lv to George
		expectedLent := []DebtC{{george, amount100 - amount20, travel}}
		actualLent := m.GetLent(peter)

		// => Peter does not owe any money to George
		expectedEmpty := []DebtC{}
		actualEmptyOwed := m.GetOwed(peter)
		actualEmptyLent := m.GetLent(george)

		assert.Equal(t, expectedOwed, actualOwed)
		assert.Equal(t, expectedLent, actualLent)
		assert.Equal(t, expectedEmpty, actualEmptyOwed)
		assert.Equal(t, expectedEmpty, actualEmptyLent)

	})
	t.Run("when the creditor owes the same to the debtor", func(t *testing.T) {
		m := getMoneyExchange(peter, george)
		// Peter takes 20lv from George (in the past)
		m.AddDebt(peter, george, amount20, food)
		// George takes 20lv from Peter (now)
		m.AddDebt(george, peter, amount20, travel)
		// => Friendship is restored
		expectedEmpty := []DebtC{}
		actualOwed1 := m.GetOwed(peter)
		actualOwed2 := m.GetOwed(george)
		actualLent1 := m.GetLent(peter)
		actualLent2 := m.GetLent(george)

		assert.Equal(t, expectedEmpty, actualOwed1)
		assert.Equal(t, expectedEmpty, actualOwed2)
		assert.Equal(t, expectedEmpty, actualLent1)
		assert.Equal(t, expectedEmpty, actualLent2)
	})
	t.Run("when the creditor owes less to the debtor", func(t *testing.T) {
		m := getMoneyExchange(peter, george)
		// Peter takes 100lv from George(in the past)
		m.AddDebt(peter, george, amount100, food)
		// George takes 20lv from Peter (now)
		m.AddDebt(george, peter, amount20, travel)
		// => Peter owes George 80lv
		expectedOwed := []DebtC{{george, amount100 - amount20, food}}
		actualOwed := m.GetOwed(peter)
		// => George lends 80lv to Peter
		expectedLent := []DebtC{{peter, amount100 - amount20, food}}
		actualLent := m.GetLent(george)

		// => George does not owe any money to Peter
		expectedEmpty := []DebtC{}
		actualEmptyOwed := m.GetOwed(george)
		actualEmptyLent := m.GetLent(peter)

		assert.Equal(t, expectedOwed, actualOwed)
		assert.Equal(t, expectedLent, actualLent)
		assert.Equal(t, expectedEmpty, actualEmptyOwed)
		assert.Equal(t, expectedEmpty, actualEmptyLent)

	})
	t.Run("when the debtor owes more money to the creditor", func(t *testing.T) {
		m := getMoneyExchange(peter, george)
		// Peter takes 20lv from George(in the past)
		m.AddDebt(peter, george, amount20, food)
		// Peter takes 100lv from George (now)
		m.AddDebt(peter, george, amount100, travel)
		// => Peter owes George 120lv
		newReason := "" + food + ", " + travel
		expectedOwed := []DebtC{{george, amount20 + amount100, newReason}}
		actualOwed := m.GetOwed(peter)
		// => George lends 80lv to Peter
		expectedLent := []DebtC{{peter, amount100 + amount20, newReason}}
		actualLent := m.GetLent(george)

		// => George does not owe any money to Peter
		expectedEmpty := []DebtC{}
		actualEmptyOwed := m.GetOwed(george)
		actualEmptyLent := m.GetLent(peter)

		assert.Equal(t, expectedOwed, actualOwed)
		assert.Equal(t, expectedLent, actualLent)
		assert.Equal(t, expectedEmpty, actualEmptyOwed)
		assert.Equal(t, expectedEmpty, actualEmptyLent)
	})
}
