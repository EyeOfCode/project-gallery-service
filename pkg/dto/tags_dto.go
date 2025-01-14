package dto

type TagsRequest struct {
	Name string `json:"name" binding:"required"`
}
