package storage

import (
	"errors"
	"github.com/google/uuid"
)

// Groups Connects groupIDs with Group to allow easy access to groups
type Groups struct {
	Groups map[uuid.UUID]Group
}

type Group struct {
	Name string
	MoneyExchange
}

func (g *Groups) CreateGroup(groupName string, participants []string) uuid.UUID {
	newID := uuid.New()
	g.Groups[newID] = Group{
		Name: groupName,
		MoneyExchange: MoneyExchange{
			Owes:  map[string]To{},
			Lends: map[string]To{},
		},
	}

	m := g.Groups[newID].MoneyExchange

	for _, username := range participants {
		m.AddUser(username)
	}

	return newID
}

func (g *Groups) GetName(groupID uuid.UUID) string {
	return g.Groups[groupID].Name
}

func (g *Groups) GetGroup(groupID uuid.UUID) Group {
	return g.Groups[groupID]
}

// GetGroupNames returns all of the group names that have IDs = <groupIDs>
func (g *Groups) GetGroupNames(groupIDs []uuid.UUID) []string {
	groupNames := make([]string, 0, len(groupIDs))
	for _, id := range groupIDs {
		if name := g.GetName(id); name != "" {
			groupNames = append(groupNames, name)
		}
	}

	return groupNames
}

// AddDebt The <creditor> lends <amount> lv for <reason> to each participant in the group with id = <groupID>
func (g *Groups) AddDebt(creditor string, groupID uuid.UUID, participants []string, amount int, reason string) {
	moneyEx := g.Groups[groupID].MoneyExchange
	for _, debtor := range participants {
		// Do not add debt to yourself
		if debtor != creditor {
			moneyEx.AddDebt(debtor, creditor, amount, reason)
		}
	}
}

// ReturnDebt The debtor returns <amount> lv to the creditor in the group with id = <groupID>
func (g *Groups) ReturnDebt(debtor string, groupID uuid.UUID, creditor string, amount int) {
	moneyEx := g.Groups[groupID].MoneyExchange

	moneyEx.AddDebt(creditor, debtor, amount, "Repay")
}

// FindGroupID Receives the groupIDs that a certain user participates in.
// Returns the groupID of the group with name = <groupName>
// or Returns an error if the user does not participate in a group with that name
func (g *Groups) FindGroupID(groupName string, participatesIn []uuid.UUID) (uuid.UUID, error) {
	for _, groupID := range participatesIn {
		if g.GetName(groupID) == groupName {
			return groupID, nil
		}
	}
	msg := "you don`t participate in group called " + groupName
	return uuid.Nil, errors.New(msg)
}

// GetOwed Returns all of the users from group with id = <groupID>
// that have lent money to the debtor, the amount of the loan and its reason.
func (g *Groups) GetOwed(debtor string, groupIDs []uuid.UUID) map[string][]DebtC {
	owed := make(map[string][]DebtC, len(groupIDs))
	for _, groupID := range groupIDs {
		exchange := g.GetGroup(groupID).MoneyExchange
		if !ContainsAll(exchange.GetOwed(debtor), []DebtC{}) {
			owed[g.GetName(groupID)] = exchange.GetOwed(debtor)
		}
	}
	return owed
}

// GetLent Returns all of the users from group with id = <groupID>
// that owe money to the creditor, the amount of the loan and its reason
func (g *Groups) GetLent(creditor string, groupIDs []uuid.UUID) map[string][]DebtC {
	lent := make(map[string][]DebtC, len(groupIDs))
	for _, groupID := range groupIDs {
		exchange := g.GetGroup(groupID).MoneyExchange
		lent[g.GetName(groupID)] = exchange.GetLent(creditor)
	}
	return lent
}
