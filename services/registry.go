package services

import (
	"invoice-service/common/gcs"
	"invoice-service/common/sentry"
	"invoice-service/repositories"
	invoiceService "invoice-service/services/invoice"
	templateService "invoice-service/services/template"
)

type ServiceRegistry struct {
	repositoryRegistry repositories.IRepositoryRegistry
	sentry             sentry.ISentry
	gcs                gcs.IGCSClient
}

type IServiceRegistry interface {
	GetTemplate() templateService.ITemplateService
	GetInvoice() invoiceService.IInvoiceService
}

func NewServiceRegistry(
	repositoryRegistry repositories.IRepositoryRegistry,
	sentry sentry.ISentry,
	gcs gcs.IGCSClient,
) IServiceRegistry {
	return &ServiceRegistry{
		repositoryRegistry,
		sentry,
		gcs,
	}
}

func (s *ServiceRegistry) GetTemplate() templateService.ITemplateService {
	return templateService.NewTemplateService(s.repositoryRegistry, s.sentry)
}

func (s *ServiceRegistry) GetInvoice() invoiceService.IInvoiceService {
	return invoiceService.NewInvoiceService(s.repositoryRegistry, s.sentry, s.gcs)
}
