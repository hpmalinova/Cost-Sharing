package storage

import "github.com/google/uuid"

// Participates stores 1:N relationship between username and groupID
// Ex: Pesho --> 101, 102, 103
type Participates struct {
	Participates map[string]map[uuid.UUID]struct{}
}

// Add Connects <username> to a group with ID = <groupID>
func (p *Participates) Add(username string, groupID uuid.UUID) {
	if _, ok := p.Participates[username]; !ok {
		p.Participates[username] = map[uuid.UUID]struct{}{}
	}

	p.Participates[username][groupID] = struct{}{}
}

// DoesParticipate Checks if a user with username = <username> participates in a group with ID = <groupID>
func (p *Participates) DoesParticipate(username string, groupID uuid.UUID) bool {
	_, ok := p.Participates[username][groupID]
	return ok
}

// GetGroups Returns the groupIDs of all of the groups that the user with username = <username> participates in
func (p *Participates) GetGroups(username string) []uuid.UUID {
	m := p.Participates[username] // m is map[uuid.UUID]struct{}

	groupIDs := make([]uuid.UUID, len(m))
	i := 0
	for k := range m {
		groupIDs[i] = k
		i++
	}

	return groupIDs
}
