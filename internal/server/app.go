package server

import (
	"context"
	"errors"
	"fmt"
	"git.myservermanager.com/varakh/upda/internal/server/config"
	"git.myservermanager.com/varakh/upda/internal/server/handler"
	"git.myservermanager.com/varakh/upda/internal/server/repository"
	"git.myservermanager.com/varakh/upda/internal/server/service"
	ginstatic "github.com/gin-contrib/static"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func Start(c context.Context) {
	// configuration init
	cfg, db := config.LoadFromEnvironment(c)

	// adhere to GOMAXPROCS, but silence default output
	_, _ = maxprocs.Set(maxprocs.Logger(nil))
	zap.L().Sugar().Debugf("GOMAXPROCS '%d'", runtime.GOMAXPROCS(0))

	// set gin mode derived
	if cfg.App.Development {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// app init (router, services, handlers)
	router := gin.New()
	router.Use(ginzap.Ginzap(zap.L(), time.RFC3339, false))
	router.Use(ginzap.RecoveryWithZap(zap.L(), true))
	router.Use(middlewareCors(cfg.Server))
	router.Use(middlewareAppName())
	router.Use(middlewareAppVersion())
	router.Use(middlewareErrorHandler())
	router.Use(middlewareAppErrorRecoveryHandler())
	router.NoRoute(middlewareGlobalNotFound())
	router.NoMethod(middlewareGlobalMethodNotAllowed())

	var err error

	prometheusService := service.NewPrometheusService(router, cfg.Prometheus, cfg.Server)

	if cfg.Prometheus.Enabled {
		if err = prometheusService.Init(); err != nil {
			zap.L().Sugar().Fatalf("Prometheus service init failed: %s", err.Error())
		}
		router.Use(prometheusService.GetProm().Instrument())
	}

	updateRepo := repository.NewUpdateDbRepo(db)
	webhookRepo := repository.NewWebhookDbRepo(db)
	eventRepo := repository.NewEventDbRepo(db)
	secretRepo := repository.NewSecretDbRepo(db)
	constantRepo := repository.NewConstantDbRepo(db)
	actionRepo := repository.NewActionDbRepo(db)
	actionInvocationRepo := repository.NewActionInvocationDbRepo(db)

	var lockService service.LockService

	if cfg.Lock.RedisEnabled {
		var e error
		lockService, e = service.NewLockRedisService(cfg.Lock)

		if err != nil {
			zap.L().Fatal("Failed to create lock service", zap.Error(e))
		}
	} else {
		lockService = service.NewLockMemService()
	}

	eventService := service.NewEventService(eventRepo)
	updateService := service.NewUpdateService(updateRepo, eventService)
	webhookService := service.NewWebhookService(webhookRepo, cfg.Webhook)
	webhookInvocationService := service.NewWebhookInvocationService(webhookService, updateService, cfg.Webhook)

	secretService := service.NewSecretService(secretRepo)
	constantService := service.NewConstantService(constantRepo)
	actionService := service.NewActionService(actionRepo, eventService)
	actionInvocationService := service.NewActionInvocationService(actionInvocationRepo, actionService, eventService, secretService, constantService)

	var taskService *service.TaskService

	if taskService, err = service.NewTaskService(updateService, eventService, webhookService, actionService, actionInvocationService, lockService, prometheusService, cfg.App, cfg.Task, cfg.Lock, cfg.Prometheus); err != nil {
		zap.L().Sugar().Fatalf("Task service creation failed: %v", err)
	}

	if err = taskService.Init(); err != nil {
		zap.L().Sugar().Fatalf("Task service initialization failed: %v", err)
	}

	taskService.Start()

	updateHandler := handler.NewUpdateHandler(updateService, cfg.App)
	webhookHandler := handler.NewWebhookHandler(webhookService)
	webhookInvocationHandler := handler.NewWebhookInvocationHandler(webhookInvocationService, webhookService)
	eventHandler := handler.NewEventHandler(eventService)
	secretHandler := handler.NewSecretHandler(secretService)
	constantHandler := handler.NewConstantHandler(constantService)
	actionHandler := handler.NewActionHandler(actionService)
	actionInvocationHandler := handler.NewActionInvocationHandler(actionService, actionInvocationService)

	infoHandler := handler.NewInfoHandler(cfg.App)
	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler()

	// embedded frontend
	// in production mode, the frontend is embedded on / during compile time utilizing -tags prod
	// if the prod tag is missing, development setup is used and a dummy frontend is shown on /
	if cfg.EmbeddedWebInterface.Enabled {
		var targetPath string
		if cfg.App.Development {
			targetPath = "web_dev"
		} else {
			targetPath = "web/build"
		}
		var embeddedFolder ginstatic.ServeFileSystem
		if embeddedFolder, err = ginstatic.EmbedFolder(embeddedFiles, targetPath); err != nil {
			zap.L().Sugar().Fatalf("Cannot serve embedded folder: %s", err.Error())
		}
		router.Use(ginstatic.Serve(fmt.Sprintf("%s", cfg.Server.BasePath), embeddedFolder))

		if !cfg.App.Development {
			embeddedFrontendGroup := router.Group(fmt.Sprintf("%s", cfg.Server.BasePath))
			embeddedFrontendGroup.GET("/conf/runtime-config.js", func(c *gin.Context) {
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
				if cfg.EmbeddedWebInterface.DarkThemeEnabled {
					webDarkThemeEnabled = 1
				}
				webEnableFooter := 0
				if cfg.EmbeddedWebInterface.FooterEnabled {
					webEnableFooter = 1
				}

				c.Data(http.StatusOK, "text/javascript; charset=utf-8", []byte(fmt.Sprintf(runtimeConfig, cfg.EmbeddedWebInterface.ApiUrl, cfg.EmbeddedWebInterface.Title, webDarkThemeEnabled, webEnableFooter)))
			})
		}
	}

	apiPublicGroup := router.Group(fmt.Sprintf("%s/api/v1", cfg.Server.BasePath))
	apiPublicGroup.GET("/health", healthHandler.Show)
	apiPublicGroup.GET("/info", infoHandler.Show)

	apiPublicGroup.POST("/webhooks/:id", middlewareEnforceJsonContentType(), webhookInvocationHandler.Execute)

	var authMethodHandler gin.HandlerFunc

	if config.AuthModeBasicSingle == cfg.Auth.AuthMethod {
		authMethodHandler = gin.BasicAuth(gin.Accounts{
			cfg.Auth.BasicAuthUser: cfg.Auth.BasicAuthPassword,
		})
	} else if config.AuthModeBasicCredentials == cfg.Auth.AuthMethod {
		authMethodHandler = gin.BasicAuth(cfg.Auth.BasicAuthCredentials)
	} else {
		zap.L().Fatal("No valid auth mode found")
	}

	apiAuthGroup := router.Group(fmt.Sprintf("%sapi/v1", cfg.Server.BasePath), authMethodHandler)

	apiAuthGroup.GET("/login", authHandler.Login)

	apiAuthGroup.GET("/updates", updateHandler.Paginate)
	apiAuthGroup.GET("/updates/:id", updateHandler.Get)
	apiAuthGroup.PATCH("/updates/:id/state", middlewareEnforceJsonContentType(), updateHandler.UpdateState)
	apiAuthGroup.DELETE("/updates/:id", updateHandler.Delete)

	apiAuthGroup.GET("/webhooks", webhookHandler.Paginate)
	apiAuthGroup.POST("/webhooks", middlewareEnforceJsonContentType(), webhookHandler.Create)
	apiAuthGroup.GET("/webhooks/:id", webhookHandler.Get)
	apiAuthGroup.PATCH("/webhooks/:id/label", middlewareEnforceJsonContentType(), webhookHandler.UpdateLabel)
	apiAuthGroup.PATCH("/webhooks/:id/ignore-host", middlewareEnforceJsonContentType(), webhookHandler.UpdateIgnoreHost)
	apiAuthGroup.DELETE("/webhooks/:id", webhookHandler.Delete)

	apiAuthGroup.GET("/events", eventHandler.Window)
	apiAuthGroup.GET("/events/:id", eventHandler.Get)
	apiAuthGroup.DELETE("/events/:id", eventHandler.Delete)

	apiAuthGroup.GET("/secrets", secretHandler.GetAll)
	apiAuthGroup.GET("/secrets/:id", secretHandler.Get)
	apiAuthGroup.POST("/secrets", middlewareEnforceJsonContentType(), secretHandler.Create)
	apiAuthGroup.PATCH("/secrets/:id/value", middlewareEnforceJsonContentType(), secretHandler.UpdateValue)
	apiAuthGroup.DELETE("/secrets/:id", secretHandler.Delete)

	apiAuthGroup.GET("/constants", constantHandler.GetAll)
	apiAuthGroup.GET("/constants/:id", constantHandler.Get)
	apiAuthGroup.POST("/constants", middlewareEnforceJsonContentType(), constantHandler.Create)
	apiAuthGroup.PATCH("/constants/:id/value", middlewareEnforceJsonContentType(), constantHandler.UpdateValue)
	apiAuthGroup.DELETE("/constants/:id", constantHandler.Delete)

	apiAuthGroup.GET("/actions", actionHandler.Paginate)
	apiAuthGroup.POST("/actions", middlewareEnforceJsonContentType(), actionHandler.Create)
	apiAuthGroup.GET("/actions/:id", actionHandler.Get)
	apiAuthGroup.PATCH("/actions/:id/label", middlewareEnforceJsonContentType(), actionHandler.UpdateLabel)
	apiAuthGroup.PATCH("/actions/:id/match-event", middlewareEnforceJsonContentType(), actionHandler.UpdateMatchEvent)
	apiAuthGroup.PATCH("/actions/:id/match-host", middlewareEnforceJsonContentType(), actionHandler.UpdateMatchHost)
	apiAuthGroup.PATCH("/actions/:id/match-application", middlewareEnforceJsonContentType(), actionHandler.UpdateMatchApplication)
	apiAuthGroup.PATCH("/actions/:id/match-provider", middlewareEnforceJsonContentType(), actionHandler.UpdateMatchProvider)
	apiAuthGroup.PATCH("/actions/:id/payload", middlewareEnforceJsonContentType(), actionHandler.UpdatePayload)
	apiAuthGroup.PATCH("/actions/:id/enabled", middlewareEnforceJsonContentType(), actionHandler.UpdateEnabled)
	apiAuthGroup.DELETE("/actions/:id", actionHandler.Delete)
	apiAuthGroup.POST("/actions/:id/test", middlewareEnforceJsonContentType(), actionInvocationHandler.Test)

	apiAuthGroup.GET("/action-invocations", actionInvocationHandler.Paginate)
	apiAuthGroup.GET("/action-invocations/:id", actionInvocationHandler.Get)
	apiAuthGroup.DELETE("/action-invocations/:id", actionInvocationHandler.Delete)

	serverAddress := fmt.Sprintf("%s:%d", cfg.Server.Listen, cfg.Server.Port)
	srv := &http.Server{
		Addr:    serverAddress,
		Handler: router,
	}

	go func() {
		var e error

		if cfg.Server.TlsEnabled {
			e = srv.ListenAndServeTLS(cfg.Server.TlsCertPath, cfg.Server.TlsKeyPath)
		} else {
			e = srv.ListenAndServe()
		}

		if e != nil && !errors.Is(e, http.ErrServerClosed) {
			zap.L().Sugar().Fatalf("Application cannot be started: %v", e)
		}
	}()

	// gracefully handle shut down
	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of x seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but cannot be caught, thus no need to add
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Shutting down...")
	taskService.Stop()

	ctx, cancel := context.WithTimeout(c, cfg.Server.Timeout)
	defer cancel()
	if err = srv.Shutdown(ctx); err != nil {
		zap.L().Sugar().Fatalf("Shutdown failed, exited directly: %v", err)
	}
	// catching ctx.Done() for configured timeout
	select {
	case <-ctx.Done():
		zap.L().Sugar().Infof("Shutdown timeout of '%v' expired, exiting...", cfg.Server.Timeout)

	}
	zap.L().Info("Exited")
}
