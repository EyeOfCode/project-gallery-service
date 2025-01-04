package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type FileStore struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	BasePath  string             `json:"base_path" bson:"base_path"`
	Dir       string             `json:"dir" bson:"dir"`
	Extension string             `json:"extension" bson:"extension"`
	ShopID    primitive.ObjectID `json:"shop_id" bson:"shop_id"`
}