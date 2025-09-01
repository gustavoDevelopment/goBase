package app

import (
	"context"
	"fmt"

	"api-ptf-core-business-orchestrator-go-ms/internal/models"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"api-ptf-core-business-orchestrator-go-ms/internal/config"
	"api-ptf-core-business-orchestrator-go-ms/internal/infrastructure/database"
	httpServer "api-ptf-core-business-orchestrator-go-ms/internal/interfaces/http"
	"api-ptf-core-business-orchestrator-go-ms/internal/pkg/logger"

	"go.uber.org/zap"
)

type applicationWrapper struct {
	*models.Application
}

// run starts the HTTP server and keeps it running until a shutdown signal is received
func (aw *applicationWrapper) run(ctx context.Context) error {
	// Create HTTP server with our router
	router := httpServer.NewRouter(aw.Application)
	srv := &http.Server{
		Addr:    ":" + aw.Configs().HTTP.Port,
		Handler: router,
	}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				logger.Log.Fatal("Graceful shutdown timed out. Forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			logger.Log.Error("HTTP server shutdown error", zap.Error(err))
		}
		serverStopCtx()
	}()

	// Start the server
	logger.Log.Info("Starting HTTP server",
		zap.String("context_path", aw.Configs().HTTP.BasePath),
		zap.String("port", aw.Configs().HTTP.Port),
	)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
	logger.Log.Info("Server stopped")
	return nil
}

func StartUp() {
	// Inicializar la aplicación básica (esto inicializa el logger)
	if err := initializeApplication(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
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

// run starts the HTTP server and keeps it running until a shutdown signal is received
func run(ctx context.Context) error {
	// Create HTTP server with our router
	app := models.NewEmptyApplication()
	router := httpServer.NewRouter(app)
	srv := &http.Server{
		Addr:    ":" + app.Configs().HTTP.Port,
		Handler: router,
	}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				logger.Log.Fatal("Graceful shutdown timed out. Forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			logger.Log.Error("HTTP server shutdown error", zap.Error(err))
		}
		serverStopCtx()
	}()

	// Start the server
	logger.Log.Info("Starting HTTP server",
		zap.String("context_path", app.Configs().HTTP.BasePath),
		zap.String("port", app.Configs().HTTP.Port),
	)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
	logger.Log.Info("Server stopped")
	return nil
}

// cleanup realiza la limpieza de recursos de la aplicación
func (a *applicationWrapper) cleanup(ctx context.Context) {
	if a.MongoDB() != nil {
		logger.Log.Info("Disconnecting from database...")
		if err := a.MongoDB().Disconnect(ctx); err != nil {
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
func initializeApp(ctx context.Context) (*applicationWrapper, error) {
	logger.Log.Info("Starting application initialization...")

	// Get the current working directory
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filepath.Dir(filepath.Dir(filename))) // Go up two levels to the project root
	configPath := filepath.Join(dir, "configs", "config.yaml")

	config, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config from %s: %w", configPath, err)
	}

	// Ensure MongoDB URI is properly formatted
	if !strings.HasPrefix(config.MongoURI, "mongodb://") && !strings.HasPrefix(config.MongoURI, "mongodb+srv://") {
		config.MongoURI = "mongodb://" + config.MongoURI
	}

	// Inicializar conexión a MongoDB
	db, err := database.NewDatabase(config)
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

	// Crear la aplicación usando el constructor
	return &applicationWrapper{Application: models.NewApplication(config, db)}, nil
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
