#
# Build image
#
FROM node:24-alpine AS builder-frontend

RUN apk --update upgrade && \
    apk add pnpm && \
    rm -rf /var/cache/apk/*

WORKDIR /src/web
COPY internal/server/web/package*.json ./
COPY internal/server/web/pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile
COPY internal/server/web/ ./
RUN pnpm build

FROM golang:1.26-alpine AS builder-server

# Enable automatic toolchain download for Go 1.25+
ENV GOTOOLCHAIN=auto

WORKDIR /src

# Download dependencies first (better layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source and frontend build
COPY . .
COPY --from=builder-frontend /src/web/build ./internal/server/web/build
RUN CGO_ENABLED=0 GOOS=linux go build -tags prod -trimpath -ldflags="-s -w" -o /upda ./cmd/upda

#
# Actual image
#
FROM gcr.io/distroless/static-debian13:nonroot

# Copy binary
COPY --chown=65532:0 --from=builder-server /upda /usr/local/bin/upda

# Labels
LABEL maintainer="Varakh <varakh@varakh.de>" \
    description="upda" \
    org.opencontainers.image.authors="Varakh" \
    org.opencontainers.image.vendor="Varakh" \
    org.opencontainers.image.title="upda" \
    org.opencontainers.image.description="upda" \
    org.opencontainers.image.base.name="gcr.io/distroless/static-debian13:nonroot" \
    org.opencontainers.image.source="https://git.myservermanager.com/varakh/upda"

# Expose HTTP port
ENV SERVER_PORT=8080
EXPOSE ${SERVER_PORT}

# Run as non-root user (required for OpenShift restricted SCC)
USER 65532:0

# Default command
ENTRYPOINT ["/usr/local/bin/upda"]
CMD ["server", "serve"]
