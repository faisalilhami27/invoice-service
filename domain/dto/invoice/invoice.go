package dto

import (
	"github.com/google/uuid"
)

type InvoiceRequest struct {
	InvoiceNumber string    `json:"invoiceNumber"`
	Data          any       `json:"data"`
	TemplateID    uuid.UUID `json:"templateID"`
}

type InvoiceResponse struct {
	UUID string `json:"uuid"`
	URL  string `json:"url"`
}
