package repository

import (
	"context"
	"api-ptf-core-business-orchestrator-go-ms/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	FindByID(ctx context.Context, id int64) (*domain.User, error)
}

// MongoUserRepository implements UserRepository with MongoDB
type MongoUserRepository struct {
	collection *mongo.Collection
}

// NewMongoUserRepository creates a new MongoDB user repository
func NewMongoUserRepository(collection *mongo.Collection) *MongoUserRepository {
	return &MongoUserRepository{
		collection: collection,
	}
}

// FindByID finds a user by ID
func (r *MongoUserRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	var user domain.User
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // User not found
		}
		return nil, err
	}
	return &user, nil
}
