package repository

import (
	"context"

	"go.uber.org/zap"

	"go-api-mono/internal/app/user/model"
	"go-api-mono/internal/pkg/database"
	"go-api-mono/internal/pkg/logger"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context) ([]*model.User, error)
}

type userRepository struct {
	db     database.DB
	logger *logger.Logger
}

// NewUserRepository creates a new user repository
func NewUserRepository(db database.DB, logger *logger.Logger) UserRepository {
	return &userRepository{
		db:     db,
		logger: logger,
	}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	r.logger.Info("Creating user", zap.String("username", user.Username), zap.String("email", user.Email))
	err := r.db.WithContext(ctx).Create(user)
	if err != nil {
		r.logger.Error("Failed to create user", zap.Error(err))
		return err
	}
	r.logger.Info("User created successfully", zap.Uint("id", user.ID))
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	r.logger.Info("Getting user by ID", zap.Uint("id", id))
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, id); err != nil {
		r.logger.Error("Failed to get user by ID", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	r.logger.Info("Getting user by username", zap.String("username", username))
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, "username = ?", username); err != nil {
		r.logger.Error("Failed to get user by username", zap.String("username", username), zap.Error(err))
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	r.logger.Info("Updating user", zap.Uint("id", user.ID), zap.String("username", user.Username))
	err := r.db.WithContext(ctx).Debug().Model(user).Updates(map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
		"password": user.Password,
	})
	if err != nil {
		r.logger.Error("Failed to update user", zap.Uint("id", user.ID), zap.Error(err))
		return err
	}
	r.logger.Info("User updated successfully", zap.Uint("id", user.ID))
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	r.logger.Info("Deleting user", zap.Uint("id", id))
	err := r.db.WithContext(ctx).Delete(&model.User{}, id)
	if err != nil {
		r.logger.Error("Failed to delete user", zap.Uint("id", id), zap.Error(err))
		return err
	}
	r.logger.Info("User deleted successfully", zap.Uint("id", id))
	return nil
}

func (r *userRepository) List(ctx context.Context) ([]*model.User, error) {
	r.logger.Info("Listing users")
	var users []*model.User
	// TODO: Implement list users with proper pagination
	r.logger.Info("Users listed successfully", zap.Int("count", len(users)))
	return users, nil
}
