package server

const (
	envDevelopment = "DEVELOPMENT"

	envLoggingLevel        = "LOGGING_LEVEL"
	loggingLevelDefault    = "info"
	envLoggingEncoding     = "LOGGING_ENCODING"
	loggingEncodingDefault = "json"

	envLoggingDirectory    = "LOGGING_DIRECTORY"
	loggingFileNameDefault = "upda.log"

	envSecret = "SECRET"

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
	authModeDefault          = authModeBasicSingle
	authModeBasicSingle      = "basic_single"
	authModeBasicCredentials = "basic_credentials"
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
