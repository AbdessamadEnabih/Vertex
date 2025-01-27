package datastore_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/AbdessamadEnabih/Vertex/internal/datastore"
)

var dataStore datastore.DataStore = *datastore.NewDataStore()

func Test(t *testing.T) {

	// Check if the type of dataStore.data is map[string]interface{}
	if reflect.TypeOf(dataStore.Data) != reflect.TypeOf(map[string]interface{}{}) {
		t.Errorf("Expected dataStore.Data to be of type map[string]interface{}, but got %T", dataStore.Data)
	}

	// Test setting and getting values
	tests := []struct {
		key      string
		expected interface{}
	}{
		{"key1", "value1"},
		{"key2", 45},
	}

	for _, tt := range tests {
		dataStore.Set(tt.key, tt.expected)

		if got, _ := dataStore.Get(tt.key); !reflect.DeepEqual(got, tt.expected) {
			t.Errorf("Expected %s to be %v, but got %v", tt.key, tt.expected, got)
		}
	}

	// Test Empty key
	emptyKey := ""
	if got, err := dataStore.Get(emptyKey); got != nil && err.Error() == "Empty key is not allowed" {
		t.Errorf("Expected nil and error for empty key, but got %v", got)
	}

	// Test for an invalid key
	invalidKey := "invalidKey"
	if got, err := dataStore.Get(invalidKey); got != nil && err.Error() == "Key not found" {
		t.Errorf("Expected nil for key %s, but got %v", invalidKey, got)
	}

	// Call GetAll and compare results
	result := dataStore.GetAll()

	// Check if all keys exist in the result
	if len(result) != 2 {
		t.Errorf("Expected 2 items, got %d", len(result))
	}

}

func TestDataStore_GetAll_Empty(t *testing.T) {
	s := datastore.NewDataStore()

	// Call GetAll on empty DataStore
	result := s.GetAll()

	// Check if result is empty
	if len(result) != 0 {
		t.Errorf("Expected empty map, got %d items", len(result))
	}
}

func TestDataStore_Delete(t *testing.T) {
	// Add some test data
	dataStore.Set("testKey1", "value1")
	dataStore.Set("testKey2", 42)

	// Test deleting a non-existent key
	err := dataStore.Delete("nonExistentKey")
	if err == nil || err.Error() != "Key not found" {
		t.Errorf("Expected error 'Key not found' when deleting non-existent key, but got %v", err)
	}

	// Test deleting an existing key
	err = dataStore.Delete("testKey1")
	if err != nil {
		t.Errorf("Expected nil when deleting existing key, but got error: %v", err)
	}

	// Verify that the key was removed
	got, _ := dataStore.Get("testKey1")
	if got != nil {
		t.Errorf("Expected nil for deleted key, but got %v", got)
	}

	// Test deleting another existing key
	err = dataStore.Delete("testKey2")
	if err != nil {
		t.Errorf("Expected nil when deleting existing key, but got error: %v", err)
	}

	// Verify that the second key was also removed
	got, _ = dataStore.Get("testKey2")
	if got != nil {
		t.Errorf("Expected nil for deleted key, but got %v", got)
	}

	// Test deleting the same key again (should return an error)
	err = dataStore.Delete("testKey2")
	if err == nil || err.Error() != "Key not found" {
		t.Errorf("Expected error 'Key not found' when deleting already deleted key, but got %v", err)
	}
}

func TestDataStore_Update(t *testing.T) {
	// Create a new DataStore instance for this test
	var updateDataStore datastore.DataStore = *datastore.NewDataStore()

	updateDataStore.Set("key1", 22)

	value := "updatedkey"
	updateDataStore.Update("key1", value)
	if got, _ := updateDataStore.Get("key1"); got != value {
		t.Fatalf("Expected value to be %v and got %v", value, got)
	}

	// Test update of non existent key
	errExpected := updateDataStore.Update("invalidkey", value)
	if errExpected.Error() != "Key not found" {
		t.Errorf("Expected false when deleting non-existent key, got true")
	}

	// Clean up
	updateDataStore.Delete("key1")
	updateDataStore.Delete("invalidkey")
}

func TestDataStore_LargeKey(t *testing.T) {
	s := datastore.NewDataStore()
	largeKey := strings.Repeat("x", 1024*1024) // 1MB key
	value := "test"

	s.Set(largeKey, value)
	got, _ := s.Get(largeKey)
	if got != value {
		t.Errorf("Expected %s, got %v", value, got)
	}
}

func TestDataStore_LargeValue(t *testing.T) {
	s := datastore.NewDataStore()
	key := "test"
	largeValue := strings.Repeat("x", 1024*1024) // 1MB value

	s.Set(key, largeValue)
	got, _ := s.Get(key)
	if got != largeValue {
		t.Errorf("Expected %s, got %v", largeValue, got)
	}
}

func TestDataStore_EmptyKey(t *testing.T) {
	s := datastore.NewDataStore()

	err := s.Set("", "test")
	if err.Error() != "Empty key is not allowed" {
		t.Errorf("Expected Set to return Error %s", err.Error())
	}
}

func TestDataStore_EmptyValue(t *testing.T) {
	s := datastore.NewDataStore()

	s.Set("key", "")
	got, _ := s.Get("key")
	if got != "" {
		t.Errorf("Expected '', got %v", got)
	}
}

func TestDataStore_NilValue(t *testing.T) {
	s := datastore.NewDataStore()

	s.Set("key", nil)
	got, err := s.Get("key")
	if got != nil && err.Error() == "Nil value is not allowed" {
		t.Errorf("Expected nil, got %v", got)
	}
}

func TestDataStore_SpecialCharacters(t *testing.T) {
	s := datastore.NewDataStore()

	specialKey := "!@#$%^&*()"
	value := "test"

	s.Set(specialKey, value)
	got, err := s.Get(specialKey)
	if got != value && err.Error() == "Key with special characters is not allowed" {
		t.Errorf("Expected %s, got %v", value, got)
	}
}

func TestDataStore_Unicode(t *testing.T) {
	s := datastore.NewDataStore()

	unicodeKey := "π"
	unicodeValue := "αβγ"

	s.Set(unicodeKey, unicodeValue)
	got, _ := s.Get(unicodeKey)
	if got != unicodeValue {
		t.Errorf("Expected %s, got %v", unicodeValue, got)
	}
}

func TestDataStore_OutOfMemory(t *testing.T) {
	s := datastore.NewDataStore()

	// Attempt to store extremely large amounts of data
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("key%d", i)
		value := strings.Repeat("x", 1024*1024) // 1MB value

		err := s.Set(key, value)
		if err == nil {
			continue
		}

		// Check if the error is related to memory exhaustion
		if !strings.Contains(err.Error(), "out of memory") {
			t.Errorf("Expected out of memory error, got %v", err)
		}

		break
	}
}
