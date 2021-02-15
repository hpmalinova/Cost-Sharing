package storage

import "github.com/google/uuid"

// In server
type Participants struct {
	Participants map[uuid.UUID]map[string]struct{} // 101 --> Pesho, Lily, Stoyan
}

func (p *Participants) Add(groupID uuid.UUID, participants []string) {
	p.Participants[groupID] = map[string]struct{}{}

	for _, user := range participants {
		user := string(user)
		p.Participants[groupID][user] = struct{}{}
	}
}

func (p *Participants) HasParticipant(groupID uuid.UUID, user string) bool {
	_, ok := p.Participants[groupID][user]
	return ok
}

func (p *Participants) GetParticipants(groupID uuid.UUID) []string {
	m := p.Participants[groupID] // m is map[string]struct{}

	usernames := make([]string, len(m))
	i := 0
	for k := range m {
		usernames[i] = k
		i++
	}

	return usernames
}

func (p *Participants) GetNumberOfParticipants(groupID uuid.UUID) int {
	return len(p.Participants[groupID])
}
