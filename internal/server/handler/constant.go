package handler

import (
	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ConstantHandler struct {
	service service.ConstantService
}

func NewConstantHandler(s *service.ConstantService) *ConstantHandler {
	return &ConstantHandler{service: *s}
}

func (h *ConstantHandler) GetAll(c *gin.Context) {
	var constants []*model.Constant
	var err error

	if constants, err = h.service.GetAll(); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	var data []*api.ConstantResponse
	data = make([]*api.ConstantResponse, 0, len(constants))

	for _, e := range constants {
		data = append(data, &api.ConstantResponse{
			ID:        e.ID.String(),
			Key:       e.Key,
			Value:     e.Value,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, api.NewDataResponseWithPayload(api.NewConstantPageResponse(data)))
}

func (h *ConstantHandler) Create(c *gin.Context) {
	var e *model.Constant
	var err error

	var req api.CreateConstantRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	if e, err = h.service.Insert(req.Key, req.Value); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewConstantSingleResponse(e.ID.String(), e.Key, e.Value, e.CreatedAt, e.UpdatedAt))
}

func (h *ConstantHandler) UpdateValue(c *gin.Context) {
	var e *model.Constant
	var err error

	var req api.ModifyConstantValueRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	var pathParams api.IDUriRequest

	if err = c.ShouldBindUri(&pathParams); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}
	if e, err = h.service.UpdateValue(pathParams.ID, req.Value); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewConstantSingleResponse(e.ID.String(), e.Key, e.Value, e.CreatedAt, e.UpdatedAt))
}

func (h *ConstantHandler) Get(c *gin.Context) {
	var err error
	var pathParams api.IDUriRequest

	if err = c.ShouldBindUri(&pathParams); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	var e *model.Constant

	if e, err = h.service.Get(pathParams.ID); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewConstantSingleResponse(e.ID.String(), e.Key, e.Value, e.CreatedAt, e.UpdatedAt))
}

func (h *ConstantHandler) Delete(c *gin.Context) {
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
