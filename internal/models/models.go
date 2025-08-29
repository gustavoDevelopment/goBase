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
