package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-api-mono/internal/app/user/model"
	"go-api-mono/internal/app/user/service"
	"go-api-mono/internal/pkg/auth"
	"go-api-mono/internal/pkg/core"
	"go-api-mono/internal/pkg/errors"
	"go-api-mono/internal/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository 模拟用户仓储
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]model.User), args.Get(1).(int64), args.Error(2)
}

func setupTest(t *testing.T) (*UserController, *MockUserRepository, *auth.JWT, *logger.Logger) {
	mockRepo := new(MockUserRepository)
	log, err := logger.New(logger.LogConfig{
		Level:    "debug",
		Filename: "test.log",
	})
	assert.NoError(t, err)

	jwt := auth.New(auth.Config{
		SigningKey:     "test-key",
		ExpirationTime: time.Hour,
		SigningMethod:  "HS256",
		TokenPrefix:    "Bearer",
	})

	userService := service.NewUserService(mockRepo)
	controller := NewUserController(userService, jwt)
	return controller, mockRepo, jwt, log
}

func TestRegister(t *testing.T) {
	controller, mockRepo, _, log := setupTest(t)

	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func()
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "successful registration",
			requestBody: model.RegisterRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func() {
				mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, nil)
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "invalid request body",
			requestBody:    "invalid json",
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "service error",
			requestBody: model.RegisterRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func() {
				mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, nil)
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).
					Return(errors.New(errors.ErrCodeInternal, "service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock for each test case
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil

			// Setup mock expectations
			tt.setupMock()

			// Create request
			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Create context
			ctx := core.NewContext(req, w, log)

			// Execute
			controller.Register(ctx)

			// Assert response
			if tt.expectedError {
				assert.GreaterOrEqual(t, w.Code, http.StatusBadRequest)
			} else {
				assert.Equal(t, tt.expectedStatus, w.Code)
			}

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}

// Add more test functions for other controller methods...
