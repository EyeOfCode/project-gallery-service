package handlers

import (
	"pre-test-gallery-service/internal/service"
	"pre-test-gallery-service/pkg/utils"

	"github.com/gofiber/fiber/v2"
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
// @Router /api/v1/tags [get]
func (h *TagsHandler) GetAllTags(c *fiber.Ctx) error {
	tags, err := h.tagsService.GetAllTags(c.Context())
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Tags not found")
	}
	return utils.SendSuccess(c, fiber.StatusOK, tags)
}

func (h *TagsHandler) CreateTags(c *fiber.Ctx) error {
	return nil
}

func (h *TagsHandler) DeleteTags(c *fiber.Ctx) error {
	return nil
}