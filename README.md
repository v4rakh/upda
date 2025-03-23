# README

upda - **Up**date **Da**shboard in Go.

The main git repository is hosted at
_[https://git.myservermanager.com/varakh/upda](https://git.myservermanager.com/varakh/upda)_.
Other repositories are mirrors and pull requests, issues, and planning are managed there.

Contributions are very welcome, please see [Development & contribution](#development--contribution).

See **[official documentation](./_doc/Home.md)**.

## Development & contribution

There's also a [embedded frontend](#embedded-frontend).

* Pay attention to `make checkstyle` (uses `go vet ./...`); pipeline fails if issues are detected.
* Each entity has its own repository
* Each entity is only used in repository and service (otherwise, mapping happens)
* Presenter layer is constructed from the entity, e.g., in REST responses and mapped
* No entity is directly returned in any REST response
* All log calls should be handled by `zap.L()`
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
    * `log.Fatalf` or `zap.L().Fatal` is allowed in `environment.go` or `app.go`
* Look into the `_doc/` folder for [OpenAPI specification](./_doc/api.yaml) and a Postman Collection.
* Consider reading [Effective Go](https://go.dev/doc/effective_go)
* Consider reading [100 Go Mistakes and How to Avoid Them](https://100go.co/)

Be aware that some are false positives and actually required.

### Getting started

Ensure to set the following environment variables for proper debug logs during development

```shell
EMBEDDED_WEB_INTERFACE_ENABLED=false
DEVELOPMENT=true
LOGGING_ENCODING=console
LOGGING_LEVEL=debug
```

1. Run `make clean dependencies` to fetch dependencies
2. Start `git.myservermanager.com/varakh/upda/cmd/server` (or `cli`) as Go application and ensure to have _required_
   environment variables set

If you like to test with PSQL and/or REDIS for task locking, here are some useful docker commands to have containers
up and running quickly. Set necessary environment variables properly.

```shell
# postgres
docker run --rm --name=upda-db \
  -p 5432:5432 \
  --restart=unless-stopped \
  -e POSTGRES_USER=upda \
  -e POSTGRES_PASSWORD=upda \
  -e POSTGRES_DB=upda \
  postgres:17-alpine
  
# redis
docker run --rm --name some-redis \
  -p 6379:6379 \
  redis redis-server --save 60 1 --loglevel warning
```

### Embedded Frontend

_upda_ includes a frontend in a monorepo fashion inside `server/web/`. For production (binary and OCI), it's
embedded into the GoLang binary itself.

For _development_, no other steps are required. Simply follow the [frontend instructions](./server/web/README.md) and
start the frontend separately.

If you like to have a look on the _production_ experience, the frontend needs to be build first and you need to build
the Golang binary with `-tags prod`. How to properly build the frontend, please look into `build-web` of
the `Makefile` (additional `rm -rf` cmd).

### Windows hints

On Windows, you need a valid `gcc`, e.g., https://jmeubank.github.io/tdm-gcc/download/ and add the `\bin` folder to your
path.

For any `go` command you run, ensure that your `PATH` has the `gcc` binary and that you add `CGO_ENABLED=1` as
environment if go commands fail.

### Using the `lockService` correctly

The `lockService` can be used to lock resources. This works in-memory and also in a distributed fashion with REDIS.

Ensure to provide proper locking options when using, although in-memory ignores those.

Example:

```shell
# invoked from an endpoint
context := c.Request.Context()

var err error
var lock appLock

if lock, err = h.lockService.lockWithOptions(context, "TEST-LOCK", withAppLockOptionExpiry(5*time.Minute), withAppLockOptionInfiniteRetries(), withAppLockOptionRetryDelay(5*time.Second)); err != nil {
    _ = c.AbortWithError(errToHttpStatus(err), err)
    return
}
# defer to avoid leakage
defer func(lock appLock) {
    _ = lock.unlock(context)
}(lock)

# simulate long running task
time.Sleep(20 * time.Second)
```

### Git workflow

The main branch is `master`. It's protected and only eligible users can push to it. Merge requests to protected branches
are safe-guarded: they need review or at least a successful pipeline run to be merged.

- Merge request branches should start with `feat/`, `fix/`, or `chore/`
- Merge requests should be squashed and the source branch should be deleted
- Merge request commits should have a meaningful commit message
- Merge request titles should have a meaningful title which is taken as commit message once merged
    - should be prefixed with `feat: ...`, `fix(...): ...`, where the contents inside the bracket should be _one word_
      which topic/component is touched
    - should reflect a breaking change by adding a `!` before the colon, e.g., `fix(deps)!: ...`
    - should include more verbose information in the body of the git commit message (merge request description)

```
feat(security)!: add OpenID Connect authentication

- This adds a new authentication mode called oidc
- This new mode is the default, which might break existing installations
```

- Merge requests should contain an entry inside the `CHANGELOG.md` with a date to provide more information if needed
- Merge requests should contain documentation changes, so that code and documentation stays in sync

### Pipeline workflow

Pipeline runs

* on merge request change (open, new push, ...)
* on protected branches

This means you need to create a merge request to trigger a pipeline run. Without merge request, no build is triggered,
thus your code cannot be merged.

### Release preparation & workflow

Releases are done by triggering the "release" pipeline workflow **manually** on `master`.

This application uses _rolling_ releases, which means that a release resets the `latest` tag and in addition adds a date
as git tag when it was published. The latest tag is automatically replaced to newer git commits once release pipeline
finished.

There's no additional preparation needed. `master` is always the latest working state and should be in "release-able"
state any time.

### Dependency updates

Dependency updates are handled by Renovate using the `renovate.json5` file. The base branch is `master`.

Major updates undergo manual review.
