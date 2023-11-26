package controllers

import (
	"invoice-service/common/sentry"
	invoiceController "invoice-service/controllers/invoice"
	templateController "invoice-service/controllers/template"
	"invoice-service/services"
)

type ControllerRegistry struct {
	serviceRegistry services.IServiceRegistry
	sentry          sentry.ISentry
}

type IControllerRegistry interface {
	GetTemplate() templateController.ITemplateController
	GetInvoice() invoiceController.IInvoiceController
}

func NewControllerRegistry(
	serviceRegistry services.IServiceRegistry,
	sentry sentry.ISentry,
) IControllerRegistry {
	return &ControllerRegistry{
		serviceRegistry: serviceRegistry,
		sentry:          sentry,
	}
}

func (r *ControllerRegistry) GetTemplate() templateController.ITemplateController {
	return templateController.NewTemplateController(r.serviceRegistry, r.sentry)
}

func (r *ControllerRegistry) GetInvoice() invoiceController.IInvoiceController {
	return invoiceController.NewInvoiceController(r.serviceRegistry, r.sentry)
}
