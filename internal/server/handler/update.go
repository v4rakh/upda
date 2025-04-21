package handler

import (
	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/internal/server/config"
	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UpdateHandler struct {
	service   service.UpdateService
	appConfig config.App
}

func NewUpdateHandler(s *service.UpdateService, c *config.App) *UpdateHandler {
	return &UpdateHandler{service: *s, appConfig: *c}
}

func (h *UpdateHandler) Paginate(c *gin.Context) {
	var queryParams api.PaginateUpdateRequest

	if err := c.ShouldBindQuery(&queryParams); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	var updates []*model.Update
	var err error

	s, stateQueryContainsAtLeastOne := c.GetQueryArray("state")

	states := make([]api.UpdateState, 0)
	if stateQueryContainsAtLeastOne {
		for _, state := range s {
			states = append(states, api.UpdateState(state))
		}
	}

	if updates, err = h.service.Paginate(queryParams.Page, queryParams.PageSize, queryParams.OrderBy, queryParams.Order, queryParams.SearchTerm, queryParams.SearchIn, states...); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	var data []*api.UpdateResponse
	data = make([]*api.UpdateResponse, 0)

	for _, e := range updates {
		data = append(data, &api.UpdateResponse{
			ID:          e.ID,
			Application: e.Application,
			Provider:    e.Provider,
			Host:        e.Host,
			Version:     e.Version,
			State:       e.State,
			CreatedAt:   e.CreatedAt,
			UpdatedAt:   e.UpdatedAt,
		})
	}

	var totalElements int64
	if totalElements, err = h.service.Count(queryParams.SearchTerm, queryParams.SearchIn, states...); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	totalPages := (totalElements + int64(queryParams.PageSize) - 1) / int64(queryParams.PageSize)
	c.JSON(http.StatusOK, api.NewDataResponseWithPayload(api.NewUpdatePageResponse(data, queryParams.Page, queryParams.PageSize, queryParams.OrderBy, queryParams.Order, totalElements, totalPages)))
}

func (h *UpdateHandler) Get(c *gin.Context) {
	var err error
	var pathParams api.IDUriRequest

	if err = c.ShouldBindUri(&pathParams); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	var e *model.Update
	if e, err = h.service.Get(pathParams.ID); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewUpdateSingleResponse(e.ID, e.Application, e.Provider, e.Host, e.Version, e.State, e.CreatedAt, e.UpdatedAt, e.Metadata))
}

func (h *UpdateHandler) UpdateState(c *gin.Context) {
	var e *model.Update
	var err error

	var req api.ModifyUpdateStateRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	var pathParams api.IDUriRequest

	if err = c.ShouldBindUri(&pathParams); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	if e, err = h.service.UpdateState(pathParams.ID, api.UpdateState(req.State)); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewUpdateSingleResponse(e.ID, e.Application, e.Provider, e.Host, e.Version, e.State, e.CreatedAt, e.UpdatedAt, e.Metadata))
}

func (h *UpdateHandler) Delete(c *gin.Context) {
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
