package server

import (
	"context"
	"errors"
	"fmt"
	"git.myservermanager.com/varakh/upda/util"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	_ "go.uber.org/automaxprocs"
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

	// app init (router, services, handlers)
	router := gin.New()
	router.Use(ginzap.Ginzap(zap.L(), time.RFC3339, false))
	router.Use(ginzap.RecoveryWithZap(zap.L(), true))

	var err error

	ps := newPrometheusService(router, env.prometheusConfig)

	if env.prometheusConfig.enabled {
		if err = ps.init(); err != nil {
			zap.L().Sugar().Fatalf("Prometheus service init failed: %s", err.Error())
		}
		router.Use(ps.prometheus.Instrument())
	}

	updateRepo := newUpdateDbRepo(env.db)
	webhookRepo := newWebhookDbRepo(env.db)
	eventRepo := newEventDbRepo(env.db)
	secretRepo := newSecretDbRepo(env.db)
	actionRepo := newActionDbRepo(env.db)
	actionInvocationRepo := newActionInvocationDbRepo(env.db)

	var ls lockService

	if env.lockConfig.redisEnabled {
		var e error
		ls, e = newLockRedisService(env.lockConfig)

		if err != nil {
			zap.L().Fatal("Failed to create lock service", zap.Error(e))
		}
	} else {
		ls = newLockMemService()
	}

	es := newEventService(eventRepo)
	us := newUpdateService(updateRepo, es)
	ws := newWebhookService(webhookRepo, env.webhookConfig)
	wis := newWebhookInvocationService(ws, us, env.webhookConfig)

	ss := newSecretService(secretRepo)
	as := newActionService(actionRepo, es)
	ais := newActionInvocationService(actionInvocationRepo, as, es, ss)

	var ts *taskService

	if ts, err = newTaskService(us, es, ws, as, ais, ls, ps, env.appConfig, env.taskConfig, env.lockConfig, env.prometheusConfig); err != nil {
		zap.L().Sugar().Fatalf("Task service creation failed: %v", err)
	}

	if err = ts.init(); err != nil {
		zap.L().Sugar().Fatalf("Task service initialization failed: %v", err)
	}

	ts.start()

	uh := newUpdateHandler(us, env.appConfig)
	wh := newWebhookHandler(ws)
	wih := newWebhookInvocationHandler(wis, ws)
	eh := newEventHandler(es)
	sh := newSecretHandler(ss)
	ah := newActionHandler(as)
	aih := newActionInvocationHandler(as, ais)

	ih := newInfoHandler(env.appConfig)
	hh := newHealthHandler()
	authH := newAuthHandler()

	router.Use(middlewareEnforceJsonContentType())
	router.Use(middlewareAppName())
	router.Use(middlewareAppVersion())
	router.Use(middlewareAppContentType())
	router.Use(middlewareErrorHandler())
	router.Use(middlewareAppErrorRecoveryHandler())
	router.NoRoute(middlewareGlobalNotFound())
	router.NoMethod(middlewareGlobalMethodNotAllowed())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     env.serverConfig.corsAllowOrigins,
		AllowMethods:     env.serverConfig.corsAllowMethods,
		AllowHeaders:     env.serverConfig.corsAllowHeaders,
		AllowCredentials: env.serverConfig.corsAllowCredentials,
		ExposeHeaders:    env.serverConfig.corsExposeHeaders,
	}))

	apiPublicGroup := router.Group("/api/v1")
	apiPublicGroup.GET("/health", hh.show)
	apiPublicGroup.GET("/info", ih.show)

	apiPublicGroup.POST("/webhooks/:id", wih.execute)

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

	apiAuthGroup.GET("/login", authH.login)

	apiAuthGroup.GET("/updates", uh.paginate)
	apiAuthGroup.GET("/updates/:id", uh.get)
	apiAuthGroup.PATCH("/updates/:id/state", uh.updateState)
	apiAuthGroup.DELETE("/updates/:id", uh.delete)

	apiAuthGroup.GET("/webhooks", wh.paginate)
	apiAuthGroup.POST("/webhooks", wh.create)
	apiAuthGroup.GET("/webhooks/:id", wh.get)
	apiAuthGroup.PATCH("/webhooks/:id/label", wh.updateLabel)
	apiAuthGroup.PATCH("/webhooks/:id/ignore-host", wh.updateIgnoreHost)
	apiAuthGroup.DELETE("/webhooks/:id", wh.delete)

	apiAuthGroup.GET("/events", eh.window)
	apiAuthGroup.GET("/events/:id", eh.get)
	apiAuthGroup.DELETE("/events/:id", eh.delete)

	apiAuthGroup.GET("/secrets", sh.getAll)
	apiAuthGroup.GET("/secrets/:id", sh.get)
	apiAuthGroup.POST("/secrets", sh.create)
	apiAuthGroup.PATCH("/secrets/:id/value", sh.updateValue)
	apiAuthGroup.DELETE("/secrets/:id", sh.delete)

	apiAuthGroup.GET("/actions", ah.paginate)
	apiAuthGroup.POST("/actions", ah.create)
	apiAuthGroup.GET("/actions/:id", ah.get)
	apiAuthGroup.PATCH("/actions/:id/label", ah.updateLabel)
	apiAuthGroup.PATCH("/actions/:id/match-event", ah.updateMatchEvent)
	apiAuthGroup.PATCH("/actions/:id/match-host", ah.updateMatchHost)
	apiAuthGroup.PATCH("/actions/:id/match-application", ah.updateMatchApplication)
	apiAuthGroup.PATCH("/actions/:id/match-provider", ah.updateMatchProvider)
	apiAuthGroup.PATCH("/actions/:id/payload", ah.updatePayload)
	apiAuthGroup.PATCH("/actions/:id/enabled", ah.updateEnabled)
	apiAuthGroup.DELETE("/actions/:id", ah.delete)
	apiAuthGroup.POST("/actions/:id/test", aih.test)

	apiAuthGroup.GET("/action-invocations", aih.paginate)
	apiAuthGroup.GET("/action-invocations/:id", aih.get)
	apiAuthGroup.DELETE("/action-invocations/:id", aih.delete)

	serverAddress := fmt.Sprintf("%s:%d", env.serverConfig.listen, env.serverConfig.port)
	srv := &http.Server{
		Addr:    serverAddress,
		Handler: router,
	}

	go func() {
		var e error

		if env.serverConfig.tlsEnabled {
			e = srv.ListenAndServeTLS(env.serverConfig.tlsCertPath, env.serverConfig.tlsKeyPath)
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
	ts.stop()

	ctx, cancel := context.WithTimeout(context.Background(), env.serverConfig.timeout)
	defer cancel()
	if err = srv.Shutdown(ctx); err != nil {
		zap.L().Sugar().Fatalf("Shutdown failed, exited directly: %v", err)
	}
	// catching ctx.Done() for configured timeout
	select {
	case <-ctx.Done():
		zap.L().Sugar().Infof("Shutdown timeout of '%v' expired, exiting...", env.serverConfig.timeout)

	}
	zap.L().Info("Exited")
}
