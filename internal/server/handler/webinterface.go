package handler

import (
	"fmt"
	"git.myservermanager.com/varakh/upda/internal/server/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

type WebinterfaceHandler struct {
	webInterfaceConfig *config.Webinterface
	authConfig         *config.Auth
}

func NewWebinterfaceHandler(wc *config.Webinterface, ac *config.Auth) *WebinterfaceHandler {
	return &WebinterfaceHandler{webInterfaceConfig: wc, authConfig: ac}
}

func (h *WebinterfaceHandler) GetConfig(c *gin.Context) {
	runtimeConfig := `
const runtime_config = Object.freeze({
  VITE_BASE_PATH: '/',
  VITE_API_URL: '%s',
  VITE_TITLE: '%s',
  VITE_AUTH_TYPE: '%s',
  VITE_ENABLE_FOOTER: %d
});

Object.defineProperty(window, 'runtime_config', {
    value: runtime_config,
    writable: false
});
	`
	enableFooter := 0
	if h.webInterfaceConfig.FooterEnabled {
		enableFooter = 1
	}

	webinterfaceConfig := fmt.Sprintf(runtimeConfig, h.webInterfaceConfig.ApiUrl, h.webInterfaceConfig.Title, h.authConfig.Type, enableFooter)
	c.Data(http.StatusOK, "text/javascript; charset=utf-8", []byte(webinterfaceConfig))
}
