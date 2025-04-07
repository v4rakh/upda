# upda

**Up**date **Da**shboard (upda). A simple application to keep track of updates from different hosts and systems.

Managing various application or OCI container image updates can be a tedious task:

* A lot of hosts to operate with a lot of different applications being deployed
* A lot of different OCI containers to watch for updated images
* No convenient dashboard to see and manage all the available updates in one place

_upda_ manages a list of updates with attributes attached to it. For new updates to arrive, _upda_ needs to be called
via a webhook call (created within _upda_) from other applications, such as a bash script, an
application like [duin](https://crazymax.dev/diun/) or simply by using the `upda` binary command-line.

Please head over to the [Usage](./Usage.md) section for a quick _Getting Started_ once you've [deployed](./Deployment.md)
_upda_.

The code is hosted here: [upda application including frontend](https://git.myservermanager.com/varakh/upda).

## Features

_upda_ manages a list of updates with attributes attached to it. For new updates to arrive, _upda_ needs to get them
from an external source.
For this, _upda_ allows to manage webhooks, which can be called with a unique URL from any other application or even a
bash script so that upda retrieves these information.

_upda_'s main features include

* Managing [Updates](./Usage.md#manage-updates) by changing their state (pending, ignored, approved)
* Managing [Webhooks](./Usage.md#getting-updates-in-via-webhooks) which allow to get information into _upda_ regarding Updates
  and their properties (like version) you like to track
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

