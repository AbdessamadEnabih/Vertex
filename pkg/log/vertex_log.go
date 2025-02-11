package vertex_log

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
)

func get_log_path() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("unable to determine caller information")
	}
	projectRoot := filepath.Join(filepath.Dir(filename), "..", "..")

	// Return the logs folder inside the project root
	return filepath.Join(projectRoot, "logs")
}

var (
	logFileName string = "vertex.log"
	logFilePath string
)

func init() {
	// Get log folder path
	logFolderPath := get_log_path()
	logFilePath = filepath.Join(logFolderPath, logFileName)

	// Create logs directory if not exists
	if _, err := os.Stat(logFolderPath); os.IsNotExist(err) {
		err := os.Mkdir(logFolderPath, 0755)
		if err != nil {
			log.Fatalf("Error creating log directory: %v", err)
		}
	}
	// Create log file if not exists
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		os.Create(logFilePath)
		os.Chmod(logFilePath, 0644)
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
