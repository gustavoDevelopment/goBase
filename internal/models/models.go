package models

import (
	"api-ptf-core-business-orchestrator-go-ms/internal/config"
	mongoDb "api-ptf-core-business-orchestrator-go-ms/internal/infrastructure/database/mongo"
	"database/sql"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	cfg     *config.Config
	mongoDb *mongoDb.Database
	pgPool  *pgxpool.Pool
	oracle  *sql.DB
}

// NewApplication creates a new Application instance with the provided dependencies
func NewApplication(cfg *config.Config, db *mongoDb.Database, pgPool *pgxpool.Pool, oracle *sql.DB) *Application {
	return &Application{
		cfg:     cfg,
		mongoDb: db,
		pgPool:  pgPool,
		oracle:  oracle,
	}
}

// NewEmptyApplication creates a new empty Application instance
func NewEmptyApplication() *Application {
	return &Application{}
}

// PostgreSQLPool returns the PostgreSQL connection pool instance
func (a *Application) PostgreSQLPool() *pgxpool.Pool {
	return a.pgPool
}

// OraclePool returns the Oracle connection pool instance
func (a *Application) OraclePool() *sql.DB {
	return a.oracle
}

// DB returns the database instance
func (a *Application) MongoDB() *mongoDb.Database {
	return a.mongoDb
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
func (a *Application) SetDB(db *mongoDb.Database) {
	a.mongoDb = db
}
