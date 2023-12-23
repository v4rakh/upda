# CHANGELOG

Changes adhere to [semantic versioning](https://semver.org).

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

[1.0.1]: https://git.myservermanager.com/varakh/upda/releases/tag/1.0.1

[1.0.0]: https://git.myservermanager.com/varakh/upda/releases/tag/1.0.0