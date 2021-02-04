package storage

type MoneyExchange struct {
	Owes  map[Username]To
	Lends map[Username]To
}

func (m *MoneyExchange) AddUser(username Username) {
	m.Owes[username] = To{To: map[Username]Debt{}}
	m.Lends[username] = To{To: map[Username]Debt{}}
}

func (m *MoneyExchange) AddDebt(debtor, creditor Username, amount int, reason string) { // "debtor" has to give "creditor" 20lv
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
		reason = debt.Reason + ", " + reason
	}

	m.Owes[debtor].To[creditor] = Debt{amount, reason}
	m.Lends[creditor].To[debtor] = Debt{amount, reason}
}

type To struct {
	To map[Username]Debt
}

type Debt struct {
	Amount int
	Reason string
}
