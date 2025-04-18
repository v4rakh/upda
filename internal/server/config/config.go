package config

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"git.myservermanager.com/varakh/upda/internal/file"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	migratepostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"moul.io/zapgorm2"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	envDevelopment = "DEVELOPMENT"

	envLoggingLevel        = "LOGGING_LEVEL"
	loggingLevelDefault    = "info"
	envLoggingEncoding     = "LOGGING_ENCODING"
	loggingEncodingDefault = "json"

	envLoggingDirectory    = "LOGGING_DIRECTORY"
	loggingFileNameDefault = "upda.log"

	EnvSecret = "SECRET"

	envTZ     = "TZ"
	tzDefault = "Etc/UTC"

	envEmbeddedWebInterfaceEnabled     = "EMBEDDED_WEB_INTERFACE_ENABLED"
	embeddedWebInterfaceEnabledDefault = "true"

	envEmbeddedWebInterfaceApiUrl = "EMBEDDED_WEB_INTERFACE_API_URL"

	envEmbeddedWebInterfaceTitle     = "EMBEDDED_WEB_INTERFACE_TITLE"
	embeddedWebInterfaceTitleDefault = "upda"

	envEmbeddedWebInterfaceDarkThemeEnabled     = "EMBEDDED_WEB_INTERFACE_DARK_THEME_ENABLED"
	embeddedWebInterfaceDarkThemeEnabledDefault = "false"

	envEmbeddedWebInterfaceFooter            = "EMBEDDED_WEB_INTERFACE_FOOTER_ENABLED"
	embeddedWebInterfaceFooterEnabledDefault = "true"

	envAuthMode              = "AUTH_MODE"
	authModeDefault          = AuthModeBasicSingle
	AuthModeBasicSingle      = "basic_single"
	AuthModeBasicCredentials = "basic_credentials"
	envBasicAuthUser         = "BASIC_AUTH_USER"
	envBasicAuthPassword     = "BASIC_AUTH_PASSWORD"
	envBasicAuthCredentials  = "BASIC_AUTH_CREDENTIALS"

	envServerPort           = "SERVER_PORT"
	serverPortDefault       = "8080"
	envServerListen         = "SERVER_LISTEN"
	serverListenDefault     = ""
	envServerBasePath       = "SERVER_BASE_PATH"
	serverBasePathDefault   = "/"
	envServerTlsEnabled     = "SERVER_TLS_ENABLED"
	serverTlsEnabledDefault = "false"
	envServerTlsCertPath    = "SERVER_TLS_CERT_PATH"
	envServerTlsKeyPath     = "SERVER_TLS_KEY_PATH"

	envServerTimeout     = "SERVER_TIMEOUT"
	serverTimeoutDefault = "1s"

	envCorsAllowOrigins         = "CORS_ALLOW_ORIGINS"
	envCorsAllowMethods         = "CORS_ALLOW_METHODS"
	envCorsAllowHeaders         = "CORS_ALLOW_HEADERS"
	envCorsAllowCredentials     = "CORS_ALLOW_CREDENTIALS"
	envCorsExposeHeaders        = "CORS_EXPOSE_HEADERS"
	corsAllowOriginsDefault     = "*"
	corsAllowMethodsDefault     = "HEAD, GET, POST, PUT, PATCH, DELETE, OPTIONS"
	corsAllowHeadersDefault     = "Authorization, Content-Type"
	corsAllowCredentialsDefault = "true"
	corsExposeHeadersDefault    = "*"

	dbTypePostgres = "postgres"

	envDbType                 = "DB_TYPE"
	envDbPostgresHost         = "DB_POSTGRES_HOST"
	envDbPostgresPort         = "DB_POSTGRES_PORT"
	envDbPostgresName         = "DB_POSTGRES_NAME"
	envDbPostgresTimeZone     = "DB_POSTGRES_TZ"
	envDbPostgresUser         = "DB_POSTGRES_USER"
	envDbPostgresPassword     = "DB_POSTGRES_PASSWORD"
	dbTypePostgresHostDefault = "localhost"
	dbTypePostgresPortDefault = "5432"
	dbTypePostgresTZDefault   = "Etc/UTC"

	envDbMigrationEnabled     = "DB_MIGRATION_ENABLED"
	dbMigrationEnabledDefault = "true"

	envTaskPrometheusRefreshInterval = "TASK_PROMETHEUS_REFRESH_INTERVAL"
	taskPrometheusRefreshDefault     = "60s"

	envWebhooksTokenLength     = "WEBHOOKS_TOKEN_LENGTH"
	webhooksTokenLengthDefault = "32"

	envPrometheusEnabled                = "PROMETHEUS_ENABLED"
	envPrometheusMetricsPath            = "PROMETHEUS_METRICS_PATH"
	envPrometheusSecureTokenEnabled     = "PROMETHEUS_SECURE_TOKEN_ENABLED"
	envPrometheusSecureToken            = "PROMETHEUS_SECURE_TOKEN"
	prometheusEnabledDefault            = "false"
	prometheusMetricsPathDefault        = "/metrics"
	prometheusSecureTokenEnabledDefault = "true"

	envTaskUpdateCleanStaleEnabled      = "TASK_UPDATE_CLEAN_STALE_ENABLED"
	envTaskUpdateCleanStaleInterval     = "TASK_UPDATE_CLEAN_STALE_INTERVAL"
	envTaskUpdateCleanStaleMaxAge       = "TASK_UPDATE_CLEAN_STALE_MAX_AGE"
	taskUpdateCleanStaleEnabledDefault  = "false"
	taskUpdateCleanStaleIntervalDefault = "1h"
	taskUpdateCleanStaleMaxAgeDefault   = "720h"

	envTaskEventCleanStaleEnabled      = "TASK_EVENT_CLEAN_STALE_ENABLED"
	envTaskEventCleanStaleInterval     = "TASK_EVENT_CLEAN_STALE_INTERVAL"
	envTaskEventCleanStaleMaxAge       = "TASK_EVENT_CLEAN_STALE_MAX_AGE"
	taskEventCleanStaleEnabledDefault  = "false"
	taskEventCleanStaleIntervalDefault = "8h"
	taskEventCleanStaleMaxAgeDefault   = "2190h"

	envTaskActionsEnqueueEnabled       = "TASK_ACTIONS_ENQUEUE_ENABLED"
	envTaskActionsEnqueueInterval      = "TASK_ACTIONS_ENQUEUE_INTERVAL"
	envTaskActionsEnqueueBatchSize     = "TASK_ACTIONS_ENQUEUE_BATCH_SIZE"
	taskActionsEnqueueEnabledDefault   = "true"
	taskActionsEnqueueIntervalDefault  = "10s"
	taskActionsEnqueueBatchSizeDefault = "1"

	envTaskActionsInvokeEnabled        = "TASK_ACTIONS_INVOKE_ENABLED"
	envTaskActionsInvokeInterval       = "TASK_ACTIONS_INVOKE_INTERVAL"
	envTaskActionsInvokeBatchSize      = "TASK_ACTIONS_INVOKE_BATCH_SIZE"
	envTaskActionsInvokeMaxRetries     = "TASK_ACTIONS_INVOKE_MAX_RETRIES"
	taskActionsInvokeEnabledDefault    = "true"
	taskActionsInvokeIntervalDefault   = "10s"
	taskActionsInvokeBatchSizeDefault  = "1"
	taskActionsInvokeMaxRetriesDefault = "3"

	envTaskActionsCleanStaleEnabled      = "TASK_ACTIONS_CLEAN_STALE_ENABLED"
	envTaskActionsCleanStaleInterval     = "TASK_ACTIONS_CLEAN_STALE_INTERVAL"
	envTaskActionsCleanStaleMaxAge       = "TASK_ACTIONS_CLEAN_STALE_MAX_AGE"
	taskActionsCleanStaleEnabledDefault  = "true"
	taskActionsCleanStaleIntervalDefault = "12h"
	taskActionsCleanStaleMaxAgeDefault   = "720h"

	envLockRedisEnabled = "LOCK_REDIS_ENABLED"
	envLockRedisUrl     = "LOCK_REDIS_URL"
	redisEnabledDefault = "false"
)

//go:embed migrations_postgres/*.sql
var migrationPostgresFS embed.FS

type App struct {
	TimeZone      string
	IsDevelopment bool
	IsDebug       bool
}

type EmbeddedWebInterface struct {
	Enabled          bool
	ApiUrl           string
	Title            string
	DarkThemeEnabled bool
	FooterEnabled    bool
}

type Server struct {
	Port                 int
	Listen               string
	BasePath             string
	TlsEnabled           bool
	TlsCertPath          string
	TlsKeyPath           string
	Timeout              time.Duration
	CorsAllowCredentials bool
	CorsAllowOrigins     []string
	CorsAllowMethods     []string
	CorsAllowHeaders     []string
	CorsExposeHeaders    []string
}

type Auth struct {
	AuthMethod           string
	BasicAuthUser        string
	BasicAuthPassword    string
	BasicAuthCredentials map[string]string
}

type Task struct {
	UpdateCleanStaleEnabled   bool
	UpdateCleanStaleInterval  time.Duration
	UpdateCleanStaleMaxAge    time.Duration
	EventCleanStaleEnabled    bool
	EventCleanStaleInterval   time.Duration
	EventCleanStaleMaxAge     time.Duration
	ActionsEnqueueEnabled     bool
	ActionsEnqueueInterval    time.Duration
	ActionsEnqueueBatchSize   int
	ActionsInvokeEnabled      bool
	ActionsInvokeInterval     time.Duration
	ActionsInvokeBatchSize    int
	ActionsInvokeMaxRetries   int
	ActionsCleanStaleEnabled  bool
	ActionsCleanStaleInterval time.Duration
	ActionsCleanStaleMaxAge   time.Duration
	PrometheusRefreshInterval time.Duration
}

type Lock struct {
	RedisEnabled bool
	RedisUrl     string
}

type Webhook struct {
	TokenLength int
}

type Prometheus struct {
	Enabled            bool
	Path               string
	SecureTokenEnabled bool
	SecureToken        string
}

type Configuration struct {
	App                  *App
	EmbeddedWebInterface *EmbeddedWebInterface
	Auth                 *Auth
	Server               *Server
	Task                 *Task
	Lock                 *Lock
	Webhook              *Webhook
	Prometheus           *Prometheus
	Database             *gorm.DB
}

func BootstrapFromEnv() *Configuration {
	var err error

	// bootstrap logging (configured independently and required before any other action)
	loggingLevel := os.Getenv(envLoggingLevel)
	if loggingLevel == "" {
		if err = os.Setenv(envLoggingLevel, loggingLevelDefault); err != nil {
			log.Fatalf("Cannot set logging level: %v", err)
		}
		loggingLevel = os.Getenv(envLoggingLevel)
	}
	var level zap.AtomicLevel
	if level, err = zap.ParseAtomicLevel(loggingLevel); err != nil {
		log.Fatalf("Cannot parse logging level: %v", err)
	}
	loggingEncoding := os.Getenv(envLoggingEncoding)
	if loggingEncoding == "" {
		if err = os.Setenv(envLoggingEncoding, loggingEncodingDefault); err != nil {
			log.Fatalf("Cannot set logging encoding: %v", err)
		}
		loggingEncoding = os.Getenv(envLoggingEncoding)
	}
	if loggingEncoding != "json" && loggingEncoding != "console" {
		log.Fatalf("Cannot parse logging level: %v", errors.New("only 'json' and 'console' are allowed logging encodings"))
	}
	isDebug := level.Level() == zap.DebugLevel
	isDevelopment := os.Getenv(envDevelopment) == "true"
	var loggingEncoderConfig zapcore.EncoderConfig
	if loggingEncoding == "json" {
		loggingEncoderConfig = zap.NewProductionEncoderConfig()
	} else {
		loggingEncoderConfig = zap.NewDevelopmentEncoderConfig()
	}

	logPaths := []string{"stderr"}
	loggingDirectory := os.Getenv(envLoggingDirectory)

	if loggingDirectory != "" {
		logFile := filepath.Join(loggingDirectory, loggingFileNameDefault)

		if err = file.CreateFileWithParent(logFile); err != nil {
			log.Fatalf("Log file '%s' cannot be created: %v", loggingDirectory, err)
		}

		logPaths = append(logPaths, logFile)
	}

	var zapConfig *zap.Config
	if isDebug {
		zapConfig = &zap.Config{
			Level:            level,
			Development:      isDevelopment,
			Encoding:         loggingEncoding,
			EncoderConfig:    loggingEncoderConfig,
			OutputPaths:      logPaths,
			ErrorOutputPaths: logPaths,
		}
	} else {
		zapConfig = &zap.Config{
			Level:       level,
			Development: isDevelopment,
			Sampling: &zap.SamplingConfig{
				Initial:    100,
				Thereafter: 100,
			},
			Encoding:         loggingEncoding,
			EncoderConfig:    loggingEncoderConfig,
			OutputPaths:      logPaths,
			ErrorOutputPaths: logPaths,
		}
	}

	zapLogger := zap.Must(zapConfig.Build())
	defer func(zapLogger *zap.Logger) {
		_ = zapLogger.Sync()
	}(zapLogger)
	zap.ReplaceGlobals(zapLogger)

	// assign defaults from given environment variables and validate
	bootstrapFromEnvironmentAndValidate()

	// parse environment variables in actual configuration structs
	// app prometheusConfig
	ac := &App{
		TimeZone:      os.Getenv(envTZ),
		IsDebug:       isDebug,
		IsDevelopment: isDevelopment,
	}

	// embedded web interface prometheusConfig
	embeddedWebInterfaceEnabled := os.Getenv(envEmbeddedWebInterfaceEnabled) == "true"

	if embeddedWebInterfaceEnabled {
		failIfEnvKeyNotPresent(envEmbeddedWebInterfaceApiUrl)
	}

	var embeddedWebInterfaceC *EmbeddedWebInterface
	embeddedWebInterfaceC = &EmbeddedWebInterface{
		Enabled:          embeddedWebInterfaceEnabled,
		ApiUrl:           os.Getenv(envEmbeddedWebInterfaceApiUrl),
		Title:            os.Getenv(envEmbeddedWebInterfaceTitle),
		DarkThemeEnabled: os.Getenv(envEmbeddedWebInterfaceDarkThemeEnabled) == "true",
		FooterEnabled:    os.Getenv(envEmbeddedWebInterfaceFooter) == "true",
	}

	// server prometheusConfig
	var sc *Server

	var serverPort int
	if serverPort, err = strconv.Atoi(os.Getenv(envServerPort)); err != nil {
		zap.L().Sugar().Fatalf("Invalid server port. Reason: %v", err)
	}

	serverTlsEnabled := os.Getenv(envServerTlsEnabled) == "true"

	if serverTlsEnabled {
		failIfEnvKeyNotPresent(envServerTlsCertPath)
		failIfEnvKeyNotPresent(envServerTlsKeyPath)
	}

	var serverTimeout time.Duration
	var errParse error
	if serverTimeout, errParse = time.ParseDuration(os.Getenv(envServerTimeout)); errParse != nil {
		zap.L().Sugar().Fatalf("Could not parse timeout. Reason: %s", errParse.Error())
	}

	sc = &Server{
		Port:                 serverPort,
		Listen:               os.Getenv(envServerListen),
		BasePath:             os.Getenv(envServerBasePath),
		Timeout:              serverTimeout,
		TlsEnabled:           serverTlsEnabled,
		TlsCertPath:          os.Getenv(envServerTlsCertPath),
		TlsKeyPath:           os.Getenv(envServerTlsKeyPath),
		CorsAllowCredentials: os.Getenv(envCorsAllowCredentials) == "true",
		CorsExposeHeaders:    []string{os.Getenv(envCorsExposeHeaders)},
		CorsAllowOrigins:     []string{os.Getenv(envCorsAllowOrigins)},
		CorsAllowMethods:     []string{os.Getenv(envCorsAllowMethods)},
		CorsAllowHeaders:     []string{os.Getenv(envCorsAllowHeaders)},
	}

	authMode := os.Getenv(envAuthMode)

	if authMode != AuthModeBasicSingle && authMode != AuthModeBasicCredentials {
		zap.L().Sugar().Fatalln("Invalid auth mode. Reason: must be one of ['basic_single','basic_credentials'")
	}

	authC := &Auth{
		AuthMethod: authMode,
	}

	if AuthModeBasicSingle == authMode {
		failIfEnvKeyNotPresent(envBasicAuthUser)
		failIfEnvKeyNotPresent(envBasicAuthPassword)
		authC.BasicAuthUser = os.Getenv(envBasicAuthUser)
		authC.BasicAuthPassword = os.Getenv(envBasicAuthPassword)
	}
	if AuthModeBasicCredentials == authMode {
		failIfEnvKeyNotPresent(envBasicAuthCredentials)
		authC.BasicAuthCredentials = parseBasicAuthCredentials(envBasicAuthCredentials)
	}

	// task prometheusConfig
	var tc *Task

	updateCleanStaleInterval := parseDuration(envTaskUpdateCleanStaleInterval)
	updateCleanStaleMaxAge := parseDuration(envTaskUpdateCleanStaleMaxAge)
	eventCleanStaleMaxAge := parseDuration(envTaskEventCleanStaleMaxAge)
	actionsCleanStaleMaxAge := parseDuration(envTaskActionsCleanStaleMaxAge)
	eventCleanStaleInterval := parseDuration(envTaskEventCleanStaleInterval)
	actionsEnqueueInterval := parseDuration(envTaskActionsEnqueueInterval)
	actionsInvokeInterval := parseDuration(envTaskActionsInvokeInterval)
	actionsCleanStaleInterval := parseDuration(envTaskActionsCleanStaleInterval)
	prometheusRefreshInterval := parseDuration(envTaskPrometheusRefreshInterval)

	var actionsEnqueueBatchSize int
	if actionsEnqueueBatchSize, err = strconv.Atoi(os.Getenv(envTaskActionsEnqueueBatchSize)); err != nil {
		zap.L().Sugar().Fatalf("Invalid actions enqueue batch size. Reason: %v", err)
	}
	if actionsEnqueueBatchSize <= 0 {
		zap.L().Sugar().Fatalf("Invalid actions enqueue batch size, must be a positive number.")
	}

	var actionsInvokeBatchSize int
	if actionsInvokeBatchSize, err = strconv.Atoi(os.Getenv(envTaskActionsInvokeBatchSize)); err != nil {
		zap.L().Sugar().Fatalf("Invalid actions invoke batch size. Reason: %v", err)
	}
	if actionsInvokeBatchSize <= 0 {
		zap.L().Sugar().Fatalf("Invalid actions invoke batch size, must be a positive number.")
	}

	var actionsInvokeMaxRetries int
	if actionsInvokeMaxRetries, err = strconv.Atoi(os.Getenv(envTaskActionsInvokeMaxRetries)); err != nil {
		zap.L().Sugar().Fatalf("Invalid actions invoke max retries. Reason: %v", err)
	}
	if actionsInvokeMaxRetries <= 0 {
		zap.L().Sugar().Fatalf("Invalid actions invoke max retries, must be a positive number.")
	}

	tc = &Task{
		UpdateCleanStaleEnabled:   os.Getenv(envTaskUpdateCleanStaleEnabled) == "true",
		UpdateCleanStaleInterval:  updateCleanStaleInterval,
		UpdateCleanStaleMaxAge:    updateCleanStaleMaxAge,
		EventCleanStaleEnabled:    os.Getenv(envTaskEventCleanStaleEnabled) == "true",
		EventCleanStaleInterval:   eventCleanStaleInterval,
		EventCleanStaleMaxAge:     eventCleanStaleMaxAge,
		ActionsEnqueueEnabled:     os.Getenv(envTaskActionsEnqueueEnabled) == "true",
		ActionsEnqueueInterval:    actionsEnqueueInterval,
		ActionsEnqueueBatchSize:   actionsEnqueueBatchSize,
		ActionsInvokeEnabled:      os.Getenv(envTaskActionsInvokeEnabled) == "true",
		ActionsInvokeInterval:     actionsInvokeInterval,
		ActionsInvokeBatchSize:    actionsInvokeBatchSize,
		ActionsInvokeMaxRetries:   actionsInvokeMaxRetries,
		ActionsCleanStaleEnabled:  os.Getenv(envTaskActionsCleanStaleEnabled) == "true",
		ActionsCleanStaleInterval: actionsCleanStaleInterval,
		ActionsCleanStaleMaxAge:   actionsCleanStaleMaxAge,
		PrometheusRefreshInterval: prometheusRefreshInterval,
	}

	var lc *Lock
	lc = &Lock{
		RedisEnabled: os.Getenv(envLockRedisEnabled) == "true",
		RedisUrl:     os.Getenv(envLockRedisUrl),
	}

	if lc.RedisEnabled {
		failIfEnvKeyNotPresent(envLockRedisUrl)
	}

	webhookTokenLength := -1
	if webhookTokenLength, err = strconv.Atoi(os.Getenv(envWebhooksTokenLength)); err != nil {
		zap.L().Sugar().Fatalf("Invalid webhook token length. Reason: %v", err)
	}
	if webhookTokenLength <= 0 {
		zap.L().Sugar().Fatalln("Invalid webhook token length. Reason: must be a positive number")
	}

	wc := &Webhook{
		TokenLength: webhookTokenLength,
	}

	pc := &Prometheus{
		Enabled:            os.Getenv(envPrometheusEnabled) == "true",
		Path:               os.Getenv(envPrometheusMetricsPath),
		SecureTokenEnabled: os.Getenv(envPrometheusSecureTokenEnabled) == "true",
		SecureToken:        os.Getenv(envPrometheusSecureToken),
	}

	if pc.Enabled {
		failIfEnvKeyNotPresent(envPrometheusMetricsPath)
	}

	if pc.Enabled && pc.SecureTokenEnabled {
		failIfEnvKeyNotPresent(envPrometheusSecureToken)
	}

	// database setup
	gormConfig := &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)}
	if isDebug && isDevelopment {
		gormZapLogger := zap.Must(zapConfig.Build())
		defer func(gormZapLogger *zap.Logger) {
			_ = gormZapLogger.Sync()
		}(gormZapLogger)
		gormLogger := zapgorm2.New(gormZapLogger)
		gormConfig = &gorm.Config{Logger: gormLogger}
	}

	var db *gorm.DB
	var migrationDriver database.Driver
	var migrationDatabaseName string
	var migrationFS source.Driver

	zap.L().Sugar().Infof("Using database type '%s'", os.Getenv(envDbType))

	if os.Getenv(envDbType) == dbTypePostgres {
		host := os.Getenv(envDbPostgresHost)
		port := os.Getenv(envDbPostgresPort)
		dbUser := os.Getenv(envDbPostgresUser)
		dbPass := os.Getenv(envDbPostgresPassword)
		dbName := os.Getenv(envDbPostgresName)
		dbTZ := os.Getenv(envDbPostgresTimeZone)
		migrationDatabaseName = dbName

		if host == "" || port == "" || dbUser == "" || dbPass == "" || dbName == "" || dbTZ == "" {
			zap.L().Sugar().Fatalf("Some configuration for database type '%s' is missing", os.Getenv(envDbType))
		}

		dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=%v", host, dbUser, dbPass, dbName, port, dbTZ)
		if db, err = gorm.Open(postgres.Open(dsn), gormConfig); err != nil {
			zap.L().Sugar().Fatalf("Could not setup database: %v", err)
		}

		var sqlDb *sql.DB
		if sqlDb, err = db.DB(); err != nil {
			zap.L().Sugar().Fatalf("Could not retrieve database: %v", err)
		}

		if err = sqlDb.Ping(); err != nil {
			zap.L().Sugar().Fatalf("Could not connect to database: %v", err)
		}

		if migrationDriver, err = migratepostgres.WithInstance(sqlDb, &migratepostgres.Config{}); err != nil {
			zap.L().Sugar().Fatalf("Could not create migration driver: %v", err)
		}

		if migrationFS, err = iofs.New(migrationPostgresFS, "migrations_postgres"); err != nil {
			zap.L().Sugar().Fatalf("Could not create migration source: %v", err)
		}
	} else {
		zap.L().Sugar().Fatalf("Database type '%s' is required", dbTypePostgres)
	}

	if db == nil {
		zap.L().Sugar().Fatalf("Could not setup database")
	}

	env := &Configuration{App: ac,
		EmbeddedWebInterface: embeddedWebInterfaceC,
		Auth:                 authC,
		Server:               sc,
		Task:                 tc,
		Lock:                 lc,
		Webhook:              wc,
		Prometheus:           pc,
		Database:             db}

	migrationEnabled := os.Getenv(envDbMigrationEnabled) == "true"
	if !migrationEnabled {
		zap.L().Warn("Database schema migration is disabled and not executed automatically. Make sure to run them manually, otherwise the application might misbehave. You can safely ignore this warning if application is started in high availability mode and you're sure necessary database schema already exists.")
	} else {
		var migrator *migrate.Migrate
		if migrator, err = migrate.NewWithInstance("iofs", migrationFS, migrationDatabaseName, migrationDriver); err != nil {
			zap.L().Sugar().Fatalf("Could not create database migration instance: %v", err)
		}

		var migrationVersion uint
		var migrationVersionDirty bool
		if migrationVersion, migrationVersionDirty, err = migrator.Version(); err != nil {
			if errors.Is(err, migrate.ErrNilVersion) {
				zap.L().Info("Database migration schema is uninitialized")
			} else {
				zap.L().Sugar().Fatalf("Could not retrieve database migration version: %v", err)
			}
		} else {
			zap.L().Sugar().Infof("Previous database migration version is '%d' (dirty '%v')", migrationVersion, migrationVersionDirty)
		}

		zap.L().Info("Applying necessary database migration steps...")
		if err = migrator.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				zap.L().Info("No database schema changes detected")
			} else {
				zap.L().Sugar().Fatalf("Could not migrate database schema: %v", err)
			}
		}

		zap.L().Info("Applied all necessary database migration steps successfully")
	}

	zap.L().Sugar().Infof("App %+v", env.App)
	zap.L().Sugar().Infof("EmbeddedWebInterface %+v", env.EmbeddedWebInterface)
	zap.L().Info("Auth ***REDACTED***")
	zap.L().Sugar().Infof("Server %+v", env.Server)
	zap.L().Sugar().Infof("Task %+v", env.Task)
	zap.L().Info("Lock ***REDACTED***")
	zap.L().Sugar().Infof("Webhook %+v", env.Webhook)
	zap.L().Info("Prometheus ***REDACTED***")

	return env
}

func bootstrapFromEnvironmentAndValidate() {
	failIfEnvKeyNotPresent(EnvSecret)

	// auth mode
	setEnvKeyDefault(envAuthMode, authModeDefault)

	// app
	setEnvKeyDefault(envTZ, tzDefault)

	// web
	setEnvKeyDefault(envEmbeddedWebInterfaceEnabled, embeddedWebInterfaceEnabledDefault)
	setEnvKeyDefault(envEmbeddedWebInterfaceTitle, embeddedWebInterfaceTitleDefault)
	setEnvKeyDefault(envEmbeddedWebInterfaceDarkThemeEnabled, embeddedWebInterfaceDarkThemeEnabledDefault)
	setEnvKeyDefault(envEmbeddedWebInterfaceFooter, embeddedWebInterfaceFooterEnabledDefault)

	// webhook
	setEnvKeyDefault(envWebhooksTokenLength, webhooksTokenLengthDefault)

	// lock
	setEnvKeyDefault(envLockRedisEnabled, redisEnabledDefault)

	// task
	setEnvKeyDefault(envTaskUpdateCleanStaleEnabled, taskUpdateCleanStaleEnabledDefault)
	setEnvKeyDefault(envTaskUpdateCleanStaleInterval, taskUpdateCleanStaleIntervalDefault)
	setEnvKeyDefault(envTaskUpdateCleanStaleMaxAge, taskUpdateCleanStaleMaxAgeDefault)

	setEnvKeyDefault(envTaskEventCleanStaleEnabled, taskEventCleanStaleEnabledDefault)
	setEnvKeyDefault(envTaskEventCleanStaleInterval, taskEventCleanStaleIntervalDefault)
	setEnvKeyDefault(envTaskEventCleanStaleMaxAge, taskEventCleanStaleMaxAgeDefault)

	setEnvKeyDefault(envTaskActionsEnqueueEnabled, taskActionsEnqueueEnabledDefault)
	setEnvKeyDefault(envTaskActionsEnqueueInterval, taskActionsEnqueueIntervalDefault)
	setEnvKeyDefault(envTaskActionsEnqueueBatchSize, taskActionsEnqueueBatchSizeDefault)

	setEnvKeyDefault(envTaskActionsInvokeEnabled, taskActionsInvokeEnabledDefault)
	setEnvKeyDefault(envTaskActionsInvokeInterval, taskActionsInvokeIntervalDefault)
	setEnvKeyDefault(envTaskActionsInvokeBatchSize, taskActionsInvokeBatchSizeDefault)
	setEnvKeyDefault(envTaskActionsInvokeMaxRetries, taskActionsInvokeMaxRetriesDefault)

	setEnvKeyDefault(envTaskActionsCleanStaleEnabled, taskActionsCleanStaleEnabledDefault)
	setEnvKeyDefault(envTaskActionsCleanStaleInterval, taskActionsCleanStaleIntervalDefault)
	setEnvKeyDefault(envTaskActionsCleanStaleMaxAge, taskActionsCleanStaleMaxAgeDefault)

	setEnvKeyDefault(envTaskPrometheusRefreshInterval, taskPrometheusRefreshDefault)

	// prometheus
	setEnvKeyDefault(envPrometheusEnabled, prometheusEnabledDefault)
	setEnvKeyDefault(envPrometheusMetricsPath, prometheusMetricsPathDefault)
	setEnvKeyDefault(envPrometheusSecureTokenEnabled, prometheusSecureTokenEnabledDefault)

	// db
	setEnvKeyDefault(envDbType, dbTypePostgres)
	setEnvKeyDefault(envDbMigrationEnabled, dbMigrationEnabledDefault)

	if os.Getenv(envDbType) == dbTypePostgres {
		setEnvKeyDefault(envDbPostgresHost, dbTypePostgresHostDefault)
		setEnvKeyDefault(envDbPostgresPort, dbTypePostgresPortDefault)
		setEnvKeyDefault(envDbPostgresTimeZone, dbTypePostgresTZDefault)
	}

	// server
	setEnvKeyDefault(envServerPort, serverPortDefault)
	setEnvKeyDefault(envServerListen, serverListenDefault)
	setEnvKeyDefault(envServerBasePath, serverBasePathDefault)
	setEnvKeyDefault(envServerTlsEnabled, serverTlsEnabledDefault)
	setEnvKeyDefault(envCorsAllowOrigins, corsAllowOriginsDefault)
	setEnvKeyDefault(envCorsAllowMethods, corsAllowMethodsDefault)
	setEnvKeyDefault(envCorsAllowHeaders, corsAllowHeadersDefault)
	setEnvKeyDefault(envCorsAllowCredentials, corsAllowCredentialsDefault)
	setEnvKeyDefault(envCorsExposeHeaders, corsExposeHeadersDefault)
	setEnvKeyDefault(envServerTimeout, serverTimeoutDefault)
}

func failIfEnvKeyNotPresent(key string) {
	if os.Getenv(key) == "" {
		zap.L().Sugar().Fatalf("Not all required ENV variables given. Please set '%s'", key)
	}
}

func setEnvKeyDefault(key string, defaultValue string) {
	var err error
	if os.Getenv(key) == "" {
		if err = os.Setenv(key, defaultValue); err != nil {
			zap.L().Sugar().Fatalf("Could not set default value for ENV variable '%s'", key)
		}

		zap.L().Sugar().Infof("Setting default for '%s' to '%s'", key, defaultValue)
	}
}

func parseDuration(envProperty string) time.Duration {
	var duration time.Duration
	var err error

	if duration, err = time.ParseDuration(os.Getenv(envProperty)); err != nil {
		zap.L().Sugar().Fatalf("Could not parse duration for '%s'. Reason: %s", envProperty, err.Error())
	}

	return duration
}

func parseBasicAuthCredentials(envProperty string) map[string]string {
	if envProperty == "" {
		zap.L().Sugar().Fatalln("Invalid env for parsing basic auth credentials")
	}
	credentialsFromEnv := os.Getenv(envProperty)

	var credentials []string
	credentials = strings.Split(credentialsFromEnv, ",")

	basicAuthCredentials := make(map[string]string)

	for _, c := range credentials {
		pair := strings.Split(c, "=")

		if len(pair) != 2 {
			zap.L().Sugar().Fatalln("Invalid basic auth credentials. Reason: credentials must be specified with the = separator per credential entry")
		}

		if pair[0] == "" {
			zap.L().Sugar().Fatalln("Invalid basic auth credentials. Reason: username must not be blank")
		}

		if pair[1] == "" {
			zap.L().Sugar().Fatalln("Invalid basic auth credentials. Reason: password must not be blank")
		}

		basicAuthCredentials[pair[0]] = pair[1]
	}

	return basicAuthCredentials
}
