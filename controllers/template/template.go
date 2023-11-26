package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"invoice-service/common/sentry"
	dto "invoice-service/domain/dto/template"
	"invoice-service/services"
	errorValidation "invoice-service/utils/error"
	"invoice-service/utils/response"

	"net/http"
)

type TemplateController struct {
	serviceRegistry services.IServiceRegistry
	sentry          sentry.ISentry
}

type ITemplateController interface {
	StoreTemplate(*gin.Context)
}

func NewTemplateController(
	serviceRegistry services.IServiceRegistry,
	sentry sentry.ISentry,
) ITemplateController {
	return &TemplateController{
		serviceRegistry: serviceRegistry,
		sentry:          sentry,
	}
}

func (t *TemplateController) StoreTemplate(c *gin.Context) {
	const logCtx = "controllers.http.template.template.StoreTemplate"
	var (
		ctx     = c.Request.Context()
		span    = t.sentry.StartSpan(ctx, logCtx)
		request dto.TemplateRequest
	)
	ctx = t.sentry.SpanContext(span)
	defer t.sentry.Finish(span)

	err := c.ShouldBindWith(&request, binding.FormMultipart)
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

	result, err := t.serviceRegistry.GetTemplate().StoreTemplate(ctx, &request)
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
