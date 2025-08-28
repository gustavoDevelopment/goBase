package handlers

import (
	"api-ptf-core-business-orchestrator-go-ms/internal/application"
	"api-ptf-core-business-orchestrator-go-ms/internal/domain"
	"api-ptf-core-business-orchestrator-go-ms/internal/interfaces/http/utils"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

// Common error messages
var (
	ErrInvalidUserID    = errors.New("invalid user ID format")
	ErrInvalidPageParam = errors.New("invalid page parameter")
	ErrInvalidLimit     = errors.New("limit must be between 1 and 100")
	ErrInvalidEmail     = errors.New("valid email is required")
)

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
	vars := mux.Vars(r)
	userID := strings.TrimSpace(vars["id"])

	if userID == "" {
		utils.NotFound(w, "User not found")
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			_ = utils.NotFound(w, "User not found")
		} else {
			_ = utils.InternalServerError(w, "Failed to retrieve user: "+err.Error())
		}
		return
	}

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
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		_ = utils.SendError(w, http.StatusUnprocessableEntity, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Basic validation
	if user.Email == "" {
		_ = utils.BadRequest(w, "Email is required")
		return
	}

	createdUser, err := h.userService.CreateUser(r.Context(), &user)
	if err != nil {
		_ = utils.InternalServerError(w, "Failed to create user: "+err.Error())
		return
	}

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
	// Parse pagination parameters
	page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	switch {
	case err != nil || limit < 1:
		limit = 10
	case limit > 100:
		_ = utils.BadRequest(w, "Limit must be between 1 and 100")
		return
	}

	// Get users
	users, err := h.userService.ListUsers(r.Context(), page, limit)
	if err != nil {
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

	_ = utils.SendSuccess(w, "SUCCESS", "Users retrieved successfully", http.StatusOK, response)
}
