package handler

import (
	"net/http"

	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/internal/meta"
	"git.myservermanager.com/varakh/upda/internal/server/config"
	"github.com/gin-gonic/gin"
)

type InfoHandler struct {
	appConfig config.App
}

func NewInfoHandler(a *config.App) *InfoHandler {
	return &InfoHandler{appConfig: *a}
}

func (h *InfoHandler) Show(c *gin.Context) {
	c.JSON(http.StatusOK, api.NewDataResponseWithPayload(api.NewInfoResponse(meta.Name, meta.Version, h.appConfig.TimeZone)))
}
