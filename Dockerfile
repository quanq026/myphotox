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
FROM ubuntu:noble AS go-build

ENV DEBIAN_FRONTEND="noninteractive" \
    PATH="/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin" \
    CGO_ENABLED="1"

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
      ca-certificates curl git gcc g++ pkg-config libvips-dev libheif-dev && \
    rm -rf /var/lib/apt/lists/* && \
    GO_VER=$(curl -fsSL "https://go.dev/VERSION?m=text" | head -n 1) && \
    curl -fsSL "https://go.dev/dl/${GO_VER}.linux-amd64.tar.gz" | tar --no-same-owner -xz -C /usr/local && \
    go version

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=frontend-build /src/assets/static /src/assets/static

RUN sed -i 's/\r$//' scripts/build.sh scripts/dist/arch.sh && \
    go mod tidy && \
    mkdir -p /opt/myphotox/bin /opt/myphotox/lib /opt/myphotox/assets/models && \
    ./scripts/build.sh prod /opt/myphotox/bin/myphotox

RUN cp -a assets/templates assets/static /opt/myphotox/assets/ && \
    mkdir -p /opt/myphotox/assets/locales && \
    cp -a assets/locales/en assets/locales/vi /opt/myphotox/assets/locales/ && \
    if [ -f assets/locales/messages.pot ]; then cp -a assets/locales/messages.pot /opt/myphotox/assets/locales/; fi

RUN LDD_OUT=$(ldd /opt/myphotox/bin/myphotox) && \
    echo "$LDD_OUT" && \
    if echo "$LDD_OUT" | grep -qi tensorflow; then echo "ERROR: libtensorflow must not be linked" && exit 1; fi && \
    echo "$LDD_OUT" | awk '/=> \// {print $3}' | grep -E 'libvips|libheif|libde265|libx265|libaom|libjxl|libsharpyuv' | while read -r lib; do \
      cp -Lv "$lib" /opt/myphotox/lib/ || true; \
    done

##################################################### STAGE 3: RUNTIME #####################################################
FROM ubuntu:noble

LABEL org.opencontainers.image.title="MyPhotoX"
LABEL org.opencontainers.image.description="Lightweight SQLite single-container self-hosted photos app"

ENV DEBIAN_FRONTEND="noninteractive" \
    DO_NOT_TRACK=true \
    ORT_DISABLE_TELEMETRY=1 \
    MALLOC_ARENA_MAX=2 \
    HOME="/myphotox" \
    PATH="/opt/myphotox/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin" \
    LD_LIBRARY_PATH="/opt/myphotox/lib:/usr/local/lib:/usr/lib" \
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

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
      ca-certificates tzdata bash curl sqlite3 ffmpeg exiftool libvips42t64 libheif1 && \
    rm -rf /var/lib/apt/lists/* && \
    mkdir -p /myphotox/originals /myphotox/import /myphotox/storage

COPY --from=go-build /opt/myphotox /opt/myphotox

RUN mkdir -p /opt/myphotox/assets/models && /opt/myphotox/bin/myphotox version

WORKDIR /myphotox
EXPOSE 2342

HEALTHCHECK --interval=15s --timeout=5s --retries=5 --start-period=20s \
  CMD curl -fsS http://127.0.0.1:2342/api/v1/status || exit 1

ENTRYPOINT ["/opt/myphotox/bin/myphotox"]
CMD ["start"]
