# AGENT.md - Vibecast Project Guide

## Overview

Vibecast is a lightweight static site hosting platform built with Go and SQLite. Users can deploy static sites (ZIP, PDF, HTML, images, etc.) and share them via URL with optional password protection.

## Project Structure

```
cmd/server/          - Entry point (main.go)
internal/server/     - HTTP handlers, pages (inline HTML/CSS/JS), static serving
internal/db/         - SQLite schema, migrations, models
internal/auth/       - JWT auth middleware
```

- Frontend HTML/CSS/JS is **inline** in `internal/server/pages.go` (~150KB+). Edit carefully.
- DB migrations use `ALTER TABLE ADD COLUMN ... DEFAULT` with errors silently ignored. **Never use DROP TABLE** — updates must preserve existing data.

## Build & Run

```bash
# Build
go build -o bin/vibecast ./cmd/server

# Run (persistent DB — never delete on restart)
./bin/vibecast -addr :18099 -storage ./storage -db ./vibecast.db

# For local testing with cloudflare tunnel
cloudflared tunnel --protocol http2 --url http://localhost:18099
# Note: Use --protocol http2 if QUIC fails (network issues)
```

## Release Process

**Do NOT create releases manually with `gh release create` or `git tag`.**

Releases are triggered via GitHub Actions workflow:

```bash
gh workflow run release.yml
```

- Workflow file: `.github/workflows/release.yml`
- Trigger: `workflow_dispatch` (manual)
- Version is auto-generated as `YYYYMMDD` (Asia/Shanghai timezone), with `-N` suffix if tag exists
- No `v` prefix in version numbers (e.g. `20260722`, not `v20260722`)
- Builds for: linux/amd64, darwin/amd64, darwin/arm64, windows/amd64
- Generates SHA256SUMS checksums
- Creates GitHub release with auto-generated release notes

### Check release status
```bash
gh run list --workflow=release.yml --limit=1
gh release list --limit=3
```

## Key Conventions

- i18n: All UI strings must be added to both `en` and `zh` in the i18n object in pages.go
- CSS: Dashboard and admin have separate `<style>` blocks — changes may need to be applied to both
- Password toggle: Default icon is slashed-eye (hidden state); `togglePwd` swaps to open-eye when visible
- Site detail layout: Two-column grid (`detail-grid`), file list spans full width (`grid-column: 1/-1`)
- `loadFileTree` targets `.detail-grid` (not `.detail-inner`) so file section stays inside the grid

## Common Pitfalls

- **i18n braces**: When adding keys, ensure `en:{...}` closes with single `}` before `,zh:{...}`. Extra `}` causes JS syntax error → blank page.
- **String escaping in Go template literals**: Avoid complex regex/quote escaping in inline onclick handlers. Pass IDs only, fetch data via API instead.
- **CSS replace_all**: `pages.go` has duplicate CSS blocks (dashboard + admin). Use `replace_all=true` or provide unique context.
- **DB path**: Use `/home/tar/` or project directory, NOT `/tmp` (tmpfs causes SQLite I/O errors).
