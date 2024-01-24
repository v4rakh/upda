package server

import (
	"git.myservermanager.com/varakh/upda/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

type infoHandler struct {
	appConfig appConfig
}

func newInfoHandler(a *appConfig) *infoHandler {
	return &infoHandler{appConfig: *a}
}

func (h *infoHandler) show(c *gin.Context) {
	c.JSON(http.StatusOK, api.DataResponse{Data: gin.H{
		"name":     Name,
		"version":  Version,
		"timeZone": h.appConfig.timeZone,
	}})
}
