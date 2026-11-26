##################################################### STAGE 1: FRONTEND ####################################################
FROM node:24-alpine AS frontend-build

WORKDIR /src
COPY package.json package-lock.json ./
COPY frontend/package.json ./frontend/
RUN npm install --ignore-scripts --no-audit --no-fund --no-update-notifier

COPY frontend/ ./frontend/
RUN npm --workspace=frontend run build && \
    node frontend/scripts/precompress.js

##################################################### STAGE 2: BACKEND #####################################################
FROM golang:1.27-alpine AS go-build

ARG BUILD_VERSION=myphotox
ENV CGO_ENABLED="1" \
    CGO_CFLAGS="-D_LARGEFILE64_SOURCE"

RUN apk add --no-cache \
      gcc g++ musl-dev pkgconf vips-dev libheif-dev git

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=frontend-build /src/assets/static /src/assets/static

RUN go mod tidy && \
    mkdir -p /opt/myphotox/bin /opt/myphotox/assets/models /opt/myphotox/assets/locales && \
    go build -trimpath -ldflags "-s -w -X main.version=${BUILD_VERSION}" -o /opt/myphotox/bin/myphotox cmd/photoprism/photoprism.go && \
    cp -a assets/templates assets/static /opt/myphotox/assets/ && \
    cp -a assets/locales/en assets/locales/vi /opt/myphotox/assets/locales/ && \
    if [ -f assets/locales/messages.pot ]; then cp -a assets/locales/messages.pot /opt/myphotox/assets/locales/; fi && \
    ldd /opt/myphotox/bin/myphotox && \
    ! ldd /opt/myphotox/bin/myphotox | grep -qi tensorflow

##################################################### STAGE 3: RUNTIME #####################################################
FROM alpine:3.23

LABEL org.opencontainers.image.title="MyPhotoX"
LABEL org.opencontainers.image.description="Lightweight SQLite self-hosted photos app"
LABEL org.opencontainers.image.source="https://github.com/quanq026/myphotox"

ENV DO_NOT_TRACK=true \
    ORT_DISABLE_TELEMETRY=1 \
    MALLOC_ARENA_MAX=2 \
    HOME="/myphotox" \
    PATH="/opt/myphotox/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin" \
    PHOTOPRISM_ASSETS_PATH="/opt/myphotox/assets" \
    PHOTOPRISM_ORIGINALS_PATH="/myphotox/originals" \
    PHOTOPRISM_IMPORT_PATH="/myphotox/import" \
    PHOTOPRISM_STORAGE_PATH="/myphotox/storage" \
    PHOTOPRISM_BACKUP_PATH="/myphotox/storage/backups" \
    PHOTOPRISM_LOG_FILENAME="/myphotox/storage/myphotox.log" \
    PHOTOPRISM_PID_FILENAME="/myphotox/storage/myphotox.pid" \
    PHOTOPRISM_DATABASE_DRIVER="sqlite" \
    PHOTOPRISM_APP_NAME="MyPhotoX" \
    PHOTOPRISM_SITE_TITLE="MyPhotoX" \
    PHOTOPRISM_SITE_CAPTION="Self-Hosted Photos App" \
    PHOTOPRISM_DEFAULT_LOCALE="en" \
    PHOTOPRISM_HTTP_HOST="0.0.0.0" \
    PHOTOPRISM_HTTP_PORT="2342" \
    PHOTOPRISM_DISABLE_TLS="true" \
    PHOTOPRISM_DISABLE_TENSORFLOW="true" \
    PHOTOPRISM_DISABLE_FACES="true" \
    PHOTOPRISM_DISABLE_CLASSIFICATION="true" \
    PHOTOPRISM_DISABLE_PLACES="true" \
    PHOTOPRISM_DISABLE_WEBDAV="true" \
    PHOTOPRISM_DISABLE_RAW="true" \
    PHOTOPRISM_DISABLE_DARKTABLE="true" \
    PHOTOPRISM_DISABLE_RAWTHERAPEE="true" \
    PHOTOPRISM_DISABLE_IMAGEMAGICK="true" \
    PHOTOPRISM_DISABLE_VECTORS="true" \
    PHOTOPRISM_DISABLE_MCP="true" \
    PHOTOPRISM_DISABLE_OIDC="true" \
    PHOTOPRISM_DISABLE_BACKUPS="true" \
    PHOTOPRISM_DETECT_NSFW="false" \
    PHOTOPRISM_AUTO_INDEX="300" \
    PHOTOPRISM_AUTO_IMPORT="-1"

RUN apk add --no-cache \
      ca-certificates tzdata bash curl sqlite ffmpeg exiftool vips vips-heif vips-jxl libheif && \
    mkdir -p /myphotox/originals /myphotox/import /myphotox/storage

COPY --from=go-build /opt/myphotox /opt/myphotox

RUN mkdir -p /opt/myphotox/assets/models && /opt/myphotox/bin/myphotox version

WORKDIR /myphotox
EXPOSE 2342

HEALTHCHECK --interval=15s --timeout=5s --retries=5 --start-period=20s \
  CMD curl -fsS http://127.0.0.1:2342/api/v1/status || exit 1

ENTRYPOINT ["/opt/myphotox/bin/myphotox"]
CMD ["start"]
