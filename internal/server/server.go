package server

import (
	"context"
	"errors"
	"fmt"
	"git.myservermanager.com/varakh/upda/internal/meta"
	"git.myservermanager.com/varakh/upda/internal/server/config"
	"git.myservermanager.com/varakh/upda/internal/server/constant"
	"git.myservermanager.com/varakh/upda/internal/server/handler"
	"git.myservermanager.com/varakh/upda/internal/server/repository"
	"git.myservermanager.com/varakh/upda/internal/server/service"
	ginstatic "github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
	"go.uber.org/automaxprocs/maxprocs"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var ServeCmd = &cli.Command{
	Name:  "serve",
	Usage: "Starts the server and keeps it running",
	Action: func(ctx context.Context, _ *cli.Command) error {
		start(ctx)
		return nil
	},
}

func start(c context.Context) {
	var err error

	// configuration init
	cfg, db := config.LoadFromEnvironment(c)

	log.Info().Msgf("Starting %s %s", meta.Name, meta.Version)

	// adhere to GOMAXPROCS, but silence default output
	_, _ = maxprocs.Set(maxprocs.Logger(nil))
	log.Debug().Msgf("GOMAXPROCS '%d'", runtime.GOMAXPROCS(0))

	// set gin mode derived
	if cfg.App.Development {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	corsMiddleware := middlewareCors(cfg.Cors)
	loggingMiddleware := middlewareLogging(cfg.Logging)
	recoveryMiddleware := middlewarePanicRecoveryHandler(cfg.Logging)
	errorMiddleware := middlewareErrorTransformer()

	// routers init
	appRouter := newEngine(loggingMiddleware, recoveryMiddleware, corsMiddleware, middlewareAppName(), middlewareAppVersion(), errorMiddleware)
	promRouter := newEngine(loggingMiddleware, recoveryMiddleware, errorMiddleware)

	separatePromServer := cfg.Prometheus.Enabled && cfg.Prometheus.Port != cfg.Server.Port

	var prometheusService *service.PrometheusService
	if cfg.Prometheus.Enabled && separatePromServer {
		prometheusService = service.NewPrometheusService(promRouter, cfg.Prometheus)
		log.Info().Msg("Starting separate Prometheus server")
	} else if cfg.Prometheus.Enabled && !separatePromServer {
		prometheusService = service.NewPrometheusService(appRouter, cfg.Prometheus)
		log.Info().Msg("Starting embedded Prometheus server")
	}
	if cfg.Prometheus.Enabled {
		if err = prometheusService.Init(); err != nil {
			log.Fatal().Msgf("Prometheus service init failed: %s", err.Error())
		}
		// always instrument tracking for the meta router
		appRouter.Use(prometheusService.GetProm().Instrument())
	}

	updateRepo := repository.NewUpdateDbRepo(db)
	webhookRepo := repository.NewWebhookDbRepo(db)
	eventRepo := repository.NewEventDbRepo(db)
	secretRepo := repository.NewSecretDbRepo(db)
	constantRepo := repository.NewConstantDbRepo(db)
	actionRepo := repository.NewActionDbRepo(db)
	actionInvocationRepo := repository.NewActionInvocationDbRepo(db)
	filterPresetRepo := repository.NewFilterPresetDbRepo(db)
	commentRepo := repository.NewCommentDbRepo(db)

	var lockService service.LockService

	if cfg.Lock.RedisEnabled {
		var e error
		lockService, e = service.NewLockRedisService(cfg.Lock)

		if err != nil {
			log.Fatal().Msgf("Failed to create lock service %+v", e)
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
	commentService := service.NewCommentService(commentRepo)

	var taskService *service.TaskService

	if taskService, err = service.NewTaskService(updateService, eventService, webhookService, actionService, actionInvocationService, lockService, prometheusService, cfg.App, cfg.Task, cfg.Lock, cfg.Prometheus); err != nil {
		log.Fatal().Msgf("Task service creation failed: %v", err)
	}

	if err = taskService.Init(); err != nil {
		log.Fatal().Msgf("Task service initialization failed: %v", err)
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
	commentHandler := handler.NewCommentHandler(updateService, commentService)

	infoHandler := handler.NewInfoHandler(cfg.App)
	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler()

	// in production, the web interface is served on SERVER_BASE_PATH, build with go flag -tags prod is required
	if cfg.Webinterface.Enabled && !cfg.App.Development {
		cacheControl := middlewareCacheControl(cfg.WebinterfaceCacheControl)
		webinterfaceHandler := handler.NewWebinterfaceHandler(cfg.Webinterface)
		appRouter.GET(fmt.Sprintf("%sui/conf/runtime-config.js", cfg.Server.BasePath), cacheControl, webinterfaceHandler.GetConfig)

		targetFSPath := "web/build"
		var webinterfaceFolderFS ginstatic.ServeFileSystem
		if webinterfaceFolderFS, err = ginstatic.EmbedFolder(webinterfaceFS, targetFSPath); err != nil {
			log.Fatal().Msgf("Cannot serve webinterface folder: %s", err.Error())
		}

		appRouter.GET(cfg.Server.BasePath, middlewareRedirect("ui/"))
		if "/" != cfg.Server.BasePath {
			appRouter.GET("", middlewareRedirect(fmt.Sprintf("%sui/", cfg.Server.BasePath)))
		}
		appRouter.Use(middlewareFSRewrite(fmt.Sprintf("%sui", cfg.Server.BasePath), webinterfaceFolderFS, &cacheControl))
	}

	apiPublicGroup := appRouter.Group(fmt.Sprintf("%s/api/v1", cfg.Server.BasePath))
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
		log.Fatal().Msg("No valid auth mode found")
	}

	apiAuthGroup := appRouter.Group(fmt.Sprintf("%sapi/v1", cfg.Server.BasePath), authMethodHandler)
	apiAuthGroup.Use(middlewareSessionProvider())

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

	apiAuthGroup.GET("/comments/:updateId", commentHandler.GetAllByUpdateId)
	apiAuthGroup.POST("/comments/:updateId", middlewareEnforceJsonContentType(), commentHandler.Create)
	apiAuthGroup.PATCH("/comments/:id/content", middlewareEnforceJsonContentType(), commentHandler.UpdateContent)
	apiAuthGroup.DELETE("/comments/:id", commentHandler.Delete)

	// start servers (run in separate goroutines)
	appSrv := newServer(appRouter, fmt.Sprintf("%s:%d", cfg.Server.Listen, cfg.Server.Port))
	prometheusSrv := newServer(promRouter, fmt.Sprintf("%s:%d", cfg.Prometheus.Listen, cfg.Prometheus.Port))

	startServer(appSrv, cfg.Server)

	if separatePromServer {
		startServer(prometheusSrv, cfg.Server)
	}

	// gracefully handle shut down
	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of x seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but cannot be caught, thus no need to add
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down...")

	timeoutCtx, timeoutCancel := context.WithTimeout(c, cfg.Server.Timeout)
	defer timeoutCancel()

	shutdownDone := make(chan struct{})
	go func() {
		taskService.Stop()
		stopServer(c, appSrv)
		stopServer(c, prometheusSrv)
		close(shutdownDone)
	}()

	select {
	case <-shutdownDone:
		log.Info().Msg("Exited")
	case <-timeoutCtx.Done():
		log.Info().Msgf("Shutdown timeout of '%v' expired, exiting forcefully...", cfg.Server.Timeout)
		os.Exit(1)
	}
}

func newServer(r *gin.Engine, address string) *http.Server {
	if r == nil || address == "" {
		log.Fatal().Msg("Failed to create server, engine or address is nil")
		return nil
	}

	return &http.Server{
		Addr:    address,
		Handler: r,
	}
}

func startServer(s *http.Server, cfg *config.Server) {
	go func() {
		var e error
		log.Info().Msgf("Server listening on '%s'", s.Addr)

		if cfg.TlsEnabled {
			e = s.ListenAndServeTLS(cfg.TlsCertPath, cfg.TlsKeyPath)
		} else {
			e = s.ListenAndServe()
		}

		if e != nil && !errors.Is(e, http.ErrServerClosed) {
			log.Fatal().Msgf("Server cannot be started: %v", e)
		}
	}()
}

func stopServer(ctx context.Context, s *http.Server) {
	if s == nil {
		return
	}

	if err := s.Shutdown(ctx); err != nil {
		log.Fatal().Msgf("Shutdown failed, exited directly: %v", err)
	}

	log.Info().Msgf("Shutdown for '%s' complete", s.Addr)
}

func newEngine(middleware ...gin.HandlerFunc) *gin.Engine {
	r := gin.New()

	for _, m := range middleware {
		r.Use(m)
	}

	r.NoMethod(middlewareGlobalMethodNotAllowed())
	r.NoRoute(middlewareGlobalNotFound())

	return r
}
