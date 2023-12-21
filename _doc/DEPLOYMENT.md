# Deployment

Use one of the provided `docker-compose` examples, edit to your needs. Then issue `docker compose up` command.

All applications should be up and running.

As of now, the web interface and the server comes as different container images.

Default image user is `appuser` (`uid=2033`) and group is `appgroup` (`gid=2033`).

The following examples are available

## Postgres

```yaml
version: '3.9'

networks:
  internal:
    external: false
    driver: bridge
    driver_opts:
      com.docker.network.bridge.name: br-upda

services:
  ui:
    container_name: upda_ui
    image: git.myservermanager.com/varakh/upda-ui:latest
    environment:
      - VITE_API_URL=https://upda.domain.tld/api/v1/
      - VITE_APP_TITLE=upda
      - VITE_APP_DESCRIPTION=upda
    restart: unless-stopped
    networks:
      - internal
    ports:
      - "127.0.0.1:8181:80"
    depends_on:
      - api

  api:
    container_name: upda_api
    image: git.myservermanager.com/varakh/upda:latest
    environment:
      - TZ=Europe/Berlin
      - DB_POSTGRES_TZ=Europe/Berlin
      - DB_TYPE=postgres
      - DB_POSTGRES_HOST=db
      - DB_POSTGRES_PORT=5432
      - DB_POSTGRES_NAME=upda
      - DB_POSTGRES_USER=upda
      - DB_POSTGRES_PASSWORD=upda
      - ADMIN_USER=admin
      - ADMIN_PASSWORD=changeit
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

## SQLite

```yaml
version: '3.9'

networks:
  internal:
    external: false
    driver: bridge
    driver_opts:
      com.docker.network.bridge.name: br-upda

services:
  ui:
    container_name: upda_ui
    image: git.myservermanager.com/varakh/upda-ui:latest
    environment:
      - VITE_API_URL=https://upda.domain.tld/api/v1/
      - VITE_APP_TITLE=upda
      - VITE_APP_DESCRIPTION=upda
    restart: unless-stopped
    networks:
      - internal
    ports:
      - "127.0.0.1:8181:80"
    depends_on:
      - api

  api:
    container_name: upda_api
    image: git.myservermanager.com/varakh/upda:latest
    environment:
      - TZ=Europe/Berlin
      - DB_SQLITE_FILE=/data/upda.db
      - ADMIN_USER=admin
      - ADMIN_PASSWORD=changeit
    restart: unless-stopped
    networks:
      - internal
    volumes:
      - upda-app-vol:/data
    ports:
      - "127.0.0.1:8080:8080"

volumes:
  upda-app-vol:
    external: false
```

### Reverse proxy

You may want to use a proxy in front of them on your host, e.g., nginx. Here's a configuration snippet which should do
the work.

The UI and API is reachable through the same domain, e.g., `https://upda.domain.tld`. In addition, Let's Encrypt is used
for transport encryption.

```shell
server {
    listen 443 ssl http2;
    ssl_certificate /etc/letsencrypt/live/upda.domain.tld/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/upda.domain.tld/privkey.pem;
    
    # ui
    location / {
        proxy_pass http://localhost:8181;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
    
    # api
    location ~* ^/(api)/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
    # metrics
    location ~* ^/metrics {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```
