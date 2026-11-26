# MyPhotoX

Lightweight, self-hosted personal photos and videos web app running in a single SQLite container.

## Features

- Photo and video indexing, upload, import, and folder browsing
- Timeline search with view modes, sort order, year, and month filters
- Albums, favorites, archive, metadata editing, and batch editing
- Video transcoding/playback via FFmpeg (`.mp4`, `.mov`, etc.)
- Anonymous public share links (`/s/...`) and ZIP downloads
- English (`en`) and Vietnamese (`vi`) UI

## Quick Start

```bash
cp .env.example .env
docker compose up -d --build
```

Open `http://localhost:2342` (default credentials in `docker-compose.yml` / `.env`).

## Directories

- `./originals` — main photo and video library
- `./import` — staging folder for manual/automatic imports
- `./storage` — SQLite database, thumbnails, cache, and logs

## Build Architecture

Standalone multi-stage `Dockerfile`:
1. `node:24-alpine` — builds and precompresses the Vue 3 + Vuetify frontend
2. `ubuntu:noble` + official Go toolchain — compiles the `myphotox` Go binary without TensorFlow CGo
3. `ubuntu:noble` runtime — minimal image with `sqlite3`, `ffmpeg`, `exiftool`, `libvips`, and `libheif`

## License

GNU Affero General Public License v3.0 (AGPL-3.0). See [LICENSE](LICENSE).
