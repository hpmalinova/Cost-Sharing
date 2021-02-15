package storage

import "github.com/google/uuid"

// In server
type Participates struct {
	Participates map[string]map[uuid.UUID]struct{} // Pesho --> 101,102,103
}

func (p *Participates) Add(username string, groupID uuid.UUID) {
	if _, ok := p.Participates[username]; !ok {
		p.Participates[username] = map[uuid.UUID]struct{}{}
	}

	p.Participates[username][groupID] = struct{}{}
}

//func (p *Participates) Add(username string, groupIDs []uuid.UUID) {
//	p.Participates[username] = map[uuid.UUID]struct{}{}
//
//	for _, groupID := range groupIDs {
//		p.Participates[username][groupID] = struct{}{}
//	}
//}

func (p *Participates) DoesParticipate(username string, groupID uuid.UUID) bool {
	_, ok := p.Participates[username][groupID]
	return ok
}

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
