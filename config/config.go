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
		Bind string `yaml:"bind"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	Persistence struct {
		Path             string `yaml:"path"`
		SnapshotInterval int    `yaml:"snapshot_interval"`
		Enabled          bool   `yaml:"enabled"`
		AppendOnly       bool   `yaml:"append_only"`
	} `yaml:"persistence"`
}

func LoadConfig() *Config {
	configPath := filepath.Join("config", "config.yaml")
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

	v := reflect.ValueOf(config).Elem()
	field := v.FieldByName(section)

	if !field.IsValid() {
		return nil, fmt.Errorf("field %s not found", section)
	}

	return field.Interface(), nil
}
