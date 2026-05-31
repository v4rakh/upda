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

The most straight forward way to get started is by looking into available commands inside the `Makefile`.

For the full setup, you need the following tools:

- go (see minimum version in `go.mod`)
- pnpm and node (see version constraints in `package.json`)
- make to execute commands of the `Makefile`

Quick start is to run:

1. Copy `.env.development` to `.env`.
2. Start with `make run-server` and make sure the environment variables from `.env` are injected.

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

_upda_ includes a frontend in a monorepo fashion inside `internal/server/web/`. For production (binary and OCI), it's
embedded into the GoLang binary itself.

For _development_, no other steps are required. Simply follow
the [frontend instructions](./internal/server/web/README.md) and start the frontend separately.

If you like to have a look on the _production_ experience, the frontend needs to be build first, and you need to build
the Golang binary with `-tags prod`. How to properly build the frontend, please look into `build-web` of
the `Makefile` (additional `rm -rf` cmd!).

### Tests

There are multiple test targets defined in the `Makefile`

- Go: Default tests skipping integration tests if certain environment variables not provided, see `make test-server`.
- Go: Full tests including integration tests, see `make test-server-coverage`.
- Web: Executed together with Go unit tests.

If you're running on rootless docker or podman, set `DOCKER_HOST=unix:///run/user/1000/podman/podman.sock` (adapt for
your user and/or docker).

### Git workflow

The main branch is `main`. It's protected and only eligible users can push to it. Merge requests to protected branches
are safeguarded: they need review or at least a successful pipeline run to be merged.

- Use conventional commits as commit style and branch naming strategy, e.g., `feat/`, `fix/`, `refactor/`, `chore/`, or
  `ci/`
- **All** merge request commits should have a meaningful commit **title** and **message** stating the **why**
- Use atomic git commits, separate **preparatory** from **functional** commits to speed up review
- Avoid merging trunk back, use `git-rebase`

### Pipeline workflow

Pipeline runs

* on merge request change (open, new push, ...)
* on protected branches

This means you need to create a merge request to trigger a pipeline run. Without merge request, no build is triggered,
thus your code cannot be merged.

### Dependency updates

Dependency updates are handled by Renovate using the `renovate.json5` file. The base branch is `main`.

Major updates undergo manual review.

### Releases

1. Prepare a new MR to trunk with the following changes
    * Adjust and align versions
        * `flake.nix`: `version`
        * `internal/meta/pkg.go`: `Version`
        * `package.json`: `version`
    * Make sure `make clean dependencies checkstyle build test-coverage` is fine
    * Make sure `nix build` is fine (you need `nix` for it, update checksums in `flake.nix` if it fails)
      ```shell
      nix build .#packages.x86_64-linux.default -L
      nix build .#packages.aarch64-linux.default -L
      ```
    * Use `release/` as branch prefix and `release: prepare XYZ` as commit message
2. Merge to trunk
3. Trigger the release job the semantic version which is inside the main trunk
4. Generate changelog and attach it to the release (use `git-cliff`)
5. Pull changes from trunk, prepare a new MR to trunk to prepare next version
    *  Adjust and align versions to the next semantic _patch_ version
    * `flake.nix`: `version`
    * `internal/meta/pkg.go`: `Version`
    * `package.json`: `version`
    * Use `release/` as branch prefix and `release: prepare next cycle...` as commit message
6. Merge to trunk
