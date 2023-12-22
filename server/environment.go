package server

import (
	"fmt"
	"github.com/adrg/xdg"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"moul.io/zapgorm2"
	"os"
	"strconv"
	"time"
)

type appConfig struct {
	timeZone string
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
	adminUser     string
	adminPassword string
}

type taskConfig struct {
	updateCleanStaleEnabled   bool
	updateCleanStaleInterval  string
	updateCleanStaleMaxAge    time.Duration
	eventCleanStaleEnabled    bool
	eventCleanStaleInterval   string
	eventCleanStaleMaxAge     time.Duration
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
	// logging (configured independently)
	var logger *zap.Logger
	var err error
	level := zap.NewAtomicLevelAt(zapcore.InfoLevel)

	envLoggingLevel := os.Getenv(envLoggingLevel)
	if envLoggingLevel != "" {
		if level, err = zap.ParseAtomicLevel(envLoggingLevel); err != nil {
			log.Fatalf("Cannot parse logging level: %v", err)
		}
	}

	logger, err = zap.NewDevelopment(zap.IncreaseLevel(level))
	if err != nil {
		log.Fatalf("Can't initialize logger: %v", err)
	}
	// flushes buffer, if any
	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	// assign defaults from given environment variables and validate
	bootstrapFromEnvironmentAndValidate()

	// parse environment variables in actual configuration structs
	// app config
	appConfig := &appConfig{
		timeZone: os.Getenv(envTZ),
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

	authConfig := &authConfig{
		adminUser:     os.Getenv(envAdminUser),
		adminPassword: os.Getenv(envAdminPassword),
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

	tc = &taskConfig{
		updateCleanStaleEnabled:   os.Getenv(envTaskUpdateCleanStaleEnabled) == "true",
		updateCleanStaleInterval:  os.Getenv(envTaskUpdateCleanStaleInterval),
		updateCleanStaleMaxAge:    updateCleanStaleMaxAge,
		eventCleanStaleEnabled:    os.Getenv(envTaskEventCleanStaleEnabled) == "true",
		eventCleanStaleInterval:   os.Getenv(envTaskEventCleanStaleInterval),
		eventCleanStaleMaxAge:     eventCleanStaleMaxAge,
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

	webhookConfig := &webhookConfig{
		tokenLength: webhookTokenLength,
	}

	prometheusConfig := &prometheusConfig{
		enabled:            os.Getenv(envPrometheusEnabled) == "true",
		path:               os.Getenv(envPrometheusMetricsPath),
		secureTokenEnabled: os.Getenv(envPrometheusSecureTokenEnabled) == "true",
		secureToken:        os.Getenv(envPrometheusSecureToken),
	}

	if prometheusConfig.enabled && prometheusConfig.secureTokenEnabled {
		failIfEnvKeyNotPresent(envPrometheusSecureToken)
	}

	// database setup
	gormLogger := zapgorm2.New(logger)
	gormLogger.SetAsDefault()

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

		if db, err = gorm.Open(sqlite.Open(dbFile), &gorm.Config{Logger: gormLogger}); err != nil {
			zap.L().Sugar().Fatalf("Could not setup database: %v", err)
		}

		if res := db.Exec("PRAGMA foreign_keys = ON"); res.Error != nil {
			zap.L().Sugar().Fatalf("Could not invoke foreign key for SQLite: %v", res.Error)
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
		if db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gormLogger}); err != nil {
			zap.L().Sugar().Fatalf("Could not setup database: %v", err)
		}
	} else {
		zap.L().Sugar().Fatalf("Database type '%s' or '%s' is required", dbTypeSqlite, dbTypePostgres)
	}

	if db == nil {
		zap.L().Sugar().Fatalf("Could not setup database")
	}

	env := &Environment{appConfig: appConfig,
		authConfig:       authConfig,
		serverConfig:     sc,
		taskConfig:       tc,
		lockConfig:       lc,
		webhookConfig:    webhookConfig,
		prometheusConfig: prometheusConfig,
		db:               db}

	if err = env.db.AutoMigrate(&Update{}, &Webhook{}, &Event{}); err != nil {
		zap.L().Sugar().Fatalf("Could not migrate database schema: %s", err)
	}

	zap.L().Sugar().Infof("appConfig %+v", env.appConfig)
	zap.L().Sugar().Infof("serverConfig %+v", env.serverConfig)
	zap.L().Sugar().Infof("taskConfig %+v", env.taskConfig)
	zap.L().Sugar().Infof("webhookConfig %+v", env.webhookConfig)

	return env
}

func bootstrapFromEnvironmentAndValidate() {
	// app
	setEnvKeyDefault(envTZ, tzDefault)

	failIfEnvKeyNotPresent(envAdminUser)
	failIfEnvKeyNotPresent(envAdminPassword)

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
