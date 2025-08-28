package repository

import (
	"context"
	"time"

	"api-ptf-core-business-orchestrator-go-ms/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindAll(ctx context.Context, page, limit int64) ([]domain.User, error)
	Create(ctx context.Context, user *domain.User) (string, error)
	Update(ctx context.Context, id string, user *domain.User) error
	Delete(ctx context.Context, id string) error
}

// MongoUserRepository is the MongoDB implementation of UserRepository
type MongoUserRepository struct {
	repo *BaseRepository[domain.User]
}

// NewMongoUserRepository creates a new MongoDB user repository
func NewMongoUserRepository(collection *mongo.Collection) *MongoUserRepository {
	return &MongoUserRepository{
		repo: NewBaseRepository[domain.User](collection),
	}
}

// Create inserts a new user into the database
func (r *MongoUserRepository) Create(ctx context.Context, user *domain.User) (string, error) {
	// Set timestamps
	now := time.Now()
	user.DateCreated = now
	user.DateUpdated = now

	// Insert the document
	result, err := r.repo.collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	// Convert the inserted ID to string
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex(), nil
	}

	return "", nil
}

// FindByID finds a user by ID
func (r *MongoUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	return r.repo.FindByID(ctx, id)
}

// FindByEmail finds a user by email address
func (r *MongoUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.repo.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // User not found
		}
		return nil, err
	}
	return &user, nil
}

// FindAll retrieves all users with pagination
func (r *MongoUserRepository) FindAll(ctx context.Context, page, limit int64) ([]domain.User, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	opts := options.Find()
	skip := (page - 1) * limit
	opts.SetSkip(skip)
	opts.SetLimit(limit)

	cursor, err := r.repo.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []domain.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// Update updates an existing user
func (r *MongoUserRepository) Update(ctx context.Context, id string, user *domain.User) error {
	user.DateUpdated = time.Now()
	
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": user,
	}

	_, err = r.repo.collection.UpdateByID(ctx, objID, update)
	return err
}

// Delete removes a user from the database
func (r *MongoUserRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.repo.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

// Count returns the total number of users in the database
func (r *MongoUserRepository) Count(ctx context.Context) (int64, error) {
	return r.repo.collection.CountDocuments(ctx, bson.M{})
}
