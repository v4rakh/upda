# CHANGELOG

Changes adhere to [semantic versioning](https://semver.org).

## [3.0.1] - 2024/06/10

* Fixed finding proper remaining Action invocations by their state

## [3.0.0] - 2024/06/10

> This is a major version upgrade. Other versions are incompatible with this release.

* Added automatic detection of `GOMAXPROCS`
* Switched to enforce JSON as `Content-Type` for all incoming requests
* Switched to properly respond with JSON on page not found or method not allowed
* Renamed `CORS_ALLOW_ORIGIN` to `CORS_ALLOW_ORIGINS`
* Added `CORS_ALLOW_CREDENTIALS` which defaults to `true`
* Added `CORS_EXPOSE_HEADERS` which defaults to `*`
* Overhauled package visibility for server module
* Updated dependencies
* Updated OCI image base to alpine `3.20` with Go `1.22`

## [2.0.1] - 2024/05/01

* Fixed retrieval of encrypted webhook token

## [2.0.0] - 2024/04/28

> This is a major version upgrade. Other versions are incompatible with this release.

* Added _Actions_, a simple way to trigger notifications via [shoutrrr](https://containrrr.dev/shoutrrr) which supports
  secrets
* Added new auth mode which allows setting multiple basic auth credentials
    * Added `AUTH_MODE` which can be one of `basic_single` (_default_) and `basic_credentials`
    * For `basic_credentials`: added `BASIC_AUTH_CREDENTIALS` which can be used as list of `username1=password1,...` (
      comma separated)
    * For `basic_single`: renamed `ADMIN_USER` and `ADMIN_PASSWORD` to `BASIC_AUTH_USER` and `BASIC_AUTH_PASSWORD`
* Added mandatory `SECRET` environment variable to encrypt some data inside the database
* Switched to producing events only for _Updates_
* Switched to encrypting webhook tokens in database
* Adapted logging which defaults to JSON encoding
* Updated dependencies

## [1.0.3] - 2024/01/21

* Updated dependencies

## [1.0.2] - 2023/12/23

* Fix wrong event type being created for update state change

## [1.0.1] - 2023/12/23

* Disable cleaning up stale updates and events by default
* Change Prometheus exporter behavior
    * Return `-1` for deleted updates in Prometheus which are evicted on next application restart
    * Ignore `PROMETHEUS_METRICS_PATH` (defaults to `/metrics`) in application metrics
* Introduce locking for periodic background tasks
    * Rename `TASK_LOCK_REDIS_ENABLED` to `LOCK_REDIS_ENABLED` which still defaults to `false` (disabled)
    * Rename `TASK_LOCK_REDIS_URL` to `LOCK_REDIS_URL`

## [1.0.0] - 2023/12/21

* Initial release

[3.0.1]: https://git.myservermanager.com/varakh/upda/releases/tag/3.0.1

[3.0.0]: https://git.myservermanager.com/varakh/upda/releases/tag/3.0.0

[2.0.1]: https://git.myservermanager.com/varakh/upda/releases/tag/2.0.1

[2.0.0]: https://git.myservermanager.com/varakh/upda/releases/tag/2.0.0

[1.0.3]: https://git.myservermanager.com/varakh/upda/releases/tag/1.0.3

[1.0.2]: https://git.myservermanager.com/varakh/upda/releases/tag/1.0.2

[1.0.1]: https://git.myservermanager.com/varakh/upda/releases/tag/1.0.1

[1.0.0]: https://git.myservermanager.com/varakh/upda/releases/tag/1.0.0