package test

import (
	"context"
	"encoding/json"
	"errors"
	"go-fiber-api/internal/handlers"
	"go-fiber-api/internal/model"
	"go-fiber-api/internal/service"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
    args := m.Called(ctx, email)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *model.User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
    args := m.Called(ctx, id)
    return args.Error(0)
}

func (m *MockUserRepository) FindOne(ctx context.Context, query bson.M) (*model.User, error) {
    args := m.Called(ctx, query)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindAll(ctx context.Context, query bson.D, opts *options.FindOptions) ([]model.User, error) {
    args := m.Called(ctx, query, opts)
    return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserRepository) Count(ctx context.Context, query bson.D) (int64, error) {
    args := m.Called(ctx, query)
    return int64(args.Int(0)), args.Error(1)
}

func TestUserList(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    map[string]string
		setupMock      func(*MockUserRepository)
		expectedStatus int
		expectedData   struct {
			Page       int         `json:"page"`
			PageSize   int         `json:"pageSize"`
			TotalItems int64       `json:"totalItems"`
			TotalPages int         `json:"totalPages"`
			Items      []model.User `json:"items"`
		}
		expectError bool
	}{
		{
			name: "Success - No Filters",
			queryParams: map[string]string{
				"page": "1",
				"size": "10",
			},
			setupMock: func(m *MockUserRepository) {
				m.On("Count", mock.Anything, mock.Anything).Return(2, nil)
				users := []model.User{
					{ID: primitive.ObjectID{}, Name: "User1"},
					{ID: primitive.ObjectID{}, Name: "User2"},
				}
				m.On("FindAll", mock.Anything, mock.Anything, mock.Anything).Return(users, nil)
			},
			expectedStatus: 200,
			expectedData: struct {
				Page       int         `json:"page"`
				PageSize   int         `json:"pageSize"`
				TotalItems int64       `json:"totalItems"`
				TotalPages int         `json:"totalPages"`
				Items      []model.User `json:"items"`
			}{
				Page:       1,
				PageSize:   10,
				TotalItems: 2,
				TotalPages: 1,
				Items: []model.User{
					{ID: primitive.ObjectID{}, Name: "User1"},
					{ID: primitive.ObjectID{}, Name: "User2"},
				},
			},
			expectError: false,
		},
		{
			name: "Success - Empty Result",
			queryParams: map[string]string{
				"page": "1",
				"size": "10",
			},
			setupMock: func(m *MockUserRepository) {
				m.On("Count", mock.Anything, mock.Anything).Return(0, nil)
				m.On("FindAll", mock.Anything, mock.Anything, mock.Anything).Return([]model.User{}, nil)
			},
			expectedStatus: 200,
			expectedData: struct {
				Page       int         `json:"page"`
				PageSize   int         `json:"pageSize"`
				TotalItems int64       `json:"totalItems"`
				TotalPages int         `json:"totalPages"`
				Items      []model.User `json:"items"`
			}{
				Page:       1,
				PageSize:   10,
				TotalItems: 0,
				TotalPages: 0,
				Items:      []model.User{},
			},
			expectError: false,
		},
		{
			name: "Error - Service Failure",
			queryParams: map[string]string{
				"page": "1",
				"size": "10",
			},
			setupMock: func(m *MockUserRepository) {
				m.On("Count", mock.Anything, mock.Anything).Return(0, errors.New("database error"))
			},
			expectedStatus: 500,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			mockRepo := new(MockUserRepository)
			tt.setupMock(mockRepo)
			
			userService := service.NewUserService(mockRepo)
			userHandler := handlers.NewUserHandler(userService)
			app.Get("/user/list", userHandler.UserList)

			req := httptest.NewRequest("GET", "/user/list", nil)
			q := req.URL.Query()
			for key, value := range tt.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if !tt.expectError {
				var result struct {
					Data    struct {
						Page       int         `json:"page"`
						PageSize   int         `json:"pageSize"`
						TotalItems int64       `json:"totalItems"`
						TotalPages int         `json:"totalPages"`
						Items      []model.User `json:"items"`
					} `json:"data"`
					Success bool `json:"success"`
				}
				err := json.NewDecoder(resp.Body).Decode(&result)
				assert.NoError(t, err)
				assert.True(t, result.Success)
				assert.Equal(t, tt.expectedData, result.Data)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}