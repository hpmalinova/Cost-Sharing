package storage

// MoneyExchange Stores the data, that connects creditors and debtors
type MoneyExchange struct {
	Owes  map[string]To
	Lends map[string]To
}

// To Creates the relationship between a user (his username) and a Debt to that user
type To struct {
	To map[string]Debt
}

type Debt struct {
	Amount int
	Reason string
}

// AddUser Adds user with the corresponding username to the MoneyExchange struct
func (m *MoneyExchange) AddUser(username string) {
	m.Owes[username] = To{To: map[string]Debt{}}
	m.Lends[username] = To{To: map[string]Debt{}}
}

// AddDebt The <debtor> has to give the <creditor> a certain <amount>, taken for <reason>
func (m *MoneyExchange) AddDebt(debtor, creditor string, amount int, reason string) {
	// Check if creditor already owes something to debtor
	if debt, ok := m.Owes[creditor].To[debtor]; ok {
		if debt.Amount > amount {
			// Decrease the debt that the creditor has to the debtor (From the past)
			newAmount := debt.Amount - amount
			m.Owes[creditor].To[debtor] = Debt{newAmount, debt.Reason}
			m.Lends[debtor].To[creditor] = Debt{newAmount, debt.Reason}
			return
		} else if debt.Amount == amount {
			// Delete the debt, no one owes anything
			delete(m.Owes[creditor].To, debtor)
			delete(m.Lends[debtor].To, creditor)
			return
		} else {
			// Decrease the debt that the debtor has to the creditor
			amount -= debt.Amount
		}
	}

	// Check if debtor already owes something to creditor
	if debt, ok := m.Owes[debtor].To[creditor]; ok {
		// Increase the debt and concatenate the new reason
		amount += debt.Amount
		if reason != "" {
			reason = debt.Reason + ", " + reason
		} else {
			reason = debt.Reason
		}
	}

	m.Owes[debtor].To[creditor] = Debt{amount, reason}
	m.Lends[creditor].To[debtor] = Debt{amount, reason}
}

// GetOwed Returns all of the users (their usernames) that has lent money to the debtor,
// the amount of the loan and its reason
func (m *MoneyExchange) GetOwed(debtor string) []DebtC {
	return convertToClientData(m.Owes[debtor])
	//return m.Owes[debtor]
}

// GetLent Returns all of the users (their usernames) that owes money to the creditor,
// the amount of the loan and its reason
func (m *MoneyExchange) GetLent(creditor string) []DebtC {
	return convertToClientData(m.Lends[creditor])
}

func convertToClientData(input To) []DebtC {
	output := make([]DebtC, len(input.To))

	i := 0
	for user, debt := range input.To {
		output[i] = DebtC{
			To:     user,
			Amount: debt.Amount,
			Reason: debt.Reason,
		}
		i++
	}

	return output
}

type DebtC struct {
	To     string
	Amount int
	Reason string
}
