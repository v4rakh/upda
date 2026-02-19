# Configuration

The following tables describe available configuration values.

### Logging

| Env Var                         | Default    | Description                                                                  |
|---------------------------------|------------|------------------------------------------------------------------------------|
| LOGGING_ENCODING                | console    | Sets log output format (console or json).                                    |
| LOGGING_ENCODING_COLORIZE       | false      | Enables colored log output if supported.                                     |
| LOGGING_ENCODING_ERROR_KEY      | error      | Field name for error messages in logs.                                       |
| LOGGING_ENCODING_FILE_KEY       | file       | Field name for file names in logs.                                           |
| LOGGING_ENCODING_FUNC_KEY       | func       | Field name for functions in logs.                                            |
| LOGGING_ENCODING_LEVEL_KEY      | level      | Field name for log levels.                                                   |
| LOGGING_ENCODING_MESSAGE_KEY    | msg        | Field name for log messages.                                                 |
| LOGGING_ENCODING_STACKTRACE_KEY | stacktrace | Field name for stack traces.                                                 |
| LOGGING_ENCODING_TIME_ENCODER   | rfc3339    | Time format for logs (e.g. rfc3339, iso8601, epoch).                         |
| LOGGING_ENCODING_TIME_KEY       | ts         | Field name for timestamps in logs.                                           |
| LOGGING_LEVEL                   | info       | Minimum log level (trace, debug, info, warn, error, fatal, panic, disabled). |
| LOGGING_LEVEL_REQUESTS          | disabled   | Logging level for HTTP requests.                                             |

### App

| Env Var     | Default | Description                                                 |
|-------------|---------|-------------------------------------------------------------|
| TZ          | Etc/UTC | System time zone.                                           |
| DEVELOPMENT | false   | Enables development mode (more verbose logging, debugging). |

### Secret

| Env Var | Default | Description                              |
|---------|---------|------------------------------------------|
| SECRET  | -       | Secret key for cryptographic operations. |

### Server

| Env Var              | Default | Description                             |
|----------------------|---------|-----------------------------------------|
| SERVER_PORT          | 8080    | Server listening port.                  |
| SERVER_LISTEN        | -       | Network address to bind (e.g. 0.0.0.0). |
| SERVER_BASE_PATH     | /       | Base path for API endpoints.            |
| SERVER_TLS_ENABLED   | false   | Enables TLS/HTTPS.                      |
| SERVER_TLS_CERT_PATH | -       | Path to TLS certificate file.           |
| SERVER_TLS_KEY_PATH  | -       | Path to TLS key file.                   |
| SERVER_TIMEOUT       | 10s     | Request timeout duration.               |

### Cors

| Env Var                | Default                                | Description                              |
|------------------------|----------------------------------------|------------------------------------------|
| CORS_ALLOW_CREDENTIALS | true                                   | Whether cookies/credentials are allowed. |
| CORS_ALLOW_ORIGINS     | *                                      | Allowed origins (comma-separated).       |
| CORS_ALLOW_METHODS     | HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS | Allowed HTTP methods.                    |
| CORS_ALLOW_HEADERS     | Authorization,Content-Type             | Allowed request headers.                 |
| CORS_EXPOSE_HEADERS    | *                                      | Response headers to expose.              |

### Database

| Env Var              | Default   | Description                              |
|----------------------|-----------|------------------------------------------|
| DB_TYPE              | postgres  | Database type (only postgres supported). |
| DB_MIGRATION_ENABLED | true      | Whether DB migrations run on startup.    |
| DB_POSTGRES_HOST     | localhost | Postgres host address.                   |
| DB_POSTGRES_PORT     | 5432      | Postgres port.                           |
| DB_POSTGRES_NAME     | -         | Postgres database name.                  |
| DB_POSTGRES_TZ       | Etc/UTC   | Postgres time zone.                      |
| DB_POSTGRES_USER     | -         | Postgres username.                       |
| DB_POSTGRES_PASSWORD | -         | Postgres password.                       |

### Webinterface

| Env Var                          | Default                    | Description                 |
|----------------------------------|----------------------------|-----------------------------|
| WEB_INTERFACE_ENABLED            | true                       | Enables web UI.             |
| WEB_INTERFACE_API_URL            | http://localhost:8080/api/ | API URL used by the web UI. |
| WEB_INTERFACE_TITLE              | upda                       | Web UI title.               |
| WEB_INTERFACE_DARK_THEME_ENABLED | false                      | Enables dark theme.         |
| WEB_INTERFACE_FOOTER_ENABLED     | true                       | Shows or hides the footer.  |

### Webinterface Cache Control

| Env Var                                 | Default | Description                             |
|-----------------------------------------|---------|-----------------------------------------|
| WEB_INTERFACE_CC_ENABLED                | true    | Enables cache control headers.          |
| WEB_INTERFACE_CC_MUST_REVALIDATE        | true    | Requires cache to revalidate.           |
| WEB_INTERFACE_CC_NO_CACHE               | false   | Disables caching of responses.          |
| WEB_INTERFACE_CC_NO_STORE               | false   | Prevents storing responses.             |
| WEB_INTERFACE_CC_NO_TRANSFORM           | false   | Disables proxy modifications.           |
| WEB_INTERFACE_CC_PUBLIC                 | true    | Marks responses as public cacheable.    |
| WEB_INTERFACE_CC_PRIVATE                | false   | Marks responses as private cacheable.   |
| WEB_INTERFACE_CC_PROXY_REVALIDATE       | true    | Proxies must revalidate cache on reuse. |
| WEB_INTERFACE_CC_MAX_AGE                | 48h     | Maximum cache duration.                 |
| WEB_INTERFACE_CC_SMAX_AGE               | -       | Shared cache duration.                  |
| WEB_INTERFACE_CC_IMMUTABLE              | false   | Marks responses as immutable.           |
| WEB_INTERFACE_CC_STALE_WHILE_REVALIDATE | -       | Allows stale cache while revalidating.  |
| WEB_INTERFACE_CC_STALE_IF_ERROR         | -       | Allows stale cache if server errors.    |

### Auth

| Env Var                       | Default      | Description                                    |
|-------------------------------|--------------|------------------------------------------------|
| AUTH_TYPE                     | session      | Authentication type (only session supported).  |
| AUTH_SESSION_SECRET           | -            | Session secret key.                            |
| AUTH_SESSION_PROVIDER         | single       | Session provider (single user or credentials). |
| AUTH_SESSION_USER             | -            | Username for single provider.                  |
| AUTH_SESSION_PASSWORD         | -            | Password for single provider.                  |
| AUTH_SESSION_CREDENTIALS      | -            | Multiple credentials (key-value format).       |
| AUTH_SESSION_CLEANUP_ENABLED  | true         | Enables cleanup of expired sessions.           |
| AUTH_SESSION_CLEANUP_INTERVAL | 1h           | Interval for session cleanup.                  |
| AUTH_SESSION_COOKIE_MAX_AGE   | 8h           | Session cookie expiration time.                |
| AUTH_SESSION_COOKIE_NAME      | UPDA_SESSION | Session cookie name.                           |
| AUTH_SESSION_COOKIE_DOMAIN    | -            | Domain for session cookie.                     |
| AUTH_SESSION_COOKIE_PATH      | /            | Path scope for session cookie.                 |
| AUTH_SESSION_COOKIE_HTTP_ONLY | true         | Disallows JavaScript access to cookie.         |
| AUTH_SESSION_COOKIE_SECURE    | true         | Requires HTTPS for cookies.                    |
| AUTH_SESSION_COOKIE_SAME_SITE | strict       | SameSite policy (lax or strict).               |

### Task

| Env Var                           | Default | Description                          |
|-----------------------------------|---------|--------------------------------------|
| TASK_UPDATE_CLEAN_STALE_ENABLED   | false   | Enables cleaning of stale updates.   |
| TASK_EVENT_CLEAN_STALE_INTERVAL   | 1h      | Interval for cleaning stale updates. |
| TASK_UPDATE_CLEAN_STALE_MAX_AGE   | 720h    | Max age for stale updates.           |
| TASK_EVENT_CLEAN_STALE_ENABLED    | false   | Enables cleaning of stale events.    |
| TASK_EVENT_CLEAN_STALE_INTERVAL   | 8h      | Interval for cleaning stale events.  |
| TASK_EVENT_CLEAN_STALE_MAX_AGE    | 2190h   | Max age for stale events.            |
| TASK_ACTIONS_ENQUEUE_ENABLED      | true    | Enables action enqueuing.            |
| TASK_ACTIONS_ENQUEUE_INTERVAL     | 10s     | Interval for enqueueing actions.     |
| TASK_ACTIONS_ENQUEUE_BATCH_SIZE   | 1       | Batch size for action enqueueing.    |
| TASK_ACTIONS_INVOKE_ENABLED       | true    | Enables invoking actions.            |
| TASK_ACTIONS_INVOKE_INTERVAL      | 10s     | Interval for invoking actions.       |
| TASK_ACTIONS_INVOKE_BATCH_SIZE    | 1       | Batch size for action invocation.    |
| TASK_ACTIONS_INVOKE_MAX_RETRIES   | 3       | Max retries for failed actions.      |
| TASK_ACTIONS_CLEAN_STALE_ENABLED  | true    | Enables cleanup of stale actions.    |
| TASK_ACTIONS_CLEAN_STALE_INTERVAL | 12h     | Interval for cleaning stale actions. |
| TASK_ACTIONS_CLEAN_STALE_MAX_AGE  | 720h    | Max age for stale actions.           |

### Lock

| Env Var                      | Default   | Description                         |
|------------------------------|-----------|-------------------------------------|
| LOCK_REDIS_ENABLED           | false     | Enables Redis-based locks.          |
| LOCK_REDIS_HOST              | localhost | Redis host.                         |
| LOCK_REDIS_PORT              | 6379      | Redis port.                         |
| LOCK_REDIS_DB_NAME           | 0         | Redis database index.               |
| LOCK_REDIS_USERNAME          | -         | Redis username.                     |
| LOCK_REDIS_PASSWORD          | -         | Redis password.                     |
| LOCK_REDIS_TASK_LOCK_TRIES   | 1         | Number of retry attempts for locks. |
| LOCK_REDIS_TASK_LOCK_AT_MOST | 30s       | Max duration a task lock is held.   |
| LOCK_REDIS_TASK_RETRY_DELAY  | 5s        | Delay between lock retries.         |

### Webhook

| Env Var               | Default | Description                         |
|-----------------------|---------|-------------------------------------|
| WEBHOOKS_TOKEN_LENGTH | 32      | Length of generated webhook tokens. |

### Prometheus

| Env Var                         | Default  | Description                                  |
|---------------------------------|----------|----------------------------------------------|
| PROMETHEUS_ENABLED              | false    | Enables Prometheus metrics.                  |
| PROMETHEUS_PORT                 | 8080     | Prometheus server port.                      |
| PROMETHEUS_LISTEN               | -        | Network address for Prometheus server.       |
| PROMETHEUS_BASE_PATH            | /        | Base path for Prometheus server.             |
| PROMETHEUS_METRICS_PATH         | /metrics | Metrics endpoint path.                       |
| PROMETHEUS_SECURE_TOKEN_ENABLED | true     | Enables secure token protection for metrics. |
| PROMETHEUS_SECURE_TOKEN         | -        | Token required for accessing metrics.        |
| PROMETHEUS_REFRESH_INTERVAL     | 60s      | Refresh interval for metrics collection.     |
