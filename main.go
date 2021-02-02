package main

import "fmt"

type User struct {
	Username
	password string
}

type Friends map[Username][]Username // "Pesho" -> ["Iliq", "Stoqn"]

type Username string

type Debt struct {
	Amount int
	Reason string // TODO predefined reasons
}

func (d* Debt) Add(amount int, reason string){ // TODO how to update Reason
	d.Amount += amount
	d.Reason = reason
}

// type Loan?

type To struct{
	To map[Username]Debt
}

type Owes struct {
	Owes map[Username]To
}

//type Owes map[Username]To

func (o* Owes) AddUser(username Username) { // when adding new users ?
	o.Owes[username] = To{To: map[Username]Debt{}}
}
// type {Owes Lends}
func (o* Owes) AddDebt(owes Username, lends Username, amount int, reason string) { // "owes" has to give "lends" 20lv
	//o.Owes[username] = To{To: map[Username]Debt{}}
	// (*o)[username]
	if debt, ok := o.Owes[owes].To[lends]; ok {
		newAmount := debt.Amount + amount
		o.Owes[owes].To[lends] = Debt{newAmount, reason}
	} else{
		o.Owes[owes].To[lends] = Debt{Amount: amount, Reason: reason}
	}
}


type Lends map[Username]To

//type Owes map[Username] OwesTo
//
//type Lends map[Username] LendsTo

//type LendsTo struct {
//	to map[Username]Debt
//}


// Pesho --> []Struct == [{Ivana, 50, food}, {Rosi, 200, hotel}]
// {Pesho, Ivana} --> {50, food}
// {Pesho, Rosi} --> {200, hotel}

// User:
// Pesho, ****
// owesTo username --> debt, reason
// lendsTo username --> amount, reason

// Group:

func main() {
	p := Username("Pesho")
	s := Username("Silvia")
	r := Username("Rado")
	o := Owes{Owes: map[Username]To{}} // Create map
	o.AddUser(p)
	o.AddUser(s)
	o.AddUser(r)
	//fmt.Println(o)

	o.AddDebt(p,s, 20, "nz")
	o.AddDebt(p,s, 30, "new")
	o.AddDebt(p,r,10,"pa")
	fmt.Println(o)

//////////////////////////////
	//p := Username("Pesho")
	//s := Username("Silvia")
	////r := Username("Rado")
	//o := make(Owes)
	//o[p] = To{To: map[Username]Debt{}}
	//o[p].To[s] = Debt{
	//	Amount: 10,
	//	Reason: "food",
	//}
	//
	////fmt.Println(o[p])
	//
	//// add new debt
	//if debt, ok := o[p].To[s]; ok {
	//	//debt := o[p].To[s]
	//	fmt.Println("here")
	//	debt.Add(50, "transport")
	//	fmt.Println(o)
	//} else{
	//	o.Add(s)
	//	o[p].To[s] = Debt{30, "hotel"}
	//}
	//
	//fmt.Println(o)


}
