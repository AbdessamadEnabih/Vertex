package log

import (
	"log"
	"os"
	"path/filepath"
	"reflect"

	"github.com/AbdessamadEnabih/Vertex/pkg/config"
)

func get_log_path() string {
	logging_config, err := config.GetConfigByField("Logging")
	if err != nil {
		Log("Error getting logging config: "+err.Error(), "ERROR")
	}

	return reflect.ValueOf(logging_config).FieldByName("Path").String()
}

var (
	logFileName string = "vertex.log"
	logFilePath string
)

func init() {
	// Get log file path
	logFilePath = filepath.Join(get_log_path(), logFileName)

	// Create logs directory if not exists
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", 0755)
	}
	// Create log file if not exists
	if _, err := os.Stat(logFileName); os.IsNotExist(err) {
		os.Create(logFileName)
		os.Chmod(logFileName, 0644)
	}
}

func Log(message string, level string) {
	f, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("Error opening log file: %v", err)
		return
	}
	defer f.Close()

	log.SetOutput(f)
	logger(message, level)
}

func logger(message string, level string) {
	switch level {
	case "INFO":
		log.SetPrefix("INFO: ")
	case "ERROR":
		log.SetPrefix("ERROR: ")
	case "DEBUG":
		log.SetPrefix("DEBUG: ")
	default:
		log.SetPrefix("INFO: ")
	}
	log.Println(message)
}
