package state_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/AbdessamadEnabih/Vertex/pkg/state"
)

var State state.State = *state.NewState()

func Test(t *testing.T) {

	// Check if the type of State.data is map[string]interface{}
	if reflect.TypeOf(State.Data) != reflect.TypeOf(map[string]interface{}{}) {
		t.Errorf("Expected State.Data to be of type map[string]interface{}, but got %T", State.Data)
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
		State.Set(tt.key, tt.expected)

		if got := State.Get(tt.key); !reflect.DeepEqual(got, tt.expected) {
			t.Errorf("Expected %s to be %v, but got %v", tt.key, tt.expected, got)
		}
	}

	// Test for an invalid key
	invalidKey := "invalidKey"
	if got := State.Get(invalidKey); got != nil {
		t.Errorf("Expected nil for key %s, but got %v", invalidKey, got)
	}

	// Call GetAll and compare results
	result := State.GetAll()

	// Check if all keys exist in the result
	if len(result) != 2 {
		t.Errorf("Expected 2 items, got %d", len(result))
	}

}

func TestState_GetAll_Empty(t *testing.T) {
	s := state.NewState()

	// Call GetAll on empty state
	result := s.GetAll()

	// Check if result is empty
	if len(result) != 0 {
		t.Errorf("Expected empty map, got %d items", len(result))
	}
}

func TestState_Delete(t *testing.T) {
	// Add some test data
	State.Set("testKey1", "value1")
	State.Set("testKey2", 42)

	// Test deleting a non-existent key
	deleted := State.Delete("nonExistentKey")
	if deleted {
		t.Errorf("Expected false when deleting non-existent key, but got true")
	}

	// Test deleting an existing key
	deleted = State.Delete("testKey1")
	if !deleted {
		t.Errorf("Expected true when deleting existing key, but got false")
	}

	// Verify that the key was removed
	got := State.Get("testKey1")
	if got != nil {
		t.Errorf("Expected nil for deleted key, but got %v", got)
	}

	// Test deleting another existing key
	deleted = State.Delete("testKey2")
	if !deleted {
		t.Errorf("Expected true when deleting existing key, but got false")
	}

	// Verify that both keys were removed
	got = State.Get("testKey2")
	if got != nil {
		t.Errorf("Expected nil for deleted key, but got %v", got)
	}

	// Test deleting again (should return false)
	deleted = State.Delete("testKey2")
	if deleted {
		t.Errorf("Expected false when deleting already deleted key, but got true")
	}
}

func TestState_Update(t *testing.T) {
	// Create a new State instance for this test
	var updateState state.State = *state.NewState()

	updateState.Set("key1", 22)

	value := "updatedkey"
	updateState.Update("key1", value)
	if got := updateState.Get("key1"); got != value {
		t.Fatalf("Expected value to be %v and got %v", value, got)
	}

	// Test update of non existent key
	expected := updateState.Update("invalidkey", value)
	if expected != false {
		t.Errorf("Expected false when deleting non-existent key, got true")
	}

	// Clean up
	updateState.Delete("key1")
	updateState.Delete("invalidkey")
}

func TestState_LargeKey(t *testing.T) {
	s := state.NewState()
	largeKey := strings.Repeat("x", 1024*1024) // 1MB key
	value := "test"

	s.Set(largeKey, value)
	got := s.Get(largeKey)
	if got != value {
		t.Errorf("Expected %s, got %v", value, got)
	}
}

func TestState_LargeValue(t *testing.T) {
	s := state.NewState()
	key := "test"
	largeValue := strings.Repeat("x", 1024*1024) // 1MB value

	s.Set(key, largeValue)
	got := s.Get(key)
	if got != largeValue {
		t.Errorf("Expected %s, got %v", largeValue, got)
	}
}

func TestState_EmptyKey(t *testing.T) {
	s := state.NewState()

	err := s.Set("", "test")
	if err.Error() != "empty key is not allowed" {
		t.Errorf("Expected Set to return Error %s", err.Error())
	}
}

func TestState_EmptyValue(t *testing.T) {
	s := state.NewState()

	s.Set("key", "")
	got := s.Get("key")
	if got != "" {
		t.Errorf("Expected '', got %v", got)
	}
}

func TestState_NilValue(t *testing.T) {
	s := state.NewState()

	s.Set("key", nil)
	got := s.Get("key")
	if got != nil {
		t.Errorf("Expected nil, got %v", got)
	}
}

func TestState_SpecialCharacters(t *testing.T) {
	s := state.NewState()

	specialKey := "!@#$%^&*()"
	value := "test"

	s.Set(specialKey, value)
	got := s.Get(specialKey)
	if got != value {
		t.Errorf("Expected %s, got %v", value, got)
	}
}

func TestState_Unicode(t *testing.T) {
	s := state.NewState()

	unicodeKey := "π"
	unicodeValue := "αβγ"

	s.Set(unicodeKey, unicodeValue)
	got := s.Get(unicodeKey)
	if got != unicodeValue {
		t.Errorf("Expected %s, got %v", unicodeValue, got)
	}
}

func TestState_OutOfMemory(t *testing.T) {
	s := state.NewState()

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
