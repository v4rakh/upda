package server

import (
	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/commons"
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
	c.JSON(http.StatusOK, api.NewDataResponseWithPayload(api.NewInfoResponse(commons.Name, commons.Version, h.appConfig.timeZone)))
}
