package storage

import (
	"errors"
	"sync"
)

type Link struct {
	Original string
	Short    string
}

type Storage struct {
	links map[string]Link
	mu    sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		links: make(map[string]Link),
	}
}

func (s *Storage) SaveLink(original string, short string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.links[short] = Link{Original: original, Short: short}
}

func (s *Storage) GetLink(short string) (Link, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	link, exists := s.links[short]
	if !exists {
		return Link{}, errors.New("link not found")
	}
	return link, nil
}