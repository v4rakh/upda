package server

import (
	"errors"
	"git.myservermanager.com/varakh/upda/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

type webhookInvocationHandler struct {
	invocationService webhookInvocationService
	webhookService    webhookService
}

func newWebhookInvocationHandler(i *webhookInvocationService, w *webhookService) *webhookInvocationHandler {
	return &webhookInvocationHandler{invocationService: *i, webhookService: *w}
}

func (h *webhookInvocationHandler) execute(c *gin.Context) {
	tokenHeader := c.GetHeader(headerWebhookToken)
	webhookId := c.Param("id")

	var w *Webhook
	var err error

	if w, err = h.webhookService.get(webhookId); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	switch w.Type {
	case api.WebhookTypeGeneric.Value():
		var req api.WebhookGenericRequest
		if err = c.ShouldBindJSON(&req); err != nil {
			errAbortWithValidatorPayload(c, err)
			return
		}
		if err = h.invocationService.executeGeneric(webhookId, tokenHeader, req); err != nil {
			_ = c.AbortWithError(errToHttpStatus(err), err)
			return
		}
		break
	case api.WebhookTypeDiun.Value():
		var req api.WebhookDiunRequest
		if err = c.ShouldBindJSON(&req); err != nil {
			errAbortWithValidatorPayload(c, err)
			return
		}
		if err = h.invocationService.executeDiun(webhookId, tokenHeader, req); err != nil {
			_ = c.AbortWithError(errToHttpStatus(err), err)
			return
		}
		break
	default:
		err = newServiceError(illegalArgument, errors.New("no default handler for webhook type found"))
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	c.Header(headerContentType, headerContentTypeApplicationJson)
	c.Status(http.StatusNoContent)
}
