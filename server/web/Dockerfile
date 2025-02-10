#
# Build image
#
FROM node:22-alpine AS builder
LABEL maintainer="Varakh <varakh@varakh.de>"

RUN mkdir -p /app
WORKDIR /app
ENV PATH=/app/node_modules/.bin:$PATH

# build sources
ADD package.json /app/package.json
ADD pnpm-lock.yaml /app/pnpm-lock.yaml
ADD index.html /app/index.html
ADD vite-env.d.ts /app/vite-env.d.ts
ADD vite.config.ts /app/vite.config.ts
ADD .env /app/.env
ADD .npmrc /app/.npmrc
ADD eslint.config.js /app/eslint.config.js
ADD .prettierrc /app/.prettierrc
ADD .prettierignore /app/.prettierignore
ADD .stylelintrc /app/.stylelintrc
ADD tsconfig.json /app/tsconfig.json

COPY public /app/public
COPY src /app/src

RUN npm install --global pnpm@^10 && \
    pnpm install --frozen-lockfile && \
    pnpm build

#
# Actual image
#
FROM nginx:alpine

ENV NGINX_PORT=80
ENV NGINX_BASE_PATH=/

LABEL maintainer="Varakh <varakh@varakh.de>" \
    description="upda-standalone-webinterface" \
    org.opencontainers.image.authors="Varakh" \
    org.opencontainers.image.vendor="Varakh" \
    org.opencontainers.image.title="upda-standalone-webinterface" \
    org.opencontainers.image.description="upda-standalone-webinterface" \
    org.opencontainers.image.base.name="nginx"

COPY --from=builder /app/build /usr/share/nginx/html
RUN mkdir -p "/usr/share/nginx/app" && \
    mkdir -p "/usr/share/nginx/html/conf" && \
    chown -R nginx:nginx "/usr/share/nginx/app" && \
    chown -R nginx:nginx "/usr/share/nginx/html" && \
    chmod -R g+w "/usr/share/nginx/app" && \
    chmod -R g+w "/usr/share/nginx/html" && \
    apk --update upgrade && \
    apk add --no-cache bash && \
    rm -rf /var/cache/apk/*

COPY docker-nginx.conf.template /etc/nginx/conf.d/default.conf.template
COPY docker-env.sh /usr/share/nginx/app/env.sh
COPY .env /usr/share/nginx/app/.env

WORKDIR /usr/share/nginx/html
EXPOSE ${NGINX_PORT}

ENTRYPOINT ["/usr/share/nginx/app/env.sh"]
CMD ["/bin/bash", "-c", "envsubst '$${NGINX_PORT} $${NGINX_BASE_PATH}' < /etc/nginx/conf.d/default.conf.template > /etc/nginx/conf.d/default.conf && nginx -g 'daemon off;'"]
