package storage

import "github.com/google/uuid"

const (
	username  = "Peter"
	username2 = "George"
	username3 = "Lily"
	password  = "123456"
	password2 = "qwerty"

	reason     = "Food"
	reason2    = "Travel"
	amount     = 20
	amount2    = 100
	groupName1 = "Christmas party"
	groupName2 = "Traveling to Japan"
)

var groupID1 = uuid.New()
var groupID2 = uuid.New()

func getMoneyExchange(us1, us2 string) MoneyExchange {
	m := MoneyExchange{Owes: map[string]To{}, Lends: map[string]To{}}
	m.AddUser(us1)
	m.AddUser(us2)
	return m
}

func getGroups() Groups {
	return Groups{map[uuid.UUID]Group{}}
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

func equal(a map[string][]DebtC, b map[string][]DebtC) bool {
	if len(a) != len(b) {
		return false
	}

	for k, debtsA := range a {
		if debtsB, ok := b[k]; ok {
			if !containsAll(debtsA, debtsB) {
				return false
			}
		} else {
			return false
		}
	}

	return true
}
