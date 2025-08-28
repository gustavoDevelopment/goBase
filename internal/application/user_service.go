package application

import (
	"api-ptf-core-business-orchestrator-go-ms/internal/domain"
	"api-ptf-core-business-orchestrator-go-ms/internal/infrastructure/repository"
	"context"
	"time"
)

// UserService contains the business logic for users
type UserService struct {
	repo repository.UserRepository
}

// NewUserService creates a new UserService
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// GetUserByID retrieves a user by their ID
func (s *UserService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	// Add any business logic/validation here before creating the user
	user.DateCreated = time.Now()
	user.DateUpdated = time.Now()

	id, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	// Fetch the created user to return it with the generated ID
	return s.repo.FindByID(ctx, id)
}

// ListUsers retrieves all users with pagination
func (s *UserService) ListUsers(ctx context.Context, page, limit int64) ([]domain.User, error) {
	return s.repo.FindAll(ctx, page, limit)
}
