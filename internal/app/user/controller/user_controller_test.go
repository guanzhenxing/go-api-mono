package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-api-mono/internal/app/user/model"
	"go-api-mono/internal/pkg/validator"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Register(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserService) Authenticate(ctx context.Context, username, password string) (*model.User, error) {
	args := m.Called(ctx, username, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) Get(ctx context.Context, id uint) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) List(ctx context.Context) ([]*model.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.User), args.Error(1)
}

func (m *MockUserService) Update(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserService) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func init() {
	gin.SetMode(gin.TestMode)
	if err := validator.Register(); err != nil {
		panic(err)
	}
}

func TestUserController_Create(t *testing.T) {
	mockService := new(MockUserService)
	controller := NewUserController(mockService)

	router := gin.New()
	router.POST("/users", controller.Create)

	user := &model.User{
		Username: "testuser123",
		Email:    "test@gmail.com",
		Password: "Test123!@#",
	}

	mockService.On("Register", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)

	body, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestUserController_Get(t *testing.T) {
	mockService := new(MockUserService)
	controller := NewUserController(mockService)

	router := gin.New()
	router.GET("/users/:id", controller.Get)

	user := &model.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@gmail.com",
	}

	mockService.On("Get", mock.Anything, uint(1)).Return(user, nil)

	req := httptest.NewRequest("GET", "/users/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestUserController_Update(t *testing.T) {
	mockService := new(MockUserService)
	controller := NewUserController(mockService)

	router := gin.New()
	router.PUT("/users/:id", controller.Update)

	user := &model.User{
		ID:       1,
		Username: "testuser123",
		Email:    "test@gmail.com",
		Password: "Test123!@#",
	}

	mockService.On("Update", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)

	body, _ := json.Marshal(user)
	req := httptest.NewRequest("PUT", "/users/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestUserController_Delete(t *testing.T) {
	mockService := new(MockUserService)
	controller := NewUserController(mockService)

	router := gin.New()
	router.DELETE("/users/:id", controller.Delete)

	mockService.On("Delete", mock.Anything, uint(1)).Return(nil)

	req := httptest.NewRequest("DELETE", "/users/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockService.AssertExpectations(t)
}
