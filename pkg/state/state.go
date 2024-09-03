package state

import (
	"fmt"
	"sync"
)

type State struct {
	Data map[string]interface{}
	mu   sync.Mutex
}

func NewState() *State {
	return &State{
		Data: make(map[string]interface{}),
	}
}

func (s *State) Set(key string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Data[key] = value
}

func (s *State) Get(key string) interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	value, ok := s.Data[key]
	if !ok {
		fmt.Printf("%s not found", key)
		return nil
	}
	return value
}

func (s *State) GetAll() map[string]interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()
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
