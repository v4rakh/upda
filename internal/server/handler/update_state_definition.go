package handler

import (
	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/internal/server/dto"
	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UpdateStateDefinitionHandler struct {
	service service.UpdateStateDefinitionService
}

func NewUpdateStateDefinitionHandler(s *service.UpdateStateDefinitionService) *UpdateStateDefinitionHandler {
	return &UpdateStateDefinitionHandler{service: *s}
}

func (h *UpdateStateDefinitionHandler) GetAll(c *gin.Context) {
	var states []*model.UpdateStateDefinition
	var err error

	if states, err = h.service.GetAll(); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	var data []*api.UpdateStateDefinitionResponse
	data = make([]*api.UpdateStateDefinitionResponse, 0, len(states))

	for _, e := range states {
		data = append(data, &api.UpdateStateDefinitionResponse{
			ID:               e.ID.String(),
			Name:             e.Name,
			Label:            e.Label,
			Color:            e.Color,
			Icon:             e.Icon,
			Description:      e.Description,
			IsInitial:        e.IsInitial,
			SkipOnNewVersion: e.SkipOnNewVersion,
			SortOrder:        e.SortOrder,
			CreatedAt:        e.CreatedAt,
			UpdatedAt:        e.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, api.NewDataResponseWithPayload(api.NewUpdateStateDefinitionPageResponse(data)))
}

func (h *UpdateStateDefinitionHandler) Get(c *gin.Context) {
	var err error
	var pathParams api.IDUriRequest

	if err = c.ShouldBindUri(&pathParams); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	var e *model.UpdateStateDefinition
	if e, err = h.service.Get(pathParams.ID); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewUpdateStateDefinitionSingleResponse(e.ID.String(), e.Name, e.Label, e.Color, e.Icon, e.Description, e.IsInitial, e.SkipOnNewVersion, e.SortOrder, e.CreatedAt, e.UpdatedAt))
}

func (h *UpdateStateDefinitionHandler) Create(c *gin.Context) {
	var e *model.UpdateStateDefinition
	var err error

	var req api.CreateUpdateStateDefinitionRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	if e, err = h.service.Create(req.Name, req.Label, req.Color, req.Icon, req.Description, req.IsInitial, req.SkipOnNewVersion); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewUpdateStateDefinitionSingleResponse(e.ID.String(), e.Name, e.Label, e.Color, e.Icon, e.Description, e.IsInitial, e.SkipOnNewVersion, e.SortOrder, e.CreatedAt, e.UpdatedAt))
}

func (h *UpdateStateDefinitionHandler) Update(c *gin.Context) {
	var e *model.UpdateStateDefinition
	var err error

	var req api.ModifyUpdateStateDefinitionRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	var pathParams api.IDUriRequest

	if err = c.ShouldBindUri(&pathParams); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	if e, err = h.service.Update(pathParams.ID, req.Name, req.Label, req.Color, req.Icon, req.Description, req.IsInitial, req.SkipOnNewVersion, req.SortOrder); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewUpdateStateDefinitionSingleResponse(e.ID.String(), e.Name, e.Label, e.Color, e.Icon, e.Description, e.IsInitial, e.SkipOnNewVersion, e.SortOrder, e.CreatedAt, e.UpdatedAt))
}

func (h *UpdateStateDefinitionHandler) Delete(c *gin.Context) {
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

func (h *UpdateStateDefinitionHandler) Reorder(c *gin.Context) {
	var err error
	var req api.ReorderUpdateStateDefinitionsRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	items := make([]dto.UpdateStateReorderItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = dto.UpdateStateReorderItem{
			ID:        item.ID,
			SortOrder: item.SortOrder,
		}
	}

	if err = h.service.Reorder(items); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.Header(api.HeaderContentType, api.HeaderContentTypeApplicationJson)
	c.Status(http.StatusNoContent)
}
