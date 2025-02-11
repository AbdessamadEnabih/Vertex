package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Logging struct {
		Level      string `yaml:"level"`
		Path       string `yaml:"path"`
		MaxSize    string `yaml:"max_size"`
		MaxAge     string `yaml:"max_age"`
		MaxBackups int    `yaml:"max_backups"`
	} `yaml:"logging"`
	Store struct {
		MaxKeyAge         string `yaml:"max_key_age"`
		MaxSize           string `yaml:"max_size"`
		MaxAllowedEntries int    `yaml:"max_allowed_entries"`
	} `yaml:"state"`
	Server struct {
		Adress string `yaml:"adress"`
		Port   int    `yaml:"port"`
		SSL    bool   `yaml:"ssl"`
	} `yaml:"server"`
	Persistence struct {
		Path             string `yaml:"path"`
		SnapshotInterval int    `yaml:"snapshot_interval"`
		Enabled          bool   `yaml:"enabled"`
		AppendOnly       bool   `yaml:"append_only"`
	} `yaml:"persistence"`
}

func getConfigPath() string {
    // Use the config path if explicitly set
    if configPath := os.Getenv("VERTEX_CONFIG_PATH"); configPath != "" {
        return configPath
    }

    // Default environment to "development" if not set
    vertexEnv := os.Getenv("VERTEX_ENV")
    if vertexEnv == "" {
        vertexEnv = "development"
    }

    switch vertexEnv {
    case "production":
        return "/etc/vertex/config.yaml"
    case "development":
        _, filename, _, ok := runtime.Caller(0)
        if !ok {
            panic("unable to determine caller information")
        }
        // Assuming this file is in [project_root]/pkg/config, move up two directories.
        projectRoot := filepath.Join(filepath.Dir(filename), "..", "..")
        return filepath.Join(projectRoot, "configs", "config.yaml")
    default:
        return "configs/config.yaml"
    }
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
