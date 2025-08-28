package database

import (
	"context"
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
		clientOptions := options.Client().ApplyURI(cfg.MongoURI)
		
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
		defer cancel()

		client, connectErr := mongo.Connect(ctx, clientOptions)
		if connectErr != nil {
			err = connectErr
			return
		}

		// Ping the database to verify the connection
		if pingErr := client.Ping(ctx, nil); pingErr != nil {
			err = pingErr
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
