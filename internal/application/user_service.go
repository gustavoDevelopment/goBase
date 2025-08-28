package application

import (
	"api-ptf-core-business-orchestrator-go-ms/internal/domain"
	"api-ptf-core-business-orchestrator-go-ms/internal/infrastructure/repository"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// UserService contiene la lógica de negocio para los usuarios
type UserService struct {
	genericRepo *repository.GenericRepository[domain.User]
	userRepo    repository.UserRepository
}

// NewUserService crea una nueva instancia de UserService
func NewUserService(genericRepo *repository.GenericRepository[domain.User], userRepo repository.UserRepository) *UserService {
	return &UserService{
		genericRepo: genericRepo,
		userRepo:    userRepo,
	}
}

// GetUserByID obtiene un usuario por su ID
func (s *UserService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	return s.genericRepo.FindByID(ctx, id)
}

// UserRepository define la interfaz para operaciones de datos de usuarios
type UserRepository interface {
	// Métodos adicionales
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Count(ctx context.Context) (int64, error)
}

// GetUserByEmail busca un usuario por su correo electrónico
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.userRepo.FindByEmail(ctx, email)
}

// CountUsers devuelve el número total de usuarios
func (s *UserService) CountUsers(ctx context.Context) (int64, error) {
	return s.genericRepo.Count(ctx)
}

// ListUsers returns a paginated list of users
func (s *UserService) ListUsers(ctx context.Context, page, limit int64) ([]domain.User, error) {
	return s.genericRepo.FindAll(ctx, page, limit)
}

// Create creates a new user
func (s *UserService) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	// Check if user with email already exists
	existingUser, err := s.userRepo.FindByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("error checking for existing user: %w", err)
	}

	if existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", user.Email)
	}

	// Hash the password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}
	user.Password = string(hashedPassword)

	// Set timestamps
	now := time.Now()
	user.DateCreated = now
	user.DateUpdated = now

	// Create the user
	id, err := s.genericRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	// Set the ID from the created user
	user.ID = id

	return user, nil
}
