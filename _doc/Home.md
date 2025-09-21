# upda

**Up**date **Da**shboard (upda). A centralized tool for tracking and managing updates across various systems,
applications, and container images.

_upda_ provides a **single** dashboard to manage, display, and organize update information from many different hosts
and sources. Instead of checking each system or container image individually, administrators can see all available
updates in one convenient place.

Already familiar with _upda_? **[Deploy it now](./Deployment.md)** or see **[Configuration](./Configuration.md)**.

## What It Solves

- Eliminates the need to manually check for updates across many servers or container images.
- Provides clarity when managing a large, complex environment with multiple deployed applications.
- Saves time by consolidating update tracking into one dashboard instead of scattered tools or logs.

Curious about **[Use Cases](./UseCases.md)** or an **[example workflow](./ExampleWorkflow.md)**?

## Key Features

- **Centralized update management**: Collects, manages, and visualizes [update](./Usage.md#manage-updates) information
  from multiple hosts, applications, or container registries. Updates can be assigned different states such as
  *pending*, *ignored*, or *approved*, which helps track which ones still need attention.
- **Data ingestion**: Updates don't appear automatically; they must be pushed to _upda_ by other systems or scripts
  through webhooks. Each webhook has a unique URL that can be triggered by any external tool, even a simple bash script.
  Any system or script can push update data to _upda_ through its [webhook](./Usage.md#getting-updates-in-via-webhooks)
  interface.
- **Flexible data sources**:
    - Works with tools like *diun* (for Docker/Podman/OCI container updates).
    - Accepts custom scripts or CLI calls to report updates.
- **Metadata-rich updates**: Each update can include details like version number, originating host, and additional
  metadata.
- **Visualization and organization**: Provides a user-friendly interface to view, filter, and manage updates.
- **Extensibility**: Users can contribute or use scripts to simplify integration from their own systems, you just need
  to call the webhook.
- **Actions**: When an update is created or changes state, _upda_ can trigger [actions](./Usage.md#actions) that notify
  or integrate with other systems using tools like [shoutrrr](https://containrrr.dev/shoutrrr/).
- **History of actions**: You can view logs of [past actions](./Usage.md#history-of-actions) triggered by updates.
- **Events**: _upda_ records [events](./Usage.md#see-what-has-changed), such as when an update changes state or version,
  providing an audit trail of what changed and when.
- **Metrics exporter**: It provides a Prometheus endpoint so that update-related statistics can
  be [monitored](./Monitoring.md) externally.

In short, _upda_ acts as a **central orchestrator for updates**, combining data ingestion, state management, and
event-triggered integrations to external systems.

Please head over to the [Usage](./Usage.md) section to dive into using _upda_. Before, make sure that you have
_upda_ [deployed](./Deployment.md).

## What It Doesn't Solve

_upda_ is **NOT** a scraper to watch docker registries or GitHub releases, it simply tracks and consolidates updates
from different sources, but you need to feed in these information on your own, e.g., via Webhooks. If you like to watch
GitHub releases, write a scraper and use the binary `upda` to report back to _upda_.

_Though_, in the documentation section, there are some contributed scripts to report into _upda_, for example to check
public GitHub releases. See [contrib/](./contrib/).
