package stores

import (
	"sync"

	"github.com/crimsonn/zm_server/internal/models"
)

type OnlineUserStore interface {
	Set(netID string, user models.User)
	Remove(netID string)
	Get(netID string) (models.User, bool)
	List() []models.User
	Count() int
}

type InMemoryUserStore struct {
	mu    sync.RWMutex
	users map[string]models.User
}

func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		users: make(map[string]models.User),
	}
}

func (s *InMemoryUserStore) Set(netID string, user models.User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[netID] = user
}

func (s *InMemoryUserStore) Remove(netID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.users, netID)
}

func (s *InMemoryUserStore) Get(netID string) (models.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, ok := s.users[netID]
	return user, ok
}

func (s *InMemoryUserStore) List() []models.User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.users) == 0 {
		return []models.User{}
	}

	users := make([]models.User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users
}

func (s *InMemoryUserStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.users)
}
