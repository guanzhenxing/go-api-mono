package service

import (
	"fmt"

	"go-api-mono/internal/app/user/model"
	"go-api-mono/internal/app/user/repository"
	"go-api-mono/internal/pkg/core"
	"go-api-mono/internal/pkg/errors"

	"golang.org/x/crypto/bcrypt"
)

// UserService 用户服务
type UserService struct {
	repo repository.UserRepository
}

// NewUserService 创建用户服务
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// RegisterUser 注册用户
func (s *UserService) RegisterUser(ctx *core.Context, user *model.User) error {
	// 检查邮箱是否已存在
	existingUser, err := s.repo.GetByEmail(ctx.Request.Context(), user.Email)
	if err != nil {
		return fmt.Errorf("failed to check email: %w", err)
	}
	if existingUser != nil {
		return errors.ErrUserExists
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hashedPassword)

	// 创建用户
	if err := s.repo.Create(ctx.Request.Context(), user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// AuthenticateUser 用户认证
func (s *UserService) AuthenticateUser(ctx *core.Context, email, password string) (*model.User, error) {
	// 获取用户
	user, err := s.repo.GetByEmail(ctx.Request.Context(), email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, errors.ErrInvalidCredentials
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.ErrInvalidCredentials
	}

	return user, nil
}

// GetUser 获取用户
func (s *UserService) GetUser(ctx *core.Context, id uint) (*model.User, error) {
	user, err := s.repo.GetByID(ctx.Request.Context(), id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, errors.ErrUserNotFound
	}
	return user, nil
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(ctx *core.Context, user *model.User) error {
	// 检查用户是否存在
	existingUser, err := s.repo.GetByID(ctx.Request.Context(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if existingUser == nil {
		return errors.ErrUserNotFound
	}

	// 如果提供了新密码，则加密
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		user.Password = string(hashedPassword)
	} else {
		user.Password = existingUser.Password
	}

	if err := s.repo.Update(ctx.Request.Context(), user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(ctx *core.Context, id uint) error {
	// 检查用户是否存在
	user, err := s.repo.GetByID(ctx.Request.Context(), id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return errors.ErrUserNotFound
	}

	if err := s.repo.Delete(ctx.Request.Context(), id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ListUsers 获取用户列表
func (s *UserService) ListUsers(ctx *core.Context, page, pageSize int) ([]model.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	users, total, err := s.repo.List(ctx.Request.Context(), page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, total, nil
}
