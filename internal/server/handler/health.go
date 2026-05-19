package handler

import (
	"net/http"

	"git.myservermanager.com/varakh/upda/api"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Show(c *gin.Context) {
	c.JSON(http.StatusOK, api.NewDataResponseWithPayload(api.NewHealthResponse(true)))
}
