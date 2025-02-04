package persistence_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/AbdessamadEnabih/Vertex/internal/persistence"
	"github.com/AbdessamadEnabih/Vertex/pkg/datastore"
)

func setup() {
	// Set up any necessary configuration or environment variables
	os.MkdirAll("testdata", os.ModePerm)
}

func teardown() {
	// Clean up any resources or files created during tests
	os.RemoveAll("testdata")
}

func TestWriteInDataStoreFile(t *testing.T) {
	setup()
	defer teardown()

	ds := datastore.NewDataStore()
	ds.Set("key1", "test value")

	datastorePath := filepath.Join("testdata", "datastore.data")
	err := persistence.WriteInDataStoreFile(ds, datastorePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if _, err := os.Stat(datastorePath); os.IsNotExist(err) {
		t.Fatalf("Expected datastore file to exist, but it does not")
	}
}

func TestReadDataStoreFromFile(t *testing.T) {
	setup()
	defer teardown()

	var key, val = "key1", "test value"

	ds := datastore.NewDataStore()
	ds.Set(key, val)

	datastorePath := filepath.Join("testdata", "datastore.data")
	err := persistence.WriteInDataStoreFile(ds, datastorePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	loadedDS, err := persistence.ReadDataStoreFromFile(datastorePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if value, _ := loadedDS.Get(key); value != val {
		t.Fatalf("Expected %v, got %v", val, value)
	}
}
