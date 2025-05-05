package handler

import (
	"fmt"
	"git.myservermanager.com/varakh/upda/internal/server/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

type EmbeddedWebInterfaceHandler struct {
	config *config.EmbeddedWebInterface
}

func NewEmbeddedWebInterfaceHandler(c *config.EmbeddedWebInterface) *EmbeddedWebInterfaceHandler {
	return &EmbeddedWebInterfaceHandler{config: c}
}

func (h *EmbeddedWebInterfaceHandler) GetConfig(c *gin.Context) {
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
	webDarkThemeEnabled := 0
	if h.config.DarkThemeEnabled {
		webDarkThemeEnabled = 1
	}
	webEnableFooter := 0
	if h.config.FooterEnabled {
		webEnableFooter = 1
	}

	c.Data(http.StatusOK, "text/javascript; charset=utf-8", []byte(fmt.Sprintf(runtimeConfig, h.config.ApiUrl, h.config.Title, webDarkThemeEnabled, webEnableFooter)))
}
