# MyPhotoX

Lightweight, self-hosted personal photos and videos web app running in a single Alpine SQLite container (~300 MB image, ~30 MB idle RAM).

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

Open `http://localhost:2342`. Image builds run on GitHub Actions and publish to `ghcr.io/quanq026/myphotox:latest`.

## Directories

- `./originals` — photo and video library
- `./import` — staging folder
- `./storage` — SQLite database, thumbnails, cache, logs
