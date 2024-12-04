package health

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"go-api-mono/internal/pkg/database"
)

// mockDB 模拟数据库接口
type mockDB struct {
	mock.Mock
}

func (m *mockDB) WithContext(ctx context.Context) database.DB {
	return m
}

func (m *mockDB) Create(value interface{}) error {
	args := m.Called(value)
	return args.Error(0)
}

func (m *mockDB) Save(value interface{}) error {
	args := m.Called(value)
	return args.Error(0)
}

func (m *mockDB) Delete(value interface{}, where ...interface{}) error {
	args := m.Called(value, where[0])
	return args.Error(0)
}

func (m *mockDB) First(dest interface{}, conds ...interface{}) error {
	args := m.Called(dest, conds[0])
	return args.Error(0)
}

func (m *mockDB) GetDB() *gorm.DB {
	args := m.Called()
	if db := args.Get(0); db != nil {
		return db.(*gorm.DB)
	}
	return nil
}

// mockCache 模拟缓存接口
type mockCache struct {
	mock.Mock
}

func (m *mockCache) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *mockCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *mockCache) Del(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *mockCache) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *mockCache) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestChecker_Check(t *testing.T) {
	// 创建模拟对象
	mockDB := new(mockDB)
	mockCache := new(mockCache)
	version := "1.0.0"

	// 创建健康检查器
	checker := NewChecker(mockDB, mockCache, version)

	// 设置模拟行为 - 模拟正常工作的数据库
	mockGormDB := &gorm.DB{
		Config: &gorm.Config{},
	}
	mockDB.On("GetDB").Return(mockGormDB)
	mockCache.On("Ping", mock.Anything).Return(nil)

	// 执行健康检查
	ctx := context.Background()
	resp := checker.Check(ctx)

	// 验证结果
	assert.Equal(t, "unhealthy", resp.Status) // 数据库连接不完整，所以状态是 unhealthy
	assert.Equal(t, version, resp.Version)
	assert.NotEmpty(t, resp.Timestamp)
	assert.NotEmpty(t, resp.Uptime)
	assert.Len(t, resp.Components, 3) // 数据库、Redis和系统组件

	// 验证组件状态
	var dbComponent *Status
	for _, component := range resp.Components {
		if component.Component == "database" {
			dbComponent = &component
			break
		}
	}

	assert.NotNil(t, dbComponent)
	assert.Equal(t, "unhealthy", dbComponent.Status)
	assert.Contains(t, dbComponent.Message, "Failed to get sql.DB instance")

	// 验证模拟对象的调用
	mockDB.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestChecker_CheckUnhealthy(t *testing.T) {
	// 创建模拟对象
	mockDB := new(mockDB)
	mockCache := new(mockCache)
	version := "1.0.0"

	// 创建健康检查器
	checker := NewChecker(mockDB, mockCache, version)

	// 设置模拟行为 - 模拟故障情况
	mockDB.On("GetDB").Return(nil)
	mockCache.On("Ping", mock.Anything).Return(assert.AnError)

	// 执行健康检查
	ctx := context.Background()
	resp := checker.Check(ctx)

	// 验证结果
	assert.Equal(t, "unhealthy", resp.Status)
	assert.Equal(t, version, resp.Version)
	assert.NotEmpty(t, resp.Timestamp)
	assert.NotEmpty(t, resp.Uptime)

	// 验证组件状态
	var redisComponent *Status
	for _, component := range resp.Components {
		if component.Component == "redis" {
			redisComponent = &component
			break
		}
	}

	assert.NotNil(t, redisComponent)
	assert.Equal(t, "unhealthy", redisComponent.Status)
	assert.Contains(t, redisComponent.Message, "Redis ping failed")

	// 验证模拟对象的调用
	mockDB.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}
