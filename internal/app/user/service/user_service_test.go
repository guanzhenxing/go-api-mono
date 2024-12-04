package service

import (
	"context"
	"net/http"
	"testing"

	"go-api-mono/internal/app/user/model"
	"go-api-mono/internal/pkg/core"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository 是用户仓储的mock实现
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

// TestUserService_RegisterUser 测试用户注册
func TestUserService_RegisterUser(t *testing.T) {
	tests := []struct {
		name    string
		user    *model.User
		setup   func(*MockUserRepository)
		wantErr bool
	}{
		{
			name: "成功注册",
			user: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			setup: func(repo *MockUserRepository) {
				repo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, nil)
				repo.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "邮箱已存在",
			user: &model.User{
				Username: "testuser",
				Email:    "exists@example.com",
				Password: "password123",
			},
			setup: func(repo *MockUserRepository) {
				repo.On("GetByEmail", mock.Anything, "exists@example.com").Return(&model.User{}, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockUserRepository)
			if tt.setup != nil {
				tt.setup(repo)
			}

			s := NewUserService(repo)
			ctx := &core.Context{
				Context: context.Background(),
				Request: &http.Request{},
			}

			err := s.RegisterUser(ctx, tt.user)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			repo.AssertExpectations(t)
		})
	}
}

// TestUserService_AuthenticateUser 测试用户认证
func TestUserService_AuthenticateUser(t *testing.T) {
	// Generate a proper password hash
	password := "password123"
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)
	hashedPassword := string(hashedBytes)

	tests := []struct {
		name     string
		email    string
		password string
		setup    func(*MockUserRepository)
		want     *model.User
		wantErr  bool
	}{
		{
			name:     "认证成功",
			email:    "test@example.com",
			password: password,
			setup: func(repo *MockUserRepository) {
				repo.On("GetByEmail", mock.Anything, "test@example.com").Return(&model.User{
					Email:    "test@example.com",
					Password: hashedPassword,
				}, nil)
			},
			wantErr: false,
		},
		{
			name:     "用户不存在",
			email:    "notfound@example.com",
			password: password,
			setup: func(repo *MockUserRepository) {
				repo.On("GetByEmail", mock.Anything, "notfound@example.com").Return(nil, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockUserRepository)
			if tt.setup != nil {
				tt.setup(repo)
			}

			s := NewUserService(repo)
			ctx := &core.Context{
				Context: context.Background(),
				Request: &http.Request{},
			}

			user, err := s.AuthenticateUser(ctx, tt.email, tt.password)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
			}

			repo.AssertExpectations(t)
		})
	}
}

// BenchmarkUserService_RegisterUser 基准测试用户注册
func BenchmarkUserService_RegisterUser(b *testing.B) {
	repo := new(MockUserRepository)
	repo.On("GetByEmail", mock.Anything, mock.Anything).Return(nil, nil)
	repo.On("Create", mock.Anything, mock.Anything).Return(nil)

	s := NewUserService(repo)
	ctx := &core.Context{
		Context: context.Background(),
		Request: &http.Request{},
	}
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.RegisterUser(ctx, user)
	}
}

// BenchmarkUserService_AuthenticateUser 基准测试用户认证
func BenchmarkUserService_AuthenticateUser(b *testing.B) {
	hashedPassword := "$2a$10$ZWVyeSBzZWNyZXQga2V5IGNoYW5nZSBpdCBpbiBwcm9kdWN0aW9u"
	repo := new(MockUserRepository)
	repo.On("GetByEmail", mock.Anything, mock.Anything).Return(&model.User{
		Email:    "test@example.com",
		Password: hashedPassword,
	}, nil)

	s := NewUserService(repo)
	ctx := &core.Context{
		Context: context.Background(),
		Request: &http.Request{},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.AuthenticateUser(ctx, "test@example.com", "password123")
	}
}
