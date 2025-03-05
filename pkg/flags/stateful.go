package flags

import "sync"

type Store struct {
	mux   sync.Mutex
	Store map[string]any
}

func NewStore() *Store {
	return &Store{
		Store: make(map[string]any),
	}
}

func (s *Store) Set(key string, value any) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.Store[key] = value
}

func (s *Store) Get(key string) (any, bool) {
	s.mux.Lock()
	defer s.mux.Unlock()
	value, ok := s.Store[key]
	return value, ok
}
