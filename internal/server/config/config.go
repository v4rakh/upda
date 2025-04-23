package config

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"git.myservermanager.com/varakh/upda/internal/file"
	"git.myservermanager.com/varakh/upda/internal/server/constant"
	"git.myservermanager.com/varakh/upda/internal/validate"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	migratepostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"moul.io/zapgorm2"
	"path/filepath"
	"time"
)

const (
	EnvSecret = "SECRET"
)

//go:embed migrations_postgres/*.sql
var migrationPostgresFS embed.FS

type Logging struct {
	Level       string `env:"LOGGING_LEVEL,default=info" validate:"required,oneof=debug info warn error dpanic panic fatal"`
	Encoding    string `env:"LOGGING_ENCODING,default=json" validate:"required,oneof=json console"`
	Directory   string `env:"LOGGING_DIRECTORY"`
	Debug       bool   `env:"DEBUG,default=false"`
	Development bool   `env:"DEVELOPMENT,default=false"`
}

type App struct {
	TimeZone    string `env:"TZ,default=Etc/UTC" validate:"required"`
	Development bool   `env:"DEVELOPMENT,default=false"`
}

type Secret struct {
	Secret string `env:"SECRET,required" validate:"required"`
}

type Server struct {
	Port                 int           `env:"SERVER_PORT,default=8080" validate:"gte=1"`
	Listen               string        `env:"SERVER_LISTEN"`
	BasePath             string        `env:"SERVER_BASE_PATH,default=/" validate:"required"`
	TlsEnabled           bool          `env:"SERVER_TLS_ENABLED,default=false"`
	TlsCertPath          string        `env:"SERVER_TLS_CERT_PATH"`
	TlsKeyPath           string        `env:"SERVER_TLS_KEY_PATH"`
	Timeout              time.Duration `env:"SERVER_TIMEOUT,default=1s" validate:"gte=0"`
	CorsAllowCredentials bool          `env:"CORS_ALLOW_CREDENTIALS,default=true"`
	CorsAllowOrigins     []string      `env:"CORS_ALLOW_ORIGINS,default=*"`
	CorsAllowMethods     []string      `env:"CORS_ALLOW_METHODS,default=HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS"`
	CorsAllowHeaders     []string      `env:"CORS_ALLOW_HEADERS,default=Authorization,Content-Type"`
	CorsExposeHeaders    []string      `env:"CORS_EXPOSE_HEADERS,default=*"`
}

type Database struct {
	Type             string `env:"DB_TYPE,default=postgres" validate:"required,oneof=postgres"`
	MigrationEnabled bool   `env:"DB_MIGRATION_ENABLED,default=true"`
	PostgresHost     string `env:"DB_POSTGRES_HOST,default=localhost" validate:"required_if=Type postgres"`
	PostgresPort     int    `env:"DB_POSTGRES_PORT,default=5432" validate:"required_if=Type postgres"`
	PostgresName     string `env:"DB_POSTGRES_NAME" validate:"required_if=Type postgres"`
	PostgresTimeZone string `env:"DB_POSTGRES_TZ,default=Etc/UTC" validate:"required_if=Type postgres"`
	PostgresUser     string `env:"DB_POSTGRES_USER" validate:"required_if=Type postgres"`
	PostgresPassword string `env:"DB_POSTGRES_PASSWORD" validate:"required_if=Type postgres"`
}

type EmbeddedWebInterface struct {
	Enabled          bool   `env:"EMBEDDED_WEB_INTERFACE_ENABLED,default=true"`
	ApiUrl           string `env:"EMBEDDED_WEB_INTERFACE_API_URL,default=http://localhost:8080" validate:"required_if=Enabled true"`
	Title            string `env:"EMBEDDED_WEB_INTERFACE_TITLE,default=upda" validate:"required_if=Enabled true"`
	DarkThemeEnabled bool   `env:"EMBEDDED_WEB_INTERFACE_DARK_THEME_ENABLED,default=false"`
	FooterEnabled    bool   `env:"EMBEDDED_WEB_INTERFACE_FOOTER_ENABLED,default=true"`
}

type Auth struct {
	AuthMethod           constant.ConfigAuthMode `env:"AUTH_MODE,default=basic_single" validate:"required,oneof=basic_single basic_credentials"`
	BasicAuthUser        string                  `env:"BASIC_AUTH_USER" validate:"required_if=AuthMethod basic_single"`
	BasicAuthPassword    string                  `env:"BASIC_AUTH_PASSWORD" validate:"required_if=AuthMethod basic_single"`
	BasicAuthCredentials map[string]string       `env:"BASIC_AUTH_CREDENTIALS,separator=|,delimiter=;" validate:"required_if=AuthMethod basic_credentials"`
}

type Task struct {
	UpdateCleanStaleEnabled  bool          `env:"TASK_UPDATE_CLEAN_STALE_ENABLED,default=true"`
	UpdateCleanStaleInterval time.Duration `env:"TASK_EVENT_CLEAN_STALE_INTERVAL,default=1h" validate:"required_if=UpdateCleanStaleEnabled true,gt=0"`
	UpdateCleanStaleMaxAge   time.Duration `env:"TASK_UPDATE_CLEAN_STALE_MAX_AGE,default=720h" validate:"required_if=UpdateCleanStaleEnabled true,gt=0"`

	EventCleanStaleEnabled  bool          `env:"TASK_EVENT_CLEAN_STALE_ENABLED,default=false"`
	EventCleanStaleInterval time.Duration `env:"TASK_EVENT_CLEAN_STALE_INTERVAL,default=8h" validate:"required_if=EventCleanStaleEnabled true,gt=0"`
	EventCleanStaleMaxAge   time.Duration `env:"TASK_EVENT_CLEAN_STALE_MAX_AGE,default=2190h" validate:"required_if=EventCleanStaleEnabled true,gt=0"`

	ActionsEnqueueEnabled   bool          `env:"TASK_ACTIONS_ENQUEUE_ENABLED,default=true"`
	ActionsEnqueueInterval  time.Duration `env:"TASK_ACTIONS_ENQUEUE_INTERVAL,default=10s" validate:"required_if=ActionsEnqueueEnabled true,gt=0"`
	ActionsEnqueueBatchSize int           `env:"TASK_ACTIONS_ENQUEUE_BATCH_SIZE,default=1" validate:"required_if=ActionsEnqueueEnabled true,numeric,gte=1"`

	ActionsInvokeEnabled    bool          `env:"TASK_ACTIONS_INVOKE_ENABLED,default=true"`
	ActionsInvokeInterval   time.Duration `env:"TASK_ACTIONS_INVOKE_INTERVAL,default=10s" validate:"required_if=ActionsInvokeEnabled true,gt=0"`
	ActionsInvokeBatchSize  int           `env:"TASK_ACTIONS_INVOKE_BATCH_SIZE,default=1" validate:"required_if=ActionsInvokeEnabled true,numeric,gte=1"`
	ActionsInvokeMaxRetries int           `env:"TASK_ACTIONS_INVOKE_MAX_RETRIES,default=3" validate:"required_if=ActionsInvokeEnabled true,numeric,gte=1"`

	ActionsCleanStaleEnabled  bool          `env:"TASK_ACTIONS_CLEAN_STALE_ENABLED,default=true"`
	ActionsCleanStaleInterval time.Duration `env:"TASK_ACTIONS_CLEAN_STALE_INTERVAL,default=12h" validate:"required_if=ActionsCleanStaleEnabled true,gt=0"`
	ActionsCleanStaleMaxAge   time.Duration `env:"TASK_ACTIONS_CLEAN_STALE_MAX_AGE,default=720h" validate:"required_if=ActionsCleanStaleEnabled true,gt=0"`

	PrometheusRefreshInterval time.Duration `env:"TASK_PROMETHEUS_REFRESH_INTERVAL,default=60s" validate:"required,gte=0"`
}

type Lock struct {
	RedisEnabled bool   `env:"LOCK_REDIS_ENABLED,default=false"`
	RedisUrl     string `env:"LOCK_REDIS_URL" validate:"required_if=RedisEnabled true"`
}

type Webhook struct {
	TokenLength int `env:"WEBHOOKS_TOKEN_LENGTH,default=32" validate:"required,numeric,gte=4"`
}

type Prometheus struct {
	Enabled            bool   `env:"PROMETHEUS_ENABLED,default=false"`
	Path               string `env:"PROMETHEUS_METRICS_PATH,default=/metrics" validate:"required_if=Enabled true"`
	SecureTokenEnabled bool   `env:"PROMETHEUS_SECURE_TOKEN_ENABLED,default=true"`
	SecureToken        string `env:"PROMETHEUS_SECURE_TOKEN" validate:"required_if=Enabled true SecureTokenEnabled true"`
}

type Configuration struct {
	App                  *App
	Auth                 *Auth
	Database             *Database
	EmbeddedWebInterface *EmbeddedWebInterface
	Lock                 *Lock
	Prometheus           *Prometheus
	Secret               *Secret
	Server               *Server
	Task                 *Task
	Webhook              *Webhook
}

func LoadFromEnvironment(ctx context.Context) (*Configuration, *gorm.DB) {
	var err error

	// bootstrap logging (configured independently and required before any other action)
	var lc Logging
	if err = envconfig.Process(ctx, &lc); err != nil {
		zap.L().Sugar().Fatalf("Cannot load logging configuration from environment. Reason: %v", err)
	}
	if err = validate.ValidOrError(lc); err != nil {
		zap.L().Sugar().Fatalf("Cannot validate logging configuration. Reason: %s", err.Error())
	}

	var level zap.AtomicLevel
	if level, err = zap.ParseAtomicLevel(lc.Level); err != nil {
		log.Fatalf("Cannot parse logging level: %v", err)
	}

	var loggingEncoderConfig zapcore.EncoderConfig
	if "json" == lc.Encoding {
		loggingEncoderConfig = zap.NewProductionEncoderConfig()
	} else {
		loggingEncoderConfig = zap.NewDevelopmentEncoderConfig()
	}

	logPaths := []string{"stderr"}
	if lc.Directory != "" {
		logFile := filepath.Join(lc.Directory, fmt.Sprintf("%s.log", constant.AppName))

		if err = file.CreateFileWithParent(logFile); err != nil {
			log.Fatalf("Log file '%s' cannot be created: %v", lc.Directory, err)
		}

		logPaths = append(logPaths, logFile)
	}

	var zapConfig *zap.Config
	if lc.Debug {
		zapConfig = &zap.Config{
			Level:            level,
			Development:      lc.Development,
			Encoding:         lc.Encoding,
			EncoderConfig:    loggingEncoderConfig,
			OutputPaths:      logPaths,
			ErrorOutputPaths: logPaths,
		}
	} else {
		zapConfig = &zap.Config{
			Level:       level,
			Development: lc.Development,
			Sampling: &zap.SamplingConfig{
				Initial:    100,
				Thereafter: 100,
			},
			Encoding:         lc.Encoding,
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

	// load configuration and validate from environment
	var c Configuration
	if err = envconfig.Process(ctx, &c); err != nil {
		zap.L().Sugar().Fatalf("Cannot load configuration from environment. Reason: %v", err)
	}
	if err = validate.ValidOrError(c); err != nil {
		zap.L().Sugar().Fatalf("Cannot validate configuration. Reason: %s", err.Error())
	}

	gormConfig := &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)}
	if lc.Debug && c.App.Development {
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

	zap.L().Sugar().Infof("Using database type '%s'", c.Database.Type)

	if "postgres" == c.Database.Type {
		host := c.Database.PostgresHost
		port := c.Database.PostgresPort
		dbUser := c.Database.PostgresUser
		dbPass := c.Database.PostgresPassword
		dbName := c.Database.PostgresName
		dbTZ := c.Database.PostgresTimeZone
		migrationDatabaseName = dbName

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
	}

	if db == nil {
		zap.L().Sugar().Fatalf("Could not setup database")
	}

	if !c.Database.MigrationEnabled {
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

	zap.L().Sugar().Infof("Configuration: App %+v", c.App)
	zap.L().Sugar().Infof("Configuration: Auth ***REDACTED***")
	zap.L().Sugar().Infof("Configuration: Database ***REDACTED***")
	zap.L().Sugar().Infof("Configuration: EmbeddedWebInterface %+v", c.EmbeddedWebInterface)
	zap.L().Sugar().Infof("Configuration: Lock ***REDACTED***")
	zap.L().Sugar().Infof("Configuration: Logging %+v", lc)
	zap.L().Sugar().Infof("Configuration: Prometheus ***REDACTED***")
	zap.L().Sugar().Infof("Configuration: Secret ***REDACTED***")
	zap.L().Sugar().Infof("Configuration: Server %+v", c.Server)
	zap.L().Sugar().Infof("Configuration: Task %+v", c.Task)
	zap.L().Sugar().Infof("Configuration: Webhook %+v", c.Webhook)

	return &c, db
}
