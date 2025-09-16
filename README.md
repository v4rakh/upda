# README

**Up**date **Da**shboard (upda) written in Go. A centralized tool for tracking and managing updates across various
systems, applications, and container images.

The main git repository is hosted at
_[https://git.myservermanager.com/varakh/upda](https://git.myservermanager.com/varakh/upda)_.
Other repositories are mirrors and pull requests, issues, and planning are managed there.

Contributions are very welcome, please see [Development & contribution](#development--contribution).

## [Official documentation](./_doc/Home.md)

Please head over to the documentation for setup, usage, and operation. This section mainly focuses on development.

## Development & contribution

Head over to the [getting started](#getting-started) for setting up the development environment first.

### Guidelines

* Pay attention to `make checkstyle` (uses `go vet ./...`); pipeline fails if issues are detected.
* Each entity has its own repository
* Each entity is only used in repository and service (otherwise, mapping happens, latest at controller level)
* Presenter layer is constructed from the entity, e.g., in REST responses and mapped
* No entity is directly returned in any REST response
* All log calls should be handled by `log`
* Configuration is bootstrapped via separated `struct` types which are given to the service which need them
* Error handling
    * Always throw an error with `NewServiceError` for repositories, services and handlers
    * Always throw an error wrapping the cause with `fmt.Errorf`
    * Forward/bubble up the error directly, when original error is already a `NewServiceError` (most likely internal
      calls)
    * Always abort handler chain with `AbortWithError`
    * Utils can throw any error
    * Repositories, handlers and services should always properly return `error` including any `init`-like function (
      best to avoid them and call in `newXXX`). **Do not abort with `Fatalf` or similar**
    * `log.Fatalf` or `log.Fatal()` is allowed in `config.go` or `server.go`
* Look into the `_doc/` folder for [OpenAPI specification](./_doc/api.yaml).
* Consider reading [Effective Go](https://go.dev/doc/effective_go)
* Consider reading [100 Go Mistakes and How to Avoid Them](https://100go.co/)

Be aware that some are false positives and actually required.

### Getting started

The most straight forward way to get started is by looking into available commands inside the `Makefile`.

For the full setup, you need the following tools:

- go (see minimum version in `go.mod`)
- pnpm and node (see version constraints in `package.json`)
- make to execute commands of the `Makefile`

Though, when you're familiar with [direnv](https://github.com/direnv/direnv) or even the package
manager [nix](https://nixos.org/), you can achieve a full and easy setup when you go into the project's directory.

#### `direnv` / `nix-direnv`

This project can optionally use [direnv](https://github.com/direnv/direnv) (
or [nix-direnv](https://github.com/nix-community/nix-direnv) with `nix`) to automatically load environment variables
from an `.env` file. Copy `.env.example` to `.env` and adjust the values accordingly.

When you change directory into the project, the environment variables are automatically loaded after you've allowed
`direnv` with `direnv allow`.

#### Nix Flakes

_In addition_, the project hosts a `flake.nix` and a `flake.lock` file. You can safely ignore them if you don't like to
use this method of bootstrapping your environment to work with this application. This setup though allows to easily have
an environment set up for this application by installing necessary required binaries without modifying your OS
installation. This is done through the package manager [nix](https://nixos.org/) which automatically installs everything
necessary under a _devShell_ (development shell). From shell which you can enter via `nix develop` inside the project's
root directory, all the above tools are available.
To automate it even further, [nix-direnv](https://github.com/nix-community/nix-direnv) recognises when you change
directory (or open an IDE terminal) within this project (similar to `direnv` itself).

**Keep in mind that the project itself is not available as flake.**

### Pre-requisites

Ensure to set the following environment variables for proper debug logs during development

```shell
WEB_INTERFACE_ENABLED=false
DEVELOPMENT=true
LOGGING_LEVEL=debug
LOGGING_LEVEL_REQUESTS=debug
LOGGING_ENCODING_COLORIZE=true
```

1. Run `make clean dependencies` to fetch dependencies
2. Start `git.myservermanager.com/varakh/upda/cmd/main.go`

If you like to test with Postgres and/or REDIS for task locking, here are some useful docker commands to have containers
up and running quickly. Set necessary environment variables properly.

```shell
# postgres
docker run --name=upda-db \
  -p 5432:5432 \
  --restart=unless-stopped \
  -e POSTGRES_USER=upda \
  -e POSTGRES_PASSWORD=upda \
  -e POSTGRES_DB=upda \
  postgres:17-alpine
  
# redis
docker run --name some-redis \
  -p 6379:6379 \
  redis redis-server --save 60 1 --loglevel warning
```

### Web interface

_upda_ includes a frontend in a monorepo fashion inside `server/web/`. For production (binary and OCI), it's
embedded into the GoLang binary itself.

For _development_, no other steps are required. Simply follow
the [frontend instructions](./internal/server/web/README.md) and
start the frontend separately.

If you like to have a look on the _production_ experience, the frontend needs to be build first, and you need to build
the Golang binary with `-tags prod`. How to properly build the frontend, please look into `build-web` of
the `Makefile` (additional `rm -rf` cmd!).

### Windows hints

On Windows, you need a valid `gcc`, e.g., https://jmeubank.github.io/tdm-gcc/download/ and add the `\bin` folder to your
path.

For any `go` command you run, ensure that your `PATH` has the `gcc` binary and that you add `CGO_ENABLED=1` as
environment if `go` commands fail.

### enums

For new enums or when changing existing ones, use the `make generate` task which
uses [go-enum](https://github.com/abice/go-enum) to generate boilerplate code.

See example `enum.go`. Make sure to use the same `//go:generate` directives.

### Using the `lockService` correctly

The `lockService` can be used to lock resources. This works in-memory and also in a distributed fashion with REDIS.

Ensure to provide proper locking options when using, although in-memory ignores those.

Example:

```shell
# invoked from an endpoint
context := c.Request.Context()

var err error
var lock Lock

if lock, err = h.lockService.lockWithOptions(context, "TEST-LOCK", withLockOptionExpiry(5*time.Minute), withLockOptionInfiniteRetries(), withLockOptionRetryDelay(5*time.Second)); err != nil {
    _ = c.AbortWithError(errToHttpStatus(err), err)
    return
}
# defer to avoid leakage
defer func(lock Lock) {
    _ = lock.unlock(context)
}(lock)

# simulate long running task
time.Sleep(20 * time.Second)
```

### Tests

There are multiple test targets defined in the `Makefile`

- Go: Unit tests.
- Go: Integration tests require an argument `image=...`, a built OCI image reference of _upda_, to work.
- Web: Executed together with Go unit tests.

If you're running on rootless docker or podman, set `DOCKER_HOST=unix:///run/user/1000/podman/podman.sock` (adapt for
docker).

### Git workflow

The main branch is `master`. It's protected and only eligible users can push to it. Merge requests to protected branches
are safe-guarded: they need review or at least a successful pipeline run to be merged.

- Merge request branches should start with `feat/`, `fix/`, `refactor/`, `chore/`, or `ci/`
- Merge requests should be squashed and the source branch should be deleted
- Merge request commits should have a meaningful commit message
- Merge request titles should have a meaningful title which is taken as commit message once merged
    - should be prefixed with `feat: ...`, `fix(...): ...`, where the contents inside the bracket should be _one word_
      which topic/component is touched (conventional commits)
    - should reflect a breaking change by adding a `!` before the colon, e.g., `fix(deps)!: ...`
    - should include more verbose information in the body of the git commit message (merge request description)

```
feat(security)!: add OpenID Connect authentication

- This adds a new authentication mode called oidc
- This new mode is the default, which might break existing installations
```

- Merge requests should contain documentation changes, so that code and documentation stays in sync

### Pipeline workflow

Pipeline runs

* on merge request change (open, new push, ...)
* on protected branches

This means you need to create a merge request to trigger a pipeline run. Without merge request, no build is triggered,
thus your code cannot be merged.

### Release preparation & workflow

Follow these steps to release the application

* Trigger the pipeline for a commit on the `master`
    * When asked, enter a version number which should align with semantic versioning
    * The pipeline creates a git tag and a release in the VCS management system
    * Wait until the release pipeline succeeded
* (_optional_) Generate the changelog
    * Got into the git repository, make sure to fetch (including just created tag)
    * Requires [git-cliff](https://git-cliff.org/) being installed
    * Invoke from last but one release git tag to the most recent release tag (just created) with
      `git-cliff OLDTAG..NEWTAG`, e.g., `git-cliff 6.0.0..6.1.0`
    * This prints markdown to your terminal.
    * Copy the markdown and edit the release in the VCS management system

There's no additional preparation needed before invoking the release pipeline on `master` as it should always represent
a working state at any time.

### Dependency updates

Dependency updates are handled by Renovate using the `renovate.json5` file. The base branch is `master`.

Major updates undergo manual review.
