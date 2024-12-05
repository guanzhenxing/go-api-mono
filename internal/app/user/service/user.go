package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"go-api-mono/internal/app/user/model"
	"go-api-mono/internal/app/user/repository"
	"go-api-mono/pkg/database"
	"go-api-mono/pkg/logger"
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
	db         database.DB
	logger     *logger.Logger
	repository repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(db database.DB, logger *logger.Logger) UserService {
	return &userService{
		db:         db,
		logger:     logger,
		repository: repository.NewUserRepository(db, logger),
	}
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
	if err := s.repository.Create(ctx, user); err != nil {
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
	return s.repository.GetByID(ctx, id)
}

func (s *userService) List(ctx context.Context) ([]*model.User, error) {
	return s.repository.List(ctx)
}

func (s *userService) Update(ctx context.Context, user *model.User) error {
	// 获取现有用户信息
	existingUser, err := s.repository.GetByID(ctx, user.ID)
	if err != nil {
		return err
	}

	s.logger.Info("Updating user",
		zap.Uint("id", user.ID),
		zap.String("existing_username", existingUser.Username),
		zap.String("existing_email", existingUser.Email),
		zap.Time("existing_created_at", existingUser.CreatedAt),
		zap.Time("existing_updated_at", existingUser.UpdatedAt),
		zap.String("new_username", user.Username),
		zap.String("new_email", user.Email),
	)

	// 准备更新的字段
	updates := map[string]interface{}{
		"username":   user.Username,
		"email":      user.Email,
		"updated_at": time.Now(),
	}

	// 如果提供了新密码，则更新密码
	if user.Password != "" {
		s.logger.Info("Updating password")
		hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if hashErr != nil {
			s.logger.Error("Failed to hash password", zap.Error(hashErr))
			return fmt.Errorf("failed to hash password: %w", hashErr)
		}
		updates["password"] = string(hashedPassword)
	}

	s.logger.Info("Executing update query", zap.Any("updates", updates))

	// 使用链式调用
	err = s.repository.Update(ctx, user)
	if err != nil {
		s.logger.Error("Failed to update user", zap.Error(err))
		return err
	}

	s.logger.Info("User updated successfully")
	return nil
}

func (s *userService) Delete(ctx context.Context, id uint) error {
	return s.repository.Delete(ctx, id)
}
