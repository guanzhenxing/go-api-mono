package service

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"go-api-mono/internal/app/user/model"
	"go-api-mono/internal/pkg/database"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrEmailAlreadyUsed  = errors.New("email already in use")
	ErrUsernameNotFound  = errors.New("username not found")
	ErrInvalidCredential = errors.New("invalid credentials")
)

// UserService defines the interface for user operations
type UserService interface {
	Register(ctx context.Context, user *model.User) error
	Authenticate(ctx context.Context, email, password string) (*model.User, error)
	Get(ctx context.Context, id uint) (*model.User, error)
	List(ctx context.Context) ([]*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint) error
}

type userService struct {
	db database.DB
}

// NewUserService creates a new user service
func NewUserService(db database.DB) UserService {
	return &userService{db: db}
}

// Register handles user registration
func (s *userService) Register(ctx context.Context, user *model.User) error {
	// 对密码进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hashedPassword)

	// 创建用户
	if err := s.db.WithContext(ctx).Create(user); err != nil {
		return err
	}

	return nil
}

// Authenticate verifies user credentials and returns the user if valid
func (s *userService) Authenticate(ctx context.Context, email, password string) (*model.User, error) {
	var user model.User
	if err := s.db.WithContext(ctx).First(&user, "email = ?", email); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidPassword
	}

	return &user, nil
}

func (s *userService) Get(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	if err := s.db.WithContext(ctx).First(&user, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (s *userService) List(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	// TODO: Implement list users with pagination
	return users, nil
}

func (s *userService) Update(ctx context.Context, user *model.User) error {
	// 获取现有用户信息
	existingUser, err := s.Get(ctx, user.ID)
	if err != nil {
		return err
	}

	// 保留原始的创建时间
	user.CreatedAt = existingUser.CreatedAt

	// 如果提供了新密码，对其进行哈希处理
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		user.Password = string(hashedPassword)
	} else {
		user.Password = existingUser.Password
	}

	return s.db.WithContext(ctx).Save(user)
}

func (s *userService) Delete(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&model.User{}, id)
}
