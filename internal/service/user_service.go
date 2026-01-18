package service

import (
	"context"
	"fmt"

	"github.com/ihsanbudiman/gh-action/internal/domain"
)

// userService implements domain.UserService
type userService struct {
	repo domain.UserRepository
}

// NewUserService creates a new user service
func NewUserService(repo domain.UserRepository) domain.UserService {
	return &userService{repo: repo}
}

// Create creates a new user
func (s *userService) Create(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error) {
	user := &domain.User{
		Name:  req.Name,
		Email: req.Email,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("service: failed to create user: %w", err)
	}

	return user, nil
}

// GetByID retrieves a user by ID
func (s *userService) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get user: %w", err)
	}
	return user, nil
}

// GetAll retrieves all users
func (s *userService) GetAll(ctx context.Context) ([]*domain.User, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get users: %w", err)
	}
	return users, nil
}

// Update updates an existing user
func (s *userService) Update(ctx context.Context, id int64, req *domain.UpdateUserRequest) (*domain.User, error) {
	// First, get the existing user
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("service: user not found")
	}

	// Update fields
	user.Name = req.Name
	user.Email = req.Email

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("service: failed to update user: %w", err)
	}

	return user, nil
}

// Delete removes a user
func (s *userService) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("service: failed to delete user: %w", err)
	}
	return nil
}
