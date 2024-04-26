package server

import (
	"errors"
	"fmt"
	"git.myservermanager.com/varakh/upda/util"
	"github.com/adrg/xdg"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
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

type appConfig struct {
	timeZone      string
	isDevelopment bool
	isDebug       bool
}

type serverConfig struct {
	port             int
	listen           string
	tlsEnabled       bool
	tlsCertPath      string
	tlsKeyPath       string
	timeout          time.Duration
	corsAllowOrigin  []string
	corsAllowMethods []string
	corsAllowHeaders []string
}

type authConfig struct {
	authMethod           string
	basicAuthUser        string
	basicAuthPassword    string
	basicAuthCredentials map[string]string
}

type taskConfig struct {
	updateCleanStaleEnabled   bool
	updateCleanStaleInterval  string
	updateCleanStaleMaxAge    time.Duration
	eventCleanStaleEnabled    bool
	eventCleanStaleInterval   string
	eventCleanStaleMaxAge     time.Duration
	actionsEnqueueEnabled     bool
	actionsEnqueueInterval    string
	actionsEnqueueBatchSize   int
	actionsInvokeEnabled      bool
	actionsInvokeInterval     string
	actionsInvokeBatchSize    int
	actionsInvokeMaxRetries   int
	actionsCleanStaleEnabled  bool
	actionsCleanStaleInterval string
	actionsCleanStaleMaxAge   time.Duration
	prometheusRefreshInterval string
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

type Environment struct {
	appConfig        *appConfig
	authConfig       *authConfig
	serverConfig     *serverConfig
	taskConfig       *taskConfig
	lockConfig       *lockConfig
	webhookConfig    *webhookConfig
	prometheusConfig *prometheusConfig
	db               *gorm.DB
}

func bootstrapEnvironment() *Environment {
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

		if err = util.CreateFileWithParent(logFile); err != nil {
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
	defer zapLogger.Sync()
	zap.ReplaceGlobals(zapLogger)

	// assign defaults from given environment variables and validate
	bootstrapFromEnvironmentAndValidate()

	// parse environment variables in actual configuration structs
	// app config
	ac := &appConfig{
		timeZone:      os.Getenv(envTZ),
		isDebug:       isDebug,
		isDevelopment: isDevelopment,
	}

	// server config
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
		port:             serverPort,
		timeout:          serverTimeout,
		listen:           os.Getenv(envServerListen),
		tlsEnabled:       serverTlsEnabled,
		tlsCertPath:      os.Getenv(envServerTlsCertPath),
		tlsKeyPath:       os.Getenv(envServerTlsKeyPath),
		corsAllowOrigin:  []string{os.Getenv(envCorsAllowOrigin)},
		corsAllowMethods: []string{os.Getenv(envCorsAllowMethods)},
		corsAllowHeaders: []string{os.Getenv(envCorsAllowHeaders)},
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

	// task config
	var tc *taskConfig

	var updateCleanStaleMaxAge time.Duration
	if updateCleanStaleMaxAge, errParse = time.ParseDuration(os.Getenv(envTaskUpdateCleanStaleMaxAge)); errParse != nil {
		zap.L().Sugar().Fatalf("Could not parse max age for cleaning stale updates. Reason: %s", errParse.Error())
	}

	var eventCleanStaleMaxAge time.Duration
	if eventCleanStaleMaxAge, errParse = time.ParseDuration(os.Getenv(envTaskEventCleanStaleMaxAge)); errParse != nil {
		zap.L().Sugar().Fatalf("Could not parse max age for cleaning stale events. Reason: %s", errParse.Error())
	}

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

	var actionsCleanStaleMaxAge time.Duration
	if actionsCleanStaleMaxAge, errParse = time.ParseDuration(os.Getenv(envTaskActionsCleanStaleMaxAge)); errParse != nil {
		zap.L().Sugar().Fatalf("Could not parse max age for cleaning stale actions. Reason: %s", errParse.Error())
	}

	tc = &taskConfig{
		updateCleanStaleEnabled:   os.Getenv(envTaskUpdateCleanStaleEnabled) == "true",
		updateCleanStaleInterval:  os.Getenv(envTaskUpdateCleanStaleInterval),
		updateCleanStaleMaxAge:    updateCleanStaleMaxAge,
		eventCleanStaleEnabled:    os.Getenv(envTaskEventCleanStaleEnabled) == "true",
		eventCleanStaleInterval:   os.Getenv(envTaskEventCleanStaleInterval),
		eventCleanStaleMaxAge:     eventCleanStaleMaxAge,
		actionsEnqueueEnabled:     os.Getenv(envTaskActionsEnqueueEnabled) == "true",
		actionsEnqueueInterval:    os.Getenv(envTaskActionsEnqueueInterval),
		actionsEnqueueBatchSize:   actionsEnqueueBatchSize,
		actionsInvokeEnabled:      os.Getenv(envTaskActionsInvokeEnabled) == "true",
		actionsInvokeInterval:     os.Getenv(envTaskActionsInvokeInterval),
		actionsInvokeBatchSize:    actionsInvokeBatchSize,
		actionsInvokeMaxRetries:   actionsInvokeMaxRetries,
		actionsCleanStaleEnabled:  os.Getenv(envTaskActionsCleanStaleEnabled) == "true",
		actionsCleanStaleInterval: os.Getenv(envTaskActionsCleanStaleInterval),
		actionsCleanStaleMaxAge:   actionsCleanStaleMaxAge,
		prometheusRefreshInterval: os.Getenv(envTaskPrometheusRefreshInterval),
	}

	var lc *lockConfig
	lc = &lockConfig{
		redisEnabled: os.Getenv(envLockRedisEnabled) == "true",
		redisUrl:     os.Getenv(envLockRedisUrl),
	}

	webhookTokenLength := 32
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

	if pc.enabled && pc.secureTokenEnabled {
		failIfEnvKeyNotPresent(envPrometheusSecureToken)
	}

	// database setup
	gormConfig := &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)}
	if isDebug && isDevelopment {
		gormZapLogger := zap.Must(zapConfig.Build())
		defer gormZapLogger.Sync()
		gormLogger := zapgorm2.New(gormZapLogger)
		gormConfig = &gorm.Config{Logger: gormLogger}
	}

	var db *gorm.DB
	zap.L().Sugar().Infof("Using database type '%s'", os.Getenv(envDbType))

	if os.Getenv(envDbType) == dbTypeSqlite {
		if os.Getenv(envDbSqliteFile) == "" {
			var defaultDbFile string
			if defaultDbFile, err = xdg.DataFile(Name + "/" + dbTypeSqliteDbNameDefault); err != nil {
				zap.L().Sugar().Fatalf("Database file '%s' could not be created. Reason: %v", defaultDbFile, err)
			}
			setEnvKeyDefault(envDbSqliteFile, defaultDbFile)
		}

		dbFile := os.Getenv(envDbSqliteFile)
		zap.L().Sugar().Infof("Using database file '%s'", dbFile)

		if err = util.CreateFileWithParent(dbFile); err != nil {
			zap.L().Sugar().Fatalf("Database file '%s' cannot be created: %v", dbFile, err)
		}

		if db, err = gorm.Open(sqlite.Open(dbFile), gormConfig); err != nil {
			zap.L().Sugar().Fatalf("Could not setup database: %v", err)
		}

		if res := db.Exec("PRAGMA foreign_keys = ON"); res.Error != nil {
			zap.L().Sugar().Fatalf("Could not execute foreign key for SQLite: %v", res.Error)
		}

		sqlDb, _ := db.DB()
		sqlDb.SetMaxOpenConns(1)
		zap.L().Sugar().Infof("SQLite: restricting max connections to '1'")
	} else if os.Getenv(envDbType) == dbTypePostgres {
		host := os.Getenv(envDbPostgresHost)
		port := os.Getenv(envDbPostgresPort)
		dbUser := os.Getenv(envDbPostgresUser)
		dbPass := os.Getenv(envDbPostgresPassword)
		dbName := os.Getenv(envDbPostgresName)
		dbTZ := os.Getenv(envDbPostgresTimeZone)

		if host == "" || port == "" || dbUser == "" || dbPass == "" || dbName == "" || dbTZ == "" {
			zap.L().Sugar().Fatalf("Some configuration for database type '%s' is missing", os.Getenv(envDbType))
		}

		dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=%v", host, dbUser, dbPass, dbName, port, dbTZ)
		if db, err = gorm.Open(postgres.Open(dsn), gormConfig); err != nil {
			zap.L().Sugar().Fatalf("Could not setup database: %v", err)
		}
	} else {
		zap.L().Sugar().Fatalf("Database type '%s' or '%s' is required", dbTypeSqlite, dbTypePostgres)
	}

	if db == nil {
		zap.L().Sugar().Fatalf("Could not setup database")
	}

	env := &Environment{appConfig: ac,
		authConfig:       authC,
		serverConfig:     sc,
		taskConfig:       tc,
		lockConfig:       lc,
		webhookConfig:    wc,
		prometheusConfig: pc,
		db:               db}

	if err = env.db.AutoMigrate(&Update{}, &Webhook{}, &Event{}, &Secret{}, &Action{}, &ActionInvocation{}); err != nil {
		zap.L().Sugar().Fatalf("Could not migrate database schema: %s", err)
	}

	zap.L().Sugar().Infof("appConfig %+v", env.appConfig)
	zap.L().Sugar().Infof("serverConfig %+v", env.serverConfig)
	zap.L().Sugar().Infof("taskConfig %+v", env.taskConfig)
	zap.L().Sugar().Infof("webhookConfig %+v", env.webhookConfig)

	return env
}

func bootstrapFromEnvironmentAndValidate() {
	failIfEnvKeyNotPresent(envSecret)

	// auth mode
	setEnvKeyDefault(envAuthMode, authModeDefault)

	// app
	setEnvKeyDefault(envTZ, tzDefault)

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
	setEnvKeyDefault(envDbType, dbTypeSqlite)

	if os.Getenv(envDbType) == dbTypePostgres {
		setEnvKeyDefault(envDbPostgresHost, dbTypePostgresHostDefault)
		setEnvKeyDefault(envDbPostgresPort, dbTypePostgresPortDefault)
		setEnvKeyDefault(envDbPostgresTimeZone, dbTypePostgresTZDefault)
	}

	// server
	setEnvKeyDefault(envServerPort, serverPortDefault)
	setEnvKeyDefault(envServerListen, serverListenDefault)
	setEnvKeyDefault(envServerTlsEnabled, serverTlsEnabledDefault)
	setEnvKeyDefault(envCorsAllowOrigin, corsAllowOriginDefault)
	setEnvKeyDefault(envCorsAllowMethods, corsAllowMethodsDefault)
	setEnvKeyDefault(envCorsAllowHeaders, corsAllowHeadersDefault)
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

		zap.L().Sugar().Infof("Set '%s' to '%s'", key, defaultValue)
	}
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
