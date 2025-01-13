package repository

import (
	"context"
	"pre-test-gallery-service/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TagsRepository interface {
	FindAll(ctx context.Context, query bson.M) ([]model.Tags, error)
	FindOne(ctx context.Context, query bson.M) (*model.Tags, error)
	Create(ctx context.Context, tags *model.Tags) ([]*model.Tags, error)
	Delete(ctx context.Context, tags *model.Tags, id primitive.ObjectID) error
}

type tagsRepository struct {
	collection *mongo.Collection
}

func NewTagsRepository(db *mongo.Database) TagsRepository {
	return &tagsRepository{
		collection: db.Collection("tags"),
	}
}

func (r *tagsRepository) FindAll(ctx context.Context, query bson.M) ([]model.Tags, error) {
	var tags []model.Tags
	cursor, err := r.collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &tags); err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *tagsRepository) FindOne(ctx context.Context, query bson.M) (*model.Tags, error) {
	var result model.Tags
	err := r.collection.FindOne(ctx, query).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *tagsRepository) Create(ctx context.Context, tags *model.Tags) ([]*model.Tags, error) {
	if tags == nil {
		return nil, nil
	}
	result, err := r.collection.InsertOne(ctx, tags)
	if err != nil {
		return nil, err
	}
	tags.ID = result.InsertedID.(primitive.ObjectID)
	return []*model.Tags{tags}, nil
}

func (r *tagsRepository) Delete(ctx context.Context, tags *model.Tags, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}