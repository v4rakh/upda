#
# Build image
#
FROM alpine:3.20 AS builder
LABEL maintainer="Varakh <varakh@varakh.de>"

RUN apk --update upgrade && \
    apk add go gcc make sqlite && \
    apk add nodejs npm && \
    # See https://stackoverflow.com/questions/34729748/installed-go-binary-not-found-in-path-on-alpine-linux-docker
    mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 && \
    rm -rf /var/cache/apk/*

WORKDIR /app
COPY . .
RUN rm -rf bin/ && \
    npm install --global pnpm@^9 && \
    CC=gcc make ci

#
# Actual image
#
FROM alpine:3.20
LABEL maintainer="Varakh <varakh@varakh.de>" \
    description="upda" \
    org.opencontainers.image.authors="Varakh" \
    org.opencontainers.image.vendor="Varakh" \
    org.opencontainers.image.title="upda" \
    org.opencontainers.image.description="upda" \
    org.opencontainers.image.base.name="alpine:3.20"

ENV USER=appuser \
    GROUP=appuser \
    UID=2033 \
    GID=2033

RUN apk --update upgrade && \
    apk add sqlite tzdata && \
    rm -rf /var/cache/apk/* && \
    addgroup -S ${GROUP} -g ${GID} && \
    adduser -S ${USER} -G ${GROUP} -u ${UID}

COPY --from=builder /app/bin/upda-cli-linux-amd64 /usr/bin/upda-cli
COPY --from=builder /app/bin/upda-server-linux-amd64 /usr/bin/upda-server

USER ${USER}

ENV SERVER_PORT 8080
EXPOSE ${SERVER_PORT}
CMD ["/usr/bin/upda-server"]
