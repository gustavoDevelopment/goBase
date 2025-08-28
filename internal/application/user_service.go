package application

import (
	"api-ptf-core-business-orchestrator-go-ms/internal/domain"
	"api-ptf-core-business-orchestrator-go-ms/internal/infrastructure/repository"
	"context"
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
func (s *UserService) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	// Here you can add business logic before/after the repository call
	// For example: validation, logging, caching, etc.
	
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// You can add additional business logic here
	// For example: data transformation, enrichment, etc.

	return user, nil
}
