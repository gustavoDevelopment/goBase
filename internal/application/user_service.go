package application

import (
	"api-ptf-core-business-orchestrator-go-ms/internal/domain"
	"api-ptf-core-business-orchestrator-go-ms/internal/infrastructure/repository"
	"api-ptf-core-business-orchestrator-go-ms/internal/pkg/logger"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
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

// generateRandomPassword generates a random 6-character string
func generateRandomPassword(length int) string {
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
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

	// Set timestamps and ID
	now := time.Now()
	user.ID = uuid.New().String()
	user.DateCreated = now

	// Generate random password if not provided
	if user.Password == "" {
		user.Password = generateRandomPassword(6)
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}
	user.Password = string(hashedPassword)

	logger.Logger().Info("Creating user", zap.String("user_id", user.ID), zap.String("email", user.Email))

	// Create the user
	id, err := s.genericRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	// Set the ID from the created user
	user.ID = id

	return user, nil
}
