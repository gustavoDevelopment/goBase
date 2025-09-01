package models

import (
	"api-ptf-core-business-orchestrator-go-ms/internal/config"
	"api-ptf-core-business-orchestrator-go-ms/internal/infrastructure/database"
)

type Application struct {
	cfg *config.Config
	db  *database.Database
}

// NewApplication creates a new Application instance with the provided dependencies
func NewApplication(cfg *config.Config, db *database.Database) *Application {
	return &Application{
		cfg: cfg,
		db:  db,
	}
}

// NewEmptyApplication creates a new empty Application instance
func NewEmptyApplication() *Application {
	return &Application{}
}

// DB returns the database instance
func (a *Application) MongoDB() *database.Database {
	return a.db
}

// DB returns the database instance
func (a *Application) Configs() *config.Config {
	return a.cfg
}

// SetConfig sets the configuration instance
func (a *Application) SetConfig(cfg *config.Config) {
	a.cfg = cfg
}

// SetDB sets the database instance
func (a *Application) SetDB(db *database.Database) {
	a.db = db
}
