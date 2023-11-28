package dto

import (
	"github.com/google/uuid"
)

type InvoiceRequest struct {
	InvoiceNumber string    `json:"invoice_number"`
	Data          any       `json:"data"`
	TemplateID    uuid.UUID `json:"template_id"`
	CreatedBy     string    `json:"created_by"`
}

type InvoiceResponse struct {
	UUID string `json:"uuid"`
	URL  string `json:"url"`
}
