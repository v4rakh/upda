package server

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"git.myservermanager.com/varakh/upda/internal/commons"
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

//go:embed migrations_postgres/*.sql
var migrationPostgresFS embed.FS

type appConfig struct {
	timeZone      string
	isDevelopment bool
	isDebug       bool
}

type embeddedWebInterfaceConfig struct {
	enabled          bool
	apiUrl           string
	title            string
	darkThemeEnabled bool
	footerEnabled    bool
}

type serverConfig struct {
	port                 int
	listen               string
	basePath             string
	tlsEnabled           bool
	tlsCertPath          string
	tlsKeyPath           string
	timeout              time.Duration
	corsAllowCredentials bool
	corsAllowOrigins     []string
	corsAllowMethods     []string
	corsAllowHeaders     []string
	corsExposeHeaders    []string
}

type authConfig struct {
	authMethod           string
	basicAuthUser        string
	basicAuthPassword    string
	basicAuthCredentials map[string]string
}

type taskConfig struct {
	updateCleanStaleEnabled   bool
	updateCleanStaleInterval  time.Duration
	updateCleanStaleMaxAge    time.Duration
	eventCleanStaleEnabled    bool
	eventCleanStaleInterval   time.Duration
	eventCleanStaleMaxAge     time.Duration
	actionsEnqueueEnabled     bool
	actionsEnqueueInterval    time.Duration
	actionsEnqueueBatchSize   int
	actionsInvokeEnabled      bool
	actionsInvokeInterval     time.Duration
	actionsInvokeBatchSize    int
	actionsInvokeMaxRetries   int
	actionsCleanStaleEnabled  bool
	actionsCleanStaleInterval time.Duration
	actionsCleanStaleMaxAge   time.Duration
	prometheusRefreshInterval time.Duration
}

type lockConfig struct {
	redisEnabled bool
	redisUrl     string
}

type webhookConfig struct {
	tokenLength int
}

type prometheusConfig struct {
	enabled            bool
	path               string
	secureTokenEnabled bool
	secureToken        string
}

type environment struct {
	appConfig                  *appConfig
	embeddedWebInterfaceConfig *embeddedWebInterfaceConfig
	authConfig                 *authConfig
	serverConfig               *serverConfig
	taskConfig                 *taskConfig
	lockConfig                 *lockConfig
	webhookConfig              *webhookConfig
	prometheusConfig           *prometheusConfig
	db                         *gorm.DB
}

func bootstrapEnvironment() *environment {
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

		if err = commons.CreateFileWithParent(logFile); err != nil {
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
	ac := &appConfig{
		timeZone:      os.Getenv(envTZ),
		isDebug:       isDebug,
		isDevelopment: isDevelopment,
	}

	// embedded web interface prometheusConfig
	embeddedWebInterfaceEnabled := os.Getenv(envEmbeddedWebInterfaceEnabled) == "true"

	if embeddedWebInterfaceEnabled {
		failIfEnvKeyNotPresent(envEmbeddedWebInterfaceApiUrl)
	}

	var embeddedWebInterfaceC *embeddedWebInterfaceConfig
	embeddedWebInterfaceC = &embeddedWebInterfaceConfig{
		enabled:          embeddedWebInterfaceEnabled,
		apiUrl:           os.Getenv(envEmbeddedWebInterfaceApiUrl),
		title:            os.Getenv(envEmbeddedWebInterfaceTitle),
		darkThemeEnabled: os.Getenv(envEmbeddedWebInterfaceDarkThemeEnabled) == "true",
		footerEnabled:    os.Getenv(envEmbeddedWebInterfaceFooter) == "true",
	}

	// server prometheusConfig
	var sc *serverConfig

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

	sc = &serverConfig{
		port:                 serverPort,
		listen:               os.Getenv(envServerListen),
		basePath:             os.Getenv(envServerBasePath),
		timeout:              serverTimeout,
		tlsEnabled:           serverTlsEnabled,
		tlsCertPath:          os.Getenv(envServerTlsCertPath),
		tlsKeyPath:           os.Getenv(envServerTlsKeyPath),
		corsAllowCredentials: os.Getenv(envCorsAllowCredentials) == "true",
		corsExposeHeaders:    []string{os.Getenv(envCorsExposeHeaders)},
		corsAllowOrigins:     []string{os.Getenv(envCorsAllowOrigins)},
		corsAllowMethods:     []string{os.Getenv(envCorsAllowMethods)},
		corsAllowHeaders:     []string{os.Getenv(envCorsAllowHeaders)},
	}

	authMode := os.Getenv(envAuthMode)

	if authMode != authModeBasicSingle && authMode != authModeBasicCredentials {
		zap.L().Sugar().Fatalln("Invalid auth mode. Reason: must be one of ['basic_single','basic_credentials'")
	}

	authC := &authConfig{
		authMethod: authMode,
	}

	if authModeBasicSingle == authMode {
		failIfEnvKeyNotPresent(envBasicAuthUser)
		failIfEnvKeyNotPresent(envBasicAuthPassword)
		authC.basicAuthUser = os.Getenv(envBasicAuthUser)
		authC.basicAuthPassword = os.Getenv(envBasicAuthPassword)
	}
	if authModeBasicCredentials == authMode {
		failIfEnvKeyNotPresent(envBasicAuthCredentials)
		authC.basicAuthCredentials = parseBasicAuthCredentials(envBasicAuthCredentials)
	}

	// task prometheusConfig
	var tc *taskConfig

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

	tc = &taskConfig{
		updateCleanStaleEnabled:   os.Getenv(envTaskUpdateCleanStaleEnabled) == "true",
		updateCleanStaleInterval:  updateCleanStaleInterval,
		updateCleanStaleMaxAge:    updateCleanStaleMaxAge,
		eventCleanStaleEnabled:    os.Getenv(envTaskEventCleanStaleEnabled) == "true",
		eventCleanStaleInterval:   eventCleanStaleInterval,
		eventCleanStaleMaxAge:     eventCleanStaleMaxAge,
		actionsEnqueueEnabled:     os.Getenv(envTaskActionsEnqueueEnabled) == "true",
		actionsEnqueueInterval:    actionsEnqueueInterval,
		actionsEnqueueBatchSize:   actionsEnqueueBatchSize,
		actionsInvokeEnabled:      os.Getenv(envTaskActionsInvokeEnabled) == "true",
		actionsInvokeInterval:     actionsInvokeInterval,
		actionsInvokeBatchSize:    actionsInvokeBatchSize,
		actionsInvokeMaxRetries:   actionsInvokeMaxRetries,
		actionsCleanStaleEnabled:  os.Getenv(envTaskActionsCleanStaleEnabled) == "true",
		actionsCleanStaleInterval: actionsCleanStaleInterval,
		actionsCleanStaleMaxAge:   actionsCleanStaleMaxAge,
		prometheusRefreshInterval: prometheusRefreshInterval,
	}

	var lc *lockConfig
	lc = &lockConfig{
		redisEnabled: os.Getenv(envLockRedisEnabled) == "true",
		redisUrl:     os.Getenv(envLockRedisUrl),
	}

	if lc.redisEnabled {
		failIfEnvKeyNotPresent(envLockRedisUrl)
	}

	webhookTokenLength := -1
	if webhookTokenLength, err = strconv.Atoi(os.Getenv(envWebhooksTokenLength)); err != nil {
		zap.L().Sugar().Fatalf("Invalid webhook token length. Reason: %v", err)
	}
	if webhookTokenLength <= 0 {
		zap.L().Sugar().Fatalln("Invalid webhook token length. Reason: must be a positive number")
	}

	wc := &webhookConfig{
		tokenLength: webhookTokenLength,
	}

	pc := &prometheusConfig{
		enabled:            os.Getenv(envPrometheusEnabled) == "true",
		path:               os.Getenv(envPrometheusMetricsPath),
		secureTokenEnabled: os.Getenv(envPrometheusSecureTokenEnabled) == "true",
		secureToken:        os.Getenv(envPrometheusSecureToken),
	}

	if pc.enabled {
		failIfEnvKeyNotPresent(envPrometheusMetricsPath)
	}

	if pc.enabled && pc.secureTokenEnabled {
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

	env := &environment{appConfig: ac,
		embeddedWebInterfaceConfig: embeddedWebInterfaceC,
		authConfig:                 authC,
		serverConfig:               sc,
		taskConfig:                 tc,
		lockConfig:                 lc,
		webhookConfig:              wc,
		prometheusConfig:           pc,
		db:                         db}

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

	zap.L().Sugar().Infof("AppConfig %+v", env.appConfig)
	zap.L().Sugar().Infof("EmbeddedWebInterfaceConfig %+v", env.embeddedWebInterfaceConfig)
	zap.L().Info("AuthConfig ***REDACTED***")
	zap.L().Sugar().Infof("ServerConfig %+v", env.serverConfig)
	zap.L().Sugar().Infof("TaskConfig %+v", env.taskConfig)
	zap.L().Info("LockConfig ***REDACTED***")
	zap.L().Sugar().Infof("WebhookConfig %+v", env.webhookConfig)
	zap.L().Info("PrometheusConfig ***REDACTED***")

	return env
}

func bootstrapFromEnvironmentAndValidate() {
	failIfEnvKeyNotPresent(envSecret)

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
