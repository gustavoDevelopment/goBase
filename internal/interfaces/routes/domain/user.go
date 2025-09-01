package domainRoutes

import (
	"api-ptf-core-business-orchestrator-go-ms/internal/application"
	"api-ptf-core-business-orchestrator-go-ms/internal/infrastructure/repository"
	"api-ptf-core-business-orchestrator-go-ms/internal/interfaces/http/handlers"
	"api-ptf-core-business-orchestrator-go-ms/internal/models"

	"github.com/gorilla/mux"
)

func RegisterUserRoutes(router *mux.Router, a *models.Application) {
	// Inicializar repositorios
	userRepo := repository.NewMongoUserRepository(a.MongoDB())

	// Inicializar servicios
	userService := application.NewUserService(userRepo.Repo, userRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)

	// User routes
	userRouter := router.PathPrefix("/users").Subrouter()
	userRouter.HandleFunc("", userHandler.ListUsers).Methods("GET")
	userRouter.HandleFunc("", userHandler.CreateUser).Methods("POST")
	userRouter.HandleFunc("/{id}", userHandler.GetUserByID).Methods("GET")
}
