# CHANGELOG

Changes adhere to [semantic versioning](https://semver.org).

## [6.1.0] - UNRELEASED

* Add `EMBEDDED_WEB_INTERFACE_DARK_THEME_ENABLED` to enforce dark mode (defaults to `0`)
* Add `EMBEDDED_WEB_INTERFACE_FOOTER_ENABLED` to show page footer (defaults to `1`)
* Remove default `robots.txt`, add this manually in your reverse proxy if you like to have it served
* Dependency updates
* Code cleanup and refactor

## [6.0.1] - 2025/02/17

* Hotfix: Incorrect instruction to forward `VITE_TITLE` inside embedded web interface

## [6.0.0] - 2025/02/16

> This is a major version upgrade. Though not incompatible, the environment configuration has breaking changes which
> require **manual intervention** to adapt.

* Add _constants_ which can be used similarly to `<SECRET>` and `<VAR>` in actions with `<CONST>`
* Disallow _constants_ and _secrets_ with identical keys
* Switch to text area for shoutrrr action URL and body form, but disallow line breaks
* Switch to ghost button style for update cards
* Switch default time zone (`TZ` and `DB_POSTGRES_TZ`) to `Etc/UTC`
* Switch default `WEBHOOKS_TOKEN_LENGTH` from `16` to `30`
* Allow to disable the embedded frontend with `EMBEDDED_WEB_INTERFACE_ENABLED=false` (defaults to `true`)
* Rename `WEB_API_URL` configuration to `EMBEDDED_WEB_INTERFACE_API_URL` which requires the full API URL with trailing slash, e.g. `https://upda.domain.tld/api/v1/`
* Rename `WEB_TITLE` configuration to `EMBEDDED_WEB_INTERFACE_TITLE`
* Allow to set a base path with `SERVER_BASE_PATH` which defaults to `/`, all routes are prefixed with this base path (except embedded web interface if enabled)
* Dependency updates

## [5.0.0] - 2024/12/21

> This is a major version upgrade. Other versions are incompatible with this release.

* Drop support for SQLite (only Postgres is supported)
* Library updates
* Update OCI image base to alpine `3.21` with Go `1.23`
* Move away from `npm` to `pnpm`

## [4.0.0] - 2024/10/25

> This is a major version upgrade. Other versions are incompatible with this release.

* Embed frontend into Go binary and only ship _one_ OCI image
* Switch license to GPLv3

## [3.0.2] - 2024/06/15

* Don't enforce JSON content type for GET and DELETE requests
* Dependency updates
    * github.com/go-playground/validator/v10 v10.22.0
    * gorm.io/driver/postgres v1.5.9
    * gorm.io/driver/sqlite v1.5.6
* Fixed filter for Updates ignoring desired state

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

[6.0.1]: https://git.myservermanager.com/varakh/upda/releases/tag/6.0.1

[6.0.0]: https://git.myservermanager.com/varakh/upda/releases/tag/6.0.0

[5.0.0]: https://git.myservermanager.com/varakh/upda/releases/tag/5.0.0

[4.0.0]: https://git.myservermanager.com/varakh/upda/releases/tag/4.0.0

[3.0.2]: https://git.myservermanager.com/varakh/upda/releases/tag/3.0.2

[3.0.1]: https://git.myservermanager.com/varakh/upda/releases/tag/3.0.1

[3.0.0]: https://git.myservermanager.com/varakh/upda/releases/tag/3.0.0

[2.0.1]: https://git.myservermanager.com/varakh/upda/releases/tag/2.0.1

[2.0.0]: https://git.myservermanager.com/varakh/upda/releases/tag/2.0.0

[1.0.3]: https://git.myservermanager.com/varakh/upda/releases/tag/1.0.3

[1.0.2]: https://git.myservermanager.com/varakh/upda/releases/tag/1.0.2

[1.0.1]: https://git.myservermanager.com/varakh/upda/releases/tag/1.0.1

[1.0.0]: https://git.myservermanager.com/varakh/upda/releases/tag/1.0.0