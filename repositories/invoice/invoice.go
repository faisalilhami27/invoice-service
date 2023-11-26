package repositories

import (
	"context"

	uuidGenerate "github.com/google/uuid"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"invoice-service/common/sentry"
	errorHelper "invoice-service/utils/error"

	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"invoice-service/config"
	"invoice-service/domain/models"
)

type Invoice struct {
	db     *mongo.Client
	sentry sentry.ISentry
}

type IInvoice interface {
	CreateInvoice(context.Context, *models.Invoice) (*models.Invoice, error)
}

func NewInvoice(db *mongo.Client, sentry sentry.ISentry) IInvoice {
	return &Invoice{
		db:     db,
		sentry: sentry,
	}
}

func (t *Invoice) CreateInvoice(ctx context.Context, invoice *models.Invoice) (*models.Invoice, error) {
	utcTime := time.Now().UTC()
	jakartaLocation, _ := time.LoadLocation("Asia/Jakarta") //nolint:errcheck
	today := utcTime.In(jakartaLocation).Format(time.RFC3339)
	parseTime, _ := time.Parse(time.RFC3339, today)
	invoice.UUID = uuidGenerate.New().String()
	invoice.CreatedAt = parseTime
	invoice.UpdatedAt = &parseTime
	collection := t.db.Database(config.Config.Database.Name).Collection("invoices")
	result, err := collection.InsertOne(ctx, invoice)
	if err != nil {
		return nil, errorHelper.WrapError(err, t.sentry)
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errorHelper.WrapError(err, t.sentry)
	}

	response := models.Invoice{
		ID:            insertedID,
		UUID:          invoice.UUID,
		InvoiceNumber: invoice.InvoiceNumber,
		Data:          invoice.Data,
	}
	return &response, nil
}
