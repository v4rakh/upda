package server

import (
	"fmt"
	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/internal/server/config"
	"git.myservermanager.com/varakh/upda/internal/server/constant"
	"git.myservermanager.com/varakh/upda/internal/server/handler"
	"git.myservermanager.com/varakh/upda/internal/server/service_error"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func middlewareCors(c *config.Server) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     c.CorsAllowOrigins,
		AllowMethods:     c.CorsAllowMethods,
		AllowHeaders:     c.CorsAllowHeaders,
		AllowCredentials: c.CorsAllowCredentials,
		ExposeHeaders:    c.CorsExposeHeaders,
	})
}

func middlewareAppName() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header(api.HeaderAppName, constant.AppName)
		c.Next()
	}
}

func middlewareGlobalNotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusNotFound, api.NewErrorResponseWithStatusAndMessage(string(service_error.ErrCodeNotFound), "page not found"))
		return
	}
}

func middlewareGlobalMethodNotAllowed() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, api.NewErrorResponseWithStatusAndMessage(string(service_error.ErrCodeMethodNotAllowed), "method not allowed"))
		return
	}
}

func middlewareEnforceJsonContentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodOptions && !strings.HasPrefix(c.GetHeader(api.HeaderContentType), api.HeaderContentTypeApplicationJson) {
			c.AbortWithStatusJSON(http.StatusBadRequest, api.NewErrorResponseWithStatusAndMessage(string(service_error.ErrCodeIllegalArgument), "content-type must be application/json"))
			return
		}
		c.Next()
	}
}

func middlewareAppVersion() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header(api.HeaderAppVersion, constant.AppVersion)
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
			c.Header(api.HeaderContentType, api.HeaderContentTypeApplicationJson)
			c.JSON(-1, api.NewErrorResponseWithStatusAndMessage(handler.CodeToStr(c.Errors.Last()), c.Errors.Last().Error()))
			return
		}
	}
}

func middlewareAppErrorRecoveryHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, api.NewErrorResponseWithStatusAndMessage(string(service_error.ErrCodeGeneral), fmt.Sprintf("%s", err)))
			}
		}()
		c.Next()
	}
}
