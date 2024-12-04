package service

import (
	"context"

	"go-api-mono/internal/app/user/model"
	"go-api-mono/internal/pkg/database"
)

// UserService defines the interface for user operations
type UserService interface {
	Register(ctx context.Context, user *model.User) error
	Authenticate(ctx context.Context, username, password string) (*model.User, error)
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

func (s *userService) Register(ctx context.Context, user *model.User) error {
	return s.db.WithContext(ctx).Create(user)
}

func (s *userService) Authenticate(ctx context.Context, username, password string) (*model.User, error) {
	var user model.User
	if err := s.db.WithContext(ctx).First(&user, "username = ?", username); err != nil {
		return nil, err
	}
	// TODO: Implement password verification
	return &user, nil
}

func (s *userService) Get(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	if err := s.db.WithContext(ctx).First(&user, id); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *userService) List(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	// TODO: Implement list users
	return users, nil
}

func (s *userService) Update(ctx context.Context, user *model.User) error {
	return s.db.WithContext(ctx).Save(user)
}

func (s *userService) Delete(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&model.User{}, id)
}
