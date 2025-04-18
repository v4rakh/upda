package handler

import (
	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/internal/server/config"
	"git.myservermanager.com/varakh/upda/internal/server/constant"
	"github.com/gin-gonic/gin"
	"net/http"
)

type InfoHandler struct {
	appConfig config.App
}

func NewInfoHandler(a *config.App) *InfoHandler {
	return &InfoHandler{appConfig: *a}
}

func (h *InfoHandler) Show(c *gin.Context) {
	c.JSON(http.StatusOK, api.NewDataResponseWithPayload(api.NewInfoResponse(constant.AppName, constant.AppVersion, h.appConfig.TimeZone)))
}
