package handler

import (
	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/internal/server/dto"
	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ActionInvocationHandler struct {
	actionService           service.ActionService
	actionInvocationService service.ActionInvocationService
}

func NewActionInvocationHandler(as *service.ActionService, ais *service.ActionInvocationService) *ActionInvocationHandler {
	return &ActionInvocationHandler{actionService: *as, actionInvocationService: *ais}
}

func (h *ActionInvocationHandler) Test(c *gin.Context) {
	var err error
	var req api.TestActionRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	var pathParams api.IDUriRequest

	if err = c.ShouldBindUri(&pathParams); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	var e *model.Action
	if e, err = h.actionService.Get(pathParams.ID); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	err = h.actionInvocationService.Execute(e, &dto.EventPayloadInformationDto{Application: req.Application, Host: req.Host, Provider: req.Provider, Version: req.Version, State: req.State})

	isSuccess := err == nil
	var message string
	if err != nil {
		message = err.Error()
	}

	c.JSON(http.StatusOK, api.NewActionTestSingleResponse(isSuccess, message))
}

func (h *ActionInvocationHandler) Paginate(c *gin.Context) {
	var queryParams api.PaginateActionInvocationRequest
	var err error
	if err = c.ShouldBindQuery(&queryParams); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	var actionInvocations []*model.ActionInvocation
	if actionInvocations, err = h.actionInvocationService.Paginate(queryParams.Page, queryParams.PageSize, queryParams.OrderBy, queryParams.Order); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	var data []*api.ActionInvocationResponse
	data = make([]*api.ActionInvocationResponse, 0, len(actionInvocations))

	for _, e := range actionInvocations {
		data = append(data, &api.ActionInvocationResponse{
			ID:         e.ID.String(),
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
	if totalElements, err = h.actionInvocationService.Count(); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	totalPages := (totalElements + int64(queryParams.PageSize) - 1) / int64(queryParams.PageSize)
	c.JSON(http.StatusOK, api.NewDataResponseWithPayload(api.NewActionInvocationPageResponse(data, queryParams.Page, queryParams.PageSize, queryParams.OrderBy, queryParams.Order, totalElements, totalPages)))
}

func (h *ActionInvocationHandler) Get(c *gin.Context) {
	var err error
	var pathParams api.IDUriRequest

	if err = c.ShouldBindUri(&pathParams); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	var e *model.ActionInvocation
	if e, err = h.actionInvocationService.Get(pathParams.ID); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewActionInvocationSingleResponse(e.ID.String(), e.RetryCount, e.State, e.Message, e.ActionID, e.EventID, e.CreatedAt, e.UpdatedAt))
}

func (h *ActionInvocationHandler) Delete(c *gin.Context) {
	var err error
	var pathParams api.IDUriRequest

	if err = c.ShouldBindUri(&pathParams); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	if err = h.actionInvocationService.Delete(pathParams.ID); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.Header(api.HeaderContentType, api.HeaderContentTypeApplicationJson)
	c.Status(http.StatusNoContent)
}
