package handler

import (
	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/service"
	"github.com/gin-gonic/gin"
	"net/http"
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
	if events, err = h.service.Window(queryParams.Size, queryParams.Skip, queryParams.OrderBy, queryParams.Order); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	var data []*api.EventResponse
	data = make([]*api.EventResponse, 0, len(events))

	for _, e := range events {
		data = append(data, &api.EventResponse{
			ID:        e.ID,
			Name:      e.Name,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
			Payload:   e.Payload,
		})
	}

	var hasNext bool
	if hasNext, err = h.service.WindowHasNext(queryParams.Size, queryParams.Skip, queryParams.OrderBy, queryParams.Order); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewDataResponseWithPayload(api.NewEventWindowResponse(data, queryParams.Size, queryParams.Skip, queryParams.OrderBy, queryParams.Order, hasNext)))
}

func (h *EventHandler) Get(c *gin.Context) {
	e, err := h.service.Get(c.Param("id"))
	if err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewEventSingleResponse(e.ID, e.Name, e.CreatedAt, e.UpdatedAt, e.Payload))
}

func (h *EventHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Param("id")); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.Header(api.HeaderContentType, api.HeaderContentTypeApplicationJson)
	c.Status(http.StatusNoContent)
}
