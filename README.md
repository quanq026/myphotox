# MyPhotoX

Lightweight, self-hosted personal photos and videos web app running in a single Alpine SQLite container (~490 MB image, ~30-37 MB idle RAM).

## Features

- Dedicated **Pictures** & **Videos** views, Albums, Folders, Favorites, and Archive
- Minimal 2-tab **Settings** (`General` + `Account`) with built-in **Rescan Library**
- Upload, metadata edit, batch edit, anonymous public share links (`/s/...`), ZIP downloads
- Video streaming via FFmpeg (`.mp4`, `.mov`, etc.)
- English (`en`) and Vietnamese (`vi`) UI

## Quick Start (Prebuilt GHCR Image)

```bash
cp .env.example .env
docker compose pull
docker compose up -d
```

Open `http://localhost:2342`.

## CI/CD

GitHub Actions (`.github/workflows/ci-cd.yml`) runs on pushes/PRs to `main` and version tags (`v*`):

- `ci-frontend`: Node 24 XSS scan, frontend Vite build & precompress check
- `ci-backend`: Go 1.27 CGo vet, binary build & TensorFlow-free check
- `cd-docker`: Multi-stage Alpine 3.23 Docker image push to `ghcr.io/quanq026/myphotox:latest` (runs after CI passes, skipped on PRs)

## Directories

- `./originals` — photo and video library
- `./import` — staging folder
- `./storage` — SQLite database, thumbnails, cache, logs
