package postgrest

import (
	"api-ptf-core-business-orchestrator-go-ms/internal/config"
	"api-ptf-core-business-orchestrator-go-ms/internal/pkg/logger"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func InitDB(cfg *config.Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.App.Postgres.User,
		cfg.App.Postgres.Password,
		cfg.App.Postgres.Host,
		cfg.App.Postgres.Port,
		cfg.App.Postgres.DBName,
		cfg.App.Postgres.SSLMode,
	)

	config, err := pgxpool.ParseConfig(dsn)

	if err != nil {
		logger.Log.Error("❌ error al parsear la configuración", zap.Error(err))
		return nil, fmt.Errorf("pgxpool.ParseConfig: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)

	if err != nil {
		logger.Log.Error("❌ error al conectar a la BD", zap.Error(err))
		return nil, fmt.Errorf("pgxpool.NewWithConfig: %w", err)
	}

	// Verificamos conexión
	if err := pool.Ping(context.Background()); err != nil {
		logger.Log.Error("❌ error al hacer ping a la BD", zap.Error(err))
		return nil, fmt.Errorf("pool.Ping: %w", err)
	}

	logger.Log.Info("✅ Conexión a PostgreSQL establecida con pgxpool")
	validatePostgresConnection(pool)
	return pool, nil
}

// Función para cerrar la conexión al terminar la app
func CloseDB(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
		logger.Log.Info("🔒 Pool de conexiones cerrado")
	}
}

func validatePostgresConnection(pool *pgxpool.Pool) {
	var version string
	err := pool.QueryRow(context.Background(), "SELECT version();").Scan(&version)
	if err != nil {
		logger.Log.Error("❌ Error ejecutando query: %v", zap.Error(err))
	}

	logger.Log.Info("Versión de PostgreSQL:", zap.String("version", version))
}
