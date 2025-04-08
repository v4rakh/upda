package server

import (
	"git.myservermanager.com/varakh/upda/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

type webhookHandler struct {
	service webhookService
}

func newWebhookHandler(s *webhookService) *webhookHandler {
	return &webhookHandler{service: *s}
}

func (h *webhookHandler) paginate(c *gin.Context) {
	var queryParams api.PaginateWebhookRequest
	var err error
	if err = c.ShouldBindQuery(&queryParams); err != nil {
		errAbortWithValidatorPayload(c, err)
		return
	}

	var webhooks []*Webhook
	if webhooks, err = h.service.paginate(queryParams.Page, queryParams.PageSize, queryParams.OrderBy, queryParams.Order); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	var data []*api.WebhookResponse
	data = make([]*api.WebhookResponse, 0, len(webhooks))

	for _, e := range webhooks {
		data = append(data, &api.WebhookResponse{
			ID:         e.ID,
			Label:      e.Label,
			Type:       e.Type,
			IgnoreHost: e.IgnoreHost,
			CreatedAt:  e.CreatedAt,
			UpdatedAt:  e.UpdatedAt,
		})
	}

	var totalElements int64
	if totalElements, err = h.service.count(); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	totalPages := (totalElements + int64(queryParams.PageSize) - 1) / int64(queryParams.PageSize)
	c.JSON(http.StatusOK, api.NewDataResponseWithPayload(api.NewWebhookPageResponse(data, queryParams.Page, queryParams.PageSize, queryParams.OrderBy, queryParams.Order, totalElements, totalPages)))
}

func (h *webhookHandler) get(c *gin.Context) {
	e, err := h.service.get(c.Param("id"))
	if err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewWebhookSingleResponse(e.ID, e.Label, e.Type, e.IgnoreHost, "", e.CreatedAt, e.UpdatedAt))
}

func (h *webhookHandler) create(c *gin.Context) {
	var e *Webhook
	var err error

	var req api.CreateWebhookRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		errAbortWithValidatorPayload(c, err)
		return
	}

	if e, err = h.service.create(req.Label, api.WebhookType(req.Type), req.IgnoreHost); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewWebhookSingleResponse(e.ID, e.Label, e.Type, e.IgnoreHost, e.Token, e.CreatedAt, e.UpdatedAt))
}

func (h *webhookHandler) updateLabel(c *gin.Context) {
	var e *Webhook
	var err error

	var req api.ModifyWebhookLabelRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		errAbortWithValidatorPayload(c, err)
		return
	}

	if e, err = h.service.updateLabel(c.Param("id"), req.Label); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewWebhookSingleResponse(e.ID, e.Label, e.Type, e.IgnoreHost, "", e.CreatedAt, e.UpdatedAt))
}

func (h *webhookHandler) updateIgnoreHost(c *gin.Context) {
	var e *Webhook
	var err error

	var req api.ModifyWebhookIgnoreHostRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		errAbortWithValidatorPayload(c, err)
		return
	}

	if e, err = h.service.updateIgnoreHost(c.Param("id"), req.IgnoreHost); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewWebhookSingleResponse(e.ID, e.Label, e.Type, e.IgnoreHost, "", e.CreatedAt, e.UpdatedAt))
}

func (h *webhookHandler) delete(c *gin.Context) {
	if err := h.service.delete(c.Param("id")); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	c.Header(api.HeaderContentType, api.HeaderContentTypeApplicationJson)
	c.Status(http.StatusNoContent)
}
