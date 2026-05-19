package handler

import (
	"net/http"

	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/service"
	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	service service.EventService
}

func NewEventHandler(s *service.EventService) *EventHandler {
	return &EventHandler{service: *s}
}

func (h *EventHandler) Window(c *gin.Context) {
	var queryParams api.EventWindowRequest
	var err error
	if err = c.ShouldBindQuery(&queryParams); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	var events []*model.Event
	if events, err = h.service.Window(queryParams.Size, queryParams.Skip, queryParams.OrderBy, queryParams.Order, queryParams.UpdateID); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	var data []*api.EventResponse
	data = make([]*api.EventResponse, 0, len(events))

	for _, e := range events {
		data = append(data, &api.EventResponse{
			ID:        e.ID.String(),
			Name:      e.Name,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
			Payload:   e.Payload,
		})
	}

	var hasNext bool
	if hasNext, err = h.service.WindowHasNext(queryParams.Size, queryParams.Skip, queryParams.OrderBy, queryParams.Order, queryParams.UpdateID); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewDataResponseWithPayload(api.NewEventWindowResponse(data, queryParams.Size, queryParams.Skip, queryParams.OrderBy, queryParams.Order, hasNext)))
}

func (h *EventHandler) Get(c *gin.Context) {
	var err error
	var pathParams api.IDUriRequest

	if err = c.ShouldBindUri(&pathParams); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}
	var e *model.Event
	if e, err = h.service.Get(pathParams.ID); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewEventSingleResponse(e.ID.String(), e.Name, e.CreatedAt, e.UpdatedAt, e.Payload))
}

func (h *EventHandler) Delete(c *gin.Context) {
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
