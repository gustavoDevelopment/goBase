package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"api-ptf-core-business-orchestrator-go-ms/internal/application"
	"api-ptf-core-business-orchestrator-go-ms/internal/config"
	"api-ptf-core-business-orchestrator-go-ms/internal/infrastructure/database"
	"api-ptf-core-business-orchestrator-go-ms/internal/infrastructure/repository"
)

func main() {
	// Cargar configuración
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Inicializar logger
	logger := log.New(os.Stdout, "[APP] ", log.LstdFlags)

	// Configurar contexto para manejo de señales de terminación
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Configurar manejo de señales
	setupGracefulShutdown(cancel)

	// Inicializar conexión a MongoDB
	db, err := database.NewDatabase(cfg)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Disconnect(ctx); err != nil {
			logger.Printf("Error disconnecting from database: %v", err)
		}
	}()

	logger.Println("Successfully connected to MongoDB")

	// Inicializar repositorios
	userRepo := repository.NewMongoUserRepository(db.GetCollection(cfg.MongoCollection))

	// Inicializar servicios
	userService := application.NewUserService(userRepo)

	// Aquí iría la inicialización del servidor HTTP, colas, etc.
	// Por ahora, solo mostramos un mensaje de éxito
	logger.Println("Application started successfully")

	// Ejemplo de uso del servicio
	exampleUser, err := userService.GetUserByID(ctx, 1)
	if err != nil {
		logger.Printf("Error getting user: %v", err)
	} else {
		logger.Printf("Found user: %+v", exampleUser)
	}

	// Mantener la aplicación en ejecución hasta recibir señal de terminación
	<-ctx.Done()
	logger.Println("Shutting down...")
}

// setupGracefulShutdown configura el manejo de señales para un apagado controlado
func setupGracefulShutdown(cancelFunc context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal: %v", sig)
		cancelFunc()
	}()
}
