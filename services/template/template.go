package services

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"invoice-service/common/sentry"
	dto "invoice-service/domain/dto/template"
	"invoice-service/domain/models"
	"invoice-service/repositories"

	"io"
	"mime/multipart"
)

type TemplateService struct {
	repositoryRegistry repositories.IRepositoryRegistry
	sentry             sentry.ISentry
}

type ITemplateService interface {
	StoreTemplate(ctx context.Context, request *dto.TemplateRequest) (*dto.TemplateResponse, error)
}

func NewTemplateService(repositoryRegistry repositories.IRepositoryRegistry, sentry sentry.ISentry) ITemplateService {
	return &TemplateService{
		repositoryRegistry: repositoryRegistry,
		sentry:             sentry,
	}
}

func (t *TemplateService) StoreTemplate(
	ctx context.Context,
	request *dto.TemplateRequest,
) (*dto.TemplateResponse, error) {
	const logCtx = "services.template.template.StoreTemplate"
	var (
		span = t.sentry.StartSpan(ctx, logCtx)
	)
	ctx = t.sentry.SpanContext(span)
	defer t.sentry.Finish(span)

	callback := func(ctx mongo.SessionContext) (any, error) {
		file, err := request.HTML.Open()
		if err != nil {
			return nil, err
		}
		defer func(file multipart.File) {
			err = file.Close()
			if err != nil {
				return
			}
		}(file)

		htmlContent, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}

		template, err := t.repositoryRegistry.GetTemplate().CreateTemplate(ctx, &models.Template{
			HTML:      string(htmlContent),
			Category:  request.Category,
			Service:   request.Service,
			CreatedBy: &request.CreatedBy,
		})
		if err != nil {
			return nil, err
		}

		response := dto.TemplateResponse{
			UUID:     template.UUID,
			Category: template.Category,
			Service:  template.Service,
		}

		return &response, nil
	}

	result, err := t.repositoryRegistry.Transaction(ctx, callback)
	if err != nil {
		return nil, err
	}

	response, _ := result.(dto.TemplateResponse) //nolint:errcheck
	return &response, nil
}
