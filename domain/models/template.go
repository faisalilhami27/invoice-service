package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"time"
)

type Template struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UUID      string             `bson:"uuid"`
	HTML      string             `bson:"html"`
	Category  string             `bson:"category"`
	Service   string             `bson:"invoice"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt *time.Time         `bson:"updated_at"`
	CreatedBy *string            `bson:"created_by"`
	UpdatedBy *string            `bson:"updated_by"`
}
