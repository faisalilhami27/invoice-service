package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"invoice-service/common/sentry"
	dto "invoice-service/domain/dto/invoice"
	"invoice-service/services"
	errorValidation "invoice-service/utils/error"
	"invoice-service/utils/response"

	"net/http"
)

type InvoiceController struct {
	serviceRegistry services.IServiceRegistry
	sentry          sentry.ISentry
}

type IInvoiceController interface {
	StoreInvoice(*gin.Context)
}

func NewInvoiceController(
	serviceRegistry services.IServiceRegistry,
	sentry sentry.ISentry,
) IInvoiceController {
	return &InvoiceController{
		serviceRegistry: serviceRegistry,
		sentry:          sentry,
	}
}

func (t *InvoiceController) StoreInvoice(c *gin.Context) {
	const logCtx = "controllers.http.invoice.invoice.StoreInvoice"
	var (
		ctx     = c.Request.Context()
		span    = t.sentry.StartSpan(ctx, logCtx)
		request dto.InvoiceRequest
	)
	ctx = t.sentry.SpanContext(span)
	defer t.sentry.Finish(span)

	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code:   http.StatusBadRequest,
			Err:    err,
			Gin:    c,
			Sentry: t.sentry,
		})
		return
	}

	validate := validator.New()
	if err = validate.Struct(request); err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errorResponse := errorValidation.ErrorValidationResponse(err)
		response.HTTPResponse(response.ParamHTTPResp{
			Err:     err,
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errorResponse,
			Sentry:  t.sentry,
			Gin:     c,
		})
		return
	}

	result, err := t.serviceRegistry.GetInvoice().StoreInvoice(ctx, &request)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code:   http.StatusBadRequest,
			Err:    err,
			Gin:    c,
			Sentry: t.sentry,
		})
		return
	}

	response.HTTPResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: result,
		Err:  err,
		Gin:  c,
	})
}
