package server

import (
	"git.myservermanager.com/varakh/upda/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

type actionInvocationHandler struct {
	actionService           actionService
	actionInvocationService actionInvocationService
}

func newActionInvocationHandler(as *actionService, ais *actionInvocationService) *actionInvocationHandler {
	return &actionInvocationHandler{actionService: *as, actionInvocationService: *ais}
}

func (h *actionInvocationHandler) test(c *gin.Context) {
	var err error
	var req api.TestActionRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		errAbortWithValidatorPayload(c, err)
		return
	}

	var e *Action
	if e, err = h.actionService.get(c.Param("id")); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	err = h.actionInvocationService.execute(e, &eventPayloadInformationDto{Application: req.Application, Host: req.Host, Provider: req.Provider, Version: req.Version})

	isSuccess := err == nil
	var message string
	if err != nil {
		message = err.Error()
	}

	c.JSON(http.StatusOK, api.NewActionTestSingleResponse(isSuccess, message))
}

func (h *actionInvocationHandler) paginate(c *gin.Context) {
	var queryParams api.PaginateActionInvocationRequest
	var err error
	if err = c.ShouldBindQuery(&queryParams); err != nil {
		errAbortWithValidatorPayload(c, err)
		return
	}

	var actionInvocations []*ActionInvocation
	if actionInvocations, err = h.actionInvocationService.paginate(queryParams.Page, queryParams.PageSize, queryParams.OrderBy, queryParams.Order); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	var data []*api.ActionInvocationResponse
	data = make([]*api.ActionInvocationResponse, 0)

	for _, e := range actionInvocations {
		data = append(data, &api.ActionInvocationResponse{
			ID:         e.ID,
			RetryCount: e.RetryCount,
			State:      e.State,
			Message:    e.Message,
			ActionID:   e.ActionID,
			EventID:    e.EventID,
			CreatedAt:  e.CreatedAt,
			UpdatedAt:  e.UpdatedAt,
		})
	}

	var totalElements int64
	if totalElements, err = h.actionInvocationService.count(); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	totalPages := (totalElements + int64(queryParams.PageSize) - 1) / int64(queryParams.PageSize)
	c.JSON(http.StatusOK, api.NewDataResponseWithPayload(api.NewActionInvocationPageResponse(data, queryParams.Page, queryParams.PageSize, queryParams.OrderBy, queryParams.Order, totalElements, totalPages)))
}

func (h *actionInvocationHandler) get(c *gin.Context) {
	e, err := h.actionInvocationService.get(c.Param("id"))
	if err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewActionInvocationSingleResponse(e.ID, e.RetryCount, e.State, e.Message, e.ActionID, e.EventID, e.CreatedAt, e.UpdatedAt))
}

func (h *actionInvocationHandler) delete(c *gin.Context) {
	if err := h.actionInvocationService.delete(c.Param("id")); err != nil {
		_ = c.AbortWithError(errToHttpStatus(err), err)
		return
	}

	c.Header(headerContentType, headerContentTypeApplicationJson)
	c.Status(http.StatusNoContent)
}
