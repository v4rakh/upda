BIN_DIR = $(shell pwd)/bin

clean:
	rm -rf ${BIN_DIR}

dependencies:
	GO111MODULE=on go mod download

ci: clean dependencies test-ci build-server-ci build-cli-ci

build-server-ci: build-server-linux-amd64

# server requires CGO_ENABLED=1 for go-sqlite
build-server-freebsd-amd64:
	CGO_ENABLED=1 GO111MODULE=on GOOS=freebsd GOARCH=amd64 go build -o ${BIN_DIR}/upda-server-freebsd-amd64 cmd/server.go
build-server-freebsd-arm64:
	CGO_ENABLED=1 GO111MODULE=on GOOS=freebsd GOARCH=arm64 go build -o ${BIN_DIR}/upda-server-freebsd-arm64 cmd/server.go
build-server-darwin-amd64:
	CGO_ENABLED=1 GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -o ${BIN_DIR}/upda-server-darwin-amd64 cmd/server.go
build-server-darwin-arm64:
	CGO_ENABLED=1 GO111MODULE=on GOOS=darwin GOARCH=arm64 go build -o ${BIN_DIR}/upda-server-darwin-arm64 cmd/server.go
build-server-linux-amd64:
	CGO_ENABLED=1 GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o ${BIN_DIR}/upda-server-linux-amd64 cmd/server.go
build-server-linux-arm64:
	CGO_ENABLED=1 GO111MODULE=on GOOS=linux GOARCH=arm64 go build -o ${BIN_DIR}/upda-server-linux-arm64 cmd/server.go
build-server-windows-amd64:
	CGO_ENABLED=1 GO111MODULE=on GOOS=windows GOARCH=amd64 go build -o ${BIN_DIR}/upda-server-windows-amd64 cmd/server.go
build-server-windows-arm64:
	CGO_ENABLED=1 GO111MODULE=on GOOS=windows GOARCH=arm64 go build -o ${BIN_DIR}/upda-server-windows-arm64 cmd/server.go

# cli does not require CGO_ENABLED=1, cross-platform build possible
build-cli-ci: build-cli-linux-amd64

build-cli-all: build-cli-freebsd-amd64 build-cli-freebsd-arm64 build-cli-darwin-amd64 build-cli-darwin-arm64 build-cli-linux-amd64 build-cli-linux-arm64 build-cli-windows-amd64 build-cli-windows-arm64

build-cli-freebsd-amd64:
	CGO_ENABLED=0 GO111MODULE=on GOOS=freebsd GOARCH=amd64 go build -o ${BIN_DIR}/upda-cli-freebsd-amd64 cmd/cli.go
build-cli-freebsd-arm64:
	CGO_ENABLED=0 GO111MODULE=on GOOS=freebsd GOARCH=arm64 go build -o ${BIN_DIR}/upda-cli-freebsd-arm64 cmd/cli.go
build-cli-darwin-amd64:
	CGO_ENABLED=0 GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -o ${BIN_DIR}/upda-cli-darwin-amd64 cmd/cli.go
build-cli-darwin-arm64:
	CGO_ENABLED=0 GO111MODULE=on GOOS=darwin GOARCH=arm64 go build -o ${BIN_DIR}/upda-cli-darwin-arm64 cmd/cli.go
build-cli-linux-amd64:
	CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o ${BIN_DIR}/upda-cli-linux-amd64 cmd/cli.go
build-cli-linux-arm64:
	CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=arm64 go build -o ${BIN_DIR}/upda-cli-linux-arm64 cmd/cli.go
build-cli-windows-amd64:
	CGO_ENABLED=0 GO111MODULE=on GOOS=windows GOARCH=amd64 go build -o ${BIN_DIR}/upda-cli-windows-amd64 cmd/cli.go
build-cli-windows-arm64:
	CGO_ENABLED=0 GO111MODULE=on GOOS=windows GOARCH=arm64 go build -o ${BIN_DIR}/upda-cli-windows-arm64 cmd/cli.go

test: test-server test-cli test-util

test-ci: test

test-server:
	GO111MODULE=on go test ./server/...
test-cli:
	GO111MODULE=on go test ./terminal/...
test-util:
	GO111MODULE=on go test ./util/...