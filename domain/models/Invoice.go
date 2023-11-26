package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invoice struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	UUID          string             `bson:"uuid"`
	InvoiceNumber string             `bson:"invoice_number"`
	Data          interface{}        `bson:"data"`
	CreatedAt     time.Time          `bson:"created_at"`
	UpdatedAt     *time.Time         `bson:"updated_at"`
	CreatedBy     *string            `bson:"created_by"`
	UpdatedBy     *string            `bson:"updated_by"`
}
