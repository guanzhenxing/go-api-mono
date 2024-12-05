package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"go-api-mono/pkg/logger"
	"go-api-mono/pkg/middleware"
)

// Config 代表应用程序的完整配置结构
// 它包含了所有组件的配置信息
type Config struct {
	App       AppConfig                  `yaml:"app"`       // 应用基本配置
	Server    ServerConfig               `yaml:"server"`    // 服务器配置
	Log       logger.LogConfig           `yaml:"log"`       // 日志配置
	Database  DatabaseConfig             `yaml:"database"`  // 数据库配置
	Redis     RedisConfig                `yaml:"redis"`     // Redis配置
	RateLimit middleware.RateLimitConfig `yaml:"rateLimit"` // 限流配置
	JWT       JWTConfig                  `yaml:"jwt"`       // JWT配置
}

// AppConfig 代表应用的基本配置信息
type AppConfig struct {
	Name    string `yaml:"name"`    // 应用名称
	Version string `yaml:"version"` // 应用版本
	Mode    string `yaml:"mode"`    // 运行模式（development/production/testing）
}

// ServerConfig 代表HTTP服务器的配置信息
type ServerConfig struct {
	Port            int           `yaml:"port"`            // 服务端口
	Addr            string        `yaml:"addr"`            // 服务地址
	ReadTimeout     time.Duration `yaml:"readTimeout"`     // 读取超时时间
	WriteTimeout    time.Duration `yaml:"writeTimeout"`    // 写入超时时间
	ShutdownTimeout time.Duration `yaml:"shutdownTimeout"` // 关闭超时时间
}

// DatabaseConfig 代表数据库的配置信息
type DatabaseConfig struct {
	Host         string        `yaml:"host"`         // 数据库主机地址
	Port         int           `yaml:"port"`         // 数据库端口
	Username     string        `yaml:"username"`     // 数据库用户名
	Password     string        `yaml:"password"`     // 数据库密码
	Database     string        `yaml:"database"`     // 数据库名称
	MaxOpenConns int           `yaml:"maxOpenConns"` // 最大打开连接数
	MaxIdleConns int           `yaml:"maxIdleConns"` // 最大空闲连接数
	MaxLifetime  time.Duration `yaml:"maxLifetime"`  // 连接最大生命周期
	Debug        bool          `yaml:"debug"`        // 是否开启调试模式
}

// RedisConfig 代表Redis的配置信息
type RedisConfig struct {
	Host     string `yaml:"host"`     // Redis主机地址
	Port     int    `yaml:"port"`     // Redis端口
	Password string `yaml:"password"` // Redis密码
	DB       int    `yaml:"db"`       // Redis数据库索引
}

// JWTConfig 代表JWT的配置信息
type JWTConfig struct {
	SigningKey     string        `yaml:"signingKey"`     // JWT签名密钥
	ExpirationTime time.Duration `yaml:"expirationTime"` // 令牌过期时间
	SigningMethod  string        `yaml:"signingMethod"`  // 签名方法
	TokenPrefix    string        `yaml:"tokenPrefix"`    // 令牌前缀
}

// MustLoad 加载配置文件，如果失败则panic
// 这个函数通常在应用启动时调用，因为配置加载失败应该导致应用程序停止
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

// Load 从指定路径加载配置文件
// 它会读取YAML格式的配置文件并解析到Config结构体中
func Load(path string) (*Config, error) {
	// 读取配置文件内容
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析YAML内容到Config结构体
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}
