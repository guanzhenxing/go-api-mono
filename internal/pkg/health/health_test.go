package health

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"go-api-mono/internal/pkg/database"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) WithContext(ctx context.Context) database.DB {
	return m
}

func (m *MockDB) Create(value interface{}) error {
	args := m.Called(value)
	return args.Error(0)
}

func (m *MockDB) Save(value interface{}) error {
	args := m.Called(value)
	return args.Error(0)
}

func (m *MockDB) Delete(value interface{}, where ...interface{}) error {
	args := m.Called(value, where[0])
	return args.Error(0)
}

func (m *MockDB) First(dest interface{}, conds ...interface{}) error {
	args := m.Called(dest, conds[0])
	return args.Error(0)
}

func (m *MockDB) GetDB() *gorm.DB {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*gorm.DB)
}

type MockCache struct {
	mock.Mock
}

func (m *MockCache) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockCache) Del(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockCache) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCache) Close() error {
	args := m.Called()
	return args.Error(0)
}

type mockSQLDB struct {
	mock.Mock
}

func (m *mockSQLDB) PingContext(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestHealthChecker(t *testing.T) {
	version := "1.0.0"

	t.Run("All Components Healthy", func(t *testing.T) {
		ctx := context.Background()

		mockDB := &MockDB{}
		mockCache := &MockCache{}

		// Setup mock DB with working SQL connection
		sqlDB := &sql.DB{}
		mockGormDB := &gorm.DB{
			Config: &gorm.Config{
				ConnPool: sqlDB,
			},
		}
		mockDB.On("GetDB").Return(mockGormDB)

		// Setup mock cache
		mockCache.On("Ping", ctx).Return(nil)

		checker := NewChecker(mockDB, mockCache, version)
		response := checker.Check(ctx)

		assert.Equal(t, "healthy", response.Status)
		assert.Equal(t, version, response.Version)
		assert.Len(t, response.Components, 3) // database, redis, system

		// Verify component statuses
		for _, component := range response.Components {
			assert.Equal(t, "healthy", component.Status)
		}

		mockCache.AssertExpectations(t)
		mockDB.AssertExpectations(t)
	})

	t.Run("Redis Unhealthy", func(t *testing.T) {
		ctx := context.Background()

		mockDB := &MockDB{}
		mockCache := &MockCache{}

		// Setup mock DB with working SQL connection
		sqlDB := &sql.DB{}
		mockGormDB := &gorm.DB{
			Config: &gorm.Config{
				ConnPool: sqlDB,
			},
		}
		mockDB.On("GetDB").Return(mockGormDB)

		// Setup mock cache with error
		mockCache.On("Ping", ctx).Return(errors.New("redis connection error"))

		checker := NewChecker(mockDB, mockCache, version)
		response := checker.Check(ctx)

		assert.Equal(t, "unhealthy", response.Status)

		// Find Redis component
		var redisComponent *Status
		for _, component := range response.Components {
			if component.Component == "redis" {
				redisComponent = &component
				break
			}
		}

		assert.NotNil(t, redisComponent)
		assert.Equal(t, "unhealthy", redisComponent.Status)
		assert.Contains(t, redisComponent.Message, "Redis ping failed")

		mockCache.AssertExpectations(t)
		mockDB.AssertExpectations(t)
	})

	t.Run("Database Unhealthy - Nil Connection", func(t *testing.T) {
		ctx := context.Background()

		mockDB := &MockDB{}
		mockCache := &MockCache{}

		// Setup mock DB to return nil
		mockDB.On("GetDB").Return(nil)

		// Setup mock cache
		mockCache.On("Ping", ctx).Return(nil)

		checker := NewChecker(mockDB, mockCache, version)
		response := checker.Check(ctx)

		assert.Equal(t, "unhealthy", response.Status)

		// Find database component
		var dbComponent *Status
		for _, component := range response.Components {
			if component.Component == "database" {
				dbComponent = &component
				break
			}
		}

		assert.NotNil(t, dbComponent)
		assert.Equal(t, "unhealthy", dbComponent.Status)
		assert.Contains(t, dbComponent.Message, "Database connection is nil")

		mockCache.AssertExpectations(t)
		mockDB.AssertExpectations(t)
	})

	t.Run("Database Unhealthy - DB Error", func(t *testing.T) {
		ctx := context.Background()

		mockDB := &MockDB{}
		mockCache := &MockCache{}

		// Setup mock DB with error
		mockGormDB := &gorm.DB{
			Config: &gorm.Config{},
		}
		mockDB.On("GetDB").Return(mockGormDB)

		// Setup mock cache
		mockCache.On("Ping", ctx).Return(nil)

		checker := NewChecker(mockDB, mockCache, version)
		response := checker.Check(ctx)

		assert.Equal(t, "unhealthy", response.Status)

		// Find database component
		var dbComponent *Status
		for _, component := range response.Components {
			if component.Component == "database" {
				dbComponent = &component
				break
			}
		}

		assert.NotNil(t, dbComponent)
		assert.Equal(t, "unhealthy", dbComponent.Status)
		assert.Contains(t, dbComponent.Message, "Failed to get sql.DB instance")

		mockCache.AssertExpectations(t)
		mockDB.AssertExpectations(t)
	})
}
