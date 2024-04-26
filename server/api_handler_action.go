package server

import (
	"git.myservermanager.com/varakh/upda/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

type actionHandler struct {
	service actionService
}

func newActionHandler(s *actionService) *actionHandler {
	return &actionHandler{service: *s}
}

func (h *actionHandler) paginate(c *gin.Context) {
	var queryParams api.PaginateActionRequest
	var err error
	if err = c.ShouldBindQuery(&queryParams); err != nil {
		errAbortWithValidatorPayload(c, err)
		return
	}

	var actions []*Action
	if actions, err = h.service.paginate(queryParams.Page, queryParams.PageSize, queryParams.OrderBy, queryParams.Order); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	var data []*api.ActionResponse
	data = make([]*api.ActionResponse, 0)

	for _, e := range actions {
		data = append(data, &api.ActionResponse{
			ID:               e.ID,
			Label:            e.Label,
			Type:             e.Type,
			MatchEvent:       e.MatchEvent,
			MatchHost:        e.MatchHost,
			MatchApplication: e.MatchApplication,
			MatchProvider:    e.MatchProvider,
			Payload:          e.Payload,
			Enabled:          e.Enabled,
			CreatedAt:        e.CreatedAt,
			UpdatedAt:        e.UpdatedAt,
		})
	}

	var totalElements int64
	if totalElements, err = h.service.count(); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	totalPages := (totalElements + int64(queryParams.PageSize) - 1) / int64(queryParams.PageSize)
	c.JSON(http.StatusOK, api.NewDataResponseWithPayload(api.NewActionPageResponse(data, queryParams.Page, queryParams.PageSize, queryParams.OrderBy, queryParams.Order, totalElements, totalPages)))
}

func (h *actionHandler) get(c *gin.Context) {
	e, err := h.service.get(c.Param("id"))
	if err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewActionSingleResponse(e.ID, e.Label, e.Type, e.MatchEvent, e.MatchHost, e.MatchApplication, e.MatchProvider, e.Payload, e.Enabled, e.CreatedAt, e.UpdatedAt))
}

func (h *actionHandler) create(c *gin.Context) {
	var e *Action
	var err error

	var req api.CreateActionRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		errAbortWithValidatorPayload(c, err)
		return
	}

	if e, err = h.service.create(req.Label, api.ActionType(req.Type), req.MatchEvent, req.MatchHost, req.MatchApplication, req.MatchProvider, req.Payload, req.Enabled); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewActionSingleResponse(e.ID, e.Label, e.Type, e.MatchEvent, e.MatchHost, e.MatchApplication, e.MatchProvider, e.Payload, e.Enabled, e.CreatedAt, e.UpdatedAt))
}

func (h *actionHandler) updateLabel(c *gin.Context) {
	var e *Action
	var err error

	var req api.ModifyActionLabelRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		errAbortWithValidatorPayload(c, err)
		return
	}

	if e, err = h.service.updateLabel(c.Param("id"), req.Label); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewActionSingleResponse(e.ID, e.Label, e.Type, e.MatchEvent, e.MatchHost, e.MatchApplication, e.MatchProvider, e.Payload, e.Enabled, e.CreatedAt, e.UpdatedAt))
}

func (h *actionHandler) updateMatchEvent(c *gin.Context) {
	var e *Action
	var err error

	var req api.ModifyActionMatchEventRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		errAbortWithValidatorPayload(c, err)
		return
	}

	if e, err = h.service.updateMatchEvent(c.Param("id"), req.MatchEvent); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewActionSingleResponse(e.ID, e.Label, e.Type, e.MatchEvent, e.MatchHost, e.MatchApplication, e.MatchProvider, e.Payload, e.Enabled, e.CreatedAt, e.UpdatedAt))
}

func (h *actionHandler) updateMatchHost(c *gin.Context) {
	var e *Action
	var err error

	var req api.ModifyActionMatchHostRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		errAbortWithValidatorPayload(c, err)
		return
	}

	if e, err = h.service.updateMatchHost(c.Param("id"), req.MatchHost); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewActionSingleResponse(e.ID, e.Label, e.Type, e.MatchEvent, e.MatchHost, e.MatchApplication, e.MatchProvider, e.Payload, e.Enabled, e.CreatedAt, e.UpdatedAt))
}

func (h *actionHandler) updateMatchApplication(c *gin.Context) {
	var e *Action
	var err error

	var req api.ModifyActionMatchApplicationRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		errAbortWithValidatorPayload(c, err)
		return
	}

	if e, err = h.service.updateMatchApplication(c.Param("id"), req.MatchApplication); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewActionSingleResponse(e.ID, e.Label, e.Type, e.MatchEvent, e.MatchHost, e.MatchApplication, e.MatchProvider, e.Payload, e.Enabled, e.CreatedAt, e.UpdatedAt))
}

func (h *actionHandler) updateMatchProvider(c *gin.Context) {
	var e *Action
	var err error

	var req api.ModifyActionMatchProviderRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		errAbortWithValidatorPayload(c, err)
		return
	}

	if e, err = h.service.updateMatchProvider(c.Param("id"), req.MatchProvider); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewActionSingleResponse(e.ID, e.Label, e.Type, e.MatchEvent, e.MatchHost, e.MatchApplication, e.MatchProvider, e.Payload, e.Enabled, e.CreatedAt, e.UpdatedAt))
}

func (h *actionHandler) updatePayload(c *gin.Context) {
	var e *Action
	var err error

	var req api.ModifyActionTypeAndPayloadRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		errAbortWithValidatorPayload(c, err)
		return
	}

	if e, err = h.service.updateTypeAndPayload(c.Param("id"), req.Type, req.Payload); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewActionSingleResponse(e.ID, e.Label, e.Type, e.MatchEvent, e.MatchHost, e.MatchApplication, e.MatchProvider, e.Payload, e.Enabled, e.CreatedAt, e.UpdatedAt))
}

func (h *actionHandler) updateEnabled(c *gin.Context) {
	var e *Action
	var err error

	var req api.ModifyActionEnabledRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		errAbortWithValidatorPayload(c, err)
		return
	}

	if e, err = h.service.updateEnabled(c.Param("id"), req.Enabled); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewActionSingleResponse(e.ID, e.Label, e.Type, e.MatchEvent, e.MatchHost, e.MatchApplication, e.MatchProvider, e.Payload, e.Enabled, e.CreatedAt, e.UpdatedAt))
}

func (h *actionHandler) delete(c *gin.Context) {
	if err := h.service.delete(c.Param("id")); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	c.Header(headerContentType, headerContentTypeApplicationJson)
	c.Status(http.StatusNoContent)
}
