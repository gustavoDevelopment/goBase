package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BaseRepository defines the common CRUD operations for MongoDB collections
type BaseRepository[T any] struct {
	collection *mongo.Collection
}

// NewBaseRepository creates a new instance of BaseRepository
func NewBaseRepository[T any](collection *mongo.Collection) *BaseRepository[T] {
	return &BaseRepository[T]{
		collection: collection,
	}
}

// Create inserts a new document into the collection
func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) (string, error) {
	// Add timestamps
	timestamp := time.Now()
	entityMap, err := toMap(entity)
	if err != nil {
		return "", err
	}

	entityMap["date_created"] = timestamp
	entityMap["updated_created"] = timestamp

	result, err := r.collection.InsertOne(ctx, entityMap)
	if err != nil {
		return "", err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex(), nil
	}

	return "", nil
}

// FindByID finds a document by its ID
func (r *BaseRepository[T]) FindByID(ctx context.Context, id string) (*T, error) {
	var entity T
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&entity)
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

// Update updates a document by ID
func (r *BaseRepository[T]) Update(ctx context.Context, id string, entity *T) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	entityMap, err := toMap(entity)
	if err != nil {
		return err
	}

	// Update the updated_created timestamp
	entityMap["updated_created"] = time.Now()

	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": entityMap},
	)

	return err
}

// Delete removes a document by ID
func (r *BaseRepository[T]) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

// FindAll retrieves all documents with pagination
func (r *BaseRepository[T]) FindAll(ctx context.Context, page, limit int64) ([]T, error) {
	var results []T

	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// Count returns the total number of documents in the collection
func (r *BaseRepository[T]) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}

// toMap converts a struct to a map[string]interface{}
func toMap(entity interface{}) (map[string]interface{}, error) {
	data, err := bson.Marshal(entity)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = bson.Unmarshal(data, &result)
	return result, err
}
