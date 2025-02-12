package server

import (
	"git.myservermanager.com/varakh/upda/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

type constantHandler struct {
	service constantService
}

func newConstantHandler(s *constantService) *constantHandler {
	return &constantHandler{service: *s}
}

func (h *constantHandler) getAll(c *gin.Context) {
	var constants []*Constant
	var err error

	if constants, err = h.service.getAll(); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	var data []*api.ConstantResponse
	data = make([]*api.ConstantResponse, 0, len(constants))

	for _, e := range constants {
		data = append(data, &api.ConstantResponse{
			ID:        e.ID,
			Key:       e.Key,
			Value:     e.Value,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, api.NewDataResponseWithPayload(api.NewConstantPageResponse(data)))
}

func (h *constantHandler) create(c *gin.Context) {
	var e *Constant
	var err error

	var req api.CreateConstantRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		errAbortWithValidatorPayload(c, err)
		return
	}

	if e, err = h.service.insert(req.Key, req.Value); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewConstantSingleResponse(e.ID, e.Key, e.Value, e.CreatedAt, e.UpdatedAt))
}

func (h *constantHandler) updateValue(c *gin.Context) {
	var e *Constant
	var err error

	var req api.ModifyConstantValueRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		errAbortWithValidatorPayload(c, err)
		return
	}

	if e, err = h.service.updateValue(c.Param("id"), req.Value); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewConstantSingleResponse(e.ID, e.Key, e.Value, e.CreatedAt, e.UpdatedAt))
}

func (h *constantHandler) get(c *gin.Context) {
	e, err := h.service.get(c.Param("id"))
	if err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewConstantSingleResponse(e.ID, e.Key, e.Value, e.CreatedAt, e.UpdatedAt))
}

func (h *constantHandler) delete(c *gin.Context) {
	if err := h.service.delete(c.Param("id")); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	c.Header(api.HeaderContentType, api.HeaderContentTypeApplicationJson)
	c.Status(http.StatusNoContent)
}
