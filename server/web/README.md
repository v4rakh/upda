# upda-ui

Frontend for upda - **Up**date **Da**shboard in React, TypeScript, Redux.

The main git repository is hosted at
_[https://git.myservermanager.com/varakh/upda-ui](https://git.myservermanager.com/varakh/upda-ui)_. Other repositories
are mirrors and pull requests, issues, and planning are managed there.

Contributions are very welcome!

[Official documentation](https://git.myservermanager.com/varakh/upda-docs) is hosted in a separate git repository.

## Development & contribution

Contributions are very welcome!

### Prerequisites

It's probably worth checking out a node environment manager like [nvm manager](https://github.com/nvm-sh/nvm).

Required node and npm versions are outlined in the `package.json`.

### Setup instructions

Run `npm install` which should install all dependencies.

### Start

Use the `npm run start` command to start the development setup. Backend should be running.

### Configuration magic in docker

What about configuration? How does the pre-compiled set of html, js and css files know about the environment variables?

In contrast to manual build, the docker image allows dynamic override of configuration, but only those outlined in
the `.env` file.

In production, all configuration values are dynamically generated inside the Docker image with a helper script
called `docker-env.sh`:

1. During docker build a template `.env` file and the helper script are copied to the docker image
2. Before the container's nginx is started, the helper script
    1. scans the `.env` file for known configuration variables and then
    2. adds their values to `conf/runtime-config.js` which is sourced inside the application in the
       immutable `window.runtime_config` object

During development, this `runtime-config.js` file is still loaded, but empty and thus the `getConfiguration()` ignores
it and prefers values from the sourced `.env.development`.

This means that new environment variables need to be added to all `.env*` files!

### Release

Releases are handled by the SCM platform and pipeline. Creating a **new git tag**, creates a new release in the SCM
platform, uploads produced artifacts to that release and publishes docker images automatically.
**Before** doing so, please ensure that the **commit on `master`** has the **correct version settings** and has been
built successfully:

* Adapt `package.json` and change `version` to the correct version number
* Invoke `npm install` once which properly sets the version inside the lock file
* Adapt language files, e.g., `en.json` and change `version` to the correct version number
* Adapt `CHANGELOG.md` to reflect changes and ensure a date is properly set in the header, also add a reference link
  in footer
* Adapt `env: VERSION_*` in `.forgejo/workflows/release.yaml`

After the release has been created, ensure to change the following settings for the _next development cycle_:

* Adapt `package.json` and change `version` to the _next_ version number (semantic versioning applied, so `patch`
  should bump patch level version and prepare branch for `develop` should bump minor or major version)
* Invoke `npm install` for each of those branches which properly sets the version inside the lock file
* Adapt language files, e.g., `en.json` and change `version` to the _next_ version number (semantic versioning
  applied, so `patch` should bump patch level version and prepare branch for `develop` should bump minor or major
  version)
* Adapt `CHANGELOG.md` and add an _UNRELEASED_ section
* Adapt `env: VERSION_*` in `.forgejo/workflows/release.yaml` to _next_ version number
