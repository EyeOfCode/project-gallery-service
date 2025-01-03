package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShopRequest struct {
	Name   string  `json:"name" binding:"required,min=3,max=30"`
	Budget float64 `json:"budget" binding:"required"`
}

type UpdateShopRequest struct {
	Name   string  `json:"name" binding:"omitempty,min=3,max=30"`
	Budget float64 `json:"budget" binding:"omitempty"`
}

type UpdateShopResponse struct {
	ID        primitive.ObjectID `json:"id"`
	Name      string             `json:"name"`
	Budget    float64            `json:"budget"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}