# Deployment

_upda_ is a server application which embeds a webinterface directly in its binary form. This makes it easy to deploy
natively. In addition, a _upda_ docker image is provided.

_upda-cli_ which is an optional command-line helper to quickly invoke webhooks or list tracked updates in
your is also embedded into the docker image, but can also be downloaded for your operating system.

The following sections outline how to deploy _upda_ in a containerized environment and also natively.

> You need to ensure to properly replace `$GENERATED_SECURE_SECRET_OF_SIZE_32_CHARS`, `$SECURE_RANDOM_DATABASE_PASSWORD`
> and `$SECURE_ADMIN_PASSWORD` with secure randomly generated passwords.

## Container

The following outlines use of `docker-compose`, but plain docker or podman (use pod for connection to database and/or
redis) work as well. Basic usage of those tools is not explained here. Please refer to other online sources.

The default image user is `appuser` (`uid=2033`) and group is `appgroup` (`gid=2033`).

```yaml
networks:
    internal:
        external: false
        driver: bridge
        driver_opts:
            com.docker.network.bridge.name: br-upda

services:
    app:
        container_name: upda_app
        image: git.myservermanager.com/varakh/upda:latest
        environment:
            - WEB_API_URL=https://upda.domain.tld
            - WEB_TITLE=upda
            - TZ=Europe/Berlin
            - DB_POSTGRES_TZ=Europe/Berlin
            - DB_POSTGRES_HOST=db
            - DB_POSTGRES_PORT=5432
            - DB_POSTGRES_NAME=upda
            - DB_POSTGRES_USER=upda
            - DB_POSTGRES_PASSWORD=$SECURE_RANDOM_DATABASE_PASSWORD
            - BASIC_AUTH_USER=admin
            - BASIC_AUTH_PASSWORD=$SECURE_ADMIN_PASSWORD
            # generate 32 character long secret, e.g., with "openssl rand -hex 16"
            - SECRET=$GENERATED_SECURE_SECRET_OF_SIZE_32_CHARS
        restart: unless-stopped
        networks:
            - internal
        ports:
            - "127.0.0.1:8080:8080"
        depends_on:
            - db

    db:
        container_name: upda_db
        image: docker.io/postgres:17
        restart: unless-stopped
        environment:
            - POSTGRES_USER=upda
            - POSTGRES_PASSWORD=$SECURE_RANDOM_DATABASE_PASSWORD
            - POSTGRES_DB=upda
        networks:
            - internal
        volumes:
            - upda-db-vol:/var/lib/postgresql/data

volumes:
    upda-db-vol:
        external: false
```

## High availability

For high availability, add [REDIS](https://redis.io/) to support proper distributed locking.

Make changes to your docker-compose deployment similar to the following:

```yaml
    # the existing app service - add these changes to all instances, so they all use the same redis instance
    # make sure that all of them can connect to the redis instance
    # ...
    app:
        environment:
            - LOCK_REDIS_ENABLED=true
            - LOCK_REDIS_URL=redis://redis:6379/0

    # the new redis service            
    redis:
        container_name: upda_redis
        image: redis
        restart: unless-stopped
        networks:
            - internal
        volumes:
            - redis-data-vol:/var/redis/data
        # optionally expose port depending on your setup
        ports:
            - "127.0.0.1:6379:6379"

    volumes:
        redis-data-vol:
            external: false
```

In addition, you need a proper load balancer which routes incoming traffic to all of your instances.

Furthermore, you can also decide to have the frontend in a high-availability setup.

## Reverse proxy

You may want to use a proxy in front of them on your host, e.g., nginx. Here's a configuration snippet which should do
the work.

The UI and API (backend/server) is reachable through the same domain, e.g., `https://upda.domain.tld`. In addition,
Let's Encrypt is used for transport encryption.

```shell
server {
    listen 443 ssl http2;
    ssl_certificate /etc/letsencrypt/live/upda.domain.tld/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/upda.domain.tld/privkey.pem;
    
    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## Native

Deploying _upda_ natively is also possible.

First, download the binary for your operating system, make it executable, e.g., with `chmod +x upda-server`, then
place it into the directory you want, e.g., `/usr/local/bin`. Afterward, run the binary with `./upda-server`.

For a native deployment, it's recommended to use a service orchestrator like systemd on UNIX/Linux machines. Here's an
example file `upda.service` which you can put into `/etc/systemd/system` or alike, then reload available systemd
services with `systemctl daemon-reload` to make it available.

Make sure that your `/etc/upda.conf` has all necessary `DB_POSTGRES_*` environment variables set to configure the
database connection.

Afterward, start and enable it with `systemctl enable --now upda.service`.

```shell
[Unit]
Description=upda
After=network.target

[Service]
Type=simple
# Using a dynamic user drops privileges and sets some security defaults
# See https://www.freedesktop.org/software/systemd/man/latest/systemd.exec.html
DynamicUser=yes
# All environment variables for upda can be put into this file
# upda picks them up (on each restart)
EnvironmentFile=/etc/upda.conf
# Requires upda' binary to be installed at this location, e.g., via package manager or copying it over manually
ExecStart=/usr/local/bin/upda-server
```

For a full set of available configuration, look into the [Configuration](./Configuration.md) section. Furthermore,
it's recommended to set up proper [Monitoring](./Monitoring.md).