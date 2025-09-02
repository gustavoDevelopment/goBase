package oracle

import (
	"api-ptf-core-business-orchestrator-go-ms/internal/config"
	"api-ptf-core-business-orchestrator-go-ms/internal/pkg/logger"
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/godror/godror"
	"go.uber.org/zap"
)

// InitDB inicializa la conexión a Oracle con godror
func InitDB(cfg *config.Config) (*sql.DB, error) {
	// DSN: usuario/contraseña@host:puerto/servicio
	dsn := fmt.Sprintf("%s/%s@%s:%d/%s",
		cfg.App.Oracle.User,
		cfg.App.Oracle.Password,
		cfg.App.Oracle.Host,
		cfg.App.Oracle.Port,
		cfg.App.Oracle.DBName,
	)

	db, err := sql.Open("godror", dsn)
	if err != nil {
		logger.Log.Error("❌ error al abrir conexión con Oracle", zap.Error(err))
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	// Configuración del pool
	db.SetMaxOpenConns(50)                  // Máx conexiones abiertas
	db.SetMaxIdleConns(10)                  // Conexiones en idle
	db.SetConnMaxLifetime(30 * time.Minute) // Tiempo máximo de vida de una conexión

	// Verificamos conexión con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		logger.Log.Error("❌ error al hacer ping a Oracle", zap.Error(err))
		return nil, fmt.Errorf("db.PingContext: %w", err)
	}

	logger.Log.Info("✅ Conexión a Oracle establecida")
	validateOracleConnection(db)
	return db, nil
}

// CloseDB cierra el pool de conexiones
func CloseDB(db *sql.DB) {
	if db != nil {
		if err := db.Close(); err != nil {
			logger.Log.Error("⚠️ error cerrando la conexión a Oracle", zap.Error(err))
			return
		}
		logger.Log.Info("🔒 Pool de conexiones Oracle cerrado")
	}
}

// validateOracleConnection ejecuta una query de validación
func validateOracleConnection(db *sql.DB) {
	var sysdate string
	err := db.QueryRow("SELECT TO_CHAR(SYSDATE, 'YYYY-MM-DD HH24:MI:SS') FROM dual").Scan(&sysdate)
	if err != nil {
		logger.Log.Error("❌ Error ejecutando query de validación", zap.Error(err))
		return
	}
	logger.Log.Info("🕒 Oracle SYSDATE:", zap.String("sysdate", sysdate))
}
