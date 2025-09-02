package repository

import (
	mongoDb "api-ptf-core-business-orchestrator-go-ms/internal/infrastructure/database/mongo"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GenericRepository es un repositorio genérico para operaciones CRUD
// con cualquier tipo de entidad que implemente la interfaz Model
// y tenga los campos base de fecha de creación y actualización
type GenericRepository[T any] struct {
	collection *mongo.Collection
}

// NewGenericRepository crea una nueva instancia de GenericRepository
func NewGenericRepository[T any](db *mongoDb.Database, collectionName string) *GenericRepository[T] {
	return &GenericRepository[T]{
		collection: db.GetCollection(collectionName),
	}
}

// Create inserta un nuevo documento en la colección
func (r *GenericRepository[T]) Create(ctx context.Context, entity *T) (string, error) {
	// Agregar timestamps
	timestamp := time.Now()
	entityMap, err := toMap(entity)
	if err != nil {
		return "", err
	}

	entityMap["date_created"] = timestamp

	result, err := r.collection.InsertOne(ctx, entityMap)
	if err != nil {
		return "", err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex(), nil
	}

	return "", nil
}

// FindByID busca un documento por su ID
func (r *GenericRepository[T]) FindByID(ctx context.Context, id string) (*T, error) {
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

// Update actualiza un documento existente
func (r *GenericRepository[T]) Update(ctx context.Context, id string, entity *T) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	entityMap, err := toMap(entity)
	if err != nil {
		return err
	}

	// Actualizar la fecha de modificación
	entityMap["date_updated"] = time.Now()

	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": entityMap},
	)

	return err
}

// Delete elimina un documento por su ID
func (r *GenericRepository[T]) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

// FindAll obtiene todos los documentos con paginación
func (r *GenericRepository[T]) FindAll(ctx context.Context, page, limit int64) ([]T, error) {
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

// Count devuelve el número total de documentos en la colección
func (r *GenericRepository[T]) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}
