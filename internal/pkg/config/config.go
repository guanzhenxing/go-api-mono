package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"go-api-mono/configs"

	"gopkg.in/yaml.v3"
)

// Config 应用程序配置
type Config struct {
	App       AppConfig       `yaml:"app"`
	Server    ServerConfig    `yaml:"server"`
	Log       LogConfig       `yaml:"log"`
	Database  DatabaseConfig  `yaml:"database"`
	JWT       JWTConfig       `yaml:"jwt"`
	RateLimit RateLimitConfig `yaml:"rateLimit"`
}

// AppConfig 应用程序基本配置
type AppConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Mode    string `yaml:"mode"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port            int           `yaml:"port"`
	ReadTimeout     time.Duration `yaml:"readTimeout"`
	WriteTimeout    time.Duration `yaml:"writeTimeout"`
	ShutdownTimeout time.Duration `yaml:"shutdownTimeout"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `yaml:"level"`
	Filename   string `yaml:"filename"`
	MaxSize    int    `yaml:"maxSize"`
	MaxBackups int    `yaml:"maxBackups"`
	MaxAge     int    `yaml:"maxAge"`
	Compress   bool   `yaml:"compress"`
}

// DatabaseConfig 数据库配置
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

// JWTConfig JWT配置
type JWTConfig struct {
	SigningKey     string        `yaml:"signingKey"`
	ExpirationTime time.Duration `yaml:"expirationTime"`
	SigningMethod  string        `yaml:"signingMethod"`
	TokenPrefix    string        `yaml:"tokenPrefix"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Requests int `yaml:"requests"`
	Burst    int `yaml:"burst"`
}

// Load 加载配置
func Load() (*Config, error) {
	// 获取环境
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	// 读取配置文件
	data, err := configs.GetConfigFile(env)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析配置
	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// 设置默认值
	setDefaults(config)

	// 验证配置
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

// MustLoad 加载配置，如果出错则panic
func MustLoad() *Config {
	config, err := Load()
	if err != nil {
		panic(err)
	}
	return config
}

// setDefaults 设置默认配置值
func setDefaults(config *Config) {
	if config.App.Name == "" {
		config.App.Name = "go-api-mono"
	}
	if config.App.Version == "" {
		config.App.Version = "v0.1.0"
	}
	if config.App.Mode == "" {
		config.App.Mode = "development"
	}

	if config.Server.Port == 0 {
		config.Server.Port = 8080
	}
	if config.Server.ReadTimeout == 0 {
		config.Server.ReadTimeout = 10 * time.Second
	}
	if config.Server.WriteTimeout == 0 {
		config.Server.WriteTimeout = 10 * time.Second
	}
	if config.Server.ShutdownTimeout == 0 {
		config.Server.ShutdownTimeout = 30 * time.Second
	}

	if config.Log.Level == "" {
		config.Log.Level = "debug"
	}
	if config.Log.Filename == "" {
		config.Log.Filename = "logs/app.log"
	}

	if config.JWT.SigningMethod == "" {
		config.JWT.SigningMethod = "HS256"
	}
	if config.JWT.TokenPrefix == "" {
		config.JWT.TokenPrefix = "Bearer"
	}
	if config.JWT.ExpirationTime == 0 {
		config.JWT.ExpirationTime = 24 * time.Hour
	}

	if config.RateLimit.Requests == 0 {
		config.RateLimit.Requests = 100
	}
	if config.RateLimit.Burst == 0 {
		config.RateLimit.Burst = 200
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	// 应用配置验证
	if err := c.validateApp(); err != nil {
		return fmt.Errorf("app config validation failed: %w", err)
	}

	// 服务器配置验证
	if err := c.validateServer(); err != nil {
		return fmt.Errorf("server config validation failed: %w", err)
	}

	// 数据库配置验证
	if err := c.validateDatabase(); err != nil {
		return fmt.Errorf("database config validation failed: %w", err)
	}

	// JWT配置验证
	if err := c.validateJWT(); err != nil {
		return fmt.Errorf("jwt config validation failed: %w", err)
	}

	// 速率限制配置验证
	if err := c.validateRateLimit(); err != nil {
		return fmt.Errorf("rate limit config validation failed: %w", err)
	}

	return nil
}

func (c *Config) validateApp() error {
	if c.App.Name == "" {
		return errors.New("app name is required")
	}
	if c.App.Version == "" {
		return errors.New("app version is required")
	}
	if c.App.Mode == "" {
		return errors.New("app mode is required")
	}
	if c.App.Mode != "development" && c.App.Mode != "production" && c.App.Mode != "testing" {
		return errors.New("app mode must be one of: development, production, testing")
	}
	return nil
}

func (c *Config) validateServer() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	if c.Server.ReadTimeout <= 0 {
		return errors.New("server read timeout must be positive")
	}
	if c.Server.WriteTimeout <= 0 {
		return errors.New("server write timeout must be positive")
	}
	if c.Server.ShutdownTimeout <= 0 {
		return errors.New("server shutdown timeout must be positive")
	}
	return nil
}

func (c *Config) validateDatabase() error {
	if c.Database.Host == "" {
		return errors.New("database host is required")
	}
	if c.Database.Port <= 0 || c.Database.Port > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}
	if c.Database.Username == "" {
		return errors.New("database username is required")
	}
	if c.Database.Database == "" {
		return errors.New("database name is required")
	}
	if c.Database.MaxOpenConns <= 0 {
		return errors.New("database max open connections must be positive")
	}
	if c.Database.MaxIdleConns <= 0 {
		return errors.New("database max idle connections must be positive")
	}
	if c.Database.MaxIdleConns > c.Database.MaxOpenConns {
		return errors.New("database max idle connections cannot be greater than max open connections")
	}
	if c.Database.MaxLifetime <= 0 {
		return errors.New("database connection max lifetime must be positive")
	}
	return nil
}

func (c *Config) validateJWT() error {
	if c.JWT.SigningKey == "" {
		return errors.New("jwt signing key is required")
	}
	if len(c.JWT.SigningKey) < 32 {
		return errors.New("jwt signing key must be at least 32 characters")
	}
	if c.JWT.ExpirationTime <= 0 {
		return errors.New("jwt expiration time must be positive")
	}
	if c.JWT.SigningMethod == "" {
		return errors.New("jwt signing method is required")
	}
	if c.JWT.SigningMethod != "HS256" && c.JWT.SigningMethod != "HS384" && c.JWT.SigningMethod != "HS512" {
		return errors.New("jwt signing method must be one of: HS256, HS384, HS512")
	}
	if c.JWT.TokenPrefix == "" {
		return errors.New("jwt token prefix is required")
	}
	return nil
}

func (c *Config) validateRateLimit() error {
	if c.RateLimit.Requests <= 0 {
		return errors.New("rate limit requests must be positive")
	}
	if c.RateLimit.Burst <= 0 {
		return errors.New("rate limit burst must be positive")
	}
	if c.RateLimit.Burst < c.RateLimit.Requests {
		return errors.New("rate limit burst must be greater than or equal to requests")
	}
	return nil
}
