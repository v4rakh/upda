package server

const (
	envTZ     = "TZ"
	tzDefault = "Europe/Berlin"

	envAdminUser     = "ADMIN_USER"
	envAdminPassword = "ADMIN_PASSWORD"

	envLoggingLevel = "LOGGING_LEVEL"

	envServerPort           = "SERVER_PORT"
	envServerListen         = "SERVER_LISTEN"
	envServerTlsEnabled     = "SERVER_TLS_ENABLED"
	envServerTlsCertPath    = "SERVER_TLS_CERT_PATH"
	envServerTlsKeyPath     = "SERVER_TLS_KEY_PATH"
	envServerTimeout        = "SERVER_TIMEOUT"
	serverListenDefault     = ""
	serverPortDefault       = "8080"
	serverTlsEnabledDefault = "false"
	serverTimeoutDefault    = "1s"

	envCorsAllowOrigin      = "CORS_ALLOW_ORIGIN"
	envCorsAllowMethods     = "CORS_ALLOW_METHODS"
	envCorsAllowHeaders     = "CORS_ALLOW_HEADERS"
	corsAllowOriginDefault  = "*"
	corsAllowMethodsDefault = "HEAD, GET, POST, PUT, PATCH, DELETE, OPTIONS"
	corsAllowHeadersDefault = "Authorization, Content-Type"

	dbTypeSqlite   = "sqlite"
	dbTypePostgres = "postgres"

	envDbType                 = "DB_TYPE"
	envDbSqliteFile           = "DB_SQLITE_FILE"
	envDbPostgresHost         = "DB_POSTGRES_HOST"
	envDbPostgresPort         = "DB_POSTGRES_PORT"
	envDbPostgresName         = "DB_POSTGRES_NAME"
	envDbPostgresTimeZone     = "DB_POSTGRES_TZ"
	envDbPostgresUser         = "DB_POSTGRES_USER"
	envDbPostgresPassword     = "DB_POSTGRES_PASSWORD"
	dbTypeSqliteDbNameDefault = "upda.db"
	dbTypePostgresHostDefault = "localhost"
	dbTypePostgresPortDefault = "5432"
	dbTypePostgresTZDefault   = "Europe/Berlin"

	envTaskPrometheusRefreshInterval = "TASK_PROMETHEUS_REFRESH_INTERVAL"
	taskPrometheusRefreshDefault     = "60s"

	envWebhooksTokenLength     = "WEBHOOKS_TOKEN_LENGTH"
	webhooksTokenLengthDefault = "16"

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
	taskUpdateCleanStaleMaxAgeDefault   = "168h"

	envTaskEventCleanStaleEnabled      = "TASK_EVENT_CLEAN_STALE_ENABLED"
	envTaskEventCleanStaleInterval     = "TASK_EVENT_CLEAN_STALE_INTERVAL"
	envTaskEventCleanStaleMaxAge       = "TASK_EVENT_CLEAN_STALE_MAX_AGE"
	taskEventCleanStaleEnabledDefault  = "false"
	taskEventCleanStaleIntervalDefault = "8h"
	taskEventCleanStaleMaxAgeDefault   = "2190h"

	envLockRedisEnabled = "LOCK_REDIS_ENABLED"
	envLockRedisUrl     = "LOCK_REDIS_URL"
	redisEnabledDefault = "false"
)
