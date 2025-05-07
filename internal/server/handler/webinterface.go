package handler

import (
	"fmt"
	"git.myservermanager.com/varakh/upda/internal/server/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

type WebinterfaceHandler struct {
	config *config.Webinterface
}

func NewWebinterfaceHandler(c *config.Webinterface) *WebinterfaceHandler {
	return &WebinterfaceHandler{config: c}
}

func (h *WebinterfaceHandler) GetConfig(c *gin.Context) {
	runtimeConfig := `
const runtime_config = Object.freeze({
  VITE_BASE_PATH: '/',
  VITE_API_URL: '%s',
  VITE_TITLE: '%s',
  VITE_ENABLE_DARK_THEME: %d,
  VITE_ENABLE_FOOTER: %d
});

Object.defineProperty(window, 'runtime_config', {
    value: runtime_config,
    writable: false
});
	`
	darkThemeEnabled := 0
	if h.config.DarkThemeEnabled {
		darkThemeEnabled = 1
	}
	enableFooter := 0
	if h.config.FooterEnabled {
		enableFooter = 1
	}

	c.Data(http.StatusOK, "text/javascript; charset=utf-8", []byte(fmt.Sprintf(runtimeConfig, h.config.ApiUrl, h.config.Title, darkThemeEnabled, enableFooter)))
}
