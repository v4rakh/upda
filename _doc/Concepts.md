# Concepts, a deeper dive

The following section goes into a deeper look into upda's internals.

_upda_ retrieves new updates when webhooks of _upda_ are invoked, e.g., [duin](https://crazymax.dev/diun/) invokes it or
any other application/script which can reach the instance. Tracked updates are unique for the
attributes `(application,provider,host)` which means that subsequent updates for an identical _application_, _provider_
and _host_ simply updates the `version` and `metadata` attributes for that tracked _update_ (regardless if the version
or metadata payload _actually_ changed - reasoning behind this is to get reflected metadata updates independent if
version attribute has changed).

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