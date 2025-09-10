# Monitoring

_upda_ exposes managed _updates_ as [prometheus](https://prometheus.io) metrics, so that you can easily build a
dashboard in [Grafana](https://grafana.com), or even attach alerts to pending updates
via [alertmanager](https://prometheus.io/docs/alerting/latest/alertmanager/).

To enable prometheus exporters of _upda_, please visit the [Configuration](./Configuration.md) section.

When `PROMETHEUS_ENABLED` is set to `true`, default metrics about memory utilization, but also custom metrics specific
to _upda_ are exposed under the `PROMETHEUS_METRICS_PATH` endpoint.

A Prometheus scrape configuration might look like the following if `PROMETHEUS_SECURE_TOKEN_ENABLED` is set to `true`.

```shell
scrape_configs:
  - job_name: 'upda'
    static_configs:
      - targets: ['<ip address of upda>:8080']
    bearer_token: 'VALUE_OF_PROMETHEUS_SECURE_TOKEN'
```

Prometheus can also be spawned independently of the main application server by setting `PROMETHEUS_PORT` to a different
port than `SERVER_PORT`.

There are two type of [Grafana](https://grafana.com) dashboards:

* [upda specific with custom metrics and](./grafana/upda.json)
* [general Go resource and HTTP request/response](./grafana/go_ginprom.json)

> Instead of importing the dashboard directly with Grafana, create a new Dashboard from scratch, directly go into
> Settings and inside the JSON Model, replace all contents with the contents of either of the above files. This keeps
> the generic "datasource" selection.
