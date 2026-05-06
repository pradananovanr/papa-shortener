package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	App      AppConfig      `yaml:"app"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	SSLMode  string `yaml:"sslmode"`
}

type AppConfig struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	BaseURL string `yaml:"base_url"`
}

var AppCfg *Config

func Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	AppCfg = &Config{}
	return yaml.Unmarshal(data, AppCfg)
}