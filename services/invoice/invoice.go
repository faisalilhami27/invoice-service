package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"

	"invoice-service/common/gcs"
	"invoice-service/common/sentry"
	dto "invoice-service/domain/dto/invoice"
	"invoice-service/domain/models"
	"invoice-service/repositories"
	"invoice-service/utils/helper"
)

type InvoiceService struct {
	repositoryRegistry repositories.IRepositoryRegistry
	sentry             sentry.ISentry
	gcs                gcs.IGCSClient
}

type IInvoiceService interface {
	StoreInvoice(ctx context.Context, request *dto.InvoiceRequest) (*dto.InvoiceResponse, error)
}

func NewInvoiceService(
	repositoryRegistry repositories.IRepositoryRegistry,
	sentry sentry.ISentry,
	gcs gcs.IGCSClient,
) IInvoiceService {
	return &InvoiceService{
		repositoryRegistry: repositoryRegistry,
		sentry:             sentry,
		gcs:                gcs,
	}
}

func (t *InvoiceService) StoreInvoice(
	ctx context.Context,
	request *dto.InvoiceRequest,
) (*dto.InvoiceResponse, error) {
	const logCtx = "services.template.template.StoreInvoice"
	var (
		span = t.sentry.StartSpan(ctx, logCtx)
		url  string
	)
	ctx = t.sentry.SpanContext(span)
	defer t.sentry.Finish(span)

	callback := func(ctx mongo.SessionContext) (any, error) {
		invoice, txErr := t.repositoryRegistry.GetInvoice().CreateInvoice(ctx, &models.Invoice{
			InvoiceNumber: request.InvoiceNumber,
			Data:          request.Data,
		})
		if txErr != nil {
			return nil, txErr
		}

		var data map[string]interface{}
		jsonData, _ := json.Marshal(request.Data) //nolint:errcheck
		txErr = json.Unmarshal(jsonData, &data)
		if txErr != nil {
			return nil, txErr
		}

		templateResult, txErr := t.repositoryRegistry.GetTemplate().FindOneByUUID(ctx, request.TemplateID)
		if txErr != nil {
			return nil, txErr
		}

		generatePDF, txErr := helper.GeneratePDF(ctx, templateResult.HTML, data)
		if txErr != nil {
			return nil, txErr
		}

		invoiceNumber := strings.ToLower(strings.ReplaceAll(request.InvoiceNumber, "/", "-"))
		filename := fmt.Sprintf("%s.pdf", invoiceNumber)
		url, txErr = t.gcs.UploadFileInByte(ctx, filename, generatePDF)
		if txErr != nil {
			return nil, txErr
		}

		response := dto.InvoiceResponse{
			UUID: invoice.UUID,
			URL:  url,
		}

		return response, nil
	}

	result, err := t.repositoryRegistry.Transaction(ctx, callback)
	if err != nil {
		return nil, err
	}

	invoice, _ := result.(dto.InvoiceResponse) //nolint:errcheck
	return &invoice, nil
}
