package handler

import (
	"errors"
	"net/http"

	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/internal/server/constant"
	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/service"
	"git.myservermanager.com/varakh/upda/internal/server/service_error"
	"github.com/gin-gonic/gin"
)

type WebhookInvocationHandler struct {
	invocationService service.WebhookInvocationService
	webhookService    service.WebhookService
}

func NewWebhookInvocationHandler(i *service.WebhookInvocationService, w *service.WebhookService) *WebhookInvocationHandler {
	return &WebhookInvocationHandler{invocationService: *i, webhookService: *w}
}

func (h *WebhookInvocationHandler) Execute(c *gin.Context) {
	tokenHeader := c.GetHeader(api.HeaderWebhookToken)

	var err error
	var pathParams api.IDUriRequest

	if err = c.ShouldBindUri(&pathParams); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	webhookId := pathParams.ID
	var w *model.Webhook

	if w, err = h.webhookService.Get(webhookId); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	switch w.Type {
	case constant.WebhookTypeGeneric.String():
		var req api.WebhookGenericRequest
		if err = c.ShouldBindJSON(&req); err != nil {
			AbortWithValidatorPayload(c, err)
			return
		}
		if err = h.invocationService.ExecuteGeneric(webhookId, tokenHeader, req); err != nil {
			_ = c.AbortWithError(ToHttpStatus(err), err)
			return
		}
	case constant.WebhookTypeDiun.String():
		var req api.WebhookDiunRequest
		if err = c.ShouldBindJSON(&req); err != nil {
			AbortWithValidatorPayload(c, err)
			return
		}
		if err = h.invocationService.ExecuteDiun(webhookId, tokenHeader, req); err != nil {
			_ = c.AbortWithError(ToHttpStatus(err), err)
			return
		}
	default:
		err = service_error.NewServiceError(service_error.ErrCodeIllegalArgument, errors.New("no default handler for webhook type found"))
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.Header(api.HeaderContentType, api.HeaderContentTypeApplicationJson)
	c.Status(http.StatusNoContent)
}
