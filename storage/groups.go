package storage

import "github.com/google/uuid"

// In server
type Groups struct {
	Groups map[uuid.UUID]Group
}

type Group struct {
	Name string
	MoneyExchange
}

func (g *Groups) CreateGroup(name string, participants ...Username) {
	newID := uuid.New()
	g.Groups[newID] = Group{
		Name:          name,
		MoneyExchange: MoneyExchange{},
	}

	//m := g.Groups[newID].MoneyExchange // TODO fix: add debt, add people to moneyExchange
	m := g.Groups[newID].MoneyExchange

	for _, user := range participants {
		m.AddUser(user)
	}
}

func (g *Groups) GetName(groupID uuid.UUID) string {
	return g.Groups[groupID].Name
}

func (g *Groups) GetGroup(groupID uuid.UUID) Group {
	return g.Groups[groupID]
}

func (g *Groups) GetGroupNames(groupIDs ...uuid.UUID) []string {
	groupNames := make([]string, len(groupIDs))
	for i, id := range groupIDs {
		groupNames[i] = g.GetName(id)
	}

	return groupNames
}

func (g *Groups) AddDebt(creditor Username, groupID uuid.UUID, participants []Username, amount int, reason string) {
	moneyEx := g.Groups[groupID].MoneyExchange
	for _, debtor := range participants {
		if debtor != creditor {
			moneyEx.AddDebt(debtor, creditor, amount, reason)
		}
	}
}

// TODO
// create group // TODO Now only unique group names per person
// show group(username)
// show participants(groupid)
// add debt(id, amount, reason)
// add person to group/ remove person from group?
// leave group - you cant leave the group, you owe money to ***
// kick person? (group owner/admin)
