BIN_DIR = $(shell pwd)/bin
WEB_DIR = $(shell pwd)/server/web
WEB_BUILD_DIR = $(shell pwd)/server/web/build
WEB_NODE_DIR = $(shell pwd)/server/web/node_modules
WEB_COVERAGE_DIR = $(shell pwd)/server/web/coverage
WEB_CI_DIR = $(shell pwd)/server/web/ci/*.xml

VERSION ?= rolling
LDFLAGS := -X 'git.myservermanager.com/varakh/upda/commons.Version=$(VERSION)'

# cleanup steps
clean: clean-server clean-web
clean-server:
	rm -rf ${BIN_DIR}
clean-web:
	rm -rf ${WEB_BUILD_DIR} ${WEB_NODE_DIR} ${WEB_COVERAGE_DIR} ${WEB_CI_DIR}

# dependencies steps
dependencies: dependencies-web dependencies-server
dependencies-server:
	GO111MODULE=on go mod download
dependencies-web:
	cd ${WEB_DIR}; pnpm install --frozen-lockfile

# checkstyle steps
checkstyle: checkstyle-web checkstyle-server
checkstyle-server:
	go vet ./...
checkstyle-web:
	cd ${WEB_DIR}; pnpm run checkstyle

# test steps
test: test-web test-server
test-server:
	go test -race ./...
test-web:
	cd ${WEB_DIR}; pnpm run test:ci

# build steps

build-server-all: build-server-freebsd-amd64 build-server-freebsd-arm64 build-server-darwin-amd64 build-server-darwin-arm64 build-server-linux-amd64 build-server-linux-arm64 build-server-windows-amd64 build-server-windows-arm64

build-server-freebsd-amd64:
	CGO_ENABLED=0 GO111MODULE=on GOOS=freebsd GOARCH=amd64 go build -tags prod -ldflags="$(LDFLAGS)" -o ${BIN_DIR}/upda-server-freebsd-amd64 cmd/server/main.go
build-server-freebsd-arm64:
	CGO_ENABLED=0 GO111MODULE=on GOOS=freebsd GOARCH=arm64 go build -tags prod -ldflags="$(LDFLAGS)" -o ${BIN_DIR}/upda-server-freebsd-arm64 cmd/server/main.go
build-server-darwin-amd64:
	CGO_ENABLED=0 GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -tags prod -ldflags="$(LDFLAGS)" -o ${BIN_DIR}/upda-server-darwin-amd64 cmd/server/main.go
build-server-darwin-arm64:
	CGO_ENABLED=0 GO111MODULE=on GOOS=darwin GOARCH=arm64 go build -tags prod -ldflags="$(LDFLAGS)" -o ${BIN_DIR}/upda-server-darwin-arm64 cmd/server/main.go
build-server-linux-amd64:
	CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=amd64 go build -tags prod -ldflags="$(LDFLAGS)" -o ${BIN_DIR}/upda-server-linux-amd64 cmd/server/main.go
build-server-linux-arm64:
	CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=arm64 go build -tags prod -ldflags="$(LDFLAGS)" -o ${BIN_DIR}/upda-server-linux-arm64 cmd/server/main.go
build-server-windows-amd64:
	CGO_ENABLED=0 GO111MODULE=on GOOS=windows GOARCH=amd64 go build -tags prod -ldflags="$(LDFLAGS)" -o ${BIN_DIR}/upda-server-windows-amd64 cmd/server/main.go
build-server-windows-arm64:
	CGO_ENABLED=0 GO111MODULE=on GOOS=windows GOARCH=arm64 go build -tags prod -ldflags="$(LDFLAGS)" -o ${BIN_DIR}/upda-server-windows-arm64 cmd/server/main.go

build-cli-all: build-cli-freebsd-amd64 build-cli-freebsd-arm64 build-cli-darwin-amd64 build-cli-darwin-arm64 build-cli-linux-amd64 build-cli-linux-arm64 build-cli-windows-amd64 build-cli-windows-arm64

build-cli-freebsd-amd64:
	CGO_ENABLED=0 GO111MODULE=on GOOS=freebsd GOARCH=amd64 go build -tags prod -ldflags="$(LDFLAGS)" -o ${BIN_DIR}/upda-cli-freebsd-amd64 cmd/cli/main.go
build-cli-freebsd-arm64:
	CGO_ENABLED=0 GO111MODULE=on GOOS=freebsd GOARCH=arm64 go build -tags prod -ldflags="$(LDFLAGS)" -o ${BIN_DIR}/upda-cli-freebsd-arm64 cmd/cli/main.go
build-cli-darwin-amd64:
	CGO_ENABLED=0 GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -tags prod -ldflags="$(LDFLAGS)" -o ${BIN_DIR}/upda-cli-darwin-amd64 cmd/cli/main.go
build-cli-darwin-arm64:
	CGO_ENABLED=0 GO111MODULE=on GOOS=darwin GOARCH=arm64 go build -tags prod -ldflags="$(LDFLAGS)" -o ${BIN_DIR}/upda-cli-darwin-arm64 cmd/cli/main.go
build-cli-linux-amd64:
	CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=amd64 go build -tags prod -ldflags="$(LDFLAGS)" -o ${BIN_DIR}/upda-cli-linux-amd64 cmd/cli/main.go
build-cli-linux-arm64:
	CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=arm64 go build -tags prod -ldflags="$(LDFLAGS)" -o ${BIN_DIR}/upda-cli-linux-arm64 cmd/cli/main.go
build-cli-windows-amd64:
	CGO_ENABLED=0 GO111MODULE=on GOOS=windows GOARCH=amd64 go build -tags prod -ldflags="$(LDFLAGS)" -o ${BIN_DIR}/upda-cli-windows-amd64 cmd/cli/main.go
build-cli-windows-arm64:
	CGO_ENABLED=0 GO111MODULE=on GOOS=windows GOARCH=arm64 go build -tags prod -ldflags="$(LDFLAGS)" -o ${BIN_DIR}/upda-cli-windows-arm64 cmd/cli/main.go

# remove built build/conf directory to be served live from the running binary
build-web:
	cd ${WEB_DIR}; pnpm run build; rm -rf build/conf

# ci
clean-ci: clean
dependencies-ci: dependencies
checkstyle-ci: checkstyle
test-ci: test
build-server-ci: build-server-all
build-cli-ci: build-cli-all
build-web-ci: build-web
ci: clean-ci dependencies-ci checkstyle-ci test-ci build-web-ci build-server-ci build-cli-ci
ci-oci: clean-ci dependencies-ci build-web-ci build-server-linux-amd64 build-cli-linux-amd64
