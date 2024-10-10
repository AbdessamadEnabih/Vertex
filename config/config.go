package config

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Port int    `yaml:"port"`
		Bind string `yaml:"bind"`
	} `yaml:"server"`

	Persistence struct {
		Enabled          bool `yaml:"enabled"`
		SnapshotInterval int  `yaml:"snapshot_interval"`
		AppendOnly       bool `yaml:"append_only"`
	} `yaml:"persistence"`

	State struct {
		MaxKeyAge         string `yaml:"max_key_age"`
		MaxSize           string `yaml:"max_size"`
		MaxAllowedEntries int    `yaml:"max_allowed_entries"`
	} `yaml:"state"`

	Logging struct {
		Level      string `yaml:"level"`
		File       string `yaml:"file"`
		MaxSize    string `yaml:"max_size"`
		MaxBackups int    `yaml:"max_backups"`
		MaxAge     string `yaml:"max_age"`
	} `yaml:"logging"`
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
