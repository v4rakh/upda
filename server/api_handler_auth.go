package server

import (
	"git.myservermanager.com/varakh/upda/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

type authHandler struct {
}

func newAuthHandler() *authHandler {
	return &authHandler{}
}

func (h *authHandler) login(c *gin.Context) {
	c.Header(api.HeaderContentType, api.HeaderContentTypeApplicationJson)
	c.Status(http.StatusNoContent)
}
