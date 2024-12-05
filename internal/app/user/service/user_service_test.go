package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"go-api-mono/internal/app/user/model"
	"go-api-mono/internal/pkg/database"
	"go-api-mono/internal/pkg/logger"
)

// MockUserRepository is a mock implementation of repository.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	args := m.Called(ctx, id)
	if user := args.Get(0); user != nil {
		return user.(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	args := m.Called(ctx, username)
	if user := args.Get(0); user != nil {
		return user.(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context) ([]*model.User, error) {
	args := m.Called(ctx)
	if users := args.Get(0); users != nil {
		return users.([]*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}

// MockDB is a mock implementation of database.DB
type MockDB struct {
	mock.Mock
}

func (m *MockDB) WithContext(ctx context.Context) database.DB {
	args := m.Called(ctx)
	if db := args.Get(0); db != nil {
		return db.(database.DB)
	}
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
	if fn, ok := args.Get(0).(func(interface{})); ok {
		fn(dest)
	}
	return args.Error(1)
}

func (m *MockDB) Model(value interface{}) database.DB {
	args := m.Called(value)
	if db := args.Get(0); db != nil {
		return db.(database.DB)
	}
	return m
}

func (m *MockDB) Updates(values interface{}) error {
	args := m.Called(values)
	return args.Error(0)
}

func (m *MockDB) Debug() database.DB {
	args := m.Called()
	if db := args.Get(0); db != nil {
		return db.(database.DB)
	}
	return m
}

func (m *MockDB) GetDB() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func TestUserService_Register(t *testing.T) {
	mockDB := &MockDB{}
	mockRepo := &MockUserRepository{}

	// 创建一个真实的 logger 实例用于测试
	testLogger, err := logger.New(logger.LogConfig{
		Level:      "debug",
		Filename:   "test.log",
		MaxSize:    1,
		MaxBackups: 1,
		MaxAge:     1,
		Compress:   false,
	})
	assert.NoError(t, err)

	service := &userService{
		db:         mockDB,
		logger:     testLogger,
		repository: mockRepo,
	}

	user := &model.User{
		Username: "testuser",
		Password: "password",
	}

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)

	err = service.Register(context.Background(), user)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Get(t *testing.T) {
	mockDB := &MockDB{}
	mockRepo := &MockUserRepository{}

	// 创建一个真实的 logger 实例用于测试
	testLogger, err := logger.New(logger.LogConfig{
		Level:      "debug",
		Filename:   "test.log",
		MaxSize:    1,
		MaxBackups: 1,
		MaxAge:     1,
		Compress:   false,
	})
	assert.NoError(t, err)

	service := &userService{
		db:         mockDB,
		logger:     testLogger,
		repository: mockRepo,
	}

	user := &model.User{
		ID:       1,
		Username: "testuser",
	}

	mockRepo.On("GetByID", mock.Anything, uint(1)).Return(user, nil)

	result, err := service.Get(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, user, result)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Update(t *testing.T) {
	mockDB := &MockDB{}
	mockRepo := &MockUserRepository{}

	// 创建一个真实的 logger 实例用于测试
	testLogger, err := logger.New(logger.LogConfig{
		Level:      "debug",
		Filename:   "test.log",
		MaxSize:    1,
		MaxBackups: 1,
		MaxAge:     1,
		Compress:   false,
	})
	assert.NoError(t, err)

	service := &userService{
		db:         mockDB,
		logger:     testLogger,
		repository: mockRepo,
	}

	existingUser := &model.User{
		ID:        1,
		Username:  "testuser",
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	updatedUser := &model.User{
		ID:       1,
		Username: "newusername",
		Password: "newpassword",
	}

	mockRepo.On("GetByID", mock.Anything, uint(1)).Return(existingUser, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)

	err = service.Update(context.Background(), updatedUser)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Delete(t *testing.T) {
	mockDB := &MockDB{}
	mockRepo := &MockUserRepository{}

	// 创建一个真实的 logger 实例用于测试
	testLogger, err := logger.New(logger.LogConfig{
		Level:      "debug",
		Filename:   "test.log",
		MaxSize:    1,
		MaxBackups: 1,
		MaxAge:     1,
		Compress:   false,
	})
	assert.NoError(t, err)

	service := &userService{
		db:         mockDB,
		logger:     testLogger,
		repository: mockRepo,
	}

	mockRepo.On("Delete", mock.Anything, uint(1)).Return(nil)

	err = service.Delete(context.Background(), 1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
