package server

import (
	"git.myservermanager.com/varakh/upda/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

type updateHandler struct {
	service   updateService
	appConfig appConfig
}

func newUpdateHandler(s *updateService, c *appConfig) *updateHandler {
	return &updateHandler{service: *s, appConfig: *c}
}

func (h *updateHandler) paginate(c *gin.Context) {
	var queryParams api.PaginateUpdateRequest

	if err := c.ShouldBindQuery(&queryParams); err != nil {
		errAbortWithValidatorPayload(c, err)
		return
	}

	var updates []*Update
	var err error

	s, stateQueryContainsAtLeastOne := c.GetQueryArray("state")

	states := make([]api.UpdateState, 0)
	if stateQueryContainsAtLeastOne {
		for _, state := range s {
			states = append(states, api.UpdateState(state))
		}
	}

	if updates, err = h.service.paginate(queryParams.Page, queryParams.PageSize, queryParams.OrderBy, queryParams.Order, queryParams.SearchTerm, queryParams.SearchIn, states...); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
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
	if totalElements, err = h.service.count(queryParams.SearchTerm, queryParams.SearchIn, states...); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	totalPages := (totalElements + int64(queryParams.PageSize) - 1) / int64(queryParams.PageSize)
	c.JSON(http.StatusOK, api.NewDataResponseWithPayload(api.NewUpdatePageResponse(data, queryParams.Page, queryParams.PageSize, queryParams.OrderBy, queryParams.Order, totalElements, totalPages)))
}

func (h *updateHandler) get(c *gin.Context) {
	e, err := h.service.get(c.Param("id"))
	if err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewUpdateSingleResponse(e.ID, e.Application, e.Provider, e.Host, e.Version, e.State, e.CreatedAt, e.UpdatedAt, e.Metadata))
}

func (h *updateHandler) updateState(c *gin.Context) {
	var e *Update
	var err error

	var req api.ModifyUpdateStateRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		errAbortWithValidatorPayload(c, err)
		return
	}

	if e, err = h.service.updateState(c.Param("id"), api.UpdateState(req.State)); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewUpdateSingleResponse(e.ID, e.Application, e.Provider, e.Host, e.Version, e.State, e.CreatedAt, e.UpdatedAt, e.Metadata))
}

func (h *updateHandler) delete(c *gin.Context) {
	if err := h.service.delete(c.Param("id")); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	c.Header(api.HeaderContentType, api.HeaderContentTypeApplicationJson)
	c.Status(http.StatusNoContent)
}
