package repositories

import (
	"context"

	uuidGenerate "github.com/google/uuid"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"invoice-service/common/sentry"
	errorHelper "invoice-service/utils/error"

	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"invoice-service/config"
	"invoice-service/domain/models"
)

type Template struct {
	db     *mongo.Client
	sentry sentry.ISentry
}

type ITemplate interface {
	CreateTemplate(ctx context.Context, template *models.Template) (*models.Template, error)
	FindOneByUUID(ctx context.Context, uuid uuidGenerate.UUID) (*models.Template, error)
}

func NewTemplate(db *mongo.Client, sentry sentry.ISentry) ITemplate {
	return &Template{
		db:     db,
		sentry: sentry,
	}
}

func (t *Template) CreateTemplate(ctx context.Context, template *models.Template) (*models.Template, error) {
	const logCtx = "repositories.template.template.CreateTemplate"
	var (
		span = t.sentry.StartSpan(ctx, logCtx)
	)
	ctx = t.sentry.SpanContext(span)
	defer t.sentry.Finish(span)

	utcTime := time.Now().UTC()
	jakartaLocation, _ := time.LoadLocation("Asia/Jakarta") //nolint:errcheck
	today := utcTime.In(jakartaLocation).Format(time.RFC3339)
	parseTime, _ := time.Parse(time.RFC3339, today)
	template.UUID = uuidGenerate.New().String()
	template.CreatedAt = parseTime
	template.UpdatedAt = &parseTime
	collection := t.db.Database(config.Config.Database.Name).Collection("templates")
	result, err := collection.InsertOne(ctx, template)
	if err != nil {
		return nil, errorHelper.WrapError(err, t.sentry)
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errorHelper.WrapError(err, t.sentry)
	}

	response := models.Template{
		ID:       insertedID,
		UUID:     template.UUID,
		HTML:     template.HTML,
		Category: template.Category,
		Service:  template.Service,
	}
	return &response, nil
}

func (t *Template) FindOneByUUID(ctx context.Context, uuid uuidGenerate.UUID) (*models.Template, error) {
	const logCtx = "repositories.template.template.FindOneByUUID"
	var (
		span = t.sentry.StartSpan(ctx, logCtx)
	)
	ctx = t.sentry.SpanContext(span)
	defer t.sentry.Finish(span)

	var template models.Template
	where := bson.D{{"uuid", uuid.String()}} //nolint:govet
	collection := t.db.Database(config.Config.Database.Name).Collection("templates")
	err := collection.FindOne(ctx, where).Decode(&template)
	if err != nil {
		return nil, errorHelper.WrapError(err, t.sentry)
	}
	return &template, nil
}
