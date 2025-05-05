VERSION ?= rolling
LDFLAGS := -X 'git.myservermanager.com/varakh/upda/internal/server/constant.AppVersion=$(VERSION)'

GO ?= GO111MODULE=on CGO_ENABLED=0 go
GO_TEST ?= CGO_ENABLED=1 go
GOOS ?= $(shell $(GO) version | cut -d' ' -f4 | cut -d'/' -f1)
GOARCH ?= $(shell $(GO) version | cut -d' ' -f4 | cut -d'/' -f2)

CMD_GO_FILES ?= ./cmd/upda/main.go

export GO111MODULE=on

PNPM ?= pnpm

BIN_DIR = $(shell pwd)/bin
WEB_DIR = $(shell pwd)/internal/server/web
WEB_BUILD_DIR = $(shell pwd)/internal/server/web/build
WEB_NODE_DIR = $(shell pwd)/internal/server/web/node_modules
WEB_COVERAGE_DIR = $(shell pwd)/internal/server/web/coverage
WEB_CI_DIR = $(shell pwd)/internal/server/web/ci/*.xml

clean: clean-server clean-web
clean-server:
	@rm -rf ${BIN_DIR}
	@$(GO) clean -testcache
clean-web:
	@rm -rf ${WEB_BUILD_DIR}
	@rm -rf ${WEB_NODE_DIR}
	@rm -rf ${WEB_COVERAGE_DIR}
	@rm -rf ${WEB_CI_DIR}
	@rm -f ${WEB_DIR}/.eslintcache
	@rm -f ${WEB_DIR}/.stylelintcache

dependencies: dependencies-web dependencies-server
dependencies-server:
	@$(GO) mod download
dependencies-web:
	@cd ${WEB_DIR}; $(PNPM) install --frozen-lockfile

checkstyle: checkstyle-web checkstyle-server
checkstyle-server:
	@$(GO) vet ./...
checkstyle-web:
	@cd ${WEB_DIR}; $(PNPM) run checkstyle

generate: generate-server
generate-server:
	@$(GO) generate ./...

test-unit: test-web-unit test-server-unit
test-server-unit:
	@$(GO_TEST) test -race -shuffle on ./...
test-server-integration:
	@IMAGE=$(image) $(GO_TEST) test -tags=integration ./...
test-web-unit:
	@cd ${WEB_DIR}; $(PNPM) run test:coverage

run-server:
	@$(GO) run -ldflags="$(LDFLAGS)" ${CMD_GO_FILES} server serve
run-web:
	@cd ${WEB_DIR}; $(PNPM) start

build: build-web build-server-all

build-server:
	@$(GO) build -tags prod -ldflags="$(LDFLAGS)" -o ${BIN_DIR}/upda-${GOOS}-${GOARCH} ${CMD_GO_FILES}

build-server-all: build-server-freebsd-amd64 build-server-freebsd-arm64 build-server-darwin-amd64 build-server-darwin-arm64 build-server-linux-amd64 build-server-linux-arm64 build-server-windows-amd64 build-server-windows-arm64

build-server-freebsd-amd64:
	@GOOS=freebsd GOARCH=amd64 $(GO) build -tags prod -ldflags="$(LDFLAGS)" -o ${BIN_DIR}/upda-freebsd-amd64 ${CMD_GO_FILES}
build-server-freebsd-arm64:
	@GOOS=freebsd GOARCH=arm64 $(GO) build -tags prod -ldflags="$(LDFLAGS)" -o ${BIN_DIR}/upda-freebsd-arm64 ${CMD_GO_FILES}
build-server-darwin-amd64:
	@GOOS=darwin GOARCH=amd64 $(GO) build -tags prod -ldflags="$(LDFLAGS)" -o ${BIN_DIR}/upda-darwin-amd64 ${CMD_GO_FILES}
build-server-darwin-arm64:
	@GOOS=darwin GOARCH=arm64 $(GO) build -tags prod -ldflags="$(LDFLAGS)" -o ${BIN_DIR}/upda-darwin-arm64 ${CMD_GO_FILES}
build-server-linux-amd64:
	@GOOS=linux GOARCH=amd64 $(GO) build -tags prod -ldflags="$(LDFLAGS)" -o ${BIN_DIR}/upda-linux-amd64 ${CMD_GO_FILES}
build-server-linux-arm64:
	@GOOS=linux GOARCH=arm64 $(GO) build -tags prod -ldflags="$(LDFLAGS)" -o ${BIN_DIR}/upda-linux-arm64 ${CMD_GO_FILES}
build-server-windows-amd64:
	@GOOS=windows GOARCH=amd64 $(GO) build -tags prod -ldflags="$(LDFLAGS)" -o ${BIN_DIR}/upda-windows-amd64 ${CMD_GO_FILES}
build-server-windows-arm64:
	@GOOS=windows GOARCH=arm64 $(GO) build -tags prod -ldflags="$(LDFLAGS)" -o ${BIN_DIR}/upda-windows-arm64 ${CMD_GO_FILES}

# remove built build/conf directory to be served live from the running binary
build-web:
	@cd ${WEB_DIR}; $(PNPM) run build; rm -rf build/conf

.PHONY: clean clean-server clean-web dependencies dependencies-server dependencies-web checkstyle checkstyle-server checkstyle-web build build-server build-server-all build-server-darwin-amd64 build-server-darwin-arm64 build-server-freebsd-amd64 build-server-freebsd-arm64 build-server-linux-amd64 build-server-linux-arm64 build-server-windows-amd64 build-server-windows-arm64 build-web run run-server run-web