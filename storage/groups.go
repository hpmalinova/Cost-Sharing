package storage

import (
	"errors"
	"github.com/google/uuid"
)

// In server
type Groups struct {
	Groups map[uuid.UUID]Group
}

type Group struct {
	Name string
	MoneyExchange
}

func (g *Groups) CreateGroup(groupName string, participants []string) uuid.UUID {
	// TODO check if group is taken?
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

func (g *Groups) GetGroupNames(groupIDs []uuid.UUID) []string {
	groupNames := make([]string, len(groupIDs))
	for i, id := range groupIDs {
		groupNames[i] = g.GetName(id)
	}

	return groupNames
}

func (g *Groups) AddDebt(creditor string, groupID uuid.UUID, participants []string, amount int, reason string) {
	moneyEx := g.Groups[groupID].MoneyExchange // TODO
	for _, debtor := range participants {
		if debtor != creditor {
			moneyEx.AddDebt(debtor, creditor, amount, reason)
		}
	}
}

func (g *Groups) ReturnDebt(debtor string, groupID uuid.UUID, creditor string, amount int) {
	moneyEx := g.Groups[groupID].MoneyExchange // TODO

	moneyEx.AddDebt(debtor, creditor, amount, "")
}

func (g *Groups) FindGroupID(groupName string, participatesIn []uuid.UUID) (uuid.UUID, error) {
	for _, groupID := range participatesIn {
		if g.GetName(groupID) == groupName {
			return groupID, nil
		}
	}
	msg := "you don`t participate in group called " + groupName
	return uuid.Nil, errors.New(msg)
}

func (g *Groups) GetOwed(debtor string, groupIDs []uuid.UUID) map[string][]DebtC {
	owed := make(map[string][]DebtC, len(groupIDs))
	for _, groupID := range groupIDs {
		exchange := g.GetGroup(groupID).MoneyExchange
		owed[g.GetName(groupID)] = exchange.GetOwed(debtor)
	}
	return owed
}

func (g *Groups) GetLent(creditor string, groupIDs []uuid.UUID) map[string][]DebtC {
	lent := make(map[string][]DebtC, len(groupIDs))
	for _, groupID := range groupIDs {
		exchange := g.GetGroup(groupID).MoneyExchange
		lent[g.GetName(groupID)] = exchange.GetLent(creditor)
	}
	return lent
}