VERSION ?= rolling
LDFLAGS := -X 'git.myservermanager.com/varakh/upda/internal/meta.Version=$(VERSION)'

GO ?= GO111MODULE=on CGO_ENABLED=0 go
GO_TEST ?= CGO_ENABLED=1 go
GOOS ?= $(shell $(GO) version | cut -d' ' -f4 | cut -d'/' -f1)
GOARCH ?= $(shell $(GO) version | cut -d' ' -f4 | cut -d'/' -f2)

CMD_GO_FILES ?= ./cmd/upda/main.go

export GO111MODULE=on

PNPM ?= pnpm
GOSEC ?= gosec
GRYPE ?= grype

BIN_DIR = $(shell pwd)/bin
TEST_DIR = "$(shell pwd)/coverage"
WEB_DIR = $(shell pwd)/internal/server/web
WEB_BUILD_DIR = $(shell pwd)/internal/server/web/build
WEB_NODE_DIR = $(shell pwd)/internal/server/web/node_modules
WEB_COVERAGE_DIR = $(shell pwd)/internal/server/web/coverage
WEB_CI_DIR = $(shell pwd)/internal/server/web/ci/*.xml

clean: clean-server clean-web
clean-server:
	@rm -rf ${BIN_DIR}
	@rm -rf ${TEST_DIR}
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
	$(GO) mod download
dependencies-web:
	cd ${WEB_DIR}; $(PNPM) install --frozen-lockfile

checkstyle: checkstyle-web checkstyle-server
checkstyle-server:
	$(GO) vet ./...
checkstyle-web:
	cd ${WEB_DIR}; $(PNPM) run checkstyle

generate: generate-server
generate-server:
	$(GO) generate ./...

test: test-web test-server
test-web:
	cd ${WEB_DIR}; $(PNPM) run test
test-server:
	@echo "⚠️ Skipping tests requiring environment variables, use 'make test-coverage' to run full test suite."
	$(GO_TEST) test -race -shuffle on -v ./...

test-coverage: test-web-coverage test-server-coverage
test-server-coverage:
	@make clean
	@mkdir -p ${TEST_DIR}
	$(GO_TEST) build -cover -o ${BIN_DIR}/testapp ./cmd/upda
	TEST_DIR=${TEST_DIR} TEST_BINARY=${BIN_DIR}/testapp $(GO_TEST) test -coverprofile ${TEST_DIR}/coverage.unit.out -race -shuffle on -v ./...
	$(GO_TEST) tool covdata textfmt -i=${TEST_DIR} -o=${TEST_DIR}/coverage.integration.out
	@cat ${TEST_DIR}/coverage.unit.out > ${TEST_DIR}/coverage.out
	@tail -n +2 ${TEST_DIR}/coverage.integration.out >> ${TEST_DIR}/coverage.out
	@grep -v -E "_generated.go|_test.go|main.go" ${TEST_DIR}/coverage.out > ${TEST_DIR}/coverage.final.out || true
	$(GO_TEST) tool cover -func=${TEST_DIR}/coverage.final.out
test-web-coverage:
	cd ${WEB_DIR}; $(PNPM) run test:coverage

run-server:
	$(GO) run -ldflags="$(LDFLAGS)" ${CMD_GO_FILES} server serve
run-web:
	cd ${WEB_DIR}; $(PNPM) start

audit: audit-web audit-server

audit-web:
	cd ${WEB_DIR}; $(PNPM) audit -P --audit-level high;

audit-server:
	$(GOSEC) -quiet -sort -severity medium -confidence high ./...

scan:
	@NO_COLOR=1 $(GRYPE) -v -o table --file bin/grype.txt --fail-on critical bin/ || true
	@cat ./bin/grype.txt

build: build-web build-server-all

build-server:
	$(GO) build -tags prod -ldflags="$(LDFLAGS)" -o ${BIN_DIR}/upda-${GOOS}-${GOARCH} ${CMD_GO_FILES}

build-server-all: build-server-freebsd-amd64 build-server-freebsd-arm64 build-server-darwin-amd64 build-server-darwin-arm64 build-server-linux-amd64 build-server-linux-arm64 build-server-windows-amd64 build-server-windows-arm64

build-server-freebsd-amd64:
	GOOS=freebsd GOARCH=amd64 $(GO) build -tags prod -ldflags="$(LDFLAGS)" -o ${BIN_DIR}/upda-freebsd-amd64 ${CMD_GO_FILES}
build-server-freebsd-arm64:
	GOOS=freebsd GOARCH=arm64 $(GO) build -tags prod -ldflags="$(LDFLAGS)" -o ${BIN_DIR}/upda-freebsd-arm64 ${CMD_GO_FILES}
build-server-darwin-amd64:
	GOOS=darwin GOARCH=amd64 $(GO) build -tags prod -ldflags="$(LDFLAGS)" -o ${BIN_DIR}/upda-darwin-amd64 ${CMD_GO_FILES}
build-server-darwin-arm64:
	GOOS=darwin GOARCH=arm64 $(GO) build -tags prod -ldflags="$(LDFLAGS)" -o ${BIN_DIR}/upda-darwin-arm64 ${CMD_GO_FILES}
build-server-linux-amd64:
	GOOS=linux GOARCH=amd64 $(GO) build -tags prod -ldflags="$(LDFLAGS)" -o ${BIN_DIR}/upda-linux-amd64 ${CMD_GO_FILES}
build-server-linux-arm64:
	GOOS=linux GOARCH=arm64 $(GO) build -tags prod -ldflags="$(LDFLAGS)" -o ${BIN_DIR}/upda-linux-arm64 ${CMD_GO_FILES}
build-server-windows-amd64:
	GOOS=windows GOARCH=amd64 $(GO) build -tags prod -ldflags="$(LDFLAGS)" -o ${BIN_DIR}/upda-windows-amd64 ${CMD_GO_FILES}
build-server-windows-arm64:
	GOOS=windows GOARCH=arm64 $(GO) build -tags prod -ldflags="$(LDFLAGS)" -o ${BIN_DIR}/upda-windows-arm64 ${CMD_GO_FILES}

# remove built build/conf directory to be served live from the running binary
build-web:
	cd ${WEB_DIR}; $(PNPM) run build; rm -rf build/conf

.PHONY: clean clean-server clean-web dependencies dependencies-server dependencies-web generate generate-server build build-server build-server-all build-server-darwin-amd64 build-server-darwin-arm64 build-server-freebsd-amd64 build-server-freebsd-arm64 build-server-linux-amd64 build-server-linux-arm64 build-server-windows-amd64 build-server-windows-arm64 build-web checkstyle checkstyle-server checkstyle-web audit audit-web audit-server scan test test-server test-web test-coverage test-server-coverage test-web-coverage run run-server run-web