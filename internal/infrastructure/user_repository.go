package infrastructure

import (
	"api-ptf-core-business-orchestrator-go-ms/internal/domain"
	"context"
	"fmt"
)

// InMemoryUserRepository is a mock implementation of the UserRepository interface.
// It stores users in memory.
type InMemoryUserRepository struct {
	users map[int64]*domain.User
}

// NewInMemoryUserRepository creates a new InMemoryUserRepository.
func NewInMemoryUserRepository() *InMemoryUserRepository {
	// Initialize with some dummy data
	users := make(map[int64]*domain.User)
	return &InMemoryUserRepository{users: users}
}

// FindByID retrieves a user by their ID from the in-memory store.
func (r *InMemoryUserRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	if user, ok := r.users[id]; ok {
		return user, nil
	}
	return nil, fmt.Errorf("user with id %d not found", id)
}
