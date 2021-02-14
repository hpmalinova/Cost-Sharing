package storage

import "github.com/google/uuid"

// In server
type Participates struct {
	Participates map[Username]map[uuid.UUID]struct{} // Pesho --> 101,102,103
}

func (p *Participates) Add(user Username, groupIDs ...uuid.UUID) {
	p.Participates[user] = map[uuid.UUID]struct{}{}

	for _, groupID := range groupIDs {
		p.Participates[user][groupID] = struct{}{}
	}
}

func (p *Participates) DoesParticipate(user Username, groupID uuid.UUID) bool {
	_, ok := p.Participates[user][groupID]
	return ok
}

func (p *Participates) GetGroups(username Username) []uuid.UUID {
	m := p.Participates[username] // m is map[uuid.UUID]struct{}

	groupIDs := make([]uuid.UUID, len(m))
	i := 0
	for k := range m {
		groupIDs[i] = k
		i++
	}

	return groupIDs
}
