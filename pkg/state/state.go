package state

import (
	"fmt"
	"sync"
)

type State struct {
	data map[string]interface{} // To allow any type
	mu   sync.Mutex
}

func newState() *State  {
	return &State{
		data : make(map[string]interface{}, 1000)
	}
}

func (s *State) Set(key string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}
