package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"api-ptf-core-business-orchestrator-go-ms/internal/application"
	"api-ptf-core-business-orchestrator-go-ms/internal/config"
	"api-ptf-core-business-orchestrator-go-ms/internal/infrastructure/database"
	"api-ptf-core-business-orchestrator-go-ms/internal/infrastructure/repository"
	"api-ptf-core-business-orchestrator-go-ms/internal/pkg/logger"
)

// Application contiene las dependencias de la aplicación
type Application struct {
	cfg         *config.Config
	db          *database.Database
	userService *application.UserService
}

func main() {
	// Inicialización básica
	if err := initializeApplication(); err != nil {
		logger.Log.Fatal("Failed to initialize application", zap.Error(err))
	}
	defer logger.Sync()

	// Configuración del contexto
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Configuración de señales
	setupGracefulShutdown(cancel)

	// Inicialización de la aplicación
	app, err := initializeApp(ctx)
	if err != nil {
		logger.Log.Fatal("Failed to initialize application services", zap.Error(err))
	}
	defer app.cleanup(ctx)

	// Iniciar la aplicación
	if err := app.run(ctx); err != nil {
		logger.Log.Fatal("Application error", zap.Error(err))
	}
}

// run inicia la ejecución de la aplicación
func (a *Application) run(ctx context.Context) error {
	logger.Log.Info("Application started successfully")

	// Ejemplo de uso del servicio
	logger.Log.Debug("Fetching example user...")
	exampleUser, err := a.userService.GetUserByID(ctx, 1)
	if err != nil {
		logger.Log.Error("Error getting user", zap.Error(err))
	} else {
		logger.Log.Info("Found user", zap.Any("user", exampleUser))
	}

	// Mantener la aplicación en ejecución hasta recibir señal de terminación
	<-ctx.Done()
	logger.Log.Info("Shutting down...")
	return nil
}

// cleanup realiza la limpieza de recursos de la aplicación
func (a *Application) cleanup(ctx context.Context) {
	if a.db != nil {
		logger.Log.Info("Disconnecting from database...")
		if err := a.db.Disconnect(ctx); err != nil {
			logger.Log.Error("Error disconnecting from database", zap.Error(err))
		}
	}
}

// initializeApplication realiza la inicialización básica de la aplicación
func initializeApplication() error {
	// Inicializar logger
	if err := logger.InitLogger(true); err != nil { // true para modo desarrollo
		return err
	}
	return nil
}

// initializeApp inicializa y configura la aplicación
func initializeApp(ctx context.Context) (*Application, error) {
	logger.Log.Info("Starting application initialization...")

	// Cargar configuración
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		return nil, err
	}

	// Inicializar conexión a MongoDB
	logger.Log.Info("Connecting to MongoDB...",
		zap.String("uri", cfg.MongoURI),
		zap.String("database", cfg.MongoDB),
	)

	db, err := database.NewDatabase(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create database client: %w", err)
	}

	// Create a new context with timeout for the ping operation
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Verify database connection with retry logic
	var pingErr error
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		if pingErr = db.GetClient().Ping(pingCtx, nil); pingErr == nil {
			break
		}
		if i < maxRetries-1 {
			time.Sleep(time.Duration(i+1) * time.Second) // Exponential backoff
		}
	}

	if pingErr != nil {
		return nil, fmt.Errorf("failed to ping database after %d attempts: %w", maxRetries, pingErr)
	}

	logger.Log.Info("Successfully connected to MongoDB")

	// Inicializar repositorios
	userRepo := repository.NewMongoUserRepository(db.GetCollection(cfg.MongoCollection))

	// Inicializar servicios
	userService := application.NewUserService(userRepo)

	return &Application{
		cfg:         cfg,
		db:          db,
		userService: userService,
	}, nil
}

// setupGracefulShutdown configura el manejo de señales para un apagado controlado
func setupGracefulShutdown(cancelFunc context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		logger.Log.Info("Received signal, shutting down...",
			zap.String("signal", sig.String()))
		cancelFunc()
	}()
}
