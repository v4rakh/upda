package server

import (
	"fmt"
	"git.myservermanager.com/varakh/upda/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

func middlewareAppName() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header(HeaderAppName, Name)
		c.Next()
	}
}

func middlewareAppVersion() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header(HeaderAppVersion, Version)
		c.Next()
	}
}

func middlewareAppContentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header(headerContentType, headerContentTypeApplicationJson)
		c.Next()
	}
}

// middlewareErrorHandler handles global error handling, does not overwrite any given status (see -1)
func middlewareErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// call next first, so this is the last in chain
		c.Next()

		if len(c.Errors) > 0 {
			// status -1 doesn't overwrite existing status code
			c.Header(headerContentType, headerContentTypeApplicationJson)
			c.JSON(-1, api.NewErrorResponseWithStatusAndMessage(errCodeToStr(c.Errors.Last()), c.Errors.Last().Error()))
			return
		}
	}
}

func middlewareAppErrorRecoveryHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, api.NewErrorResponseWithStatusAndMessage(string(General), fmt.Sprintf("%s", err)))
			}
		}()
		c.Next()
	}
}
