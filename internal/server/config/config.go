package config

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	golog "log"
	"net/http"
	"os"
	"time"

	"git.myservermanager.com/varakh/upda/internal/server/constant"
	"git.myservermanager.com/varakh/upda/internal/validate"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	migratepostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-envconfig"
	zerologgorm "github.com/skynet2/zerolog-gorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	EnvSecret = "SECRET"
)

//go:embed migrations_postgres/*.sql
var migrationPostgresFS embed.FS

type Logging struct {
	Encoding              constant.ConfigLogEncoding    `env:"LOGGING_ENCODING,default=console"                   validate:"required,oneof=json console"`
	EncodingColorize      bool                          `env:"LOGGING_ENCODING_COLORIZE,default=false"`
	EncodingErrorKey      string                        `env:"LOGGING_ENCODING_ERROR_KEY,default=error"           validate:"required"`
	EncodingFileKey       string                        `env:"LOGGING_ENCODING_FILE_KEY,default=file"             validate:"required"`
	EncodingFuncKey       string                        `env:"LOGGING_ENCODING_FUNC_KEY,default=func"             validate:"required"`
	EncodingLevelKey      string                        `env:"LOGGING_ENCODING_LEVEL_KEY,default=level"           validate:"required"`
	EncodingMessageKey    string                        `env:"LOGGING_ENCODING_MESSAGE_KEY,default=msg"           validate:"required"`
	EncodingStacktraceKey string                        `env:"LOGGING_ENCODING_STACKTRACE_KEY,default=stacktrace" validate:"required"`
	EncodingTimeEncoder   constant.ConfigLogTimeEncoder `env:"LOGGING_ENCODING_TIME_ENCODER,default=rfc3339"      validate:"required,oneof=epoch epochmillis epochnanos iso8601 rfc3339 rfc3339nano"`
	EncodingTimeKey       string                        `env:"LOGGING_ENCODING_TIME_KEY,default=ts"               validate:"required"`
	Level                 string                        `env:"LOGGING_LEVEL,default=info"                         validate:"required,oneof=trace debug info warn error fatal panic disabled"`
	LevelRequests         string                        `env:"LOGGING_LEVEL_REQUESTS,default=disabled"            validate:"required,oneof=trace debug info warn error fatal panic disabled"`
}

type App struct {
	TimeZone    string `env:"TZ,default=Etc/UTC"        validate:"required"`
	Development bool   `env:"DEVELOPMENT,default=false"`
}

type Secret struct {
	Secret string `env:"SECRET,required" validate:"required"`
}

type Server struct {
	Port        int           `env:"SERVER_PORT,default=8080"         validate:"gte=1"`
	Listen      string        `env:"SERVER_LISTEN"`
	BasePath    string        `env:"SERVER_BASE_PATH,default=/"       validate:"required"`
	TlsEnabled  bool          `env:"SERVER_TLS_ENABLED,default=false"`
	TlsCertPath string        `env:"SERVER_TLS_CERT_PATH"`
	TlsKeyPath  string        `env:"SERVER_TLS_KEY_PATH"`
	Timeout     time.Duration `env:"SERVER_TIMEOUT,default=10s"       validate:"gte=0"`
}

type Cors struct {
	AllowCredentials bool     `env:"CORS_ALLOW_CREDENTIALS,default=true"`
	AllowOrigins     []string `env:"CORS_ALLOW_ORIGINS,default=*"`
	AllowMethods     []string `env:"CORS_ALLOW_METHODS,default=HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS"`
	AllowHeaders     []string `env:"CORS_ALLOW_HEADERS,default=Authorization,Content-Type"`
	ExposeHeaders    []string `env:"CORS_EXPOSE_HEADERS,default=*"`
}

type Database struct {
	Type             constant.ConfigDatabaseType `env:"DB_TYPE,default=postgres"           validate:"required,oneof=postgres"`
	MigrationEnabled bool                        `env:"DB_MIGRATION_ENABLED,default=true"`
	PostgresHost     string                      `env:"DB_POSTGRES_HOST,default=localhost" validate:"required_if=Type postgres"`
	PostgresPort     int                         `env:"DB_POSTGRES_PORT,default=5432"      validate:"required_if=Type postgres"`
	PostgresName     string                      `env:"DB_POSTGRES_NAME"                   validate:"required_if=Type postgres"`
	PostgresTimeZone string                      `env:"DB_POSTGRES_TZ,default=Etc/UTC"     validate:"required_if=Type postgres"`
	PostgresUser     string                      `env:"DB_POSTGRES_USER"                   validate:"required_if=Type postgres"`
	PostgresPassword string                      `env:"DB_POSTGRES_PASSWORD"               validate:"required_if=Type postgres"`
}

type Webinterface struct {
	Enabled       bool   `env:"WEB_INTERFACE_ENABLED,default=true"`
	ApiUrl        string `env:"WEB_INTERFACE_API_URL,default=http://localhost:8080/api/" validate:"required_if=Enabled true"`
	Title         string `env:"WEB_INTERFACE_TITLE,default=upda"                         validate:"required_if=Enabled true"`
	FooterEnabled bool   `env:"WEB_INTERFACE_FOOTER_ENABLED,default=true"`
}

type WebinterfaceCacheControl struct {
	Enabled              bool           `env:"WEB_INTERFACE_CC_ENABLED,default=true"`
	MustRevalidate       bool           `env:"WEB_INTERFACE_CC_MUST_REVALIDATE,default=true"`
	NoCache              bool           `env:"WEB_INTERFACE_CC_NO_CACHE,default=false"`
	NoStore              bool           `env:"WEB_INTERFACE_CC_NO_STORE,default=false"`
	NoTransform          bool           `env:"WEB_INTERFACE_CC_NO_TRANSFORM,default=false"`
	Public               bool           `env:"WEB_INTERFACE_CC_PUBLIC,default=true"`
	Private              bool           `env:"WEB_INTERFACE_CC_PRIVATE,default=false"`
	ProxyRevalidate      bool           `env:"WEB_INTERFACE_CC_PROXY_REVALIDATE,default=true"`
	MaxAge               *time.Duration `env:"WEB_INTERFACE_CC_MAX_AGE,noinit,default=48h"`
	SMaxAge              *time.Duration `env:"WEB_INTERFACE_CC_SMAX_AGE,noinit"`
	Immutable            bool           `env:"WEB_INTERFACE_CC_IMMUTABLE,default=false"`
	StaleWhileRevalidate *time.Duration `env:"WEB_INTERFACE_CC_STALE_WHILE_REVALIDATE,noinit"`
	StaleIfError         *time.Duration `env:"WEB_INTERFACE_CC_STALE_IF_ERROR,noinit"`
}

type Auth struct {
	Type                   constant.ConfigAuthType            `env:"AUTH_TYPE,default=session"                        validate:"required,oneof=session"`
	SessionSecret          string                             `env:"AUTH_SESSION_SECRET,required"                     validate:"required_if=Type session"`
	SessionProvider        constant.ConfigAuthSessionProvider `env:"AUTH_SESSION_PROVIDER,default=single"             validate:"required_if=Type session,oneof=single credentials"`
	SessionUser            string                             `env:"AUTH_SESSION_USER"                                validate:"required_if=SessionProvider single"`
	SessionPassword        string                             `env:"AUTH_SESSION_PASSWORD"                            validate:"required_if=SessionProvider single"`
	SessionCredentials     map[string]string                  `env:"AUTH_SESSION_CREDENTIALS,separator=|,delimiter=;" validate:"required_if=SessionProvider credentials"`
	SessionCleanupEnabled  bool                               `env:"AUTH_SESSION_CLEANUP_ENABLED,default=true"`
	SessionCleanupInterval time.Duration                      `env:"AUTH_SESSION_CLEANUP_INTERVAL,default=1h"         validate:"required_if=Type session,gte=0"`
	SessionCookieMaxAge    time.Duration                      `env:"AUTH_SESSION_COOKIE_MAX_AGE,default=8h"           validate:"required_if=Type session"`
	SessionCookieName      string                             `env:"AUTH_SESSION_COOKIE_NAME,default=UPDA_SESSION"    validate:"required_if=Type session"`
	// if set to non-blank, subdomains send it
	SessionCookieDomain *string `env:"AUTH_SESSION_COOKIE_DOMAIN"`
	// cookie's scope (/ meaning when browser adds it automatically)
	SessionCookiePath string `env:"AUTH_SESSION_COOKIE_PATH,default=/" validate:"required_if=Type session"`
	// true means JavaScript cannot access it (recommended)
	SessionCookieHttpOnly bool `env:"AUTH_SESSION_COOKIE_HTTP_ONLY,default=true"`
	// true means SSL is required, otherwise cookie is not sent
	SessionCookieSecure bool `env:"AUTH_SESSION_COOKIE_SECURE,default=true"`
	// force mode when browser stores (recommended is strict)
	SessionCookieSameSite     constant.ConfigAuthSessionSameSite `env:"AUTH_SESSION_COOKIE_SAME_SITE,default=strict" validate:"required_if=Type session,oneof=lax strict"`
	SessionCookieSameSiteMode http.SameSite
}

type Task struct {
	EventCleanStaleEnabled  bool          `env:"TASK_EVENT_CLEAN_STALE_ENABLED,default=false"`
	EventCleanStaleInterval time.Duration `env:"TASK_EVENT_CLEAN_STALE_INTERVAL,default=8h"   validate:"required_if=EventCleanStaleEnabled true,gt=0"`
	EventCleanStaleMaxAge   time.Duration `env:"TASK_EVENT_CLEAN_STALE_MAX_AGE,default=2190h" validate:"required_if=EventCleanStaleEnabled true,gt=0"`

	ActionsEnqueueEnabled   bool          `env:"TASK_ACTIONS_ENQUEUE_ENABLED,default=true"`
	ActionsEnqueueInterval  time.Duration `env:"TASK_ACTIONS_ENQUEUE_INTERVAL,default=10s" validate:"required_if=ActionsEnqueueEnabled true,gt=0"`
	ActionsEnqueueBatchSize int           `env:"TASK_ACTIONS_ENQUEUE_BATCH_SIZE,default=1" validate:"required_if=ActionsEnqueueEnabled true,numeric,gte=1"`

	ActionsInvokeEnabled    bool          `env:"TASK_ACTIONS_INVOKE_ENABLED,default=true"`
	ActionsInvokeInterval   time.Duration `env:"TASK_ACTIONS_INVOKE_INTERVAL,default=10s"  validate:"required_if=ActionsInvokeEnabled true,gt=0"`
	ActionsInvokeBatchSize  int           `env:"TASK_ACTIONS_INVOKE_BATCH_SIZE,default=1"  validate:"required_if=ActionsInvokeEnabled true,numeric,gte=1"`
	ActionsInvokeMaxRetries int           `env:"TASK_ACTIONS_INVOKE_MAX_RETRIES,default=3" validate:"required_if=ActionsInvokeEnabled true,numeric,gte=1"`

	ActionsCleanStaleEnabled  bool          `env:"TASK_ACTIONS_CLEAN_STALE_ENABLED,default=true"`
	ActionsCleanStaleInterval time.Duration `env:"TASK_ACTIONS_CLEAN_STALE_INTERVAL,default=12h" validate:"required_if=ActionsCleanStaleEnabled true,gt=0"`
	ActionsCleanStaleMaxAge   time.Duration `env:"TASK_ACTIONS_CLEAN_STALE_MAX_AGE,default=720h" validate:"required_if=ActionsCleanStaleEnabled true,gt=0"`
}

type Lock struct {
	RedisEnabled        bool          `env:"LOCK_REDIS_ENABLED,default=false"`
	RedisHost           string        `env:"LOCK_REDIS_HOST,default=localhost"        validate:"required_if=RedisEnabled true"`
	RedisPort           int           `env:"LOCK_REDIS_PORT,default=6379"             validate:"required_if=RedisEnabled true,numeric,gte=1"`
	RedisDbName         int           `env:"LOCK_REDIS_DB_NAME,default=0"             validate:"numeric,gte=0"`
	RedisUsername       string        `env:"LOCK_REDIS_USERNAME"`
	RedisPassword       string        `env:"LOCK_REDIS_PASSWORD"`
	RedisTaskTries      int           `env:"LOCK_REDIS_TASK_LOCK_TRIES,default=1"     validate:"required_if=RedisEnabled true,numeric,gte=1"`
	RedisTaskLockAtMost time.Duration `env:"LOCK_REDIS_TASK_LOCK_AT_MOST,default=30s" validate:"required_if=RedisEnabled true,gte=0"`
	RedisTaskRetryDelay time.Duration `env:"LOCK_REDIS_TASK_RETRY_DELAY,default=5s"   validate:"required_if=RedisEnabled true,gte=0"`
	RedisUrl            string
}

type Webhook struct {
	TokenLength int `env:"WEBHOOKS_TOKEN_LENGTH,default=32" validate:"required,numeric,gte=4"`
}

type Prometheus struct {
	Enabled            bool          `env:"PROMETHEUS_ENABLED,default=false"`
	Port               int           `env:"PROMETHEUS_PORT,default=8080"                 validate:"required_if=Enabled true,gte=1"`
	Listen             string        `env:"PROMETHEUS_LISTEN"`
	BasePath           string        `env:"PROMETHEUS_BASE_PATH,default=/"               validate:"required_if=Enabled true"`
	Path               string        `env:"PROMETHEUS_METRICS_PATH,default=/metrics"     validate:"required_if=Enabled true"`
	SecureTokenEnabled bool          `env:"PROMETHEUS_SECURE_TOKEN_ENABLED,default=true"`
	SecureToken        string        `env:"PROMETHEUS_SECURE_TOKEN"                      validate:"required_if=Enabled true SecureTokenEnabled true"`
	RefreshInterval    time.Duration `env:"PROMETHEUS_REFRESH_INTERVAL,default=60s"      validate:"required_if=Enabled true,gte=0"`
}

type Configuration struct {
	App                      *App
	Auth                     *Auth
	Cors                     *Cors
	Database                 *Database
	Lock                     *Lock
	Logging                  *Logging
	Prometheus               *Prometheus
	Secret                   *Secret
	Server                   *Server
	Task                     *Task
	Webhook                  *Webhook
	Webinterface             *Webinterface
	WebinterfaceCacheControl *WebinterfaceCacheControl
}

func LoadFromEnvironment(ctx context.Context) (*Configuration, *gorm.DB) {
	var err error

	// bootstrap logging (configured independently and required before any other action)
	var lc Logging
	if err = envconfig.Process(ctx, &lc); err != nil {
		golog.Fatalf("Cannot load logging configuration from environment. Reason: %v", err)
	}
	if err = validate.ValidOrError(lc); err != nil {
		golog.Fatalf("Cannot validate logging configuration. Reason: %s", err)
	}

	configureLogger(&lc)

	// load configuration and validate from environment
	var c Configuration
	if err = envconfig.Process(ctx, &c); err != nil {
		log.Fatal().Msgf("Cannot load configuration from environment. Reason: %v", err)
	}
	if err = validate.ValidOrError(c); err != nil {
		log.Fatal().Msgf("Cannot validate configuration. Reason: %s", err.Error())
	}

	var db *gorm.DB
	var migrationDriver database.Driver
	var migrationDatabaseName string
	var migrationFS source.Driver

	log.Info().Msgf("Using database type '%s'", c.Database.Type)

	if constant.ConfigDatabaseTypePostgres == c.Database.Type {
		host := c.Database.PostgresHost
		port := c.Database.PostgresPort
		dbUser := c.Database.PostgresUser
		dbPass := c.Database.PostgresPassword
		dbName := c.Database.PostgresName
		dbTZ := c.Database.PostgresTimeZone
		migrationDatabaseName = dbName

		dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=%v", host, dbUser, dbPass, dbName, port, dbTZ)

		gormLog := zerologgorm.NewLogger(
			zerologgorm.WithDefaultLogLevel(zerolog.DebugLevel),
			zerologgorm.WithSlowThreshold(500*time.Millisecond),
			zerologgorm.WithLogParams(),
			zerologgorm.WithIgnoreNotFoundError(),
		)

		if db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gormLog}); err != nil {
			log.Fatal().Msgf("Could not setup database: %v", err)
		}

		var sqlDb *sql.DB
		if sqlDb, err = db.DB(); err != nil {
			log.Fatal().Msgf("Could not retrieve database: %v", err)
		}

		if err = sqlDb.Ping(); err != nil {
			log.Fatal().Msgf("Could not connect to database: %v", err)
		}

		if migrationDriver, err = migratepostgres.WithInstance(sqlDb, &migratepostgres.Config{}); err != nil {
			log.Fatal().Msgf("Could not create migration driver: %v", err)
		}

		if migrationFS, err = iofs.New(migrationPostgresFS, "migrations_postgres"); err != nil {
			log.Fatal().Msgf("Could not create migration source: %v", err)
		}
	}

	if db == nil {
		log.Fatal().Msgf("Could not setup database")
	}

	if !c.Database.MigrationEnabled {
		log.Warn().Msg("Database schema migration is disabled and not executed automatically. Make sure to run them manually, otherwise the application might misbehave. You can safely ignore this warning if application is started in high availability mode and you're sure necessary database schema already exists.")
	} else {
		var migrator *migrate.Migrate
		if migrator, err = migrate.NewWithInstance("iofs", migrationFS, migrationDatabaseName, migrationDriver); err != nil {
			log.Fatal().Msgf("Could not create database migration instance: %v", err)
		}

		var migrationVersion uint
		var migrationVersionDirty bool
		if migrationVersion, migrationVersionDirty, err = migrator.Version(); err != nil {
			if errors.Is(err, migrate.ErrNilVersion) {
				log.Info().Msg("Database migration schema is uninitialized")
			} else {
				log.Fatal().Msgf("Could not retrieve database migration version: %v", err)
			}
		} else {
			log.Info().Msgf("Previous database migration version is '%d' (dirty '%v')", migrationVersion, migrationVersionDirty)
		}

		log.Info().Msg("Applying necessary database migration steps...")
		if err = migrator.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				log.Info().Msg("No database schema changes detected")
			} else {
				log.Fatal().Msgf("Could not migrate database schema: %v", err)
			}
		}

		log.Info().Msg("Applied all necessary database migration steps successfully")
	}

	// custom defaults and validation
	if c.Lock.RedisEnabled {
		if c.Lock.RedisUsername != "" && c.Lock.RedisPassword != "" {
			c.Lock.RedisUrl = fmt.Sprintf("redis://%s:%s@%s:%d/%d", c.Lock.RedisUsername, c.Lock.RedisPassword, c.Lock.RedisHost, c.Lock.RedisPort, c.Lock.RedisDbName)
		} else {
			c.Lock.RedisUrl = fmt.Sprintf("redis://%s:%d/%d", c.Lock.RedisHost, c.Lock.RedisPort, c.Lock.RedisDbName)
		}
	}

	c.Auth.SessionCookieSameSiteMode = convertSessionCookieSameSite(c.Auth.SessionCookieSameSite)

	log.Info().Msgf("Configuration: App %+v", c.App)
	log.Info().Msg("Configuration: Auth ***REDACTED***")
	log.Info().Msgf("Configuration: Cors %+v", c.Cors)
	log.Info().Msg("Configuration: Database ***REDACTED***")
	log.Info().Msg("Configuration: Lock ***REDACTED***")
	log.Info().Msgf("Configuration: Logging %+v", lc)
	log.Info().Msg("Configuration: Prometheus ***REDACTED***")
	log.Info().Msg("Configuration: Secret ***REDACTED***")
	log.Info().Msgf("Configuration: Server %+v", c.Server)
	log.Info().Msgf("Configuration: Task %+v", c.Task)
	log.Info().Msgf("Configuration: Webhook %+v", c.Webhook)
	log.Info().Msgf("Configuration: Webinterface %+v", c.Webinterface)
	log.Info().Msgf("Configuration: WebinterfaceCacheControl %+v", c.WebinterfaceCacheControl)

	return &c, db
}

func configureLogger(cfg *Logging) {
	var level zerolog.Level
	var err error
	if level, err = zerolog.ParseLevel(cfg.Level); err != nil {
		golog.Fatalf("Cannot parse logging level: %v", err)
	}
	zerolog.SetGlobalLevel(level)

	zerolog.CallerFieldName = cfg.EncodingFuncKey
	zerolog.ErrorFieldName = cfg.EncodingErrorKey
	zerolog.ErrorStackFieldName = cfg.EncodingStacktraceKey
	zerolog.LevelFieldName = cfg.EncodingLevelKey
	zerolog.MessageFieldName = cfg.EncodingMessageKey
	zerolog.TimestampFieldName = cfg.EncodingTimeKey

	var timeEncoders = map[constant.ConfigLogTimeEncoder]string{
		constant.ConfigLogTimeEncoderEpoch:       zerolog.TimeFormatUnix,
		constant.ConfigLogTimeEncoderEpochmillis: zerolog.TimeFormatUnixMs,
		constant.ConfigLogTimeEncoderEpochnanos:  zerolog.TimeFormatUnixNano,
		constant.ConfigLogTimeEncoderIso8601:     "2006-01-02T15:04:05-0700",
		constant.ConfigLogTimeEncoderRfc3339:     time.RFC3339,
		constant.ConfigLogTimeEncoderRfc3339nano: time.RFC3339Nano,
	}
	if enc, ok := timeEncoders[cfg.EncodingTimeEncoder]; ok {
		zerolog.TimeFieldFormat = enc
	}

	if constant.ConfigLogEncodingJson == cfg.Encoding {
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
	} else {
		log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: zerolog.TimeFieldFormat, NoColor: !cfg.EncodingColorize}).With().Timestamp().Caller().Logger()
	}
}

func convertSessionCookieSameSite(site constant.ConfigAuthSessionSameSite) http.SameSite {
	switch site {
	case constant.ConfigAuthSessionSameSiteLax:
		return http.SameSiteLaxMode
	case constant.ConfigAuthSessionSameSiteStrict:
		return http.SameSiteStrictMode
	default:
		return http.SameSiteDefaultMode
	}
}
