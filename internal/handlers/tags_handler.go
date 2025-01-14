package handlers

import (
	"context"
	"fmt"
	"pre-test-gallery-service/internal/service"
	"pre-test-gallery-service/pkg/dto"
	"pre-test-gallery-service/pkg/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type TagsHandler struct {
	tagsService *service.TagsService
}

func NewTagsHandler(tagsService *service.TagsService) *TagsHandler {
	return &TagsHandler{
		tagsService:      tagsService,
	}
}

// @Summary Get all tags
// @Description Get all tags
// @Tags tags
// @Produce json
// @Success 200 {object} []model.Tags
// @Router /tags [get]
func (h *TagsHandler) GetAllTags(c *fiber.Ctx) error {
	tags, err := h.tagsService.GetAllTags(c.Context())
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Tags not found")
	}
	return utils.SendSuccess(c, fiber.StatusOK, tags)
}

// @Summary Create a new tag
// @Description Create a new tag
// @Tags tags
// @Accept json
// @Produce json
// @Param tags body dto.TagsRequest true "Tags request"
// @Success 200 {object} model.Tags
// @Router /tags [post]
func (h *TagsHandler) CreateTags(c *fiber.Ctx) error {
	var req dto.TagsRequest

	if err := c.BodyParser(&req); err != nil {
        return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body")
    }

    if err := utils.ValidateStruct(&req); err != nil {
        return utils.SendValidationError(c, err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

	existingTag, err := h.tagsService.FindOneTags(ctx, bson.M{"name": req.Name})
    if err != nil {
        return utils.SendError(c, fiber.StatusBadRequest, err.Error())
    }
	
	fmt.Println(existingTag)

	if existingTag != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Tags already exist")
	}

	tag, err := h.tagsService.CreateTags(ctx, req)
    if err != nil {
        return utils.SendError(c, fiber.StatusInternalServerError, err.Error())
    }

	return utils.SendSuccess(c, fiber.StatusOK, tag)
}

// @Summary Delete a tag
// @Description Delete a tag
// @Tags tags
// @Produce json
// @Param id path string true "Tag ID"
// @Success 200 {object} nil
// @Router /tags/{id} [delete]
func (h *TagsHandler) DeleteTags(c *fiber.Ctx) error {
	tag := c.Params("tag")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tagRes, err := h.tagsService.FindOneTags(ctx, bson.M{"name": tag})
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	if tagRes == nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Tags not found")
	}

	if err := h.tagsService.DeleteTags(ctx, tagRes.ID); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, nil, "Tags deleted successfully")
}