package server

import (
	"git.myservermanager.com/varakh/upda/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

type healthHandler struct {
}

func newHealthHandler() *healthHandler {
	return &healthHandler{}
}

func (h *healthHandler) show(c *gin.Context) {
	c.JSON(http.StatusOK, api.DataResponse{Data: gin.H{
		"healthy": true,
	}})
}
