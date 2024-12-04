package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"go-api-mono/internal/pkg/logger"
	"go-api-mono/internal/pkg/middleware"
)

// Config represents the application configuration
type Config struct {
	App       AppConfig                  `yaml:"app"`
	Server    ServerConfig               `yaml:"server"`
	Log       logger.LogConfig           `yaml:"log"`
	Database  DatabaseConfig             `yaml:"database"`
	RateLimit middleware.RateLimitConfig `yaml:"rateLimit"`
}

// AppConfig represents the application configuration
type AppConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Mode    string `yaml:"mode"`
}

// ServerConfig represents the server configuration
type ServerConfig struct {
	Port            int           `yaml:"port"`
	Addr            string        `yaml:"addr"`
	ReadTimeout     time.Duration `yaml:"readTimeout"`
	WriteTimeout    time.Duration `yaml:"writeTimeout"`
	ShutdownTimeout time.Duration `yaml:"shutdownTimeout"`
}

// DatabaseConfig represents the database configuration
type DatabaseConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	Username     string        `yaml:"username"`
	Password     string        `yaml:"password"`
	Database     string        `yaml:"database"`
	MaxOpenConns int           `yaml:"maxOpenConns"`
	MaxIdleConns int           `yaml:"maxIdleConns"`
	MaxLifetime  time.Duration `yaml:"maxLifetime"`
	Debug        bool          `yaml:"debug"`
}

// MustLoad loads the configuration from file
func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}

	config, err := Load(configPath)
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	return config
}

// Load loads the configuration from file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}
