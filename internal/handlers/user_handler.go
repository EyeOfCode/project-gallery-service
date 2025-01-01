package handlers

import (
	"context"
	"go-fiber-api/internal/service"
	"go-fiber-api/pkg/dto"
	"go-fiber-api/pkg/utils"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserHandler struct {
    userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
    return &UserHandler{
        userService: userService,
    }
}

// @Summary User list endpoint
// @Description Get the API's user list
// @Tags admin
// @Accept json
// @Produce json
// @Router /user/list [get]
func (u *UserHandler) UserList(c *fiber.Ctx) error {
    page, pageSize := utils.PaginationParams(c)

    var filter dto.UserFilter
    if err := c.QueryParser(&filter); err != nil {
        utils.SendError(c, http.StatusBadRequest, "Invalid filter parameters")
        return nil
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    mongoFilter := bson.D{}
    if filter.Name != "" {
        mongoFilter = append(mongoFilter, bson.E{
            Key: "name", 
            Value: bson.D{{
                Key: "$regex", 
                Value: primitive.Regex{Pattern: filter.Name, Options: "i"},
            }},
        })
    }

    total, err := u.userService.Count(ctx, mongoFilter)
    if err != nil {
        utils.SendError(c, http.StatusInternalServerError, "Failed to count users: "+err.Error())
        return nil
    }

    opts := options.Find().
        SetSkip(int64((page - 1) * pageSize)).
        SetLimit(int64(pageSize)).
        SetSort(bson.D{{Key: "created_at", Value: -1}})

    users, err := u.userService.FindAll(ctx, mongoFilter, opts)
    if err != nil {
        utils.SendError(c, http.StatusInternalServerError, err.Error())
        return nil
    }

    response := utils.CreatePagination(page, pageSize, total, users)
    return utils.SendSuccess(c, http.StatusOK, response)
}