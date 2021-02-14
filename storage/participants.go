package storage

import "github.com/google/uuid"

// In server
type Participants struct {
	Participants map[uuid.UUID]map[Username]struct{} // 101 --> Pesho, Lily, Stoyan
}

func (p *Participants) Add(groupID uuid.UUID, participants ...Username) {
	p.Participants[groupID] = map[Username]struct{}{}

	for _, user := range participants {
		p.Participants[groupID][user] = struct{}{}
	}
}

func (p *Participants) HasParticipant(groupID uuid.UUID, user Username) bool {
	_, ok := p.Participants[groupID][user]
	return ok
}

func (p *Participants) GetParticipants (groupID uuid.UUID) []Username {
	m := p.Participants[groupID] // m is map[Username]struct{}

	usernames := make([]Username, len(m))
	i:=0
	for k := range m{
		usernames[i] = k
		i++
	}

	return usernames
}

func (p *Participants) GetNumberOfParticipants(groupID uuid.UUID) int {
	return len(p.Participants[groupID])
}
