# README

upda - **Up**date **Da**shboard in Go. Please see [motivation](#motivation) and [concepts](#concepts) what this
application does.

There's also a [upda web interface](https://git.myservermanager.com/varakh/upda-ui). It's recommended to take a look (at
least at the screenshots).

In addition, there's a commandline tool called `upda-cli`. For more information, download it and run `./upda-cli help`
for further instructions. This is especially useful, if you have an `upda` (server) running and like to invoke webhooks
from CLI. `upda-cli` is also bundled in the docker images.

**See the [deployment instructions](./_doc/DEPLOYMENT.md) for examples on how to deploy upda and upda-ui**

The main git repository is hosted at
_[https://git.myservermanager.com/varakh/upda](https://git.myservermanager.com/varakh/upda)_.
Other repositories are mirrors and pull requests, issues, and planning are managed there.

Contributions are very welcome!

* [Motivation](#motivation)
* [Concepts](#concepts)
* [Configuration](#configuration)
* [3rd party integrations](#3rd-party-integrations)
    * [Webhooks](#webhooks)
    * [Prometheus Metrics](#prometheus-metrics)
* [Deployment](#deployment)
    * [Native](#native)
    * [Docker](#docker)
        * [Build docker image](#build-docker-image)
* [Development & contribution](#development--contribution)
    * [Getting started](#getting-started)
        * [Windows hints](#windows-hints)
    * [Release](#release)

## Motivation

> [duin](https://crazymax.dev/diun/) can determine which OCI images have updates
> available. [Argus](https://release-argus.io) can query other sources like GitHub and even invoke actions when an
> update
> has been found, but there's no _convenient_ way of having **one** dashboard or source of truth for all of them across
> different hosts without tinkering with collecting them somewhere in one place. This application is the result of that
> tinkering. :-)

Managing various application or OCI container image updates can be a tedious task:

* A lot of hosts to operate with a lot of different applications being deployed
* A lot of different OCI containers to watch for updated images
* No convenient dashboard to see and manage all the available updates in one place

_upda_ manages a list of updates with attributes attached to it. For new updates to arrive, _upda_ needs to be called
via a webhook call (created within _upda_) from other applications, such as a bash script, an
application like [duin](https://crazymax.dev/diun/) or simply by using the `upda-cli`.

After an update is being tracked, _upda_ provides a convenient way to have everything in one place. In addition, it
exposes managed _updates_ as [prometheus](https://prometheus.io) metrics, so that you can easily build a dashboard
in [Grafana](https://grafana.com), or even attach alerts to pending updates
via [alertmanager](https://prometheus.io/docs/alerting/latest/alertmanager/).

In addition, you can use _upda_'s UI to manage updates, e.g. _approve_ them when they have been rolled out to a host.

Important to note:

* _upda_ is **NOT a scraper** to watch docker registries or GitHub releases, it simply collects and consolidates updates
  from different sources via _webhooks_. If you like to watch GitHub releases, write a scraper and use `upda-cli` to
  report back to _upda_.
* _upda_ uses basic auth for administrative tasks like viewing available updates or setting up the initial webhooks.

## Concepts

_upda_ retrieves new updates when webhooks of upda are invoked, e.g., [duin](https://crazymax.dev/diun/) invokes it
or any other application which can reach the instance.
Tracked updates are unique for the attributes `(application,provider,host)` which means that subsequent updates for an
identical _application_, _provider_ and _host_ simply updates the `version` and `metadata` attributes for that tracked
_update_ (regardless if the version or metadata payload _actually_ changed - reasoning behind this is to get reflected
metadata updates independent if version attribute has changed).

State management of tracked updates:

* On first creation, state is set to _pending_.
* When an _update_ is in _approved_ state, an invocation for it resets its state to _pending_.
* _Ignored_ updates are skipped entirely and no attribute is updated.

##### The `application` attribute

The _application_ attribute is an arbitrary identifier, name or label of a subject you like to track,
e.g., `docker.io/varakh/upda` for an OCI image.

##### The `provider` attribute

The _provider_ attribute is an arbitrary name or label. During webhook invocation the provider attribute is derived in
priority:

For the _generic_ webhook:

1. If the incoming payload contains a non-blank `provider` attribute, it's taken from the request.
2. If the incoming payload contains a blank or missing `provider` attribute, the issuing webhook's label is taken.

For the _diun_ webhook:

1. If the issuing webhook's label is blank, then `oci` is used.
2. In any other case, the webhook's label is used.

Because the first priority is the issuing webhook's label, setting the _same_ label for all webhooks results in a
grouping. Also see the _ignore host_ setting for `host` below.

_Remember that changing a webhook's label won't be reflected in already created/tracked updates!_

##### The `host` attribute

_host_ should be set to the originating host name a webhook has been issued from. The _host_
attribute can also be "ignored" (a setting in each webhook). If set to ignored, _upda_ sets _host_ to _global_, thus
update versions can be grouped independent of the originating host. If set for all webhooks, you'll end up with a host
independent update dashboard.

##### The `version` attribute

The _version_ attribute is an arbitrary name or label and subject to change across invocations of webhooks. This can be
a version number, a number of total updates, anything.

##### The `metadata` attribute

An update can hold any additional metadata information provided by request payload `metadata`. Metadata can be inspected
via web interface or API.

## Configuration

The following environment variables can be used to modify application behavior.

| Variable                           | Purpose                                                                                                                                                                                                       | Default/Description                                                                                                                  |
|:-----------------------------------|:--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|:-------------------------------------------------------------------------------------------------------------------------------------|
| `TZ`                               | The time zone (**recommended** to set it properly, background tasks depend on it)                                                                                                                             | Defaults to `Europe/Berlin`, can be any time zone according to _tz database_                                                         |
| `ADMIN_USER`                       | Admin user name for login                                                                                                                                                                                     | Not set by default, you need to explicitly set it to user name                                                                       |
| `ADMIN_PASSWORD`                   | Admin password for login                                                                                                                                                                                      | Not set by default, you need to explicitly set it to a secure random                                                                 |
|                                    |                                                                                                                                                                                                               |                                                                                                                                      |
| `DB_TYPE`                          | The database type (Postgres is **recommended**)                                                                                                                                                               | Defaults to `sqlite`, possible values are `sqlite` or `postgres`                                                                     |
| `DB_SQLITE_FILE`                   | Path to the SQLITE file                                                                                                                                                                                       | Defaults to `<XDG_DATA_DIR>/upda/upda.db`, e.g. `~/.local/share/upda/upda.db`                                                        |
| `DB_POSTGRES_HOST`                 | The postgres host                                                                                                                                                                                             | Postgres host address, defaults to `localhost`                                                                                       |
| `DB_POSTGRES_PORT`                 | The postgres port                                                                                                                                                                                             | Postgres port, defaults to `5432`                                                                                                    |
| `DB_POSTGRES_NAME`                 | The postgres database name                                                                                                                                                                                    | Postgres database name, needs to be set                                                                                              |
| `DB_POSTGRES_TZ`                   | The postgres time zone                                                                                                                                                                                        | Postgres time zone settings, defaults to `Europe/Berlin`                                                                             |
| `DB_POSTGRES_USER`                 | The postgres user                                                                                                                                                                                             | Postgres user name, needs to be set                                                                                                  |
| `DB_POSTGRES_PASSWORD`             | The postgres password                                                                                                                                                                                         | Postgres user password, needs to be set                                                                                              |
|                                    |                                                                                                                                                                                                               |                                                                                                                                      |
| `SERVER_PORT`                      | Port                                                                                                                                                                                                          | Defaults to `8080`                                                                                                                   |
| `SERVER_LISTEN`                    | Server's listen address                                                                                                                                                                                       | Defaults to empty which equals `0.0.0.0`                                                                                             |
| `SERVER_TLS_ENABLED`               | If server uses TLS                                                                                                                                                                                            | Defaults `false`                                                                                                                     |
| `SERVER_TLS_CERT_PATH`             | When TLS enabled, provide the certificate path                                                                                                                                                                |                                                                                                                                      |
| `SERVER_TLS_KEY_PATH`              | When TLS enabled, provide the key path                                                                                                                                                                        |                                                                                                                                      |
| `SERVER_TIMEOUT`                   | Timeout the server waits before shutting down to end any pending tasks                                                                                                                                        | Defaults to `1s` (1 second), qualifier can be `s = second`, `m = minute`, `h = hour` prefixed with a positive number                 |
| `CORS_ALLOW_ORIGIN`                | CORS configuration                                                                                                                                                                                            | Defaults to `*`                                                                                                                      |
| `CORS_ALLOW_METHODS`               | CORS configuration                                                                                                                                                                                            | Defaults to `GET, POST, PUT, PATCH, DELETE, OPTIONS`                                                                                 |
| `CORS_ALLOW_HEADERS`               | CORS configuration                                                                                                                                                                                            | Defaults to `Authorization, Content-Type`                                                                                            |
|                                    |                                                                                                                                                                                                               |                                                                                                                                      |
| `LOGGING_LEVEL`                    | Logging level. Possible are `debug`, `info`, `warn`, `error`, `dpanic`, `panic`, `fatal`. Setting to `debug` enables high verbosity output.                                                                   | Defaults to `info`                                                                                                                   |
| `LOGGING_ENCODING`                 | Logging encoding. Possible are `console` and `json`                                                                                                                                                           | Defaults to `json`                                                                                                                   |
| `LOGGING_DIRECTORY`                | Logging directory. When set, logs will be added to a file called `upda.log` in addition to the standard output. Ensure that upda has access permissions. Use an external program for log rotation if desired. |                                                                                                                                      |
|                                    |                                                                                                                                                                                                               |                                                                                                                                      |
| `WEBHOOKS_TOKEN_LENGTH`            | The length of the token                                                                                                                                                                                       | Defaults to `16`, positive number                                                                                                    |
|                                    |                                                                                                                                                                                                               |                                                                                                                                      |
| `TASK_UPDATE_CLEAN_STALE_ENABLED`  | If background task should run to do housekeeping of stale (ignored/approved) updates from the database                                                                                                        | Defaults to `false`                                                                                                                  | 
| `TASK_UPDATE_CLEAN_STALE_INTERVAL` | Interval at which a background task does housekeeping by deleting stale (ignored/approved) updates from the database                                                                                          | Defaults to `1h` (1 hour), qualifier can be `s = second`, `m = minute`, `h = hour` prefixed with a positive number                   | 
| `TASK_UPDATE_CLEAN_STALE_MAX_AGE`  | Number defining at which age stale (ignored/approved) updates are deleted by the background task (_updatedAt_ attribute decides)                                                                              | Defaults to `168h` (168 hours = 1 week), qualifier can be `s = second`, `m = minute`, `h = hour` prefixed with a positive number     |
| `TASK_EVENT_CLEAN_STALE_ENABLED`   | If background task should run to do housekeeping of stale (old) events from the database                                                                                                                      | Defaults to `false`                                                                                                                  |
| `TASK_EVENT_CLEAN_STALE_INTERVAL`  | Interval at which a background task does housekeeping by deleting stale (old) events from the database                                                                                                        | Defaults to `8h` (8 hours), qualifier can be `s = second`, `m = minute`, `h = hour` prefixed with a positive number                  |
| `TASK_EVENT_CLEAN_STALE_MAX_AGE`   | Number defining at which age stale (old) events are deleted by the background task (_updatedAt_ attribute decides)                                                                                            | Defaults to `2190h` (2190 hours = 3 months), qualifier can be `s = second`, `m = minute`, `h = hour` prefixed with a positive number |
| `TASK_PROMETHEUS_REFRESH_INTERVAL` | Interval at which a background task updates custom metrics                                                                                                                                                    | Defaults to `60s` (60 seconds), qualifier can be `s = second`, `m = minute`, `h = hour` prefixed with a positive number              |
|                                    |                                                                                                                                                                                                               |                                                                                                                                      |
| `LOCK_REDIS_ENABLED`               | If locking via REDIS (multiple instances) is enabled. Requires REDIS. Otherwise uses in-memory locks.                                                                                                         | Defaults to `false`                                                                                                                  |
| `LOCK_REDIS_URL`                   | If locking via REDIS is enabled, this should point to a resolvable REDIS instance, e.g. `redis://<user>:<pass>@localhost:6379/<db>`.                                                                          |                                                                                                                                      |
|                                    |                                                                                                                                                                                                               |                                                                                                                                      |
| `PROMETHEUS_ENABLED`               | If Prometheus metrics are exposed                                                                                                                                                                             | Defaults to `false`                                                                                                                  |
| `PROMETHEUS_METRICS_PATH`          | Defines the metrics endpoint path                                                                                                                                                                             | Defaults to `/metrics`                                                                                                               |
| `PROMETHEUS_SECURE_TOKEN_ENABLED`  | If Prometheus metrics endpoint is protected by a token when enabled (**recommended**)                                                                                                                         | Defaults to `true`                                                                                                                   |
| `PROMETHEUS_SECURE_TOKEN`          | The token securing the metrics endpoint when enabled (**recommended**)                                                                                                                                        | Not set by default, you need to explicitly set it to a secure random                                                                 |

## 3rd party integrations

### Webhooks

This is the core mechanism of _upda_ and why it exists. Webhooks are the central piece of how _upda_ gets notified about
updates.

In order to configure a 3rd party application like [duin](https://crazymax.dev/diun/) to send updates to _upda_ with
the [duin webhook notification configuration](https://crazymax.dev/diun/notif/webhook/), **create** a new _upda_ webhook
token via _upda_'s web interface or via API call. This gives you

* a unique _upda_ URL to configure in the notification part of [duin](https://crazymax.dev/diun/),
  e.g., `/api/v1/webhooks/<a unique identifier>`
* a corresponding token for the URL which must be sent as `X-Webhook-Token` header when calling _upda_'s URL

Expected payload is derived from the _type_ of the webhook which has been created in _upda_.

Example for [duin Webhook notification](https://crazymax.dev/diun/notif/webhook/) `notif`:

```yaml
notif:
    webhook:
        endpoint: https://upda.domain.tld/api/v1/webhooks/ee03cd9e-04d0-4c7f-9866-efe219c2501e
        method: POST
        headers:
            content-type: application/json
            X-Webhook-Token: <the token from webhook creation in upda>
        timeout: 10s
```

### Prometheus Metrics

When `PROMETHEUS_ENABLED` is set to `true`, default metrics about memory utilization, but also custom metrics specific
to _upda_ are exposed under the `PROMETHEUS_METRICS_PATH` endpoint.

A Prometheus scrape configuration might look like the following if `PROMETHEUS_SECURE_TOKEN_ENABLED` is set to `true`.

```shell
scrape_configs:
  - job_name: 'upda'
    static_configs:
      - targets: ['upda:8080']
    bearer_token: 'VALUE_OF_PROMETHEUS_SECURE_TOKEN'
```

Custom exposed metrics are exposed under the `upda_` namespace.

Examples:

```shell
# HELP upda_updates details for all updates, -1=deleted (deleted next restart), 0=pending, 1=approved, 2=ignored
upda_updates{application="codeberg.org/forgejo/forgejo",host="myserver",provider="oci"} 0
upda_updates{application="docker.io/library/mysql",host="myserver",provider="oci"} 2
upda_updates{application="quay.io/navidys/prometheus-podman-exporter",host="myserver",provider="oci"} 1
upda_updates{application="quay.io/navidys/prometheus-podman-exporter",host="myserver2",provider="oci"} 1
# HELP upda_updates_all amount of all updates
upda_updates_all 4
# HELP upda_updates_approved amount of all updates in approved state
upda_updates_approved 2
# HELP upda_updates_ignored amount of all updates in ignored state
upda_updates_ignored 1
# HELP upda_updates_pending amount of all updates in pending state
upda_updates_pending 1
# HELP upda_webhooks amount of all webhooks
upda_webhooks 2
# HELP upda_events amount of all events
upda_events 146
```

There's an example [Grafana](https://grafana.com) dashboard in the `_doc/` folder.

[Alertmanager](https://prometheus.io/docs/alerting/latest/alertmanager/) could check for the following:

```yaml
- name: update_checks
  rules:
    - alert: UpdatesAvailable
      expr: upda_updates == 0 and upda_updates_pending > 0
      for: 4w
      labels:
        severity: high
        class: update
      annotations:
        summary: "Updates available from upda for {{ $labels.job }}"
        description: "Updates available from upda for {{ $labels.job }}"
```

## Deployment

### Native

Use the released binary for your platform or run `make clean build-server-{your-platform}` and the binary will be placed
into the `bin/` folder.

### Docker

For examples how to run, look into [deployment instructions](./_doc/DEPLOYMENT.md) which contains examples
for `docker-compose` files.

#### Build docker image

To build docker images, do the following

```shell
docker build --rm --no-cache -t upda:latest .
```

## Development & contribution

* Ensure to set the following environment variables for proper debug logs during development

```shell
DEVELOPMENT=true
LOGGING_ENCODING=console
LOGGING_LEVEL=debug
```

* Code guidelines
    * Each entity has its own repository
    * Each entity is only used in repository and service (otherwise, mapping happens)
    * Presenter layer is constructed from the entity, e.g., in REST responses and mapped
    * No entity is directly returned in any REST response
    * All log calls should be handled by `zap.L()`
    * Configuration is bootstrapped via separated `struct` types which are given to the service which need them
    * Error handling
        * Always throw an error with `NewServiceError`
        * Always wrap the cause error with `fmt.Errorf`
        * Forward/bubble up the error directly, when original error is already a `NewServiceError` (most likely internal
          calls)
        * Always abort handler chain with `AbortWithError`
        * Utils can throw any error

Please look into the `_doc/` folder for [OpenAPI specification](./_doc/api.yaml) and a Postman Collection.

### Getting started

1. Run `make clean dependencies` to fetch dependencies
2. Start `git.myservermanager.com/varakh/upda/cmd/server` (or `cli`) as Go application and ensure to have _required_
   environment variables set

If you like to test with PSQL and/or REDIS for task locking, here are some useful docker commands to have containers
up and running quickly. Set necessary environment variables properly.

```shell
# postgres
docker run --rm --name=upda-db \
  -p 5432:5432 \
  --restart=unless-stopped \
  -e POSTGRES_USER=upda \
  -e POSTGRES_PASSWORD=upda \
  -e POSTGRES_DB=upda \
  postgres:16-alpine
  
# redis
docker run --rm --name some-redis \
  -p 6379:6379 \
  redis redis-server --save 60 1 --loglevel warning
```

#### Windows hints

On Windows, you need a valid `gcc`, e.g., https://jmeubank.github.io/tdm-gcc/download/ and add the `\bin` folder to your
path.

For any `go` command you run, ensure that your `PATH` has the `gcc` binary and that you add `CGO_ENABLED=1` as
environment.

### Release

Releases are handled by the SCM platform and pipeline. Creating a **new git tag**, creates a new release in the SCM
platform, uploads produced artifacts to that release and publishes docker images automatically.
**Before** doing so, please ensure that the **commit on `master`** has the **correct version settings** and has been
built successfully:

* Adapt `constants_app.go` and change `Version` to the correct version number
* Adapt `CHANGELOG.md` to reflect changes and ensure a date is properly set in the header, also add a reference link
  in footer (link to scm git tag source)
* Adapt `api.yaml`: `version` attribute must reflect the to be released version
* Adapt `env: VERSION_*` in `.forgejo/workflows/release.yaml`

After the release has been created, ensure to change the following settings for the _next development cycle_:

* Adapt `constants_app.go` and change `Version` to the _next_ version number
* Adapt `CHANGELOG.md` and add an _UNRELEASED_ section
* Adapt `api.yaml`: `version` attribute must reflect the _next_ version number
* Adapt `env: VERSION_*` in `.forgejo/workflows/release.yaml` to _next_ version number

