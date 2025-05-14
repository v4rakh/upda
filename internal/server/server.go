package server

import (
	"context"
	"errors"
	"fmt"
	"git.myservermanager.com/varakh/upda/internal/server/config"
	"git.myservermanager.com/varakh/upda/internal/server/constant"
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
	router.Use(ginzap.Ginzap(zap.L(), "", cfg.Logging.UTC))
	router.Use(ginzap.RecoveryWithZap(zap.L(), true))
	router.Use(middlewareCors(cfg.Cors))
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
	filterPresetRepo := repository.NewFilterPresetDbRepo(db)

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
	filterPresetService := service.NewFilterPresetService(filterPresetRepo)

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
	filterPresetHandler := handler.NewFilterPresetHandler(filterPresetService)

	infoHandler := handler.NewInfoHandler(cfg.App)
	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler()

	// in production, the web interface is served on SERVER_BASE_PATH, build with go flag -tags prod is required
	if cfg.Webinterface.Enabled && !cfg.App.Development {
		cacheControl := middlewareCacheControl(cfg.WebinterfaceCacheControl)
		webinterfaceHandler := handler.NewWebinterfaceHandler(cfg.Webinterface)
		router.GET(fmt.Sprintf("%sui/conf/runtime-config.js", cfg.Server.BasePath), cacheControl, webinterfaceHandler.GetConfig)

		targetFSPath := "web/build"
		var webinterfaceFolderFS ginstatic.ServeFileSystem
		if webinterfaceFolderFS, err = ginstatic.EmbedFolder(webinterfaceFS, targetFSPath); err != nil {
			zap.L().Sugar().Fatalf("Cannot serve webinterface folder: %s", err.Error())
		}

		router.GET(cfg.Server.BasePath, middlewareRedirect("ui/"))
		if "/" != cfg.Server.BasePath {
			router.GET("", middlewareRedirect(fmt.Sprintf("%sui/", cfg.Server.BasePath)))
		}
		router.Use(middlewareFSRewrite(fmt.Sprintf("%sui", cfg.Server.BasePath), webinterfaceFolderFS, &cacheControl))
	}

	apiPublicGroup := router.Group(fmt.Sprintf("%s/api/v1", cfg.Server.BasePath))
	apiPublicGroup.GET("/health", healthHandler.Show)
	apiPublicGroup.GET("/info", infoHandler.Show)

	apiPublicGroup.POST("/webhooks/:id", middlewareEnforceJsonContentType(), webhookInvocationHandler.Execute)

	var authMethodHandler gin.HandlerFunc

	if constant.ConfigAuthModeBasicSingle == cfg.Auth.AuthMethod {
		authMethodHandler = gin.BasicAuth(gin.Accounts{
			cfg.Auth.BasicAuthUser: cfg.Auth.BasicAuthPassword,
		})
	} else if constant.ConfigAuthModeBasicCredentials == cfg.Auth.AuthMethod {
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

	apiAuthGroup.GET("/filter-presets/:type", filterPresetHandler.GetByType)
	apiAuthGroup.POST("/filter-presets", middlewareEnforceJsonContentType(), filterPresetHandler.Create)
	apiAuthGroup.DELETE("/filter-presets/:id", filterPresetHandler.Delete)

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
	select {
	case <-ctx.Done():
		zap.L().Sugar().Infof("Shutdown timeout of '%v' expired, exiting...", cfg.Server.Timeout)

	}
	zap.L().Info("Exited")
}
