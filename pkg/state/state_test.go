package state_test

import (
	"reflect"
	"testing"

	"github.com/AbdessamadEnabih/Vertex/pkg/state"
)

func Test(t *testing.T) {

	var State state.State = *state.NewState()

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
	// Create a new State instance
	state := state.NewState()

	// Call GetAll on empty state
	result := state.GetAll()

	// Check if result is empty
	if len(result) != 0 {
		t.Errorf("Expected empty map, got %d items", len(result))
	}
}
