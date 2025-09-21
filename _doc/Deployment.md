# Deployment

_upda_ is a server application which embeds a web interface directly in its binary form (can be disabled). This makes it
easy to deploy natively. Besides native binaries, _upda_ is published as docker image. The `upda` binary can manage the
server and is at the same time a command-line utility to quickly invoke webhooks or list tracked updates in your
instance.

Depending on **how you like to reach _upda_** (reverse proxy setup with a (sub)domain or reverse proxy setup on sub
path of your existing domain), pick one of the below **deployment** options.

The following sections outline how to deploy _upda_ in a containerized environment and also natively.

## Container

In addition to native binaries for your operating system, _upda_ is published as docker image. The default container
image user is `appuser` (`uid=2033`). The group is `appgroup` (`gid=2033`).

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
            - WEB_INTERFACE_API_URL=https://upda.domain.tld/api/v1/
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
`docker-compose.yaml` file accordingly. Let's assume you like to deploy under the `/upda-app` base path, then add
`SERVER_BASE_PATH=/upda-app/`. Make sure to adapt `WEB_INTERFACE_API_URL=https://domain.tld/upda-app/api/v1/`
as well.

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
            - LOCK_REDIS_HOST=redis
            - LOCK_REDIS_PORT=6379
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

The following examples use `nginx` as reverse proxy and Let's Encrypt for transport encryption (https).

You probably want to set the `gzip on;` directive.

### (Sub)Domain

Most likely, this is the default setup and used for the majority of deployments. _upda_ is deployed as a single
container (excluding database) or [natively](#native-deployment).

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

This requires to set `SERVER_BASE_PATH=/upda-app/` as outlined in the deployment section above.

```shell
server {
    # ... your other domain setup

    # forward matching requests to the main upda application
    # make sure that SERVER_BASE_PATH is the same as the path inside the location (except for trailing slash)
    location /upda-meta {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### `robots.txt`

By default, _upda_ does not serve a `robots.txt`. If you like to have it served, you can use the following snippet in
your reverse proxy.

Remember to adapt to your liking before.

```shell
location = /robots.txt {
  add_header Content-Type text/plain;
  return 200 "User-agent: *\nDisallow:\n";
}
```

## Native deployment

Deploying _upda_ natively is also possible.

First, download the binary for your operating system, make it executable, e.g., with `chmod +x upda`, then
place it into the directory you want, e.g., `/usr/local/bin`. Afterward, run the binary with `./upda server serve`.

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
# Requires upda binary to be installed at this location, e.g., via package manager or copying it over manually
ExecStart=/usr/local/bin/upda server serve
```

For a full set of available configuration, look into the [Configuration](./Configuration.md) section. Furthermore,
it's recommended to set up proper [Monitoring](./Monitoring.md).