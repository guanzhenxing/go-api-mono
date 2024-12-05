package database

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB defines the interface for database operations
type DB interface {
	WithContext(ctx context.Context) DB
	Create(value interface{}) error
	Save(value interface{}) error
	Delete(value interface{}, where ...interface{}) error
	First(dest interface{}, conds ...interface{}) error
	Model(value interface{}) DB
	Updates(values interface{}) error
	Debug() DB
	GetDB() *gorm.DB
}

// DatabaseConfig defines the configuration for database connection
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

type db struct {
	*gorm.DB
}

// New creates a new database connection
func New(config DatabaseConfig) (DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.MaxLifetime)

	if config.Debug {
		gormDB = gormDB.Debug()
	}

	return &db{gormDB}, nil
}

func (d *db) WithContext(ctx context.Context) DB {
	return &db{d.DB.WithContext(ctx)}
}

func (d *db) Create(value interface{}) error {
	return d.DB.Create(value).Error
}

func (d *db) Save(value interface{}) error {
	return d.DB.Save(value).Error
}

func (d *db) Delete(value interface{}, where ...interface{}) error {
	return d.DB.Delete(value, where...).Error
}

func (d *db) First(dest interface{}, conds ...interface{}) error {
	return d.DB.First(dest, conds...).Error
}

func (d *db) Model(value interface{}) DB {
	return &db{DB: d.DB.Model(value)}
}

func (d *db) Updates(values interface{}) error {
	return d.DB.Updates(values).Error
}

func (d *db) Debug() DB {
	return &db{DB: d.DB.Debug()}
}

func (d *db) GetDB() *gorm.DB {
	return d.DB
}
