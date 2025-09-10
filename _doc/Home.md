# upda

**Up**date **Da**shboard (upda). A simple application to keep track of updates from different hosts and systems.

Managing various application updates or OCI container image updates can be a tedious task, especially if you maintain a
lot of machines:

* Large set of different deployed applications to operate and maintain
* Staying up-to-date with Docker/Podman/OCI container updates
* Staying up-to-date with any other update your systems need

There's no central space to view, manage, and organize all available updates in one location, effortlessly and with a
convenient user interface.

_upda_ steps in here. From any machine, you can send "updates" to _upda_ which _upda_ then manages and visualizes. Each
update has certain attributes, like version, metadata, or the host it came from. For new updates to arrive, a client
needs to call _upda_ via its webhook functionality. Webhooks can be created in _upda's_ user interface.

This basically means _upda_ can retrieve new information about updates from anywhere, even a simple bash script you have
flying around. For Docker/Podman/OCI images, there's also an application called [duin](https://crazymax.dev/diun/) which
supports sending information via webhooks and _upda_ can be configured to retrieve them. You can also decide to send new
updates with your own script, just by calling the `upda` binary in your command-line.

Please head over to the [Usage](./Usage.md) section to learn more what _upda_ can do (like Actions, Update state
management), or jump to the _Getting Started_ once you've [deployed](./Deployment.md) _upda_. Also make sure to read
the [What it is not](#what-it-is-not) section. There are user contributed scripts already to ease sending update
information from a host to your _upda_ instance.

## Features

_upda_ manages a list of updates with attributes attached to it, like version or host. For new updates to arrive, _upda_
needs to get them from an external source.
For this, _upda_ allows to manage webhooks, which can be called with a unique URL from any other application or even a
bash script so that _upda_ retrieves these information.

_upda_'s main features include

* Managing [Updates](./Usage.md#manage-updates) by changing their state (pending, ignored, approved)
* Managing [Webhooks](./Usage.md#getting-updates-in-via-webhooks) which allow to get information into _upda_ regarding
  Updates and their properties (like version) you like to track
* Managing [Actions](./Usage.md#actions) which allow you to further process changes made to an Update (created, state
  changed, version
  changed,), basically allowing you to invoke other systems with the help
  of [shoutrrr](https://containrrr.dev/shoutrrr/)
* View [past invocation of Actions](./Usage.md#history-of-actions)
* Viewing [events](./Usage.md#see-what-has-changed) which allow you to see what has changed and how Updates
* [Metrics exporter](./Monitoring.md) via prometheus

_upda_ is designed to be simple. Only supported authorization mechanism is basic.

## What it is not

_upda_ is **NOT** a scraper to watch docker registries or GitHub releases, it simply tracks and consolidates updates
from different sources, but you need to feed in these information on your own, e.g., via Webhooks. If you like to watch
GitHub releases, write a scraper and use the binary `upda` to report back to _upda_.

Though, in the documentation section, there are some contributed scripts to report into _upda_, for example to check
public GitHub releases. See [contrib/](./contrib/).

