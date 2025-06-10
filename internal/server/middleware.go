package server

import (
	"fmt"
	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/internal/app"
	"git.myservermanager.com/varakh/upda/internal/server/config"
	"git.myservermanager.com/varakh/upda/internal/server/constant"
	"git.myservermanager.com/varakh/upda/internal/server/dto"
	"git.myservermanager.com/varakh/upda/internal/server/handler"
	"git.myservermanager.com/varakh/upda/internal/server/service_error"
	"github.com/gin-contrib/cors"
	ginstatic "github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"go.eigsys.de/gin-cachecontrol/v2"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

// middlewareCors applies CORS configuration
func middlewareCors(c *config.Cors) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     c.AllowOrigins,
		AllowMethods:     c.AllowMethods,
		AllowHeaders:     c.AllowHeaders,
		AllowCredentials: c.AllowCredentials,
		ExposeHeaders:    c.ExposeHeaders,
	})
}

// middlewareCacheControl applies cache control settings
func middlewareCacheControl(c *config.WebinterfaceCacheControl) gin.HandlerFunc {
	if !c.Enabled {
		return cachecontrol.New(cachecontrol.NoCachePreset)
	}

	return cachecontrol.New(
		cachecontrol.Config{
			MustRevalidate:       c.MustRevalidate,
			NoCache:              c.NoCache,
			NoStore:              c.NoStore,
			NoTransform:          c.NoTransform,
			Public:               c.Public,
			Private:              c.Private,
			ProxyRevalidate:      c.ProxyRevalidate,
			MaxAge:               c.MaxAge,
			SMaxAge:              c.SMaxAge,
			Immutable:            c.Immutable,
			StaleWhileRevalidate: c.StaleWhileRevalidate,
			StaleIfError:         c.StaleIfError,
		},
	)
}

// middlewareAppName adds custom HTTP header to each request
func middlewareAppName() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header(api.HeaderAppName, app.Name)
		c.Next()
	}
}

// middlewareAppVersion adds custom HTTP header to each request
func middlewareAppVersion() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header(api.HeaderAppVersion, app.Version)
		c.Next()
	}
}

// middlewareGlobalNotFound adds a global not found in the same style as the API responds
func middlewareGlobalNotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusNotFound, api.NewErrorResponseWithStatusAndMessage(string(service_error.ErrCodeNotFound), "page not found"))
		return
	}
}

// middlewareGlobalMethodNotAllowed adds a global method not allowed in the same style as the API responds
func middlewareGlobalMethodNotAllowed() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, api.NewErrorResponseWithStatusAndMessage(string(service_error.ErrCodeMethodNotAllowed), "method not allowed"))
		return
	}
}

// middlewareEnforceJsonContentType enforces JSON content type
func middlewareEnforceJsonContentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodOptions && !strings.HasPrefix(c.GetHeader(api.HeaderContentType), api.HeaderContentTypeApplicationJson) {
			c.AbortWithStatusJSON(http.StatusBadRequest, api.NewErrorResponseWithStatusAndMessage(string(service_error.ErrCodeIllegalArgument), "content-type must be application/json"))
			return
		}
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

// middlewareErrorRecoveryHandler recovers from panics, returning a 500 error
func middlewareErrorRecoveryHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, api.NewErrorResponseWithStatusAndMessage(string(service_error.ErrCodeGeneral), fmt.Sprintf("%s", err)))
			}
		}()
		c.Next()
	}
}

// middlewareFSRewrite rewrites a path to a static file system allowing to provide additional gin handlers to be applied on success
func middlewareFSRewrite(basePath string, fs ginstatic.ServeFileSystem, handlerFuncs ...*gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := c.Request.URL.Path
		if strings.HasPrefix(p, basePath) {
			relPath := strings.TrimPrefix(p, basePath)
			if "" == relPath || "/" == relPath {
				relPath = "/index.html"
			}
			relPath = strings.TrimPrefix(relPath, "/")

			var err error
			var f http.File
			if f, err = fs.Open(relPath); err != nil {
				c.Status(http.StatusNotFound)
				c.Abort()
				return
			}
			defer func(f http.File) {
				_ = f.Close()
			}(f)

			c.Status(http.StatusOK)
			c.Header(api.HeaderContentType, mime.TypeByExtension(filepath.Ext(relPath)))

			for _, h := range handlerFuncs {
				if h != nil {
					i := *h
					i(c)
				}
			}

			if _, err = io.Copy(c.Writer, f); err != nil {
				c.Status(http.StatusInternalServerError)
				c.Abort()
			}

			c.Abort()
			return
		}
		c.Next()
	}
}

// middlewareRedirect redirects when not /
func middlewareRedirect(targetPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := c.Request.URL.Path
		var t string
		if "" == p || "/" == p {
			t = targetPath
		} else {
			t = fmt.Sprintf("%s/%s", p, targetPath)
		}
		c.Redirect(http.StatusMovedPermanently, t)
		c.Abort()
		return
	}
}

// middlewareSessionProvider apply middleware to necessary routes only to provide session information.
func middlewareSessionProvider() gin.HandlerFunc {
	return func(c *gin.Context) {
		var session *dto.ContextSession

		user, exists := c.Get(gin.AuthUserKey)
		if !exists || user == nil {
			e := service_error.ErrUnauthorized
			_ = c.AbortWithError(handler.ToHttpStatus(e), e)
			return
		}

		session = &dto.ContextSession{
			User: user.(string),
		}

		c.Set(constant.GinContextSession, *session)
		c.Next()
	}
}
