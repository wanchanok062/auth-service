package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetUser struct {
	Id   primitive.ObjectID `json:"id" bson:"_id"`
	Name string             `json:"name"`
}

type RegisterRequest struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	UserName  string             `json:"user_name" bson:"user_name" validate:"required"`
	Name      string             `json:"name" validate:"required"`
	Password  string             `json:"password" validate:"required"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
}
