package handlers

import (
	"context"
	"go-fiber-api/internal/service"
	"go-fiber-api/pkg/utils"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ShopHandler struct {
    shopService *service.ShopService
}

func NewShopHandler(shopService *service.ShopService) *ShopHandler {
    return &ShopHandler{
        shopService: shopService,
    }
}

// @Summary List shops
// @Description Get paginated list of shops with optional filtering 
// @Tags shop
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "Page number" default(1) 
// @Param page_size query int false "Page size" default(10)
// @Param name query string false "Filter by name"
// @Success 200
// @Router /shop/list [get]
func (s *ShopHandler) ShopList(c *fiber.Ctx) error {
    page, pageSize := utils.PaginationParams(c)

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    total, err := s.shopService.Count(ctx, bson.D{})
    if err != nil {
        return utils.SendError(c, http.StatusInternalServerError, "Failed to count shops: "+err.Error())
    }

    opts := options.Find().
        SetSkip(int64((page - 1) * pageSize)).
        SetLimit(int64(pageSize)).
        SetSort(bson.D{{Key: "created_at", Value: -1}})

    users, err := s.shopService.FindAll(ctx, bson.D{}, opts)
    if err != nil {
        return utils.SendError(c, http.StatusInternalServerError, err.Error())
    }

    response := utils.CreatePagination(page, pageSize, total, users)
    return utils.SendSuccess(c, http.StatusOK, response)
}

func (s *ShopHandler) CreateShop(c *fiber.Ctx) error {
    return nil
}