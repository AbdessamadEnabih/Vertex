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

func validateKey(key string) error {
	re := regexp.MustCompile(`^[a-zA-Z0-9[\x80-\xFF]\s-_]+$`)

	// ^ Matches the start of the string
	// [a-zA-Z0-9] Matches any letter or digit
	// [\x80-\xFF] Matches any Unicode character between U+0080 and U+00FF
	// \s Matches any whitespace character
	// _ Matches underscore
	// - Matches hyphen
	// + Means one or more occurrences of the preceding element

	// Check for empty key
	if key == "" {
		return ErrEmptyKey
	}

	if re.MatchString(key) {
		return ErrSpecialCharactersKey
	}

	return nil
}

func (s *State) Set(key string, value interface{}) error {
	maxAllowedEntries := 100000

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := validateKey(key); err != nil {
		return err
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

	if err := validateKey(key); err != nil {
		return nil, err
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

func (s *State) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := validateKey(key); err != nil {
		return err
	}

	if _, ok := s.Data[key]; !ok {
		return ErrKeyNotFound
	}
	delete(s.Data, key)
	return nil
}

func (s *State) Update(key string, value interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := validateKey(key); err != nil {
		return err
	}

	if _, ok := s.Data[key]; !ok {
		return ErrKeyNotFound
	}

	s.Data[key] = value

	return nil
}
