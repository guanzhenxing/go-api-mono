package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"go-api-mono/internal/app/user/model"
	"go-api-mono/internal/pkg/database"
)

type MockDB struct {
	mock.Mock
	db *gorm.DB
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
	return m.db
}

func TestUserService_Register(t *testing.T) {
	mockDB := &MockDB{}
	service := NewUserService(mockDB)

	user := &model.User{
		Username: "testuser",
		Password: "password",
	}

	mockDB.On("Create", user).Return(nil)

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

	mockDB.On("First", mock.Anything, uint(1)).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*model.User)
		*arg = *user
	})

	result, err := service.Get(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, user, result)
	mockDB.AssertExpectations(t)
}

func TestUserService_Update(t *testing.T) {
	mockDB := &MockDB{}
	service := NewUserService(mockDB)

	user := &model.User{
		ID:       1,
		Username: "testuser",
	}

	mockDB.On("Save", user).Return(nil)

	err := service.Update(context.Background(), user)
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
