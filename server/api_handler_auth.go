package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type authHandler struct {
}

func newAuthHandler() *authHandler {
	return &authHandler{}
}

func (h *authHandler) login(c *gin.Context) {
	c.Header(headerContentType, headerContentTypeApplicationJson)
	c.Status(http.StatusNoContent)
}
