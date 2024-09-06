package state

import (
	"errors"
	"fmt"
	"sync"
)

type State struct {
	Data map[string]interface{}
	mu   sync.RWMutex
}

func NewState() *State {
	return &State{
		Data: make(map[string]interface{}),
	}
}

func (s *State) Set(key string, value interface{}) error {
	maxAllowedEntries := 100000

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check for empty key
	if key == "" {
		return errors.New("empty key is not allowed")
	}

	// Check for nil value
	if value == nil {
		return errors.New("nil value is not allowed")
	}

	s.Data[key] = value

	// Check for out-of-memory error
	if len(s.Data) > maxAllowedEntries {
		delete(s.Data, key)
		return fmt.Errorf("out of memory: maximum entries (%d) exceeded", maxAllowedEntries)
	}

	return nil
}

func (s *State) Get(key string) interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.Data[key]
	if !ok {
		fmt.Printf("%s not found", key)
		return nil
	}
	return value
}

func (s *State) GetAll() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Data
}

func (s *State) Delete(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.Data[key]; !ok {
		return false
	}
	delete(s.Data, key)
	return true
}

func (s *State) Update(key string, value interface{}) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.Data[key]; !ok {
		return false
	}

	s.Data[key] = value

	return true
}
