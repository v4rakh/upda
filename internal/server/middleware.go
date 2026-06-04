package server

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"runtime/debug"
	"strings"

	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/internal/meta"
	"git.myservermanager.com/varakh/upda/internal/server/config"
	"git.myservermanager.com/varakh/upda/internal/server/handler"
	"git.myservermanager.com/varakh/upda/internal/server/service_error"
	"github.com/gin-contrib/cors"
	ginstatic "github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	cachecontrol "go.eigsys.de/gin-cachecontrol/v2"
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

// middlewareLogging logs access
func middlewareLogging(lc *config.Logging) gin.HandlerFunc {
	var err error
	var logLevel zerolog.Level
	if logLevel, err = zerolog.ParseLevel(lc.LevelRequests); err != nil {
		logLevel = zerolog.Disabled
	}
	return func(c *gin.Context) {
		c.Next()
		log.WithLevel(logLevel).Msgf("Handled request %s %s: %d", c.Request.Method, c.Request.URL.Path, c.Writer.Status())
	}
}

// middlewarePanicRecoveryHandler recovers meta from panics, logs them and returns proper response
// logs the error and stack trace using zerolog.Logger, and returns a 500 response.
func middlewarePanicRecoveryHandler(lc *config.Logging) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Error().Str(lc.EncodingStacktraceKey, string(debug.Stack())).Msgf("panic recovered: %v", err)
				c.Header(api.HeaderContentType, api.HeaderContentTypeApplicationJson)
				c.AbortWithStatusJSON(http.StatusInternalServerError, api.NewErrorResponseWithStatusAndMessage(string(service_error.ErrCodeGeneral), fmt.Sprintf("%s", err)))
			}
		}()

		c.Next()
	}
}

// middlewareErrorTransformer transforms errors into proper responses (does not overwrite any given status)
func middlewareErrorTransformer() gin.HandlerFunc {
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
		c.Header(api.HeaderAppName, meta.Name)
		c.Next()
	}
}

// middlewareAppVersion adds custom HTTP header to each request
func middlewareAppVersion() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header(api.HeaderAppVersion, meta.Version)
		c.Next()
	}
}

// middlewareGlobalNotFound adds a global not found in the same style as the API responds
func middlewareGlobalNotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusNotFound, api.NewErrorResponseWithStatusAndMessage(string(service_error.ErrCodeNotFound), "page not found"))
	}
}

// middlewareGlobalMethodNotAllowed adds a global method not allowed in the same style as the API responds
func middlewareGlobalMethodNotAllowed() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, api.NewErrorResponseWithStatusAndMessage(string(service_error.ErrCodeMethodNotAllowed), "method not allowed"))
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

// middlewareFSRewrite rewrites a path to a static file system allowing to provide additional gin handlers to be applied on success
func middlewareFSRewrite(basePath string, fs ginstatic.ServeFileSystem, handlerFuncs ...*gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := c.Request.URL.Path
		if strings.HasPrefix(p, basePath) {
			relPath := strings.TrimPrefix(p, basePath)
			if relPath == "" || relPath == "/" {
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

// middlewareSecurityHeaders sets Content-Security-Policy and Strict-Transport-Security
// headers for the web interface when the respective options are enabled.
func middlewareSecurityHeaders(c *config.WebinterfaceSecurityHeaders) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.Writer.Header()
		if c.CspEnabled {
			header.Set(api.HeaderContentSecurityPolicy, c.CspValue)
		}
		if c.HstsEnabled {
			hstsValue := fmt.Sprintf("max-age=%d", int(c.HstsMaxAge.Seconds()))
			if c.HstsIncludeSubDomains {
				hstsValue += "; includeSubDomains"
			}
			if c.HstsPreload {
				hstsValue += "; preload"
			}
			header.Set(api.HeaderStrictTransportSecurity, hstsValue)
		}
	}
}

// middlewareRedirect redirects when not /
func middlewareRedirect(targetPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := c.Request.URL.Path
		var t string
		if p == "" || p == "/" {
			t = targetPath
		} else {
			t = fmt.Sprintf("%s/%s", p, targetPath)
		}
		c.Redirect(http.StatusMovedPermanently, t)
		c.Abort()
	}
}
