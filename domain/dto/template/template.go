package dto

import (
	"mime/multipart"
)

type TemplateRequest struct {
	HTML      *multipart.FileHeader `form:"html" validate:"required"`
	Category  string                `form:"category" validate:"required"`
	Service   string                `form:"service" validate:"required"`
	CreatedBy string                `form:"createdBy"`
	UpdatedBy string                `form:"updatedBy"`
}

type TemplateQueryParamRequest struct {
	Category string `form:"category"`
	Service  string `form:"service"`
}

type TemplateResponse struct {
	UUID     string `json:"uuid"`
	Category string `json:"category"`
	Service  string `json:"invoice"`
}
