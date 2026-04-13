package handler

import (
	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UpdateStateTransitionHandler struct {
	service service.UpdateStateTransitionService
}

func NewUpdateStateTransitionHandler(s *service.UpdateStateTransitionService) *UpdateStateTransitionHandler {
	return &UpdateStateTransitionHandler{service: *s}
}

func (h *UpdateStateTransitionHandler) GetAll(c *gin.Context) {
	var transitions []*model.UpdateStateTransition
	var err error

	if transitions, err = h.service.GetAll(); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	var data []*api.UpdateStateTransitionResponse
	data = make([]*api.UpdateStateTransitionResponse, 0, len(transitions))

	for _, e := range transitions {
		data = append(data, h.toResponse(e))
	}

	c.JSON(http.StatusOK, api.NewDataResponseWithPayload(api.NewUpdateStateTransitionPageResponse(data)))
}

func (h *UpdateStateTransitionHandler) GetByFromStateId(c *gin.Context) {
	var err error
	var pathParams api.StateIdUriRequest

	if err = c.ShouldBindUri(&pathParams); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	var transitions []*model.UpdateStateTransition
	if transitions, err = h.service.GetByFromStateId(pathParams.StateID); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	var data []*api.UpdateStateTransitionResponse
	data = make([]*api.UpdateStateTransitionResponse, 0, len(transitions))

	for _, e := range transitions {
		data = append(data, h.toResponse(e))
	}

	c.JSON(http.StatusOK, api.NewDataResponseWithPayload(api.NewUpdateStateTransitionPageResponse(data)))
}

func (h *UpdateStateTransitionHandler) Create(c *gin.Context) {
	var e *model.UpdateStateTransition
	var err error

	var req api.CreateUpdateStateTransitionRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	if e, err = h.service.Create(req.FromStateId, req.ToStateId); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewUpdateStateTransitionSingleResponse(
		e.ID.String(),
		api.UpdateStateDefinitionResponse{
			ID:               e.FromState.ID.String(),
			Name:             e.FromState.Name,
			Label:            e.FromState.Label,
			Color:            e.FromState.Color,
			Icon:             e.FromState.Icon,
			Description:      e.FromState.Description,
			IsInitial:        e.FromState.IsInitial,
			SkipOnNewVersion: e.FromState.SkipOnNewVersion,
			SortOrder:        e.FromState.SortOrder,
			CreatedAt:        e.FromState.CreatedAt,
			UpdatedAt:        e.FromState.UpdatedAt,
		},
		api.UpdateStateDefinitionResponse{
			ID:               e.ToState.ID.String(),
			Name:             e.ToState.Name,
			Label:            e.ToState.Label,
			Color:            e.ToState.Color,
			Icon:             e.ToState.Icon,
			Description:      e.ToState.Description,
			IsInitial:        e.ToState.IsInitial,
			SkipOnNewVersion: e.ToState.SkipOnNewVersion,
			SortOrder:        e.ToState.SortOrder,
			CreatedAt:        e.ToState.CreatedAt,
			UpdatedAt:        e.ToState.UpdatedAt,
		},
		e.CreatedAt,
		e.UpdatedAt,
	))
}

func (h *UpdateStateTransitionHandler) Delete(c *gin.Context) {
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

func (h *UpdateStateTransitionHandler) toResponse(e *model.UpdateStateTransition) *api.UpdateStateTransitionResponse {
	return &api.UpdateStateTransitionResponse{
		ID: e.ID.String(),
		FromState: api.UpdateStateDefinitionResponse{
			ID:               e.FromState.ID.String(),
			Name:             e.FromState.Name,
			Label:            e.FromState.Label,
			Color:            e.FromState.Color,
			Icon:             e.FromState.Icon,
			Description:      e.FromState.Description,
			IsInitial:        e.FromState.IsInitial,
			SkipOnNewVersion: e.FromState.SkipOnNewVersion,
			SortOrder:        e.FromState.SortOrder,
			CreatedAt:        e.FromState.CreatedAt,
			UpdatedAt:        e.FromState.UpdatedAt,
		},
		ToState: api.UpdateStateDefinitionResponse{
			ID:               e.ToState.ID.String(),
			Name:             e.ToState.Name,
			Label:            e.ToState.Label,
			Color:            e.ToState.Color,
			Icon:             e.ToState.Icon,
			Description:      e.ToState.Description,
			IsInitial:        e.ToState.IsInitial,
			SkipOnNewVersion: e.ToState.SkipOnNewVersion,
			SortOrder:        e.ToState.SortOrder,
			CreatedAt:        e.ToState.CreatedAt,
			UpdatedAt:        e.ToState.UpdatedAt,
		},
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
