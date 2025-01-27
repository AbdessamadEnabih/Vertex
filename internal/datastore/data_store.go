package datastore

import (
	"regexp"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

const defaultTTL = time.Minute * 60
const maxAllowedEntries = 100000

type DataStore struct {
	Data   map[string]interface{}
	cache  *cache.Cache
	ttlMap map[string]time.Time
	mu     sync.RWMutex
}
type DataStoreError struct {
	Cause   error
	Message string
}

// Declared Errors
var (
	ErrOutOfMemory          = &DataStoreError{Message: "Out of memory"}
	ErrEmptyKey             = &DataStoreError{Message: "Empty key is not allowed"}
	ErrNilValue             = &DataStoreError{Message: "Nil value is not allowed"}
	ErrKeyNotFound          = &DataStoreError{Message: "Key not found"}
	ErrDuplicateKey         = &DataStoreError{Message: "Key already exists"}
	ErrSpecialCharactersKey = &DataStoreError{Message: "Key with special characters is not allowed"}
)

func (e *DataStoreError) Error() string { return e.Message }
func NewDataStore() *DataStore {
	return &DataStore{
		Data:   make(map[string]interface{}),
		cache:  cache.New(5*time.Minute, 30*time.Minute),
		ttlMap: make(map[string]time.Time),
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
func (s *DataStore) Set(key string, value interface{}) error {
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

	if s.Data == nil {
        s.Data = make(map[string]interface{})
    }
	if s.ttlMap == nil {
		s.ttlMap = make(map[string]time.Time)
	}
	
	s.Data[key] = value
	s.ttlMap[key] = time.Now().Add(defaultTTL)
	// Check for out-of-memory error
	if len(s.Data) > maxAllowedEntries {
		delete(s.Data, key)
		return ErrOutOfMemory
	}
	return nil
}
func (s *DataStore) Get(key string) (interface{}, error) {
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
func (s *DataStore) GetAll() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Data
}
func (s *DataStore) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := validateKey(key); err != nil {
		return err
	}
	if _, ok := s.Data[key]; !ok {
		return ErrKeyNotFound
	}
	delete(s.Data, key)

	if s.ttlMap != nil {
		delete(s.ttlMap, key)
	}
	return nil
}
func (s *DataStore) Update(key string, value interface{}) error {
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
func (s *DataStore) FlushAll() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Data = make(map[string]interface{})

	if s.cache != nil {
		s.cache.Flush()
	}

	if s.ttlMap != nil {
		s.ttlMap = make(map[string]time.Time)
	}

	return nil
}
