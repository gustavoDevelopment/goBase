package database

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"sync"
	"time"

	"api-ptf-core-business-orchestrator-go-ms/internal/config"
	"api-ptf-core-business-orchestrator-go-ms/internal/pkg/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const (
	defaultMinPoolSize = 5
	defaultMaxPoolSize = 100
	defaultTimeout     = "10s"
)

type Database struct {
	client *mongo.Client
	db     *mongo.Database
	config *config.Config
	mu     sync.RWMutex
}

var (
	once        sync.Once
	instance    *Database
	instanceErr error
)

// NewDatabase creates a new database connection instance
func NewDatabase(cfg *config.Config) (*Database, error) {
	once.Do(func() {
		mongoCfg := cfg.App.MongoDB
		if mongoCfg.URI == "" {
			instanceErr = errors.New("MongoDB URI is not configured")
			return
		}

		// Parse timeout duration
		timeout, err := time.ParseDuration(mongoCfg.Timeout)
		if err != nil {
			timeout, _ = time.ParseDuration(defaultTimeout)
		}

		logger.Log.Info("Connecting to MongoDB...",
			zap.String("uri", mongoCfg.URI),
			zap.String("database", mongoCfg.Database),
			zap.String("timeout", mongoCfg.Timeout),
		)
		clientOptions := options.Client().
			ApplyURI(mongoCfg.URI).
			SetMinPoolSize(defaultMinPoolSize).
			SetMaxPoolSize(defaultMaxPoolSize).
			SetServerSelectionTimeout(timeout).
			SetConnectTimeout(timeout).
			SetTLSConfig(&tls.Config{
				InsecureSkipVerify: true, // Keep TLS verification disabled as per original
			})

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			instanceErr = fmt.Errorf("failed to connect to MongoDB: %w", err)
			return
		}

		// Verify the connection with timeout
		pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer pingCancel()

		if err := client.Ping(pingCtx, nil); err != nil {
			_ = client.Disconnect(context.Background())
			instanceErr = fmt.Errorf("failed to ping MongoDB: %w", err)
			return
		}

		instance = &Database{
			client: client,
			db:     client.Database(mongoCfg.Database),
			config: cfg,
		}
	})

	if instanceErr != nil {
		return nil, instanceErr
	}

	if instance == nil {
		return nil, errors.New("database instance is nil")
	}

	return instance, nil
}

// getUintConfig safely gets a uint value from config with a default fallback
func getUintConfig(value, defaultValue uint) uint {
	if value == 0 {
		return defaultValue
	}
	return value
}

// GetCollection returns a handle for a specific collection
func (d *Database) GetCollection(name string) *mongo.Collection {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.db.Collection(name)
}

// Disconnect closes the database connection
func (d *Database) Disconnect(ctx context.Context) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.client != nil {
		return d.client.Disconnect(ctx)
	}
	return nil
}

// GetClient returns the underlying MongoDB client
func (d *Database) GetClient() *mongo.Client {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.client
}
