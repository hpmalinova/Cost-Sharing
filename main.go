package main

import "fmt"

type User struct {
	Username
	password string
}

type Friends map[Username][]Username // "Pesho" -> ["Iliq", "Stoqn"]

type Username string

type MoneyExchange struct {
	Owes map[Username]To
	Lends map[Username]To
}

func (m* MoneyExchange) AddUser(username Username) {
	m.Owes[username] = To{To: map[Username]Debt{}}
	m.Lends[username] = To{To: map[Username]Debt{}}
}

func (m* MoneyExchange) AddDebt(owes Username, lends Username, amount int, reason string) { // "owes" has to give "lends" 20lv
	if debt, ok := m.Owes[owes].To[lends]; ok {
		newAmount := debt.Amount + amount
		newReason := debt.Reason + ", " + reason
		m.Owes[owes].To[lends] = Debt{newAmount, newReason}
	} else{
		m.Owes[owes].To[lends] = Debt{Amount: amount, Reason: reason}
	}
	if debt, ok := m.Lends[lends].To[owes]; ok {
		newAmount := debt.Amount + amount
		m.Lends[lends].To[owes] = Debt{newAmount, reason}
	} else{
		m.Lends[lends].To[owes] = Debt{Amount: amount, Reason: reason}
	}

	//
}

type To struct{
	To map[Username]Debt
}

type Debt struct {
	Amount int
	Reason string // TODO predefined reasons? // TODO concatenate strings
}

// Pesho --> []Struct == [{Ivana, 50, food}, {Rosi, 200, hotel}]
// {Pesho, Ivana} --> {50, food}
// {Pesho, Rosi} --> {200, hotel}

// User:
// Pesho, ****
// owesTo username --> debt, reason
// lendsTo username --> amount, reason

// Group:

func main() {
	m := MoneyExchange{
		Owes:  map[Username]To{},
		Lends: map[Username]To{},
	}
	p := Username("Pesho")
	m.AddUser(p)
	s := Username("Silvia")
	m.AddUser(s)
	m.AddDebt(p, s, 20, "hotel")
	m.AddDebt(p,s, 10, "idk")
	fmt.Println(m)
}
