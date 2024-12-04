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
	if fn, ok := args.Get(0).(func(interface{})); ok {
		fn(dest)
	}
	return args.Error(1)
}

func (m *MockDB) GetDB() *gorm.DB {
	return nil
}

func TestUserService_Register(t *testing.T) {
	mockDB := &MockDB{}
	service := NewUserService(mockDB)

	user := &model.User{
		Username: "testuser",
		Password: "password",
	}

	mockDB.On("Create", mock.AnythingOfType("*model.User")).Return(nil)

	err := service.Register(context.Background(), user)
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestUserService_Get(t *testing.T) {
	mockDB := &MockDB{}
	service := NewUserService(mockDB)

	user := &model.User{
		ID:       1,
		Username: "testuser",
	}

	mockDB.On("First", mock.AnythingOfType("*model.User"), uint(1)).Return(func(dest interface{}) {
		u := dest.(*model.User)
		*u = *user
	}, nil)

	result, err := service.Get(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, user, result)
	mockDB.AssertExpectations(t)
}

func TestUserService_Update(t *testing.T) {
	mockDB := &MockDB{}
	service := NewUserService(mockDB)

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

	// Mock First call to get existing user
	mockDB.On("First", mock.AnythingOfType("*model.User"), uint(1)).Return(func(dest interface{}) {
		u := dest.(*model.User)
		*u = *existingUser
	}, nil)

	// Mock Save call
	mockDB.On("Save", mock.AnythingOfType("*model.User")).Return(nil)

	err := service.Update(context.Background(), updatedUser)
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestUserService_Delete(t *testing.T) {
	mockDB := &MockDB{}
	service := NewUserService(mockDB)

	mockDB.On("Delete", &model.User{}, uint(1)).Return(nil)

	err := service.Delete(context.Background(), 1)
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}
