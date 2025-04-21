package handler

import (
	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type WebhookHandler struct {
	service service.WebhookService
}

func NewWebhookHandler(s *service.WebhookService) *WebhookHandler {
	return &WebhookHandler{service: *s}
}

func (h *WebhookHandler) Paginate(c *gin.Context) {
	var queryParams api.PaginateWebhookRequest
	var err error
	if err = c.ShouldBindQuery(&queryParams); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	var webhooks []*model.Webhook
	if webhooks, err = h.service.Paginate(queryParams.Page, queryParams.PageSize, queryParams.OrderBy, queryParams.Order); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
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
	if totalElements, err = h.service.Count(); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	totalPages := (totalElements + int64(queryParams.PageSize) - 1) / int64(queryParams.PageSize)
	c.JSON(http.StatusOK, api.NewDataResponseWithPayload(api.NewWebhookPageResponse(data, queryParams.Page, queryParams.PageSize, queryParams.OrderBy, queryParams.Order, totalElements, totalPages)))
}

func (h *WebhookHandler) Get(c *gin.Context) {
	var err error
	var pathParams api.IDUriRequest

	if err = c.ShouldBindUri(&pathParams); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	var e *model.Webhook
	if e, err = h.service.Get(pathParams.ID); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewWebhookSingleResponse(e.ID, e.Label, e.Type, e.IgnoreHost, "", e.CreatedAt, e.UpdatedAt))
}

func (h *WebhookHandler) Create(c *gin.Context) {
	var e *model.Webhook
	var err error

	var req api.CreateWebhookRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	if e, err = h.service.Create(req.Label, api.WebhookType(req.Type), req.IgnoreHost); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewWebhookSingleResponse(e.ID, e.Label, e.Type, e.IgnoreHost, e.Token, e.CreatedAt, e.UpdatedAt))
}

func (h *WebhookHandler) UpdateLabel(c *gin.Context) {
	var e *model.Webhook
	var err error

	var req api.ModifyWebhookLabelRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	var pathParams api.IDUriRequest

	if err = c.ShouldBindUri(&pathParams); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	if e, err = h.service.UpdateLabel(pathParams.ID, req.Label); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewWebhookSingleResponse(e.ID, e.Label, e.Type, e.IgnoreHost, "", e.CreatedAt, e.UpdatedAt))
}

func (h *WebhookHandler) UpdateIgnoreHost(c *gin.Context) {
	var e *model.Webhook
	var err error

	var req api.ModifyWebhookIgnoreHostRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	var pathParams api.IDUriRequest

	if err = c.ShouldBindUri(&pathParams); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	if e, err = h.service.UpdateIgnoreHost(pathParams.ID, req.IgnoreHost); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewWebhookSingleResponse(e.ID, e.Label, e.Type, e.IgnoreHost, "", e.CreatedAt, e.UpdatedAt))
}

func (h *WebhookHandler) Delete(c *gin.Context) {
	var err error
	var pathParams api.IDUriRequest

	if err = c.ShouldBindUri(&pathParams); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	if err = h.service.Delete(pathParams.ID); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.Header(api.HeaderContentType, api.HeaderContentTypeApplicationJson)
	c.Status(http.StatusNoContent)
}
