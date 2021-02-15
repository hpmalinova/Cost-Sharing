package storage

import "github.com/google/uuid"

// Participants stores 1:N relationship between groupID and username
// Ex: 101 --> Pesho, Lily, Stoyan
type Participants struct {
	Participants map[uuid.UUID]map[string]struct{}
}

// Add Connects <participants> to a group with ID = <groupID>
func (p *Participants) Add(groupID uuid.UUID, participants []string) {
	p.Participants[groupID] = map[string]struct{}{}

	for _, user := range participants {
		user := user
		p.Participants[groupID][user] = struct{}{}
	}
}

// HasParticipant Checks if a group with ID = <groupID> has a member = <user>
func (p *Participants) HasParticipant(groupID uuid.UUID, user string) bool {
	_, ok := p.Participants[groupID][user]
	return ok
}

// GetParticipants Returns the usernames of all of the participants in group with ID = <groupID>
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

// GetNumberOfParticipants Returns the number of participants in group with ID = <groupID>
func (p *Participants) GetNumberOfParticipants(groupID uuid.UUID) int {
	return len(p.Participants[groupID])
}
