package handler

import (
	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SecretHandler struct {
	service service.SecretService
}

func NewSecretHandler(s *service.SecretService) *SecretHandler {
	return &SecretHandler{service: *s}
}

func (h *SecretHandler) GetAll(c *gin.Context) {
	var secrets []*model.Secret
	var err error

	if secrets, err = h.service.GetAll(); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
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

func (h *SecretHandler) Create(c *gin.Context) {
	var e *model.Secret
	var err error

	var req api.CreateSecretRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	if e, err = h.service.Insert(req.Key, req.Value); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewSecretSingleResponse(e.ID, e.Key, e.Value, e.CreatedAt, e.UpdatedAt))
}

func (h *SecretHandler) UpdateValue(c *gin.Context) {
	var e *model.Secret
	var err error

	var req api.ModifySecretValueRequest

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

	c.JSON(http.StatusOK, api.NewSecretSingleResponse(e.ID, e.Key, e.Value, e.CreatedAt, e.UpdatedAt))
}

func (h *SecretHandler) Get(c *gin.Context) {
	var err error
	var pathParams api.IDUriRequest

	if err = c.ShouldBindUri(&pathParams); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	var e *model.Secret
	if e, err = h.service.Get(pathParams.ID); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewSecretSingleResponse(e.ID, e.Key, "", e.CreatedAt, e.UpdatedAt))
}

func (h *SecretHandler) Delete(c *gin.Context) {
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
