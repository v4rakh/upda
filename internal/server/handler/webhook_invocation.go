package handler

import (
	"errors"
	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/service"
	"git.myservermanager.com/varakh/upda/internal/server/service_error"
	"github.com/gin-gonic/gin"
	"net/http"
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
	webhookId := c.Param("id")

	var w *model.Webhook
	var err error

	if w, err = h.webhookService.Get(webhookId); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	switch w.Type {
	case api.WebhookTypeGeneric.Value():
		var req api.WebhookGenericRequest
		if err = c.ShouldBindJSON(&req); err != nil {
			AbortWithValidatorPayload(c, err)
			return
		}
		if err = h.invocationService.ExecuteGeneric(webhookId, tokenHeader, req); err != nil {
			_ = c.AbortWithError(ToHttpStatus(err), err)
			return
		}
		break
	case api.WebhookTypeDiun.Value():
		var req api.WebhookDiunRequest
		if err = c.ShouldBindJSON(&req); err != nil {
			AbortWithValidatorPayload(c, err)
			return
		}
		if err = h.invocationService.ExecuteDiun(webhookId, tokenHeader, req); err != nil {
			_ = c.AbortWithError(ToHttpStatus(err), err)
			return
		}
		break
	default:
		err = service_error.NewServiceError(service_error.ErrCodeIllegalArgument, errors.New("no default handler for webhook type found"))
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.Header(api.HeaderContentType, api.HeaderContentTypeApplicationJson)
	c.Status(http.StatusNoContent)
}
