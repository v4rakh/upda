package handler

import (
	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/internal/server/constant"
	"git.myservermanager.com/varakh/upda/internal/server/dto"
	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CommentHandler struct {
	updateService  service.UpdateService
	commentService service.CommentService
}

func NewCommentHandler(us *service.UpdateService, cs *service.CommentService) *CommentHandler {
	return &CommentHandler{updateService: *us, commentService: *cs}
}

func (h *CommentHandler) GetAllByUpdateId(c *gin.Context) {
	var err error

	var pathParams api.UpdateIDUriRequest
	if err = c.ShouldBindUri(&pathParams); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	var queryParams api.PaginateCommentRequest
	if err = c.ShouldBindQuery(&queryParams); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	var update *model.Update
	if update, err = h.updateService.Get(pathParams.ID); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	var comments []*model.Comment
	if comments, err = h.commentService.GetByUpdateId(update.ID.String(), queryParams.Page, queryParams.PageSize); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	var data []*api.CommentResponse
	data = make([]*api.CommentResponse, 0, len(comments))

	for _, e := range comments {
		data = append(data, &api.CommentResponse{
			ID:        e.ID.String(),
			Author:    e.Author,
			Content:   e.Content,
			UpdateID:  e.UpdateID,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		})
	}

	var totalElements int64
	if totalElements, err = h.commentService.CountByUpdateId(update.ID.String()); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}
	totalPages := (totalElements + int64(queryParams.PageSize) - 1) / int64(queryParams.PageSize)

	c.JSON(http.StatusOK, api.NewDataResponseWithPayload(api.NewCommentPageResponse(data, queryParams.Page, queryParams.PageSize, totalElements, totalPages)))
}

func (h *CommentHandler) Create(c *gin.Context) {
	var err error
	var pathParams api.UpdateIDUriRequest

	if err = c.ShouldBindUri(&pathParams); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	var update *model.Update
	if update, err = h.updateService.Get(pathParams.ID); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	var req api.CreateCommentRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	session := c.MustGet(constant.GinContextSession).(dto.ContextSession)

	var e *model.Comment
	if e, err = h.commentService.Create(session.User, req.Content, update); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewCommentSingleResponse(e.ID.String(), e.Author, e.Content, e.UpdateID, e.CreatedAt, e.UpdatedAt))
}

func (h *CommentHandler) UpdateContent(c *gin.Context) {
	var e *model.Comment
	var err error

	var req api.ModifyCommentContentRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	var pathParams api.IDUriRequest

	if err = c.ShouldBindUri(&pathParams); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}
	if e, err = h.commentService.UpdateContent(pathParams.ID, req.Content); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, api.NewCommentSingleResponse(e.ID.String(), e.Author, e.Content, e.UpdateID, e.CreatedAt, e.UpdatedAt))
}

func (h *CommentHandler) Delete(c *gin.Context) {
	var err error
	var pathParams api.IDUriRequest

	if err = c.ShouldBindUri(&pathParams); err != nil {
		AbortWithValidatorPayload(c, err)
		return
	}

	if err = h.commentService.Delete(pathParams.ID); err != nil {
		_ = c.AbortWithError(ToHttpStatus(err), err)
		return
	}

	c.Header(api.HeaderContentType, api.HeaderContentTypeApplicationJson)
	c.Status(http.StatusNoContent)
}
