# VibeShare

> Build with vibe. Share instantly.

A self-hosted, multi-user static site hosting platform built in pure Go.  
No Nginx, no external web server — the application server handles everything: 
authentication, site management, ZIP deployment, and static file serving.

## Features

- **User System** — Register, login, session-based auth (bcrypt + secure cookies)
- **Multi-Site** — Each user can create multiple sites with custom slugs
- **ZIP Deploy** — Upload a ZIP file, auto-extract and deploy instantly
- **Password Protection** — Optional per-site password gate with session cookies
- **Pure Go File Server** — Custom static file handler with MIME detection, path traversal protection, SPA fallback, and conditional requests
- **Zero Dependencies** — Only SQLite (pure-Go driver) + Go standard library
- **Single Binary** — Build and run, that's it

## Quick Start

```bash
# Build
make build

# Run
./bin/vibeshare

# Or with custom config
./bin/vibeshare --addr :3000 --storage ./data/sites --db ./data/vibeshare.db
```

Then open `http://localhost:8080/dashboard` to register and start deploying.

## Configuration

| Flag | Env Var | Default | Description |
|------|---------|---------|-------------|
| `--addr` | `VIBESHARE_ADDR` | `:8080` | Listen address |
| `--storage` | `VIBESHARE_STORAGE` | `./data/sites` | Site files storage directory |
| `--db` | `VIBESHARE_DB` | `./data/vibeshare.db` | SQLite database path |

## API

### Auth

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/auth/register` | Register (`{email, password}`) |
| POST | `/api/auth/login` | Login (`{email, password}`) |
| POST | `/api/auth/logout` | Logout |
| GET | `/api/auth/me` | Current user |

### Sites (auth required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/sites` | List user's sites |
| POST | `/api/sites` | Create site (`{name, slug, password}`) |
| GET | `/api/sites/{id}` | Get site detail |
| PUT | `/api/sites/{id}` | Update site (`{name, password}`) |
| DELETE | `/api/sites/{id}` | Delete site |
| POST | `/api/sites/{id}/deploy` | Deploy ZIP (multipart `file` field) |

### Static Site Access

- `/s/{slug}/` — Access deployed site
- `/p/{slug}` — Password gate (if site is protected)

## Architecture

```
cmd/server/main.go          — Entry point, CLI flags, graceful shutdown
internal/db/                 — SQLite schema, migrations, data models
internal/auth/               — bcrypt hashing, session tokens, middleware
internal/storage/            — ZIP extraction with path traversal protection
internal/server/             — HTTP handlers, routing, static file serving, UI
```

## Tech Stack

- **Language**: Go 1.23+
- **Database**: SQLite (via modernc.org/sqlite — pure Go, no CGO)
- **Auth**: bcrypt password hashing, random session tokens
- **File Serving**: Custom handler (no http.FileServer) for full control over MIME, caching, and security
- **Frontend**: Single-page dashboard, vanilla JS, no build step

## License

MIT
