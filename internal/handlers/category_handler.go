package handlers

import (
	"context"
	"go-fiber-api/internal/service"
	"go-fiber-api/pkg/dto"
	"go-fiber-api/pkg/utils"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoryHandler struct {
	categoryService *service.CategoryService
	shopService *service.ShopService
}

func NewCategoryHandler(categoryService *service.CategoryService, shopService *service.ShopService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
		shopService: shopService,
	}
}

// @Summary Create Category endpoint
// @Description Post the API's create category
// @Tags category
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.CategoryRequest true "Category details"
// @Router /category [post]
func (h *CategoryHandler) Create(c *fiber.Ctx) error {
	var req dto.CategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendValidationError(c, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	shopId, err := primitive.ObjectIDFromHex(req.ShopId)
	if err != nil {
		return utils.SendError(c, http.StatusBadRequest, "Invalid shop ID format")
	}

	shop, err := h.shopService.FindByID(ctx, shopId)
	if err != nil || shop == nil {
		return utils.SendError(c, http.StatusNotFound, "Failed to find shop")
	}

	category, err := h.categoryService.Create(ctx, &req, shop)
	if err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "Failed to create category")
	}

	return utils.SendSuccess(c, http.StatusCreated, category, "Category created successfully")
}

// @Summary Get All Categories endpoint
// @Description Get the API's get all categories
// @Tags category
// @Accept json
// @Produce json
// @Security Bearer
// @Router /category/list [get]
func (h *CategoryHandler) GetAll(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	categories, err := h.categoryService.FindAll(ctx)
	if err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "Failed to fetch categories")
	}

	return utils.SendSuccess(c, http.StatusOK, categories, "Categories fetched successfully")
}

func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	var req dto.UpdateCategoryRequest
	userId := c.Locals("userID").(string)
	id := c.Params("id")

	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendValidationError(c, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	shopId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.SendError(c, http.StatusBadRequest, "Invalid shop ID format")
	}

	shop, err := h.shopService.FindByID(ctx, shopId)
	if err != nil || shop == nil {
		return utils.SendError(c, http.StatusNotFound, "Failed to find shop")
	}

	category, err := h.categoryService.Update(ctx, &req)
	if err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "Failed to update category")
	}

	return utils.SendSuccess(c, http.StatusOK, category, "Category updated successfully")
}