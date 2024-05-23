# README

Backend for upda - **Up**date **Da**shboard in Go.

The main git repository is hosted at
_[https://git.myservermanager.com/varakh/upda](https://git.myservermanager.com/varakh/upda)_.
Other repositories are mirrors and pull requests, issues, and planning are managed there.

Contributions are very welcome!

[Official documentation](https://git.myservermanager.com/varakh/upda-docs) is hosted in a separate git repository.

## Development & contribution

* Ensure to set the following environment variables for proper debug logs during development

```shell
DEVELOPMENT=true
LOGGING_ENCODING=console
LOGGING_LEVEL=debug
```

* Code guidelines
    * Each entity has its own repository
    * Each entity is only used in repository and service (otherwise, mapping happens)
    * Presenter layer is constructed from the entity, e.g., in REST responses and mapped
    * No entity is directly returned in any REST response
    * All log calls should be handled by `zap.L()`
    * Configuration is bootstrapped via separated `struct` types which are given to the service which need them
    * Error handling
        * Always throw an error with `NewServiceError`
        * Always wrap the cause error with `fmt.Errorf`
        * Forward/bubble up the error directly, when original error is already a `NewServiceError` (most likely internal
          calls)
        * Always abort handler chain with `AbortWithError`
        * Utils can throw any error

Please look into the `_doc/` folder for [OpenAPI specification](./_doc/api.yaml) and a Postman Collection.

### Lock service

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

### Getting started

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
  postgres:16-alpine
  
# redis
docker run --rm --name some-redis \
  -p 6379:6379 \
  redis redis-server --save 60 1 --loglevel warning
```

#### Windows hints

On Windows, you need a valid `gcc`, e.g., https://jmeubank.github.io/tdm-gcc/download/ and add the `\bin` folder to your
path.

For any `go` command you run, ensure that your `PATH` has the `gcc` binary and that you add `CGO_ENABLED=1` as
environment.

### Release

Releases are handled by the SCM platform and pipeline. Creating a **new git tag**, creates a new release in the SCM
platform, uploads produced artifacts to that release and publishes docker images automatically.
**Before** doing so, please ensure that the **commit on `master`** has the **correct version settings** and has been
built successfully:

* Adapt `constants_app.go` and change `Version` to the correct version number
* Adapt `CHANGELOG.md` to reflect changes and ensure a date is properly set in the header, also add a reference link
  in footer (link to scm git tag source)
* Adapt `api.yaml`: `version` attribute must reflect the to be released version
* Adapt `env: VERSION_*` in `.forgejo/workflows/release.yaml`

After the release has been created, ensure to change the following settings for the _next development cycle_:

* Adapt `constants_app.go` and change `Version` to the _next_ version number
* Adapt `CHANGELOG.md` and add an _UNRELEASED_ section
* Adapt `api.yaml`: `version` attribute must reflect the _next_ version number
* Adapt `env: VERSION_*` in `.forgejo/workflows/release.yaml` to _next_ version number

