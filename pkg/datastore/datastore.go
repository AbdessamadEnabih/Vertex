package datastore

import "github.com/AbdessamadEnabih/Vertex/internal/datastore"

type DataStore struct {
	InternalDataStore *datastore.DataStore
}

func NewDataStore() *DataStore {
	return &DataStore{
		InternalDataStore: datastore.NewDataStore(),
	}
}

func (s *DataStore) Set(key string, value interface{}) error {
	return s.InternalDataStore.Set(key, value)
}

func (s *DataStore) Get(key string) (interface{}, error) {
	return s.InternalDataStore.Get(key)
}

func (s *DataStore) Delete(key string) error {
	return s.InternalDataStore.Delete(key)
}

func (s *DataStore) GetAll() map[string]interface{} {
	return s.InternalDataStore.GetAll()
}

func (s *DataStore) Update(key string, value interface{}) error {
	return s.InternalDataStore.Update(key, value)
}

func (s *DataStore) FlushAll() error {
	return s.InternalDataStore.FlushAll()
}
