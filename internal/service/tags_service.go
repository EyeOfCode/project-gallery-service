package service

import (
	"context"
	"pre-test-gallery-service/internal/model"
	"pre-test-gallery-service/internal/repository"
	"pre-test-gallery-service/pkg/dto"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TagsService struct {
	tagsRepo repository.TagsRepository
}

func NewTagsService(tagsRepo repository.TagsRepository) *TagsService {
	return &TagsService{
		tagsRepo: tagsRepo,
	}
}

func (s *TagsService) GetAllTags(ctx context.Context) ([]model.Tags, error) {
	return s.tagsRepo.FindAll(ctx, bson.M{})
}

func (s *TagsService) CreateTags(ctx context.Context, req dto.TagsRequest) (*model.Tags, error) {
	now := time.Now()
    tag := &model.Tags{
        ID:        primitive.NewObjectID(),
        Name:      req.Name,
        CreatedAt: now,
        UpdatedAt: now,
    }
	if err := s.tagsRepo.Create(ctx, tag); err != nil {
        return nil, err
    }
	return tag, nil
}

func (s *TagsService) FindOneTags(ctx context.Context, query bson.M) (*model.Tags, error) {
	return s.tagsRepo.FindOne(ctx, query)
}

func (s *TagsService) DeleteTags(ctx context.Context, id primitive.ObjectID) error {
	return s.tagsRepo.Delete(ctx, id)
}