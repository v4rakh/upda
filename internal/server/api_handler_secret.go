package server

import (
	"git.myservermanager.com/varakh/upda/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

type secretHandler struct {
	service secretService
}

func newSecretHandler(s *secretService) *secretHandler {
	return &secretHandler{service: *s}
}

func (h *secretHandler) getAll(c *gin.Context) {
	var secrets []*Secret
	var err error

	if secrets, err = h.service.getAll(); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	var data []*api.SecretResponse
	data = make([]*api.SecretResponse, 0, len(secrets))

	for _, e := range secrets {
		data = append(data, &api.SecretResponse{
			ID:        e.ID,
			Key:       e.Key,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, api.NewDataResponseWithPayload(api.NewSecretPageResponse(data)))
}

func (h *secretHandler) create(c *gin.Context) {
	var e *Secret
	var err error

	var req api.CreateSecretRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		errAbortWithValidatorPayload(c, err)
		return
	}

	if e, err = h.service.insert(req.Key, req.Value); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewSecretSingleResponse(e.ID, e.Key, e.Value, e.CreatedAt, e.UpdatedAt))
}

func (h *secretHandler) updateValue(c *gin.Context) {
	var e *Secret
	var err error

	var req api.ModifySecretValueRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		errAbortWithValidatorPayload(c, err)
		return
	}

	if e, err = h.service.updateValue(c.Param("id"), req.Value); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewSecretSingleResponse(e.ID, e.Key, e.Value, e.CreatedAt, e.UpdatedAt))
}

func (h *secretHandler) get(c *gin.Context) {
	e, err := h.service.get(c.Param("id"))
	if err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewSecretSingleResponse(e.ID, e.Key, "", e.CreatedAt, e.UpdatedAt))
}

func (h *secretHandler) delete(c *gin.Context) {
	if err := h.service.delete(c.Param("id")); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	c.Header(api.HeaderContentType, api.HeaderContentTypeApplicationJson)
	c.Status(http.StatusNoContent)
}
