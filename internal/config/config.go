package config

import (
	"os"
	"strconv"

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
	Host             string `yaml:"host"`
	Port             int    `yaml:"port"`
	BaseURL          string `yaml:"base_url"`
	DefaultAdminUser string `yaml:"default_admin_user"`
	DefaultAdminPass string `yaml:"default_admin_pass"`
}

var AppCfg *Config

func Load(path string) error {
	AppCfg = &Config{}

	// 1. Try to load from yaml file if it exists, but don't error out if it doesn't
	if data, err := os.ReadFile(path); err == nil {
		_ = yaml.Unmarshal(data, AppCfg)
	}

	// 2. Override or fill values from environment variables
	if host := os.Getenv("DB_HOST"); host != "" {
		AppCfg.Database.Host = host
	}
	if portStr := os.Getenv("DB_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			AppCfg.Database.Port = port
		}
	}
	if user := os.Getenv("DB_USER"); user != "" {
		AppCfg.Database.User = user
	}
	if pass := os.Getenv("DB_PASSWORD"); pass != "" {
		AppCfg.Database.Password = pass
	}
	if name := os.Getenv("DB_NAME"); name != "" {
		AppCfg.Database.Name = name
	}
	if sslmode := os.Getenv("DB_SSLMODE"); sslmode != "" {
		AppCfg.Database.SSLMode = sslmode
	}

	if host := os.Getenv("APP_HOST"); host != "" {
		AppCfg.App.Host = host
	}
	if portStr := os.Getenv("APP_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			AppCfg.App.Port = port
		}
	}
	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		AppCfg.App.BaseURL = baseURL
	}
	if adminUser := os.Getenv("DEFAULT_ADMIN_USER"); adminUser != "" {
		AppCfg.App.DefaultAdminUser = adminUser
	}
	if adminPass := os.Getenv("DEFAULT_ADMIN_PASS"); adminPass != "" {
		AppCfg.App.DefaultAdminPass = adminPass
	}

	// 3. Fallback defaults if they are still empty
	if AppCfg.Database.Host == "" {
		AppCfg.Database.Host = "localhost"
	}
	if AppCfg.Database.Port == 0 {
		AppCfg.Database.Port = 5432
	}
	if AppCfg.Database.User == "" {
		AppCfg.Database.User = "postgres"
	}
	if AppCfg.Database.Password == "" {
		AppCfg.Database.Password = "postgres"
	}
	if AppCfg.Database.Name == "" {
		AppCfg.Database.Name = "urlshortener"
	}
	if AppCfg.Database.SSLMode == "" {
		AppCfg.Database.SSLMode = "disable"
	}

	if AppCfg.App.Host == "" {
		AppCfg.App.Host = "0.0.0.0"
	}
	if AppCfg.App.Port == 0 {
		AppCfg.App.Port = 8080
	}

	return nil
}