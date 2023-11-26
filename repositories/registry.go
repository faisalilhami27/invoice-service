package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"invoice-service/common/sentry"
	invoiceRepo "invoice-service/repositories/invoice"
	templateRepo "invoice-service/repositories/template"
)

type RepositoryRegistry struct {
	client *mongo.Client
	sentry sentry.ISentry
}

type IRepositoryRegistry interface {
	GetTx() *mongo.Client
	GetTemplate() templateRepo.ITemplate
	GetInvoice() invoiceRepo.IInvoice
	Transaction(ctx context.Context, callback func(ctx mongo.SessionContext) (any, error)) (any, error)
}

func NewRepositoryRegistry(client *mongo.Client, sentry sentry.ISentry) IRepositoryRegistry {
	return &RepositoryRegistry{
		client: client,
		sentry: sentry,
	}
}

func (r *RepositoryRegistry) Transaction(
	ctx context.Context,
	callback func(ctx mongo.SessionContext) (any, error),
) (any, error) {
	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := r.client.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	res, err := session.WithTransaction(ctx, callback, txnOptions)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *RepositoryRegistry) GetTx() *mongo.Client {
	return r.client
}

func (r *RepositoryRegistry) GetTemplate() templateRepo.ITemplate {
	return templateRepo.NewTemplate(r.client, r.sentry)
}

func (r *RepositoryRegistry) GetInvoice() invoiceRepo.IInvoice {
	return invoiceRepo.NewInvoice(r.client, r.sentry)
}
