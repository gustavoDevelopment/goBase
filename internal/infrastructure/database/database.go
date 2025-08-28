package database

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"

	"api-ptf-core-business-orchestrator-go-ms/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	client *mongo.Client
	db     *mongo.Database
	config *config.Config
}

var (
	once     sync.Once
	instance *Database
)

// NewDatabase creates a new database connection instance
func NewDatabase(cfg *config.Config) (*Database, error) {
	var err error
	once.Do(func() {
		// Parse the connection string
		clientOptions := options.Client().ApplyURI(cfg.MongoURI)

		// Configure TLS settings
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true, // Skip certificate verification
		}

		// Apply TLS configuration
		clientOptions.SetTLSConfig(tlsConfig)
		clientOptions.SetServerSelectionTimeout(cfg.Timeout)
		clientOptions.SetConnectTimeout(cfg.Timeout)

		// Create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
		defer cancel()

		// Connect to MongoDB
		client, connectErr := mongo.Connect(ctx, clientOptions)
		if connectErr != nil {
			err = fmt.Errorf("failed to connect to MongoDB: %w", connectErr)
			return
		}

		// Verify the connection
		if pingErr := client.Ping(ctx, nil); pingErr != nil {
			err = fmt.Errorf("failed to ping MongoDB: %w", pingErr)
			return
		}

		instance = &Database{
			client: client,
			db:     client.Database(cfg.MongoDB),
			config: cfg,
		}
	})

	if err != nil {
		return nil, err
	}

	return instance, nil
}

// GetCollection returns a handle for a specific collection
func (d *Database) GetCollection(name string) *mongo.Collection {
	return d.db.Collection(name)
}

// Disconnect closes the database connection
func (d *Database) Disconnect(ctx context.Context) error {
	if d.client != nil {
		return d.client.Disconnect(ctx)
	}
	return nil
}

// GetClient returns the underlying MongoDB client
func (d *Database) GetClient() *mongo.Client {
	return d.client
}
