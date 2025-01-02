package state

import "github.com/AbdessamadEnabih/Vertex/internal/state"

type State struct {
	InternalState *state.State
}

func NewState() *State {
	return &State{
		InternalState: state.NewState(),
	}
}

func (s *State) Set(key string, value interface{}) error {
	return s.InternalState.Set(key, value)
}

func (s *State) Get(key string) (interface{}, error) {
	return s.InternalState.Get(key)
}

func (s *State) Delete(key string) error {
	return s.InternalState.Delete(key)
}

func (s *State) GetAll() map[string]interface{} {
	return s.InternalState.GetAll()
}

func (s *State) Update(key string, value interface{}) error {
	return s.InternalState.Update(key, value)
}

func (s *State) FlushAll() error {
	return s.InternalState.FlushAll()
}
