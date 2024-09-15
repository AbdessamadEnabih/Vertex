package state

import (
	"regexp"
	"sync"
)

type State struct {
	Data map[string]interface{}
	mu   sync.RWMutex
}

type StateError struct {
	Message string
	Cause   error
}

// Declared Errors
var (
	ErrOutOfMemory          = &StateError{Message: "Out of memory"}
	ErrEmptyKey             = &StateError{Message: "Empty key is not allowed"}
	ErrNilValue             = &StateError{Message: "Nil value is not allowed"}
	ErrKeyNotFound          = &StateError{Message: "Key not found"}
	ErrDuplicateKey         = &StateError{Message: "Key already exists"}
	ErrSpecialCharactersKey = &StateError{Message: "Key with special characters is not allowed"}
)

func (e *StateError) Error() string { return e.Message }

func NewState() *State {
	return &State{
		Data: make(map[string]interface{}),
	}
}

// isKeyValid checks if a given key is valid according to the following rules:
//
//  1. The key must be non-empty.
//  2. The key must only contain ASCII letters, digits, whitespace, underscores (_),
//     hyphens (-), and Unicode characters in the range U+0080 to U+00FF.
//
// This validation helps prevent potential issues with keys containing special
// characters or invalid Unicode sequences.
func isKeyValid(key string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9[\x80-\xFF]\s-_]+$`)
	return !re.MatchString(key)
}

func (s *State) Set(key string, value interface{}) error {
	maxAllowedEntries := 100000

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check for empty key
	if key == "" {
		return ErrEmptyKey
	}

	if !isKeyValid(key) {
		return ErrSpecialCharactersKey
	}

	// Check for nil value
	if value == nil {
		return ErrNilValue
	}

	if _, exists := s.Data[key]; exists {
		return ErrDuplicateKey
	}

	s.Data[key] = value

	// Check for out-of-memory error
	if len(s.Data) > maxAllowedEntries {
		delete(s.Data, key)
		return ErrOutOfMemory
	}

	return nil
}

func (s *State) Get(key string) (interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if key == "" {
		return nil, ErrEmptyKey
	}

	value, ok := s.Data[key]
	if !ok {
		return nil, ErrKeyNotFound
	}
	return value, nil
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
