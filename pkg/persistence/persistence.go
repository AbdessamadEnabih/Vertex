package persistence

import (
	"errors"

	"github.com/AbdessamadEnabih/Vertex/internal/persistence"
	"github.com/AbdessamadEnabih/Vertex/pkg/datastore"
	vertex_log "github.com/AbdessamadEnabih/Vertex/pkg/log"
)

func Save(datastore *datastore.DataStore) error {
	if err := persistence.Save(datastore); err != nil {
		vertex_log.Log("Error saving datastore: "+err.Error(), "Error")
		return err
	}
	return nil
}

func Load() (*datastore.DataStore, error) {
	DataStore, err := persistence.Load()
	if err != nil {
		vertex_log.Log("Error loading datastore: "+err.Error(), "Error")
		return nil, err
	}

	if DataStore == nil {
		vertex_log.Log("Error: DataStore is nil", "Error")
		return nil, errors.New("datastore is nil")
	}

	return DataStore, nil
}
