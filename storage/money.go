package storage

type MoneyExchange struct {
	Owes  map[string]To
	Lends map[string]To
}

func (m *MoneyExchange) AddUser(username string) {
	m.Owes[username] = To{To: map[string]Debt{}}
	m.Lends[username] = To{To: map[string]Debt{}}
}

// <debtor> has to give <creditor> <amount>lv for <reason>
func (m *MoneyExchange) AddDebt(debtor, creditor string, amount int, reason string) {
	// Check if creditor owes something to debtor (In the past)
	if debt, ok := m.Owes[creditor].To[debtor]; ok {
		if debt.Amount > amount {
			newAmount := debt.Amount - amount
			m.Owes[creditor].To[debtor] = Debt{newAmount, debt.Reason}
			m.Lends[debtor].To[creditor] = Debt{newAmount, debt.Reason}
			return
		} else if debt.Amount == amount {
			delete(m.Owes[creditor].To, debtor)
			delete(m.Lends[debtor].To, creditor)
			return
		} else {
			amount -= debt.Amount
		}
	}

	// Check if debtor already owes something to creditor
	if debt, ok := m.Owes[debtor].To[creditor]; ok {
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

func (m *MoneyExchange) GetOwed(debtor string) []DebtC {
	//return m.Owes[debtor]
	return convertToClientData(m.Owes[debtor])
}

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

type To struct {
	To map[string]Debt
}

type Debt struct {
	Amount int
	Reason string
}

type DebtC struct {
	To     string
	Amount int
	Reason string
}
