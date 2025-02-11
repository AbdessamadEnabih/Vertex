package persistence

import (
    "bytes"
    "compress/gzip"
    "encoding/json"
    "io"
    "os"
    "path/filepath"
    "reflect"

    "github.com/AbdessamadEnabih/Vertex/pkg/config"
    "github.com/AbdessamadEnabih/Vertex/pkg/datastore"
    logger "github.com/AbdessamadEnabih/Vertex/pkg/logger"
)

func get_datastore_path() string {
    persistence_config, err := config.GetConfigByField("Persistence")
    if err != nil {
        logError("persistence.get_datastore_path: Error getting persistence config", "", err)
    }

    dir, _ := os.Getwd()

    return filepath.Join(filepath.Join(dir, reflect.ValueOf(persistence_config).FieldByName("Path").String()), "datastore.data")
}

func Save(datastore *datastore.DataStore) error {
    datastorepath := get_datastore_path()

    if err := WriteInDataStoreFile(datastore, datastorepath); err != nil {
        logError("persistence.Save: Error saving datastore", datastorepath, err)
        return err
    }
    return nil
}

func WriteInDataStoreFile(datastore *datastore.DataStore, filepath string) error {
    file, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
    if err != nil {
        logError("persistence.WriteInDataStoreFile: Error opening file for writing", filepath, err)
        return err
    }
    defer file.Close()

    jsonData, err := json.Marshal(datastore)
    if err != nil {
        logError("persistence.WriteInDataStoreFile: Error marshaling datastore to JSON", filepath, err)
        return err
    }

    var compressedBuffer bytes.Buffer
    gzipWriter := gzip.NewWriter(&compressedBuffer)
    _, err = gzipWriter.Write(jsonData)
    if err != nil {
        logError("persistence.WriteInDataStoreFile: Error compressing data", filepath, err)
        return err
    }
    gzipWriter.Close()
    data := compressedBuffer.Bytes()

    _, err = file.Write(data)
    if err != nil {
        logError("persistence.WriteInDataStoreFile: Error writing compressed data to file", filepath, err)
        return err
    }

    return nil
}

func ReadDataStoreFromFile(filepath string) (*datastore.DataStore, error) {
    compressedData, err := os.ReadFile(filepath)
    if err != nil {
        logError("persistence.ReadDataStoreFromFile: Error reading file", filepath, err)
        return nil, err
    }

    compressedBuffer := bytes.NewBuffer(compressedData)
    gzipReader, err := gzip.NewReader(compressedBuffer)
    if err != nil {
        logError("persistence.ReadDataStoreFromFile: Error creating gzip reader", filepath, err)
        return nil, err
    }
    defer gzipReader.Close()

    var decompressedBuffer bytes.Buffer
    _, err = io.Copy(&decompressedBuffer, gzipReader)
    if err != nil {
        logError("persistence.ReadDataStoreFromFile: Error decompressing data", filepath, err)
        return nil, err
    }

    var savedDataStore datastore.DataStore
    err = json.Unmarshal(decompressedBuffer.Bytes(), &savedDataStore)
    if err != nil {
        logError("persistence.ReadDataStoreFromFile: Error unmarshaling JSON data", filepath, err)
        return nil, err
    }

    return &savedDataStore, nil
}

func Load() (*datastore.DataStore, error) {
    datastorepath := get_datastore_path()

    _, err := os.Stat(datastorepath)
    if os.IsNotExist(err) {
        logError("persistence.Load: DataStore not found", datastorepath, err)
        return datastore.NewDataStore(), nil
    }
    if err == nil {
        savedDataStore, err := ReadDataStoreFromFile(datastorepath)
        if err != nil {
            logError("persistence.Load: Error reading datastore file", datastorepath, err)
            return datastore.NewDataStore(), err
        }
        return savedDataStore, nil
    } else {
        logError("persistence.Load: Error checking file existence", datastorepath, err)
        return datastore.NewDataStore(), err
    }
}

func logError(message, filepath string, err error) {
    if filepath != "" {
        logger.Log(message+" at path "+filepath+": "+err.Error(), "ERROR")
    } else {
        logger.Log(message+": "+err.Error(), "ERROR")
    }
}
