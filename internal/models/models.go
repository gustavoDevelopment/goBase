package models

import (
	"api-ptf-core-business-orchestrator-go-ms/internal/application"
	"api-ptf-core-business-orchestrator-go-ms/internal/config"
	"api-ptf-core-business-orchestrator-go-ms/internal/infrastructure/database"
)

type Application struct {
	cfg         *config.Config
	db          *database.Database
	userService *application.UserService
}

// NewApplication creates a new Application instance with the provided dependencies
func NewApplication(cfg *config.Config, db *database.Database) *Application {
	return &Application{
		cfg: cfg,
		db:  db,
	}
}

// DB returns the database instance
func (a *Application) DB() *database.Database {
	return a.db
}
