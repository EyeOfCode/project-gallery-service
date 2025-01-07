package service

import (
	"context"
	"go-fiber-api/internal/model"
	"go-fiber-api/internal/repository"
	"go-fiber-api/pkg/dto"
	"go-fiber-api/pkg/utils"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
    userRepo    repository.UserRepository
    redisClient *redis.Client
}

func NewUserService(userRepo repository.UserRepository, redisClient *redis.Client) *UserService {
    return &UserService{
        userRepo: userRepo,
        redisClient: redisClient,
    }
}
func (s *UserService) FindByEmail(ctx context.Context, email string) (*model.User, error) {
    user, err := s.userRepo.FindOne(ctx, bson.M{"email": email})
    if err != nil || user == nil {
        return nil, err
    }
    return user, nil
}

func (s *UserService) Count(ctx context.Context, query bson.D) (int64, error) {
    return s.userRepo.Count(ctx, query)
}

func (s *UserService) FindAll(ctx context.Context, query bson.D, opts *options.FindOptions) ([]model.User, error) {
    return s.userRepo.FindAll(ctx, query, opts)
}

func (s *UserService) Create(ctx context.Context, payload *dto.RegisterRequest) (*model.User, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    if len(payload.Roles) == 0 {
        payload.Roles = []string{string(utils.UserRole)}
    }

    now := time.Now()
    user := &model.User{
        ID:        primitive.NewObjectID(),
        Name:      payload.Name,
        Email:     payload.Email,
        Password:  string(hashedPassword),
        Roles:     payload.Roles,
        CreatedAt: now,
        UpdatedAt: now,
    }
    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }
    return user, nil
}

func (s *UserService) Login(ctx context.Context, password string, user *model.User) (*string, error) {
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return nil, err
    }

    auth := utils.NewAuthHandler(os.Getenv("JWT_SECRET"), os.Getenv("JWT_EXPIRY"))
    token, err := auth.GenerateToken(user.ID.Hex(), user.Roles)
    if err != nil {
        return nil, err
    }
    return &token, nil
}

func (s *UserService) FindByID(ctx context.Context, id string) (*model.User, error) {
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }
    return s.userRepo.FindOne(ctx, bson.M{"_id": objID})
}
func (s *UserService) UpdateById(ctx context.Context, id primitive.ObjectID, payload *dto.UpdateUserRequest) error {
    _, err := s.userRepo.UpdateByID(ctx, id, payload)
    return err
}

func (s *UserService) Delete(ctx context.Context, id primitive.ObjectID) error {
    return s.userRepo.Delete(ctx, id)
}
