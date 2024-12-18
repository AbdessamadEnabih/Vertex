package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Logging struct {
		Level      string `yaml:"level"`
		File       string `yaml:"file"`
		MaxSize    string `yaml:"max_size"`
		MaxAge     string `yaml:"max_age"`
		MaxBackups int    `yaml:"max_backups"`
	} `yaml:"logging"`
	State struct {
		MaxKeyAge         string `yaml:"max_key_age"`
		MaxSize           string `yaml:"max_size"`
		MaxAllowedEntries int    `yaml:"max_allowed_entries"`
	} `yaml:"state"`
	Server struct {
		Adress string `yaml:"adress"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	Persistence struct {
		Path             string `yaml:"path"`
		SnapshotInterval int    `yaml:"snapshot_interval"`
		Enabled          bool   `yaml:"enabled"`
		AppendOnly       bool   `yaml:"append_only"`
	} `yaml:"persistence"`
}

func getConfigPath() string {
	// Get environment and config path variables
	vertexEnv := os.Getenv("VERTEX_ENV")
	configPath := os.Getenv("VERTEX_CONFIG_PATH")

	// Use the config path if explicitly set
	if configPath != "" {
		return configPath
	}

	// Default environment to "development" if not set
	if vertexEnv == "" {
		vertexEnv = "development"
	}

	// Determine config path based on environment
	switch vertexEnv {
	case "production":
		configPath = "/etc/vertex/config.yaml"
	case "development":
		// Use the current working directory to form the path
		if cwd, err := os.Getwd(); err == nil {
			configPath = filepath.Join(cwd, "config", "config.yaml")
		} else {
			// Fallback path if unable to get working directory
			configPath = "config/config.yaml"
		}
	default:
		// Default fallback if env is unknown
		configPath = "config/config.yaml"
	}

	return configPath
}

func LoadConfig() *Config {
	configPath := getConfigPath()
	cfg, err := os.ReadFile(configPath)
	if err != nil {
		log.Printf("Error while reading config file %v\n", err)
		return nil
	}
	var c Config
	err = yaml.Unmarshal(cfg, &c)
	if err != nil {
		log.Printf("Error while decoding config file %v\n", err)
		return nil
	}
	return &c
}

func GetConfigByField(section string) (interface{}, error) {
	var config Config = *LoadConfig()

	v := reflect.ValueOf(config)
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("config is not a struct")
	}
	f := v.FieldByName(section)
	if !f.IsValid() {
		return nil, fmt.Errorf("field %s not found", section)
	}
	if f.Kind() != reflect.Struct {
		return nil, fmt.Errorf("field %s is not a struct", section)
	}

	return f.Interface(), nil
}
