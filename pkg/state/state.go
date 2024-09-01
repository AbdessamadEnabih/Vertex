package main

import (
	"fmt"
	"sync"
)

type State struct {
	data map[string]interface{}
	mu   sync.Mutex
}

func NewState() *State {
	return &State{
		data: make(map[string]interface{}),
	}
}

func (s *State) Set(key string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *State) Get(key string) (interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	value, ok := s.data[key]
	if !ok {
		fmt.Println("Key not found")
		return
	}
	return value
}
