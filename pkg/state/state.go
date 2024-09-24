package state

import "github.com/AbdessamadEnabih/Vertex/internal/state"

type State struct {
	internalState *state.State
}

func NewState() *State {
	return &State{
		internalState: state.NewState(),
	}
}

func (s *State) Set(key string, value interface{}) error {
	return s.internalState.Set(key, value)
}

func (s *State) Get(key string) (interface{}, error) {
	return s.internalState.Get(key)
}

func (s *State) Delete(key string) error {
	return s.internalState.Delete(key)
}

func (s *State) GetAll() map[string]interface{} {
	return s.internalState.GetAll()
}

func (s *State) FlushAll() error {
	return s.internalState.FlushAll()
}
