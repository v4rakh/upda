# Deployment

_upda_ is a server application which embeds a webinterface directly in its binary form. This makes it easy to deploy
natively. In addition, a _upda_ docker image is provided to get started quickly.

_upda-cli_ which is an optional commandline helper to quickly invoke webhooks or list tracked updates in
your is also embedded into the docker image, but can also be downloaded for your operating system.

The following sections outline how to deploy _upda_ in a containerized environment and also natively.

## Container

Use one of the provided `docker-compose` examples, edit to your needs. Then issue `docker compose up -d` command and
`docker compose logs -f` to trace the log.

Default image user is `appuser` (`uid=2033`) and group is `appgroup` (`gid=2033`).

The following examples are available

### Postgres

#### docker-compose

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
            - DB_TYPE=postgres
            - DB_POSTGRES_HOST=db
            - DB_POSTGRES_PORT=5432
            - DB_POSTGRES_NAME=upda
            - DB_POSTGRES_USER=upda
            - DB_POSTGRES_PASSWORD=upda
            - BASIC_AUTH_USER=admin
            - BASIC_AUTH_PASSWORD=changeit
            # generate 32 character long secret, e.g., with "openssl rand -hex 16"
            - SECRET=generated-secure-secret-32-chars
        restart: unless-stopped
        networks:
            - internal
        ports:
            - "127.0.0.1:8080:8080"
        depends_on:
            - db

    db:
        container_name: upda_db
        image: postgres:16
        restart: unless-stopped
        environment:
            - POSTGRES_USER=upda
            - POSTGRES_PASSWORD=upda
            - POSTGRES_DB=upda
        networks:
            - internal
        volumes:
            - upda-db-vol:/var/lib/postgresql/data

volumes:
    upda-db-vol:
        external: false
```

### SQLite

#### docker-compose

You can use the following to get it up running quickly via docker compose.

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
            - BASIC_AUTH_USER=admin
            - BASIC_AUTH_PASSWORD=changeit
            # generate 32 character long secret, e.g., with "openssl rand -hex 16"
            - SECRET=generated-secure-secret-32-chars
        restart: unless-stopped
        networks:
            - internal
        volumes:
            - upda-app-vol:/home/appuser
        ports:
            - "127.0.0.1:8080:8080"

volumes:
    upda-app-vol:
        external: false
```

#### Local example

For spinning it up **locally** and without a [reverse proxy](#reverse-proxy), you can use the following simple `docker`
commands.

Make sure to adapt `DOMAIN` and pipe in your device IP address (LAN), e.g., `192.168.1.42`.

```shell
# create volume
docker volume create upda-app-vol

# run locally binding to your LAN IP address
docker run --name upda_app \
  -p 8080:8080 \
  -e TZ=Europe/Berlin \
  -e WEB_API_URL=http://192.168.1.42:8080 \
  -e BASIC_AUTH_USER=admin \
  -e BASIC_AUTH_PASSWORD=changeit \
  -v upda-app-vol:/home/appuser \
  varakh/upda:latest
```

## High availability

For high availability, pick the [Postgres setup](#postgres) and add [REDIS](https://redis.io/) to support proper
distributed locking.

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
        # other already defined volumes
        # ...
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

Native deployment is also possible.

Download the binary for your operating system. Next, use the binary or execute it locally.

See the provided systemd service example `upda.service` to deploy on a UNIX/Linux machine.

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