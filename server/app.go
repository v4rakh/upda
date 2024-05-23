package server

import (
	"context"
	"errors"
	"fmt"
	"git.myservermanager.com/varakh/upda/util"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Start() {
	// configuration init
	env := bootstrapEnvironment()

	// secure init
	util.AssertAvailablePRNG()

	// set gin mode derived
	if env.appConfig.isDevelopment {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(ginzap.Ginzap(zap.L(), time.RFC3339, false))
	router.Use(ginzap.RecoveryWithZap(zap.L(), true))

	// metrics
	prometheusService := newPrometheusService(router, env.prometheusConfig)

	if env.prometheusConfig.enabled {
		prometheusService.init()
		router.Use(prometheusService.prometheus.Instrument())
	}

	updateRepo := newUpdateDbRepo(env.db)
	webhookRepo := newWebhookDbRepo(env.db)
	eventRepo := newEventDbRepo(env.db)
	secretRepo := newSecretDbRepo(env.db)
	actionRepo := newActionDbRepo(env.db)
	actionInvocationRepo := newActionInvocationDbRepo(env.db)

	var lockService lockService

	if env.lockConfig.redisEnabled {
		var err error
		lockService, err = newLockRedisService(env.lockConfig)

		if err != nil {
			zap.L().Fatal("Failed to create lock service", zap.Error(err))
		}
	} else {
		lockService = newLockMemService()
	}

	eventService := newEventService(eventRepo)
	updateService := newUpdateService(updateRepo, eventService)
	webhookService := newWebhookService(webhookRepo, env.webhookConfig)
	webhookInvocationService := newWebhookInvocationService(webhookService, updateService, env.webhookConfig)

	secretService := newSecretService(secretRepo)
	actionService := newActionService(actionRepo, eventService)
	actionInvocationService := newActionInvocationService(actionInvocationRepo, actionService, eventService, secretService)

	taskService := newTaskService(updateService, eventService, webhookService, actionService, actionInvocationService, lockService, prometheusService, env.appConfig, env.taskConfig, env.lockConfig, env.prometheusConfig)
	taskService.init()
	taskService.start()

	updateHandler := newUpdateHandler(updateService, env.appConfig)
	webhookHandler := newWebhookHandler(webhookService)
	webhookInvocationHandler := newWebhookInvocationHandler(webhookInvocationService, webhookService)
	eventHandler := newEventHandler(eventService)
	secretHandler := newSecretHandler(secretService)
	actionHandler := newActionHandler(actionService)
	actionInvocationHandler := newActionInvocationHandler(actionService, actionInvocationService)

	infoHandler := newInfoHandler(env.appConfig)
	healthHandler := newHealthHandler()
	authHandler := newAuthHandler()

	router.Use(middlewareAppName())
	router.Use(middlewareAppVersion())
	router.Use(middlewareAppContentType())
	router.Use(middlewareErrorHandler())
	router.Use(middlewareAppErrorRecoveryHandler())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     env.serverConfig.corsAllowOrigin,
		AllowMethods:     env.serverConfig.corsAllowMethods,
		AllowHeaders:     env.serverConfig.corsAllowHeaders,
		AllowCredentials: true,
	}))

	apiPublicGroup := router.Group("/api/v1")
	apiPublicGroup.GET("/health", healthHandler.show)
	apiPublicGroup.GET("/info", infoHandler.show)

	apiPublicGroup.POST("/webhooks/:id", webhookInvocationHandler.execute)

	var authMethodHandler gin.HandlerFunc

	if authModeBasicSingle == env.authConfig.authMethod {
		authMethodHandler = gin.BasicAuth(gin.Accounts{
			env.authConfig.basicAuthUser: env.authConfig.basicAuthPassword,
		})
	} else if authModeBasicCredentials == env.authConfig.authMethod {
		authMethodHandler = gin.BasicAuth(env.authConfig.basicAuthCredentials)
	} else {
		zap.L().Fatal("No valid auth mode found")
	}

	apiAuthGroup := router.Group("/api/v1", authMethodHandler)

	apiAuthGroup.GET("/login", authHandler.login)

	apiAuthGroup.GET("/updates", updateHandler.paginate)
	apiAuthGroup.GET("/updates/:id", updateHandler.get)
	apiAuthGroup.PATCH("/updates/:id/state", updateHandler.updateState)
	apiAuthGroup.DELETE("/updates/:id", updateHandler.delete)

	apiAuthGroup.GET("/webhooks", webhookHandler.paginate)
	apiAuthGroup.POST("/webhooks", webhookHandler.create)
	apiAuthGroup.GET("/webhooks/:id", webhookHandler.get)
	apiAuthGroup.PATCH("/webhooks/:id/label", webhookHandler.updateLabel)
	apiAuthGroup.PATCH("/webhooks/:id/ignore-host", webhookHandler.updateIgnoreHost)
	apiAuthGroup.DELETE("/webhooks/:id", webhookHandler.delete)

	apiAuthGroup.GET("/events", eventHandler.window)
	apiAuthGroup.GET("/events/:id", eventHandler.get)
	apiAuthGroup.DELETE("/events/:id", eventHandler.delete)

	apiAuthGroup.GET("/secrets", secretHandler.getAll)
	apiAuthGroup.GET("/secrets/:id", secretHandler.get)
	apiAuthGroup.POST("/secrets", secretHandler.create)
	apiAuthGroup.PATCH("/secrets/:id/value", secretHandler.updateValue)
	apiAuthGroup.DELETE("/secrets/:id", secretHandler.delete)

	apiAuthGroup.GET("/actions", actionHandler.paginate)
	apiAuthGroup.POST("/actions", actionHandler.create)
	apiAuthGroup.GET("/actions/:id", actionHandler.get)
	apiAuthGroup.PATCH("/actions/:id/label", actionHandler.updateLabel)
	apiAuthGroup.PATCH("/actions/:id/match-event", actionHandler.updateMatchEvent)
	apiAuthGroup.PATCH("/actions/:id/match-host", actionHandler.updateMatchHost)
	apiAuthGroup.PATCH("/actions/:id/match-application", actionHandler.updateMatchApplication)
	apiAuthGroup.PATCH("/actions/:id/match-provider", actionHandler.updateMatchProvider)
	apiAuthGroup.PATCH("/actions/:id/payload", actionHandler.updatePayload)
	apiAuthGroup.PATCH("/actions/:id/enabled", actionHandler.updateEnabled)
	apiAuthGroup.DELETE("/actions/:id", actionHandler.delete)
	apiAuthGroup.POST("/actions/:id/test", actionInvocationHandler.test)

	apiAuthGroup.GET("/action-invocations", actionInvocationHandler.paginate)
	apiAuthGroup.GET("/action-invocations/:id", actionInvocationHandler.get)
	apiAuthGroup.DELETE("/action-invocations/:id", actionInvocationHandler.delete)

	// start server
	serverAddress := fmt.Sprintf("%s:%d", env.serverConfig.listen, env.serverConfig.port)
	srv := &http.Server{
		Addr:    serverAddress,
		Handler: router,
	}

	go func() {
		var err error

		if env.serverConfig.tlsEnabled {
			err = srv.ListenAndServeTLS(env.serverConfig.tlsCertPath, env.serverConfig.tlsKeyPath)
		} else {
			err = srv.ListenAndServe()
		}

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.L().Sugar().Fatalf("Application cannot be started: %v", err)
		}
	}()

	// gracefully handle shut down
	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of x seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but cannot be caught, thus no need to add
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Shutting down...")
	taskService.stop()

	ctx, cancel := context.WithTimeout(context.Background(), env.serverConfig.timeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Sugar().Fatalf("Shutdown failed, exited directly: %v", err)
	}
	// catching ctx.Done() for configured timeout
	select {
	case <-ctx.Done():
		zap.L().Sugar().Infof("Shutdown timeout of '%v' expired, exiting...", env.serverConfig.timeout)

	}
	zap.L().Info("Exited")
}
