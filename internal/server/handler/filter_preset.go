package handler

import (
	"net/http"

	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/internal/server/constant"
	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/service"
	"github.com/gin-gonic/gin"
)

type FilterPresetHandler struct {
	service service.FilterPresetService
}

func NewFilterPresetHandler(s *service.FilterPresetService) *FilterPresetHandler {
	return &FilterPresetHandler{service: *s}
}

func (h *FilterPresetHandler) GetByType(c *gin.Context) {
	var filters []*model.FilterPreset
	var err error

	var pathParams api.FilterPresetUriRequest
	if err = c.ShouldBindUri(&pathParams); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	if filters, err = h.service.GetByType(constant.FilterPresetType(pathParams.Type)); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	var data []*api.FilterPresetResponse
	data = make([]*api.FilterPresetResponse, 0, len(filters))

	for _, e := range filters {
		data = append(data, &api.FilterPresetResponse{
			ID:         e.ID.String(),
			Type:       e.Type,
			Label:      e.Label,
			Parameters: e.Parameters,
			Color:      e.Color,
			CreatedAt:  e.CreatedAt,
			UpdatedAt:  e.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, api.NewDataResponseWithPayload(api.NewFilterPresetPageResponse(data)))
}

func (h *FilterPresetHandler) Create(c *gin.Context) {
	var e *model.FilterPreset
	var err error

	var req api.CreateFilterPresetRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	if e, err = h.service.Create(constant.FilterPresetType(req.Type), req.Label, req.Parameters, req.Color); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewFilterPresetSingleResponse(e.ID.String(), e.Type, e.Label, e.Parameters, e.Color, e.CreatedAt, e.UpdatedAt))
}

func (h *FilterPresetHandler) Delete(c *gin.Context) {
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
