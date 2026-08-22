package presence

import "sync"

// Status describes whether the owner is available to handle conversations.
type Status string

const (
	Online  Status = "online"
	Busy    Status = "busy"
	Away    Status = "away"
	Offline Status = "offline"
	AIOnly  Status = "ai_only"
)

// Store keeps owner presence in memory for the current process.
// Persistence can be added later without changing the consumer API.
type Store struct {
	mu     sync.RWMutex
	status Status
}

func NewStore(initial Status) *Store {
	if !valid(initial) {
		initial = Offline
	}
	return &Store{status: initial}
}

func (s *Store) Get() Status {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.status
}

func (s *Store) Set(status Status) bool {
	if !valid(status) {
		return false
	}
	s.mu.Lock()
	s.status = status
	s.mu.Unlock()
	return true
}

func valid(status Status) bool {
	switch status {
	case Online, Busy, Away, Offline, AIOnly:
		return true
	default:
		return false
	}
}

// MayAutoReply determines whether an inbound client message may be answered
// automatically based solely on owner presence. Content policy is evaluated separately.
func MayAutoReply(status Status) bool {
	return status == Offline || status == AIOnly
}
