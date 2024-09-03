package state_test

import (
	"reflect"
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
	// Call GetAll on empty state
	result := State.GetAll()

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
