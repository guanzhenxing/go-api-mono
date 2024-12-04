package database

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	// ErrRecordNotFound 记录未找到错误
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

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

// DB 封装 gorm.DB
type DB struct {
	*gorm.DB
}

// New 创建数据库连接
func New(config DatabaseConfig) (*DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.MaxLifetime)

	return &DB{DB: db}, nil
}

// Close 关闭数据库连接
func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}
	return sqlDB.Close()
}

// AutoMigrate 自动迁移数据库表
func (db *DB) AutoMigrate(models ...interface{}) error {
	return db.DB.AutoMigrate(models...)
}

// Transaction 执行事务
func (db *DB) Transaction(fn func(tx *gorm.DB) error) error {
	return db.DB.Transaction(fn)
}
