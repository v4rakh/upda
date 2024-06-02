package server

import (
	"git.myservermanager.com/varakh/upda/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

type eventHandler struct {
	service eventService
}

func newEventHandler(s *eventService) *eventHandler {
	return &eventHandler{service: *s}
}

func (h *eventHandler) window(c *gin.Context) {
	var queryParams api.EventWindowRequest
	var err error
	if err = c.ShouldBindQuery(&queryParams); err != nil {
		errAbortWithValidatorPayload(c, err)
		return
	}

	var events []*Event
	if events, err = h.service.window(queryParams.Size, queryParams.Skip, queryParams.OrderBy, queryParams.Order); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
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
	if hasNext, err = h.service.windowHasNext(queryParams.Size, queryParams.Skip, queryParams.OrderBy, queryParams.Order); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewDataResponseWithPayload(api.NewEventWindowResponse(data, queryParams.Size, queryParams.Skip, queryParams.OrderBy, queryParams.Order, hasNext)))
}

func (h *eventHandler) get(c *gin.Context) {
	e, err := h.service.get(c.Param("id"))
	if err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewEventSingleResponse(e.ID, e.Name, e.CreatedAt, e.UpdatedAt, e.Payload))
}

func (h *eventHandler) delete(c *gin.Context) {
	if err := h.service.delete(c.Param("id")); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	c.Header(headerContentType, headerContentTypeApplicationJson)
	c.Status(http.StatusNoContent)
}
