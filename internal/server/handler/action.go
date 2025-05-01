package handler

import (
	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/internal/server/constant"
	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ActionHandler struct {
	service service.ActionService
}

func NewActionHandler(s *service.ActionService) *ActionHandler {
	return &ActionHandler{service: *s}
}

func (h *ActionHandler) Paginate(c *gin.Context) {
	var queryParams api.PaginateActionRequest
	var err error
	if err = c.ShouldBindQuery(&queryParams); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	var actions []*model.Action
	if actions, err = h.service.Paginate(queryParams.Page, queryParams.PageSize, queryParams.OrderBy, queryParams.Order); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	var data []*api.ActionResponse
	data = make([]*api.ActionResponse, 0, len(actions))

	for _, e := range actions {
		data = append(data, &api.ActionResponse{
			ID:               e.ID.String(),
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
	if totalElements, err = h.service.Count(); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	totalPages := (totalElements + int64(queryParams.PageSize) - 1) / int64(queryParams.PageSize)
	c.JSON(http.StatusOK, api.NewDataResponseWithPayload(api.NewActionPageResponse(data, queryParams.Page, queryParams.PageSize, queryParams.OrderBy, queryParams.Order, totalElements, totalPages)))
}

func (h *ActionHandler) Get(c *gin.Context) {
	var err error
	var pathParams api.IDUriRequest

	if err = c.ShouldBindUri(&pathParams); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	var e *model.Action
	if e, err = h.service.Get(pathParams.ID); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewActionSingleResponse(e.ID.String(), e.Label, e.Type, e.MatchEvent, e.MatchHost, e.MatchApplication, e.MatchProvider, e.Payload, e.Enabled, e.CreatedAt, e.UpdatedAt))
}

func (h *ActionHandler) Create(c *gin.Context) {
	var e *model.Action
	var err error

	var req api.CreateActionRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	if e, err = h.service.Create(req.Label, constant.ActionType(req.Type), req.MatchEvent, req.MatchHost, req.MatchApplication, req.MatchProvider, req.Payload, req.Enabled); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewActionSingleResponse(e.ID.String(), e.Label, e.Type, e.MatchEvent, e.MatchHost, e.MatchApplication, e.MatchProvider, e.Payload, e.Enabled, e.CreatedAt, e.UpdatedAt))
}

func (h *ActionHandler) UpdateLabel(c *gin.Context) {
	var e *model.Action
	var err error

	var req api.ModifyActionLabelRequest

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

	c.JSON(http.StatusOK, api.NewActionSingleResponse(e.ID.String(), e.Label, e.Type, e.MatchEvent, e.MatchHost, e.MatchApplication, e.MatchProvider, e.Payload, e.Enabled, e.CreatedAt, e.UpdatedAt))
}

func (h *ActionHandler) UpdateMatchEvent(c *gin.Context) {
	var e *model.Action
	var err error

	var req api.ModifyActionMatchEventRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	var pathParams api.IDUriRequest

	if err = c.ShouldBindUri(&pathParams); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	if e, err = h.service.UpdateMatchEvent(pathParams.ID, req.MatchEvent); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewActionSingleResponse(e.ID.String(), e.Label, e.Type, e.MatchEvent, e.MatchHost, e.MatchApplication, e.MatchProvider, e.Payload, e.Enabled, e.CreatedAt, e.UpdatedAt))
}

func (h *ActionHandler) UpdateMatchHost(c *gin.Context) {
	var e *model.Action
	var err error

	var req api.ModifyActionMatchHostRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	var pathParams api.IDUriRequest

	if err = c.ShouldBindUri(&pathParams); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	if e, err = h.service.UpdateMatchHost(pathParams.ID, req.MatchHost); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewActionSingleResponse(e.ID.String(), e.Label, e.Type, e.MatchEvent, e.MatchHost, e.MatchApplication, e.MatchProvider, e.Payload, e.Enabled, e.CreatedAt, e.UpdatedAt))
}

func (h *ActionHandler) UpdateMatchApplication(c *gin.Context) {
	var e *model.Action
	var err error

	var req api.ModifyActionMatchApplicationRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	var pathParams api.IDUriRequest

	if err = c.ShouldBindUri(&pathParams); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	if e, err = h.service.UpdateMatchApplication(pathParams.ID, req.MatchApplication); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewActionSingleResponse(e.ID.String(), e.Label, e.Type, e.MatchEvent, e.MatchHost, e.MatchApplication, e.MatchProvider, e.Payload, e.Enabled, e.CreatedAt, e.UpdatedAt))
}

func (h *ActionHandler) UpdateMatchProvider(c *gin.Context) {
	var e *model.Action
	var err error

	var req api.ModifyActionMatchProviderRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	var pathParams api.IDUriRequest

	if err = c.ShouldBindUri(&pathParams); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	if e, err = h.service.UpdateMatchProvider(pathParams.ID, req.MatchProvider); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewActionSingleResponse(e.ID.String(), e.Label, e.Type, e.MatchEvent, e.MatchHost, e.MatchApplication, e.MatchProvider, e.Payload, e.Enabled, e.CreatedAt, e.UpdatedAt))
}

func (h *ActionHandler) UpdatePayload(c *gin.Context) {
	var e *model.Action
	var err error

	var req api.ModifyActionTypeAndPayloadRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	var pathParams api.IDUriRequest

	if err = c.ShouldBindUri(&pathParams); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	if e, err = h.service.UpdateTypeAndPayload(pathParams.ID, constant.ActionType(req.Type), req.Payload); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewActionSingleResponse(e.ID.String(), e.Label, e.Type, e.MatchEvent, e.MatchHost, e.MatchApplication, e.MatchProvider, e.Payload, e.Enabled, e.CreatedAt, e.UpdatedAt))
}

func (h *ActionHandler) UpdateEnabled(c *gin.Context) {
	var e *model.Action
	var err error

	var req api.ModifyActionEnabledRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	var pathParams api.IDUriRequest

	if err = c.ShouldBindUri(&pathParams); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	if e, err = h.service.UpdateEnabled(pathParams.ID, req.Enabled); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewActionSingleResponse(e.ID.String(), e.Label, e.Type, e.MatchEvent, e.MatchHost, e.MatchApplication, e.MatchProvider, e.Payload, e.Enabled, e.CreatedAt, e.UpdatedAt))
}

func (h *ActionHandler) Delete(c *gin.Context) {
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
