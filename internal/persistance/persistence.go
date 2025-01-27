package persistance

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"

	"github.com/AbdessamadEnabih/Vertex/pkg/config"
	"github.com/AbdessamadEnabih/Vertex/pkg/datastore"
	vertex_log "github.com/AbdessamadEnabih/Vertex/pkg/log"
)



func get_datastore_path() string {
	persistence_config, err := config.GetConfigByField("Persistence")
	if err != nil {
		vertex_log.Log("Error getting persistence config: "+err.Error(), "ERROR")
	}

	dir, _ := os.Getwd()

	return filepath.Join(filepath.Join(dir, reflect.ValueOf(persistence_config).FieldByName("Path").String()), "datastore.data")

}

func Save(datastore *datastore.DataStore) {

	datastorepath := get_datastore_path()
	_, err := os.Stat(datastorepath)
	if err == nil {
		writeInDataStoreFile(datastore, datastorepath)
	} else if os.IsNotExist(err) {
		file, err := os.Create(datastorepath)
		if err != nil {
			vertex_log.Log("Error creating file: "+err.Error(), "ERROR")
			return
		}
		defer file.Close()

		writeInDataStoreFile(datastore, datastorepath)
	} else {
		vertex_log.Log("Error checking file existence: "+err.Error(), "ERROR")	
		return
	}
}

func writeInDataStoreFile(datastore *datastore.DataStore, filepath string) error {
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		vertex_log.Log("Error opening file for writing: "+err.Error(), "ERROR")
		return err
	}
	defer file.Close()

	jsonData, err := json.Marshal(datastore)
	if err != nil {
		vertex_log.Log("Error marshaling datastore to JSON: "+err.Error(), "ERROR")
		return err
	}

	data := []byte(jsonData)

	var compressedBuffer bytes.Buffer

	gzipWriter := gzip.NewWriter(&compressedBuffer)

	_, err = gzipWriter.Write(data)
	if err != nil {
		vertex_log.Log("Error compressing data: "+err.Error(), "ERROR")
		return err
	}

	err = gzipWriter.Close()
	if err != nil {
		vertex_log.Log("Error closing gzip writer: "+err.Error(), "ERROR")
		return err
	}

	encodedData := base64.StdEncoding.EncodeToString(compressedBuffer.Bytes())

	_, err = file.WriteString(encodedData)
	if err != nil {
		vertex_log.Log("Error writing datastore to file: "+err.Error(), "ERROR")
		return err
	}

	return nil
}

func readDataStoreFromFile(filepath string) (*datastore.DataStore, error) {
	encodedData, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err
	}

	compressedData, err := base64.StdEncoding.DecodeString(string(encodedData))
	if err != nil {
		vertex_log.Log("Error decoding base64 data: "+err.Error(), "ERROR")
		return nil, err
	}

	compressedBuffer := bytes.NewBuffer(compressedData)

	gzipReader, err := gzip.NewReader(compressedBuffer)
	if err != nil {
		vertex_log.Log("Error creating gzip reader: "+err.Error(), "ERROR")
		return nil, err
	}
	defer gzipReader.Close()

	// Decompress the gzip data into a byte slice
	var decompressedBuffer bytes.Buffer
	_, err = io.Copy(&decompressedBuffer, gzipReader)
	if err != nil {
		vertex_log.Log("Error decompressing data: "+err.Error(), "ERROR")
		return nil, err
	}

	var savedDataStore datastore.DataStore
	json.Unmarshal(decompressedBuffer.Bytes(), &savedDataStore)

	return &savedDataStore, nil
}

func Load() (*datastore.DataStore, error) {
	datastorepath := get_datastore_path()

	_, err := os.Stat(datastorepath)
	if os.IsNotExist(err) {
		vertex_log.Log("DataStore not found: "+err.Error(), "ERROR")
		return datastore.NewDataStore(), nil
	}
	if err == nil {
		savedDataStore, err := readDataStoreFromFile(datastorepath)
		if err != nil {
			vertex_log.Log("Error reading datastore file: "+err.Error(), "ERROR")
			return datastore.NewDataStore(), err
		}
		return savedDataStore, nil
	} else {
		vertex_log.Log("Error checking file existence: "+err.Error(), "ERROR")
		return datastore.NewDataStore(), err
	}
}
