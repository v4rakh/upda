package server

import (
	"context"
	"errors"
	"fmt"
	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/internal/meta"
	"git.myservermanager.com/varakh/upda/internal/server/auth"
	"git.myservermanager.com/varakh/upda/internal/server/auth/sessiongormstore"
	"git.myservermanager.com/varakh/upda/internal/server/config"
	"git.myservermanager.com/varakh/upda/internal/server/constant"
	"git.myservermanager.com/varakh/upda/internal/server/handler"
	"git.myservermanager.com/varakh/upda/internal/server/repository"
	"git.myservermanager.com/varakh/upda/internal/server/service"
	"git.myservermanager.com/varakh/upda/internal/server/service_error"
	"github.com/gin-contrib/sessions"
	ginstatic "github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/wader/gormstore/v2"
	"go.uber.org/automaxprocs/maxprocs"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"slices"
	"syscall"
)

type Server struct {
	ctx context.Context
}

func New(ctx *context.Context) *Server {
	s := &Server{}
	if ctx == nil {
		s.ctx = context.Background()
	} else {
		s.ctx = *ctx
	}

	return s
}

func (s *Server) Start() {
	var err error

	// configuration init
	cfg, db := config.LoadFromEnvironment(s.ctx)

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
	appRouter := s.newEngine(loggingMiddleware, recoveryMiddleware, corsMiddleware, middlewareAppName(), middlewareAppVersion(), errorMiddleware)
	promRouter := s.newEngine(loggingMiddleware, recoveryMiddleware, errorMiddleware)

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
		// always instrument tracking for the meta router
		appRouter.Use(prometheusService.GetProm().Instrument())
	}

	// repositories init
	updateRepo := repository.NewUpdateDbRepo(db)
	webhookRepo := repository.NewWebhookDbRepo(db)
	eventRepo := repository.NewEventDbRepo(db)
	secretRepo := repository.NewSecretDbRepo(db)
	constantRepo := repository.NewConstantDbRepo(db)
	actionRepo := repository.NewActionDbRepo(db)
	actionInvocationRepo := repository.NewActionInvocationDbRepo(db)
	filterPresetRepo := repository.NewFilterPresetDbRepo(db)
	commentRepo := repository.NewCommentDbRepo(db)

	// services init
	lockService := service.NewLockMemService()
	if cfg.Lock.RedisEnabled {
		var e error
		if lockService, e = service.NewLockRedisService(cfg.Lock); e != nil {
			log.Fatal().Msgf("Failed to create lock service: %+v", e)
		}
	}

	var taskService *service.TaskService
	if taskService, err = service.NewTaskService(lockService, cfg.App, cfg.Lock); err != nil {
		log.Fatal().Msgf("Task service creation failed: %v", err)
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

	// tasks init
	updatesCleanTask := service.NewUpdatesCleanTask(updateService, taskService, cfg.Task)
	if err = updatesCleanTask.Init(); err != nil {
		log.Fatal().Msgf("Task updates clean initialization failed: %v", err)
	}
	eventsCleanTask := service.NewEventsCleanTask(eventService, taskService, cfg.Task)
	if err = eventsCleanTask.Init(); err != nil {
		log.Fatal().Msgf("Task events clean initialization failed: %v", err)
	}
	actionsCleanTask := service.NewActionsCleanTask(actionInvocationService, taskService, cfg.Task)
	if err = actionsCleanTask.Init(); err != nil {
		log.Fatal().Msgf("Task actions clean initialization failed: %v", err)
	}
	actionsEnqueueTask := service.NewActionsEnqueueTask(actionInvocationService, taskService, cfg.Task)
	if err = actionsEnqueueTask.Init(); err != nil {
		log.Fatal().Msgf("Task actions enqueue initialization failed: %v", err)
	}
	actionsInvokeTask := service.NewActionsInvokeTask(actionInvocationService, taskService, cfg.Task)
	if err = actionsInvokeTask.Init(); err != nil {
		log.Fatal().Msgf("Task actions invoke initialization failed: %v", err)
	}
	prometheusTask := service.NewPrometheusTask(updateService, eventService, webhookService, actionService, prometheusService, taskService, cfg.Prometheus)
	if err = prometheusTask.Init(); err != nil {
		log.Fatal().Msgf("Task prometheus task initialization failed: %v", err)
	}

	taskService.Start()

	// handlers init
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

	// in production, the web interface is served on SERVER_BASE_PATH, build with go flag -tags prod is required
	if cfg.Webinterface.Enabled && !cfg.App.Development {
		cacheControl := middlewareCacheControl(cfg.WebinterfaceCacheControl)
		webinterfaceHandler := handler.NewWebinterfaceHandler(cfg.Webinterface, cfg.Auth)
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

	apiProtectedGroup := appRouter.Group(fmt.Sprintf("%sapi/v1", cfg.Server.BasePath))

	// authentication provider
	if !slices.Contains(constant.ConfigAuthTypeValues(), cfg.Auth.Type) {
		log.Fatal().Msg("No valid auth type found")
	}

	var authProvider auth.Provider

	if constant.ConfigAuthTypeSession == cfg.Auth.Type {
		var validCredentials map[string]string
		if constant.ConfigAuthSessionProviderSingle == cfg.Auth.SessionProvider {
			validCredentials = map[string]string{cfg.Auth.SessionUser: cfg.Auth.SessionPassword}
		} else if constant.ConfigAuthSessionProviderCredentials == cfg.Auth.SessionProvider {
			validCredentials = cfg.Auth.SessionCredentials
		}

		authSessionCleanTask := service.NewAuthSessionCleanTask(cfg.Auth, taskService)

		sessionOptions := sessions.Options{
			MaxAge:   int(cfg.Auth.SessionCookieMaxAge.Seconds()),
			SameSite: cfg.Auth.SessionCookieSameSiteMode,
			Path:     cfg.Auth.SessionCookiePath,
			HttpOnly: cfg.Auth.SessionCookieHttpOnly,
			Secure:   cfg.Auth.SessionCookieSecure,
		}
		if cfg.Auth.SessionCookieDomain != nil {
			sessionOptions.Domain = *cfg.Auth.SessionCookieDomain
		}

		gormStoreOptions := gormstore.Options{
			TableName:       "sessions",
			SkipCreateTable: true,
		}

		sessionStore := sessiongormstore.New(db,
			sessiongormstore.WithKeyPairs([]byte(cfg.Auth.SessionSecret)),
			sessiongormstore.WithSessionOpts(sessionOptions),
			sessiongormstore.WithGormStoreOpts(gormStoreOptions),
			sessiongormstore.WithCleanupFunc(authSessionCleanTask.GetCleanupFn()),
		)

		authProvider = auth.NewSessionProvider(cfg.Auth.SessionCookieName, cfg.Auth.SessionCookiePath, auth.NewStaticUserPasswordValidator(validCredentials), sessionStore)
	}

	s.applyAuthProvider(authProvider, apiPublicGroup, apiProtectedGroup)

	apiProtectedGroup.GET("/updates", updateHandler.Paginate)
	apiProtectedGroup.GET("/updates/:id", updateHandler.Get)
	apiProtectedGroup.PATCH("/updates/:id/state", middlewareEnforceJsonContentType(), updateHandler.UpdateState)
	apiProtectedGroup.DELETE("/updates/:id", updateHandler.Delete)

	apiProtectedGroup.GET("/webhooks", webhookHandler.Paginate)
	apiProtectedGroup.POST("/webhooks", middlewareEnforceJsonContentType(), webhookHandler.Create)
	apiProtectedGroup.GET("/webhooks/:id", webhookHandler.Get)
	apiProtectedGroup.PATCH("/webhooks/:id/label", middlewareEnforceJsonContentType(), webhookHandler.UpdateLabel)
	apiProtectedGroup.PATCH("/webhooks/:id/ignore-host", middlewareEnforceJsonContentType(), webhookHandler.UpdateIgnoreHost)
	apiProtectedGroup.PATCH("/webhooks/:id/ignore-host-replacement", middlewareEnforceJsonContentType(), webhookHandler.UpdateIgnoreHostReplacement)
	apiProtectedGroup.DELETE("/webhooks/:id", webhookHandler.Delete)

	apiProtectedGroup.GET("/events", eventHandler.Window)
	apiProtectedGroup.GET("/events/:id", eventHandler.Get)
	apiProtectedGroup.DELETE("/events/:id", eventHandler.Delete)

	apiProtectedGroup.GET("/secrets", secretHandler.GetAll)
	apiProtectedGroup.GET("/secrets/:id", secretHandler.Get)
	apiProtectedGroup.POST("/secrets", middlewareEnforceJsonContentType(), secretHandler.Create)
	apiProtectedGroup.PATCH("/secrets/:id/value", middlewareEnforceJsonContentType(), secretHandler.UpdateValue)
	apiProtectedGroup.DELETE("/secrets/:id", secretHandler.Delete)

	apiProtectedGroup.GET("/constants", constantHandler.GetAll)
	apiProtectedGroup.GET("/constants/:id", constantHandler.Get)
	apiProtectedGroup.POST("/constants", middlewareEnforceJsonContentType(), constantHandler.Create)
	apiProtectedGroup.PATCH("/constants/:id/value", middlewareEnforceJsonContentType(), constantHandler.UpdateValue)
	apiProtectedGroup.DELETE("/constants/:id", constantHandler.Delete)

	apiProtectedGroup.GET("/actions", actionHandler.Paginate)
	apiProtectedGroup.POST("/actions", middlewareEnforceJsonContentType(), actionHandler.Create)
	apiProtectedGroup.GET("/actions/:id", actionHandler.Get)
	apiProtectedGroup.PATCH("/actions/:id/label", middlewareEnforceJsonContentType(), actionHandler.UpdateLabel)
	apiProtectedGroup.PATCH("/actions/:id/match-event", middlewareEnforceJsonContentType(), actionHandler.UpdateMatchEvent)
	apiProtectedGroup.PATCH("/actions/:id/match-host", middlewareEnforceJsonContentType(), actionHandler.UpdateMatchHost)
	apiProtectedGroup.PATCH("/actions/:id/match-application", middlewareEnforceJsonContentType(), actionHandler.UpdateMatchApplication)
	apiProtectedGroup.PATCH("/actions/:id/match-provider", middlewareEnforceJsonContentType(), actionHandler.UpdateMatchProvider)
	apiProtectedGroup.PATCH("/actions/:id/payload", middlewareEnforceJsonContentType(), actionHandler.UpdatePayload)
	apiProtectedGroup.PATCH("/actions/:id/enabled", middlewareEnforceJsonContentType(), actionHandler.UpdateEnabled)
	apiProtectedGroup.DELETE("/actions/:id", actionHandler.Delete)
	apiProtectedGroup.POST("/actions/:id/test", middlewareEnforceJsonContentType(), actionInvocationHandler.Test)

	apiProtectedGroup.GET("/action-invocations", actionInvocationHandler.Paginate)
	apiProtectedGroup.GET("/action-invocations/:id", actionInvocationHandler.Get)
	apiProtectedGroup.DELETE("/action-invocations/:id", actionInvocationHandler.Delete)

	apiProtectedGroup.GET("/filter-presets/:type", filterPresetHandler.GetByType)
	apiProtectedGroup.POST("/filter-presets", middlewareEnforceJsonContentType(), filterPresetHandler.Create)
	apiProtectedGroup.DELETE("/filter-presets/:id", filterPresetHandler.Delete)

	apiProtectedGroup.GET("/comments/:updateId", commentHandler.GetAllByUpdateId)
	apiProtectedGroup.POST("/comments/:updateId", middlewareEnforceJsonContentType(), commentHandler.Create)
	apiProtectedGroup.PATCH("/comments/:id/content", middlewareEnforceJsonContentType(), commentHandler.UpdateContent)
	apiProtectedGroup.DELETE("/comments/:id", commentHandler.Delete)

	// start servers (run in separate goroutines)
	appSrv := s.newServer(appRouter, fmt.Sprintf("%s:%d", cfg.Server.Listen, cfg.Server.Port))
	prometheusSrv := s.newServer(promRouter, fmt.Sprintf("%s:%d", cfg.Prometheus.Listen, cfg.Prometheus.Port))

	s.startServer(appSrv, cfg.Server)

	if separatePromServer {
		s.startServer(prometheusSrv, cfg.Server)
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

	timeoutCtx, timeoutCancel := context.WithTimeout(s.ctx, cfg.Server.Timeout)
	defer timeoutCancel()

	shutdownDone := make(chan struct{})
	go func() {
		taskService.Stop()
		s.stopServer(s.ctx, appSrv)
		s.stopServer(s.ctx, prometheusSrv)
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

func (s *Server) newServer(r *gin.Engine, address string) *http.Server {
	if r == nil || address == "" {
		log.Fatal().Msg("Failed to create server, engine or address is nil")
		return nil
	}

	return &http.Server{
		Addr:    address,
		Handler: r,
	}
}

func (s *Server) startServer(h *http.Server, cfg *config.Server) {
	go func() {
		var e error
		log.Info().Msgf("Server listening on '%s'", h.Addr)

		if cfg.TlsEnabled {
			e = h.ListenAndServeTLS(cfg.TlsCertPath, cfg.TlsKeyPath)
		} else {
			e = h.ListenAndServe()
		}

		if e != nil && !errors.Is(e, http.ErrServerClosed) {
			log.Fatal().Msgf("Server cannot be started: %v", e)
		}
	}()
}

func (s *Server) stopServer(ctx context.Context, h *http.Server) {
	if h == nil {
		return
	}

	if err := h.Shutdown(ctx); err != nil {
		log.Fatal().Msgf("Shutdown failed, exited directly: %v", err)
	}

	log.Info().Msgf("Shutdown for '%s' complete", h.Addr)
}

func (s *Server) newEngine(middleware ...gin.HandlerFunc) *gin.Engine {
	r := gin.New()

	for _, m := range middleware {
		r.Use(m)
	}

	r.NoMethod(middlewareGlobalMethodNotAllowed())
	r.NoRoute(middlewareGlobalNotFound())

	return r
}

func (s *Server) applyAuthProvider(p auth.Provider, public *gin.RouterGroup, protected *gin.RouterGroup) {
	s.applyGroupMiddleware(public, p.PublicMiddleware()...)
	s.applyGroupMiddleware(protected, p.ProtectedMiddleware()...)

	protected.Use(func(c *gin.Context) {
		pass := p.IsAuthenticated(c)

		if !pass {
			c.Header(api.HeaderContentType, api.HeaderContentTypeApplicationJson)
			_ = c.AbortWithError(handler.ToHttpStatus(service_error.ErrUnauthorized), service_error.ErrUnauthorized)
			return
		}

		c.Next()
	})

	if p.Config().HasLoginRoute {
		public.Handle(p.Config().LoginMethod, p.Config().LoginPath, func(c *gin.Context) {
			var loginErr error
			var loginHttpStatus int

			if loginErr, loginHttpStatus = p.Login(c); loginErr != nil {
				httpErr := service_error.NewServiceErrorHttp(loginHttpStatus, loginErr)
				c.Header(api.HeaderContentType, api.HeaderContentTypeApplicationJson)
				_ = c.AbortWithError(handler.ToHttpStatus(httpErr), httpErr)
				return
			}

			c.Header(api.HeaderContentType, api.HeaderContentTypeApplicationJson)
			c.Status(loginHttpStatus)
		})
	}
	if p.Config().HasLogoutRoute {
		protected.Handle(p.Config().LogoutMethod, p.Config().LogoutPath, func(c *gin.Context) {
			var logoutErr error
			var logoutHttpStatus int

			if logoutErr, logoutHttpStatus = p.Logout(c); logoutErr != nil {
				httpErr := service_error.NewServiceErrorHttp(logoutHttpStatus, logoutErr)
				c.Header(api.HeaderContentType, api.HeaderContentTypeApplicationJson)
				_ = c.AbortWithError(handler.ToHttpStatus(httpErr), httpErr)
				return
			}

			c.Header(api.HeaderContentType, api.HeaderContentTypeApplicationJson)
			c.Status(logoutHttpStatus)
		})
	}
	if p.Config().HasProfileRoute {
		protected.Handle(p.Config().ProfileMethod, p.Config().ProfilePath, func(c *gin.Context) {
			var profile *auth.Profile
			var profileError error
			var profileHttpStatus int
			if profile, profileError, profileHttpStatus = p.Profile(c); profileError != nil {
				httpErr := service_error.NewServiceErrorHttp(profileHttpStatus, profileError)
				_ = c.AbortWithError(handler.ToHttpStatus(httpErr), httpErr)
				return
			}

			c.Header(api.HeaderContentType, api.HeaderContentTypeApplicationJson)
			c.JSON(profileHttpStatus, api.NewDataResponseWithPayload(profile))
		})
	}

	protected.Use(func(c *gin.Context) {
		var err error
		var profile *auth.Profile

		if profile, err = p.RouteContext(c); err != nil {
			log.Warn().Msg("Cannot pass route context")
		} else {
			c.Set(auth.RouteContext, profile)
		}

		c.Next()
	})
}

func (s *Server) applyGroupMiddleware(g *gin.RouterGroup, middleware ...gin.HandlerFunc) {
	for _, m := range middleware {
		g.Use(m)
	}
}
