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
