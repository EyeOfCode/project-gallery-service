package repository

import (
	"context"
	"go-fiber-api/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ShopRepository interface {
	Create(ctx context.Context, shop *model.Shop) error
	FindOne(ctx context.Context, query bson.M) (*model.Shop, error)
	FindAll(ctx context.Context, query bson.D, opts *options.FindOptions) ([]model.Shop, error)
	Count(ctx context.Context, query bson.D) (int64, error)
}

type shopRepository struct {
	collection *mongo.Collection
}

func NewShopRepository(db *mongo.Database) ShopRepository {
	return &shopRepository{
		collection: db.Collection("shops"),
	}
}

func (r *shopRepository) Create(ctx context.Context, shop *model.Shop) error {
	shop.ID = primitive.NewObjectID()
	shop.CreatedAt = time.Now()
	shop.UpdatedAt = time.Now()
	
	_, err := r.collection.InsertOne(ctx, shop)
	return err
}

func (r *shopRepository) FindOne(ctx context.Context, query bson.M) (*model.Shop, error) {
	var shop model.Shop
	err := r.collection.FindOne(ctx, query).Decode(&shop)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil 
		}
		return nil, err
	}
	return &shop, nil
}

func (r *shopRepository) FindAll(ctx context.Context, query bson.D, opts *options.FindOptions) ([]model.Shop, error) {
	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var shops []model.Shop
	if err := cursor.All(ctx, &shops); err != nil {
		return nil, err
	}
	return shops, nil
}

func (r *shopRepository) Count(ctx context.Context, query bson.D) (int64, error) {
	return r.collection.CountDocuments(ctx, query)
}