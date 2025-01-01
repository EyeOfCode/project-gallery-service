package service

import (
	"context"
	"go-fiber-api/internal/model"
	"go-fiber-api/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserService struct {
    userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
    return &UserService{
        userRepo: userRepo,
    }
}
func (s *UserService) FindByEmail(ctx context.Context, email string) ([]model.User, error) {
    user, err := s.userRepo.FindByEmail(ctx, email)
    if err != nil {
        return nil, err
    }
    return []model.User{*user}, nil
}

func (s *UserService) Count(ctx context.Context, query bson.D) (int64, error) {
    return s.userRepo.Count(ctx, query)
}

func (s *UserService) FindAll(ctx context.Context, query bson.D, opts *options.FindOptions) ([]model.User, error) {
    return s.userRepo.FindAll(ctx, query, opts)
}