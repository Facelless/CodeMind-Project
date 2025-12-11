package usecase

import (
	"miservicegolang/core/domain/match"
	"miservicegolang/core/domain/user"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MatchUsecase struct {
	mu      sync.Mutex
	queue   []user.User
	clients map[primitive.ObjectID]chan match.MatchEvent
}

func NewMatchUsecase() *MatchUsecase {
	return &MatchUsecase{
		queue:   []user.User{},
		clients: make(map[primitive.ObjectID]chan match.MatchEvent),
	}
}

func (m *MatchUsecase) RegisterClient(id primitive.ObjectID) chan match.MatchEvent {
	ch := make(chan match.MatchEvent, 10)
	m.clients[id] = ch
	return ch
}

func (m *MatchUsecase) UnregisterClient(id primitive.ObjectID) {
	if ch, ok := m.clients[id]; ok {
		close(ch)
		delete(m.clients, id)
	}
}

func (m *MatchUsecase) sendTo(id primitive.ObjectID, evt match.MatchEvent) {
	if ch, ok := m.clients[id]; ok {
		ch <- evt
	}
}

func (m *MatchUsecase) PlayerLookingForMatch(u user.User) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.queue = append(m.queue, u)
	m.sendTo(u.ID, match.MatchEvent{
		Type:    "queue_join",
		Message: "You entered the queue",
		User:    u,
	})

	if len(m.queue) >= 2 {
		p1 := m.queue[0]
		p2 := m.queue[1]
		m.queue = m.queue[2:]
		go m.startMatch(p1, p2)
	}
}

func (m *MatchUsecase) startMatch(p1, p2 user.User) {
	m.sendTo(p1.ID, match.MatchEvent{
		Type:    "start_match",
		Message: "Match started",
		User:    p2,
	})

	m.sendTo(p2.ID, match.MatchEvent{
		Type:    "start_match",
		Message: "Match started",
		User:    p1,
	})
}
