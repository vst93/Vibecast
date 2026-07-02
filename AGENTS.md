# Vibecast — Agent Guide

Pure-Go static site hosting. SQLite + filesystem, no Nginx. Upload ZIP → get live URL.

## Quick Start

```bash
make build              # → bin/vibecast
./bin/vibecast          # listens :8080, data in ./data/
./bin/vibecast -addr :9090 -storage ./data/sites -db ./data/vibecast.db
```

Default admin: `admin@test.com` / `admin123` (auto-created on first run if no users exist — actually first registered user becomes admin).

**Build env:** `GOPROXY=https://goproxy.cn,direct` (set in Makefile already).

## Project Structure

```
cmd/server/main.go          # Entry point, flags, graceful shutdown
internal/
  server/
    server.go               # Server struct, Config, Router(), middleware, slug helpers
    api.go                  # Auth + Sites API handlers (register, login, CRUD, deploy, files)
    admin.go                # Admin API (stats, users, sites, settings, cleanup, system-info)
    static.go               # Static file serving (path safety, MIME, dir listing, SPA fallback)
    password.go             # Password gate page handler (/p/{slug})
    captcha.go              # SVG captcha generation/verification
    update.go               # Self-update system (check/apply/restart)
    pages.go                # All HTML templates (landing, dashboard, admin, password gate) — ~740 lines
    messages.go             # i18n message map + tMsg(r,key) / tStatic(key)
  auth/
    auth.go                 # bcrypt, token gen, RequireAuth/RequireAdmin middleware, GetSessionToken/GetSiteToken
    admin.go                # Admin-specific auth helpers
  db/
    db.go                   # SQLite open + migration (schema, seed settings)
    models.go               # User/Site/Settings structs, all DB queries
  storage/
    zip.go                  # ZIP extraction (dangerous file blocking, zip bomb protection, atomic swap)
```

## Key Architecture Decisions

### Database
- **modernc.org/sqlite** (pure-Go, no CGO). DSN with WAL + busy_timeout + foreign_keys.
- `db.SetMaxOpenConns(1)` — SQLite single-writer.
- Tables: `users`, `sessions`, `sites`, `site_sessions`, `settings`.
- Migration in `db.go:migrate()` — `CREATE TABLE IF NOT EXISTS` + `ALTER TABLE ADD COLUMN` (ignores errors for existing columns).
- Sites store `password` (bcrypt hash) + `password_plain` (for admin display in self-hosted context).

### Routing
- `http.ServeMux` — **no method routing**. Each path → one handler, dispatch via `switch r.Method`.
- Site serving: `/s/{slug}/...` → `staticHandler`.
- Password gate: `/p/{slug}` → `passwordPageHandler`.
- API: `/api/...` (auth via `Authorization: Bearer <token>` header).
- Slugs are **random 12-char** `[a-z0-9]`, never user-input.

### Static File Serving (static.go)
- Path traversal protection: `filepath.Clean` + prefix check.
- Dotfiles rejected (`.env`, `.git`, `.htaccess`).
- MIME: custom map in `extraMimeTypes` + Go built-in + `application/octet-stream` fallback.
- `http.ServeContent` for conditional requests, range, ETag.
- Directory listing (nginx-style) when no `index.html`.
- Password protection: cookie-based site session (`site_token` cookie, 7-day expiry).

### Upload / Deploy
- `POST /api/sites/{id}/deploy` — multipart form, field `file`.
- Currently **ZIP only**. Extract via `storage.ExtractZip()`:
  - Strips common top-level prefix if all entries share one (skips dotfiles for detection).
  - Blocks dangerous extensions (`.php`, `.exe`, `.sh`, etc. — see `blockedExtensions` in zip.go).
  - Zip bomb protection: 500MB total uncompressed, 10000 files max, 100MB per file.
  - Atomic swap: extract to `.tmp` dir → `os.RemoveAll` old → `os.Rename`.
- Frontend uses raw `fetch()` for FormData (bypasses `api()` wrapper which hardcodes JSON content-type).

### Frontend (pages.go)
- All HTML/CSS/JS embedded as Go raw string literals in `pages.go`.
- **⚠️ Go raw strings use backticks — NO JS template literals (`` ` ``) allowed.** Use string concatenation instead.
- **⚠️ Single quotes in inline handlers**: use `\\'` in Go raw strings for JS single-quoted strings (e.g. `onclick="foo(\\'bar\\')"`). NOT `\\\\'` (double-escaped).
- **⚠️ `fmt.Sprintf` with `%` in CSS**: escape as `%%` (e.g. `width:100%%`).
- Light/dark theme: `:root` = light vars, `[data-theme="dark"]` overrides. Anti-flash inline script in `<head>`.
- Bilingual EN/ZH: `i18n` object + `t(key)` function, `localStorage` persistence.
- Dashboard: two-column grid (320px sidebar + 1fr main), `max-width:1280px`.
- Admin: sidebar tab layout (180px nav + main content), 6 tabs.
- **No browser `confirm()`/`alert()`** — custom `customConfirm()` modal.
- `api()` wrapper: JSON content-type, Bearer token, `Accept-Language` header. 401 → clear token + reload (with auth-page guard).
- `siteUrl(u)` wrapper prepends `BASE` for sub-path deployment compatibility.

### Auth
- User auth: Bearer token, 7-day session, stored in `sessions` table.
- Site auth: cookie-based (`site_token`), 7-day, `Path:/s/{slug}/`, `HttpOnly`, `SameSite:Lax`.
- First registered user → auto admin.
- Password gate page uses **fixed English** (not `tMsg()`) — public-facing, international audience.

### Backend i18n
- `tMsg(r, key)` reads `Accept-Language` header → `zh` or `en`.
- `tStatic(key)` returns English (for CLI context, no request).
- Messages in `messages.go` `messageMap`.

### Self-Update System
- CLI: `vibecast update` subcommand (intercepted before `flag.Parse()`).
- API: `GET /api/admin/update/check`, `POST /api/admin/update/apply`, `POST /api/admin/update/restart`.
- GitHub mirror proxies for China: `ghfast.top` → `gh-proxy.com` → direct.
- Binary self-replacement: `os.Rename` → `copyFile` fallback (cross-device).
- SHA256 verification. Concurrency protection via atomic CAS flag.
- Restart: `syscall.Exec` (requires `net.Listen` + `http.Server` refactor for graceful shutdown).

### Version System
- `var version = "dev"` in main.go, overridden by `-ldflags "-X main.version=X.Y.Z"`.
- `--version` CLI flag.
- `GET /api/version` public endpoint.
- Admin UI version badge.

## Common Pitfalls

1. **Stale process after rebuild**: `kill -9` old process before starting new one. Verify via `curl` output, not just build success.
2. **PTY for long-running server**: use `terminal(pty=true)` — bash `&`/`nohup` gets SIGKILL'd.
3. **`patch` tool + Go raw strings with backslashes**: may double-escape `\\'` → `\\\\'`. Verify with `grep -c "\\\\'" file.go`. Use `write_file` for entire regions if patching is problematic.
4. **`patch` `replace_all=true` + multi-line JS**: corrupts output when same pattern appears in multiple templates. Use targeted single-occurrence patches.
5. **FormData uploads**: do NOT set `Content-Type` header — browser sets `multipart/form-data; boundary=...` automatically. Only set `Authorization` + `Accept-Language`.
6. **Hyphenated element IDs**: never reference as JS globals if ID contains hyphen. Always use `document.getElementById()`.
7. **`fmt.Sprintf` placeholder/arg count mismatch**: `go build` doesn't catch this. `go vet` does.
8. **JS validation after subagent edits**: run `node --check` on extracted `<script>` content. `go build` passing ≠ JS valid.

## Build & Test

```bash
make build                    # build
go vet ./...                  # lint (catches fmt.Sprintf mismatches)
node -e "..."                 # validate embedded JS (extract <script> from pages.go)

# Start server (PTY mode to avoid SIGKILL):
./bin/vibecast -addr :8080

# Test API:
curl -s http://localhost:8080/api/version | jq .
curl -s -X POST http://localhost:8080/api/auth/login -H 'Content-Type: application/json' -d '{"email":"admin@test.com","password":"admin123"}' | jq .
```

## Env Vars

| Var | Default | Description |
|-----|---------|-------------|
| `VIBECAST_ADDR` | `:8080` | Listen address |
| `VIBECAST_STORAGE` | `./data/sites` | Site files storage |
| `VIBECAST_DB` | `./data/vibecast.db` | SQLite path |

## Conventions

- Tag: no `v` prefix (e.g. `1.0.0`, not `v1.0.0`).
- Release title: add `Beta` suffix, mark as pre-release.
- Code version: inject via ldflags `-X main.version=X.Y.Z`.
- Branch: `main` is source of truth. Feature branches from latest `main`, delete after merge.
- README: bilingual EN first, ZH second. Update in same commit as feature changes.
