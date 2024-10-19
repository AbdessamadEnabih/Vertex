package persistance

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"github.com/AbdessamadEnabih/Vertex/config"
	"github.com/AbdessamadEnabih/Vertex/pkg/state"
)

func Save(state *state.State) {
	persistence_config, err := config.GetConfigByField("Persistence")
	if err != nil {
		log.Printf("Error while loading Persistence configuration : %s", err)
	}

	dir, _ := os.Getwd()

	datapath := filepath.Join(dir, reflect.ValueOf(persistence_config).FieldByName("Path").String())
	statepath := filepath.Join(datapath, "state.data")
	_, err = os.Stat(statepath)
	if err == nil {
		writeInStateFile(state, statepath)
	} else if os.IsNotExist(err) {
		file, err := os.Create(statepath)
		if err != nil {
			log.Printf("Error creating file: %s", err)
			return
		}
		defer file.Close()

		writeInStateFile(state, statepath)
	} else {
		log.Printf("Error checking file existence: %s", err)
		return
	}
}

func writeInStateFile(state *state.State, filepath string) error {
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		log.Printf("Error opening file for writing: %s", err)
		return err
	}
	defer file.Close()

	jsonData, err := json.Marshal(state.GetAll())
	if err != nil {
		log.Printf("Error marshaling state to JSON: %s", err)
		return err
	}

	data := []byte(jsonData)

	var compressedBuffer bytes.Buffer

	gzipWriter := gzip.NewWriter(&compressedBuffer)

	_, err = gzipWriter.Write(data)
	if err != nil {
		log.Printf("Error compressing data: %s", err)
		return err
	}

	err = gzipWriter.Close()
	if err != nil {
		log.Printf("Error closing gzip writer: %s", err)
		return err
	}

	encodedData := base64.StdEncoding.EncodeToString(compressedBuffer.Bytes())

	_, err = file.WriteString(encodedData)
	if err != nil {
		log.Printf("Error writing state to file: %s", err)
		return err
	}

	return nil
}

func readStateFromFile(filepath string) {
	encodedData, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	compressedData, err := base64.StdEncoding.DecodeString(string(encodedData))
	if err != nil {
		log.Printf("Error decoding base64 data: %s", err)
		return
	}

	compressedBuffer := bytes.NewBuffer(compressedData)

	gzipReader, err := gzip.NewReader(compressedBuffer)
	if err != nil {
		log.Printf("Error creating gzip reader: %s", err)
		return
	}
	defer gzipReader.Close()

	// Decompress the gzip data into a byte slice
	var decompressedBuffer bytes.Buffer
	_, err = io.Copy(&decompressedBuffer, gzipReader)
	if err != nil {
		log.Printf("Error decompressing data: %s", err)
		return
	}

	fmt.Print(string(decompressedBuffer.Bytes()))
}

func Load() {
	persistence_config, err := config.GetConfigByField("Persistence")
	if err != nil {
		log.Printf("Error while loading Persistence configuration : %s", err)
	}

	dir, _ := os.Getwd()

	datapath := filepath.Join(dir, reflect.ValueOf(persistence_config).FieldByName("Path").String())
	statepath := filepath.Join(datapath, "state.data")

	_, err = os.Stat(statepath)
	if os.IsNotExist(err) {
		log.Printf("State not found: %s", err)
	}
	if err == nil {

		readStateFromFile(statepath)

	} else {
		log.Printf("Error checking file existence: %s", err)
		return
	}

}
