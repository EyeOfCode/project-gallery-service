package service

import (
	"context"
	"go-fiber-api/internal/model"
	"go-fiber-api/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ShopService struct {
	shopRepo repository.ShopRepository
}

func NewShopService(shopRepo repository.ShopRepository) *ShopService {
	return &ShopService{
		shopRepo: shopRepo,
	}
}
func (s *ShopService) FindAll(ctx context.Context, query bson.D, opts *options.FindOptions) ([]model.Shop, error) {
    return s.shopRepo.FindAll(ctx, query, opts)
}

func (s *ShopService) Count(ctx context.Context, query bson.D) (int64, error) {
	return s.shopRepo.Count(ctx, query)
}