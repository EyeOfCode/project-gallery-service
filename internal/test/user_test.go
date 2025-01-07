package test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"go-fiber-api/internal/config"
	"go-fiber-api/internal/handlers"
	"go-fiber-api/internal/model"
	"go-fiber-api/internal/service"
	"go-fiber-api/pkg/dto"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Mock Implementations
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

func (m *MockUserRepository) FindOne(ctx context.Context, filter bson.M) (*model.User, error) {
    args := m.Called(ctx, filter)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindAll(ctx context.Context, filter bson.D, opts *options.FindOptions) ([]model.User, error) {
    args := m.Called(ctx, filter, opts)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserRepository) UpdateByID(ctx context.Context, id primitive.ObjectID, payload *dto.UpdateUserRequest) (*model.User, error) {
    args := m.Called(ctx, id, payload)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
    args := m.Called(ctx, id)
    return args.Error(0)
}

func (m *MockUserRepository) Count(ctx context.Context, filter bson.D) (int64, error) {
    args := m.Called(ctx, filter)
    return int64(args.Int(0)), args.Error(1)
}

type MockRedisClient struct {
    mock.Mock
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
    args := m.Called(ctx, key, value, expiration)
    return args.Get(0).(*redis.StatusCmd)
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
    args := m.Called(ctx, key)
    return args.Get(0).(*redis.StringCmd)
}

func (m *MockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
    args := m.Called(ctx, keys)
    return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClient) Pipeline() redis.Pipeliner {
    args := m.Called()
    return args.Get(0).(redis.Pipeliner)
}

type MockPipeliner struct {
    mock.Mock
}

func (m *MockPipeliner) Exec(ctx context.Context) ([]redis.Cmder, error) {
    args := m.Called(ctx)
    return args.Get(0).([]redis.Cmder), args.Error(1)
}

func (m *MockPipeliner) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
    args := m.Called(ctx, key, value, expiration)
    return args.Get(0).(*redis.StatusCmd)
}

func (m *MockPipeliner) Del(ctx context.Context, keys ...string) *redis.IntCmd {
    args := m.Called(ctx, keys)
    return args.Get(0).(*redis.IntCmd)
}

// Test Setup
var testConfig = &config.Config{
    JWTSecretKey:  "test-secret",
    JWTRefreshKey: "test-refresh-secret",
    JWTExpiresIn:  "24h",
    JWTRefreshIn:  "168h",
}

// Helper Functions
func setupTestApp(mockRepo *MockUserRepository, mockRedis *MockRedisClient) *fiber.App {
    app := fiber.New()
    userService := service.NewUserService(mockRepo, mockRedis, testConfig)
    userHandler := handlers.NewUserHandler(userService)
    
    auth := app.Group("/auth")
    auth.Post("/register", userHandler.Register)
    auth.Post("/login", userHandler.Login)
    auth.Post("/refresh", userHandler.RefreshToken)
    auth.Get("/logout", userHandler.Logout)
    
    return app
}

func TestLogin(t *testing.T) {
    tests := []struct {
        name           string
        payload        dto.LoginRequest
        setupMock      func(*MockUserRepository, *MockRedisClient)
        expectedStatus int
        expectError    bool
    }{
        {
            name: "Success",
            payload: dto.LoginRequest{
                Email:    "test@example.com",
                Password: "password123",
            },
            setupMock: func(mr *MockUserRepository, mrc *MockRedisClient) {
                user := &model.User{
                    ID:       primitive.NewObjectID(),
                    Email:    "test@example.com",
                    Password: "$2a$10$abcdefghijklmnopqrstuvwxyz", // hashed password
                    Roles:    []string{"user"},
                }
                mr.On("FindOne", mock.Anything, mock.Anything).Return(user, nil)
                
                // Mock Redis Set for token storage
                cmd := redis.NewStatusCmd(context.Background())
                cmd.SetVal("OK")
                mrc.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(cmd)
            },
            expectedStatus: 200,
            expectError:    false,
        },
        {
            name: "Invalid Credentials",
            payload: dto.LoginRequest{
                Email:    "wrong@example.com",
                Password: "wrongpass",
            },
            setupMock: func(mr *MockUserRepository, mrc *MockRedisClient) {
                mr.On("FindOne", mock.Anything, mock.Anything).Return(nil, errors.New("user not found"))
            },
            expectedStatus: 404,
            expectError:    true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            app := fiber.New()
            mockRepo := new(MockUserRepository)
            mockRedis := new(MockRedisClient)
            tt.setupMock(mockRepo, mockRedis)

            userService := service.NewUserService(mockRepo, mockRedis, testConfig)
            userHandler := handlers.NewUserHandler(userService)
            app.Post("/auth/login", userHandler.Login)

            payloadBytes, _ := json.Marshal(tt.payload)
            req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(payloadBytes))
            req.Header.Set("Content-Type", "application/json")

            resp, err := app.Test(req)
            assert.NoError(t, err)
            assert.Equal(t, tt.expectedStatus, resp.StatusCode)

            if !tt.expectError {
                var result struct {
                    Data struct {
                        AccessToken  string `json:"access_token"`
                        RefreshToken string `json:"refresh_token"`
                    } `json:"data"`
                    Success bool `json:"success"`
                }
                err := json.NewDecoder(resp.Body).Decode(&result)
                assert.NoError(t, err)
                assert.True(t, result.Success)
                assert.NotEmpty(t, result.Data.AccessToken)
                assert.NotEmpty(t, result.Data.RefreshToken)
            }

            mockRepo.AssertExpectations(t)
            mockRedis.AssertExpectations(t)
        })
    }
}

func TestRefreshToken(t *testing.T) {
    tests := []struct {
        name           string
        payload        dto.RefreshTokenRequest
        setupMock      func(*MockUserRepository, *MockRedisClient)
        expectedStatus int
        expectError    bool
    }{
        {
            name: "Success",
            payload: dto.RefreshTokenRequest{
                RefreshToken: "valid.refresh.token",
            },
            setupMock: func(mr *MockUserRepository, mrc *MockRedisClient) {
                user := &model.User{
                    ID:    primitive.NewObjectID(),
                    Roles: []string{"user"},
                }
                mr.On("FindOne", mock.Anything, mock.Anything).Return(user, nil)
                
                // Mock Redis operations
                getCmd := redis.NewStringCmd(context.Background())
                getCmd.SetVal("")
                mrc.On("Get", mock.Anything, mock.Anything).Return(getCmd)
                
                setCmd := redis.NewStatusCmd(context.Background())
                setCmd.SetVal("OK")
                mrc.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(setCmd)
            },
            expectedStatus: 200,
            expectError:    false,
        },
        {
            name: "Invalid Refresh Token",
            payload: dto.RefreshTokenRequest{
                RefreshToken: "invalid.token",
            },
            setupMock: func(mr *MockUserRepository, mrc *MockRedisClient) {
                getCmd := redis.NewStringCmd(context.Background())
                getCmd.SetErr(redis.Nil)
                mrc.On("Get", mock.Anything, mock.Anything).Return(getCmd)
            },
            expectedStatus: 401,
            expectError:    true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            app := fiber.New()
            mockRepo := new(MockUserRepository)
            mockRedis := new(MockRedisClient)
            tt.setupMock(mockRepo, mockRedis)

            userService := service.NewUserService(mockRepo, mockRedis, testConfig)
            userHandler := handlers.NewUserHandler(userService)
            app.Post("/auth/refresh", userHandler.RefreshToken)

            payloadBytes, _ := json.Marshal(tt.payload)
            req := httptest.NewRequest("POST", "/auth/refresh", bytes.NewReader(payloadBytes))
            req.Header.Set("Content-Type", "application/json")

            resp, err := app.Test(req)
            assert.NoError(t, err)
            assert.Equal(t, tt.expectedStatus, resp.StatusCode)

            if !tt.expectError {
                var result struct {
                    Data struct {
                        AccessToken  string `json:"access_token"`
                        RefreshToken string `json:"refresh_token"`
                    } `json:"data"`
                    Success bool `json:"success"`
                }
                err := json.NewDecoder(resp.Body).Decode(&result)
                assert.NoError(t, err)
                assert.True(t, result.Success)
                assert.NotEmpty(t, result.Data.AccessToken)
                assert.NotEmpty(t, result.Data.RefreshToken)
            }

            mockRepo.AssertExpectations(t)
            mockRedis.AssertExpectations(t)
        })
    }
}

func TestLogout(t *testing.T) {
    tests := []struct {
        name           string
        setupMock      func(*MockUserRepository, *MockRedisClient)
        setupHeaders   func(*httptest.Request)
        expectedStatus int
        expectError    bool
    }{
        {
            name: "Success",
            setupMock: func(mr *MockUserRepository, mrc *MockRedisClient) {
                // Mock Redis Pipeline operations
                pipe := new(MockPipeliner)
                mrc.On("Pipeline").Return(pipe)
                
                // Mock successful pipeline execution
                pipe.On("Exec", mock.Anything).Return([]redis.Cmder{}, nil)
            },
            setupHeaders: func(req *httptest.Request) {
                req.Header.Set("Authorization", "Bearer valid.access.token")
                req.Header.Set("X-Refresh-Token", "valid.refresh.token")
            },
            expectedStatus: 200,
            expectError:    false,
        },
        {
            name: "Missing Token",
            setupMock: func(mr *MockUserRepository, mrc *MockRedisClient) {},
            setupHeaders: func(req *httptest.Request) {},
            expectedStatus: 401,
            expectError:    true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            app := fiber.New()
            mockRepo := new(MockUserRepository)
            mockRedis := new(MockRedisClient)
            tt.setupMock(mockRepo, mockRedis)

            userService := service.NewUserService(mockRepo, mockRedis, testConfig)
            userHandler := handlers.NewUserHandler(userService)
            app.Get("/auth/logout", userHandler.Logout)

            req := httptest.NewRequest("GET", "/auth/logout", nil)
            tt.setupHeaders(req)

            resp, err := app.Test(req)
            assert.NoError(t, err)
            assert.Equal(t, tt.expectedStatus, resp.StatusCode)

            mockRepo.AssertExpectations(t)
            mockRedis.AssertExpectations(t)
        })
    }
}

// Mock Redis Pipeliner
type MockPipeliner struct {
    mock.Mock
}

func (m *MockPipeliner) Exec(ctx context.Context) ([]redis.Cmder, error) {
    args := m.Called(ctx)
    return args.Get(0).([]redis.Cmder), args.Error(1)
}

func (m *MockPipeliner) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
    args := m.Called(ctx, key, value, expiration)
    return args.Get(0).(*redis.StatusCmd)
}

func (m *MockPipeliner) Del(ctx context.Context, keys ...string) *redis.IntCmd {
    args := m.Called(ctx, keys)
    return args.Get(0).(*redis.IntCmd)
}