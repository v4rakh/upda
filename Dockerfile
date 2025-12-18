#
# Build image
#
FROM alpine:3.23 AS builder
LABEL maintainer="Varakh <varakh@varakh.de>"
ARG VERSION="rolling-oci"

RUN apk --update upgrade && \
    apk add git && \
    apk add go gcc make && \
    apk add nodejs npm && \
    rm -rf /var/cache/apk/*

WORKDIR /app
COPY . .
RUN npm install --global pnpm@^10 && \
    VERSION=${VERSION} CC=gcc make clean dependencies build-web build-server-linux-amd64

#
# Actual image
#
FROM alpine:3.23
LABEL maintainer="Varakh <varakh@varakh.de>" \
    description="upda" \
    org.opencontainers.image.authors="Varakh" \
    org.opencontainers.image.vendor="Varakh" \
    org.opencontainers.image.title="upda" \
    org.opencontainers.image.description="upda" \
    org.opencontainers.image.base.name="alpine:3.23"

ENV USER=appuser
ENV GROUP=appuser
ENV UID=2033
ENV GID=2033

RUN apk --update upgrade && \
    apk add tzdata && \
    rm -rf /var/cache/apk/* && \
    addgroup -S ${GROUP} -g ${GID} && \
    adduser -S ${USER} -G ${GROUP} -u ${UID}

COPY --from=builder /app/bin/upda-linux-amd64 /usr/bin/upda

USER ${USER}

ENV SERVER_PORT=8080
EXPOSE ${SERVER_PORT}
CMD ["/usr/bin/upda", "server", "serve"]
