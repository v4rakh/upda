# Deployment

_upda_ is a server application which embeds a web interface directly in its binary form (can be disabled). This makes it
easy to deploy natively. Besides native binaries, _upda_ is published as docker image. The _upda-cli_ which is an
optional command-line helper to quickly invoke webhooks or list tracked updates in your instance is also embedded into
the docker image, but can also be downloaded for your operating system.

Depending on **how you like to reach _upda_** (reverse proxy setup with a (sub)domain or reverse proxy setup on sub
path of your existing domain), pick one of the below **deployment** options.

Keep in mind that _upda_ does not support sub path deployments with the embedded web interface.

The following sections outline how to deploy _upda_ in a containerized environment and also natively.

## Container

In addition to native binaries for your operating system, _upda_ is published as docker images:

* `upda`: This is the _"server"_ which includes the embedded web interface (**recommended**, but can be disabled). The
  default container image user is `appuser` (`uid=2033`). The group is `appgroup` (`gid=2033`).
* `upda-standalone-webinterface`: This is the standalone web interface to be used for reverse proxy sub path
  deployments **only**.

The following outlines how to deploy using `docker-compose`. If you prefer using plain `docker` or `podman` commands,
make sure to create necessary network (for podman use the _pod_ concept) and volume definitions. Please refer to online
resources if you're not familiar how to translate the docker-compose examples to plain container engine commands.

> You need to ensure to properly replace `$GENERATED_SECURE_SECRET_OF_SIZE_32_CHARS`, `$SECURE_RANDOM_DATABASE_PASSWORD`
> and `$SECURE_ADMIN_PASSWORD` with secure randomly generated passwords.

By default, the following examples only make _upda_ listen on `localhost`/`127.0.0.1` which can be used with
a [reverse proxy](#reverse-proxy) (**recommended**). For testing, you can also remove the local part in the port mapping
directives and expose _upda_ directly (not recommended).

### docker-compose: Deployment on a (sub)domain

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
            - EMBEDDED_WEB_INTERFACE_API_URL=https://upda.domain.tld/api/v1/
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

### docker-compose: Deployment on a sub path

Use the [deployment on a (sub)domain](#docker-compose-deployment-on-a-subdomain) as starting point and adapt your
`docker-compose.yaml` file accordingly. Let's assume you like to deploy under the `/upda-app` base path.

* Disable embedded web interface
* Add another container for the standalone web interface

```yaml
# ... = other defined directives

services:
    app:
        # ...
        environment:
            - EMBEDDED_WEB_INTERFACE_ENABLED=false
            - SERVER_BASE_PATH=/upda-app/
            # ...

    webinterface:
        container_name: upda_webinterface
        image: git.myservermanager.com/varakh/upda-standalone-webinterface:latest
        restart: unless-stopped
        environment:
            - NGINX_BASE_PATH=/upda-app/
            - VITE_API_URL=https://domain.tld/upda-app/api/v1/
        networks:
            - internal
        ports:
            - "127.0.0.1:8081:80"
```

Next, look into the fitting [reverse proxy setup](#reverse-proxy) or decide if you
need [high availability](#high-availability).

## High availability

For high availability, add [REDIS](https://redis.io/) to support proper distributed locking.

Make changes to your docker-compose deployment similar to the following:

```yaml
# ... = other defined directives
services:
    app:
        # ...
        environment:
            - LOCK_REDIS_ENABLED=true
            - LOCK_REDIS_URL=redis://redis:6379/0
            # ...

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
    # ...
```

You need a proper load balancer which routes incoming traffic to all of your instances.

## Reverse proxy

> A reverse proxy setup to proxy _upda_ through a sub path, e.g., `https://domain.tld/upda-app`, is only possible if
> the [embedded web interface is disabled and the standalone web interface](#docker-compose-deployment-on-a-sub-path) is
> deployed in addition.

The following examples use `nginx` as reverse proxy and Let's Encrypt for transport encryption (https).

You probably want to set the `gzip on;` directive.

### (Sub)Domain

Most likely, this is the default setup and used for the majority of deployments. _upda_ is deployed as a single
container (excluding database) or [natively](#native-deployment) utilizing the embedded web interface.

We assume your deployment works, and you like to make it available behind `https://upda.domain.tld`.

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

### Sub path

We assume your deployment works, and you like to make it available behind `https://domain.tld/upda-app`.

This requires to set the `SERVER_BASE_PATH=/upda-app/` for upda and for the web interface `NGINX_BASE_PATH=/upda-app/`
and maybe update `VITE_BASE_PATH` (depends on your routing, but likely default will do).

You can also combine

```shell
server {
    # ... your other domain setup
    
    # forward matching requests to the main upda application 
    location /upda-app/api {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
    
    # forward matching requests to the main upda application
    # comment in if prometheus metrics exporter is disabled
    location /upda-app/metrics {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
    
    # forward matching requests to upda standalone frontend
    location /upda-app {
        proxy_pass http://127.0.0.1:8081;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## Native deployment

Deploying _upda_ natively is also possible.

> A reverse proxy setup to proxy _upda_ through a sub path, e.g., `https://domain.tld/upda`, is **not** possible when
> deploying _upda_ natively. In addition, the _standalone_ web interface is **not** distributed natively. Use the
> containerized setup if you need to deploy _upda_ behind a sub path!

First, download the binary for your operating system, make it executable, e.g., with `chmod +x upda-server`, then
place it into the directory you want, e.g., `/usr/local/bin`. Afterward, run the binary with `./upda-server`.

For a native deployment, it's recommended to use a service orchestrator like systemd on UNIX/Linux machines. Here's an
example file `upda.service` which you can put into `/etc/systemd/system` or alike, then reload available systemd
services with `systemctl daemon-reload` to make it available.

Make sure that your `/etc/upda.conf` has all necessary environment variables, e.g. `DB_POSTGRES_*` and alike set to
configure the database connection.

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