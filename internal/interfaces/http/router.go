package http

import (
	"api-ptf-core-business-orchestrator-go-ms/internal/application"
	"api-ptf-core-business-orchestrator-go-ms/internal/config"
	"api-ptf-core-business-orchestrator-go-ms/internal/interfaces/http/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter creates a new HTTP router with all the routes
func NewRouter(cfg *config.Config, userService *application.UserService) *mux.Router {
	r := mux.NewRouter()

	// Use the base path from config
	api := r.PathPrefix(cfg.HTTP.BasePath).Subrouter()

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)

	// User routes
	userRouter := api.PathPrefix("/users").Subrouter()
	userRouter.HandleFunc("", userHandler.ListUsers).Methods("GET")
	userRouter.HandleFunc("", userHandler.CreateUser).Methods("POST")
	userRouter.HandleFunc("/{id}", userHandler.GetUserByID).Methods("GET")

	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Add middleware
	r.Use(loggingMiddleware)
	r.Use(mux.CORSMethodMiddleware(r))

	return r
}

// loggingMiddleware logs the HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log request details
		// You can use your logger here
		// logger.Log.Info("Request", zap.String("method", r.Method), zap.String("path", r.URL.Path))
		next.ServeHTTP(w, r)
	})
}
