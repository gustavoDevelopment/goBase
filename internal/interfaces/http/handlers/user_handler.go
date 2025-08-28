package handlers

import (
	"api-ptf-core-business-orchestrator-go-ms/internal/application"
	"api-ptf-core-business-orchestrator-go-ms/internal/domain"
	httpMiddleware "api-ptf-core-business-orchestrator-go-ms/internal/interfaces/http/middleware"
	"api-ptf-core-business-orchestrator-go-ms/internal/interfaces/http/utils"
	"api-ptf-core-business-orchestrator-go-ms/internal/pkg/logger"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// Common error messages
var (
	ErrInvalidUserID    = errors.New("invalid user ID format")
	ErrInvalidPageParam = errors.New("invalid page parameter")
	ErrInvalidLimit     = errors.New("limit must be between 1 and 100")
	ErrInvalidEmail     = errors.New("valid email is required")
)

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"pass" validate:"required,min=6"`
}

type UserHandler struct {
	userService *application.UserService
}

func NewUserHandler(userService *application.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUserByID handles GET /api/v1/users/{id}
// @Summary Get user by ID
// @Description Get a user by their ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} domain.User
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := httpMiddleware.GetRequestID(r.Context())

	// Get logger from context and add request context
	logger := logger.FromContext(r.Context()).With(
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
	)

	// Log request
	logger.Info("GetUserByID started")

	vars := mux.Vars(r)
	userID := strings.TrimSpace(vars["id"])

	logger = logger.With(zap.String("user_id", userID))
	logger.Info("Fetching user by ID")

	if userID == "" {
		logger.Warn("Empty user ID provided")
		_ = utils.BadRequest(w, "User ID is required")
		return
	}

	// Log database call
	dbStart := time.Now()
	user, err := h.userService.GetUserByID(r.Context(), userID)
	dbDuration := time.Since(dbStart)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.Warn("User not found")
			_ = utils.NotFound(w, "User not found")
			return
		}

		logger.Error("Failed to fetch user", zap.Error(err))
		_ = utils.InternalServerError(w, "Failed to fetch user: "+err.Error())
		return
	}

	// Log successful response
	logger.Info("GetUserByID completed",
		zap.String("request_id", requestID),
		zap.String("user_id", userID),
		zap.Duration("db_duration", dbDuration),
		zap.Duration("total_duration", time.Since(start)),
	)

	_ = utils.SendSuccess(w, "SUCCESS", "User retrieved successfully", http.StatusOK, user)
}

// CreateUser handles POST /api/v1/users
// @Summary Create a new user
// @Description Create a new user with the provided data
// @Tags users
// @Accept json
// @Produce json
// @Param user body domain.User true "User data"
// @Success 201 {object} domain.User
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := httpMiddleware.GetRequestID(r.Context())

	// Log request
	logger := logger.FromContext(r.Context())
	logger.Info("CreateUser started",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
	)

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Failed to decode request body", zap.Error(err))
		_ = utils.BadRequest(w, "Invalid request body")
		return
	}

	// Create domain user
	user := domain.User{
		Email:    req.Email,
		Password: req.Password,
	}

	logger = logger.With(zap.String("email", user.Email))
	logger.Info("Creating new user")

	// Log user creation attempt
	dbStart := time.Now()
	createdUser, err := h.userService.Create(r.Context(), &user)
	dbDuration := time.Since(dbStart)

	if err != nil {
		logger.Error("Failed to create user", zap.Error(err))
		_ = utils.InternalServerError(w, "Failed to create user: "+err.Error())
		return
	}

	// Log successful creation
	logger.Info("User created successfully",
		zap.String("request_id", requestID),
		zap.String("user_id", createdUser.ID),
		zap.String("email", createdUser.Email),
		zap.Duration("db_duration", dbDuration),
		zap.Duration("total_duration", time.Since(start)),
	)

	_ = utils.SendSuccess(w, "USER_CREATED", "User created successfully", http.StatusCreated, createdUser)
}

// ListUsersResponse represents the response structure for the ListUsers endpoint
type ListUsersResponse struct {
	Data       interface{} `json:"data"`
	Pagination struct {
		Page  int64 `json:"page"`
		Limit int64 `json:"limit"`
		Total int64 `json:"total,omitempty"`
	} `json:"pagination"`
}

// ListUsers handles GET /api/v1/users
// @Summary List all users with pagination
// @Description Get a paginated list of users
// @Tags users
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (max: 100, default: 10)"
// @Success 200 {object} ListUsersResponse
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users [get]
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := httpMiddleware.GetRequestID(r.Context())

	// Log request
	logger := logger.FromContext(r.Context())
	logger.Info("ListUsers started",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("query", r.URL.RawQuery),
	)

	// Parse pagination parameters
	page, limit, err := httpMiddleware.GetPaginationParams(r)
	if err != nil {
		logger.Warn("Invalid pagination parameters", zap.Error(err))
		_ = utils.BadRequest(w, err.Error())
		return
	}

	logger = logger.With(
		zap.Int64("page", page),
		zap.Int64("limit", limit),
	)
	logger.Info("Fetching users list")

	// Get users using the service layer
	dbStart := time.Now()
	users, err := h.userService.ListUsers(r.Context(), page, limit)
	dbDuration := time.Since(dbStart)

	if err != nil {
		logger.Error("Failed to fetch users", zap.Error(err))
		_ = utils.InternalServerError(w, "Failed to fetch users: "+err.Error())
		return
	}

	// Get total count (implementation depends on your service)
	total := int64(len(users)) // This should be replaced with actual count from service

	// Prepare response
	response := ListUsersResponse{
		Data: users,
	}
	response.Pagination.Page = page
	response.Pagination.Limit = limit
	response.Pagination.Total = total

	// Log successful response
	logger.Info("ListUsers completed",
		zap.String("request_id", requestID),
		zap.Int64("page", page),
		zap.Int64("limit", limit),
		zap.Int("users_count", len(users)),
		zap.Int64("total_users", total),
		zap.Duration("db_duration", dbDuration),
		zap.Duration("total_duration", time.Since(start)),
	)

	_ = utils.SendSuccess(w, "SUCCESS", "Users retrieved successfully", http.StatusOK, response)
}
