package store

import (
	"sync"

	"queue-svc/internal/model"
)

type Store struct {
	m  map[string]*model.Task
	mu sync.RWMutex
}

func New() *Store {
	return &Store{
		m: make(map[string]*model.Task),
	}
}

func (s *Store) Put(t *model.Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[t.ID] = t
}

func (s *Store) UpdateStatus(id string, status model.Status, err string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if v, ok := s.m[id]; ok {
		v.Status = status
		v.Err = err
	}
}

func (s *Store) IncAttempts(id string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	if t, ok := s.m[id]; ok {
		t.Attempts++
		return t.Attempts
	}
	return 0
}

// read data for testing
func (s *Store) Get(id string) *model.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t := s.m[id]
	return t
}

func (s *Store) Snapshot() map[string]model.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make(map[string]model.Task, len(s.m))
	for k, v := range s.m {
		out[k] = *v
	}
	return out
}
