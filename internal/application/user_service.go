package application

import (
	"api-ptf-core-business-orchestrator-go-ms/internal/domain"
	"context"
)

// UserRepository is an interface for user persistence.
// The implementation is in the infrastructure layer.
type UserRepository interface {
	FindByID(ctx context.Context, id int64) (*domain.User, error)
}

// UserService contains the business logic for users.
type UserService struct {
	userRepo UserRepository
}

// NewUserService creates a new UserService.
func NewUserService(userRepo UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// GetUserByID retrieves a user by their ID.
func (s *UserService) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	return s.userRepo.FindByID(ctx, id)
}
