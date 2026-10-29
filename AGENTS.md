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

- Page markup + inline JS lives in `internal/server/pages.go` (~200KB). Edit carefully.
- All CSS lives in `internal/server/ui_css.go`, layered as: `uiTokensCSS` (palette/type/a11y reset)
  → `uiComponentsCSS` (navbar, buttons, forms, cards, tables, modals, toasts, menus)
  → page layer (`uiDashboardCSS`, `uiAdminCSS`, `uiLandingCSS`, `uiGateCSS`).
  Read `DESIGN.md` before changing UI.
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

- i18n: all UI strings go in both `en` and `zh`. Dashboard and admin each have their own i18n
  object; add keys with `Object.assign(i18n.en,{...});Object.assign(i18n.zh,{...})` after the
  literal to avoid brace mistakes.
- CSS: never append "override layer" styles — add tokens/components in `ui_css.go` instead.
  Dashboard and admin share `uiComponentsCSS`, so a fix there applies to both.
- Fonts: keep `--ff-sans` starting with concrete families (`-apple-system`, `Segoe UI`, `Noto Sans`,
  `Arial`). Do **not** lead with `system-ui`/`ui-sans-serif` — fontconfig resolves them to a serif
  face at weight 600 on Linux and silently turns the whole UI serif.
- Password toggle: slashed eye = hidden; `togglePwd` swaps icon and keeps `aria-label`/`aria-pressed` in sync.
- Site row: `[chevron] [status + name + badges] [url + visits] [Upload] [overflow menu]`.
  Head is `role="button"` with `aria-expanded`; inner links/buttons call `stopPropagation`.
- Site detail: `.detail-clip` (overflow hidden) wraps `.detail-inner` (padding) so the
  `grid-template-rows:0fr→1fr` collapse does not leak padding. Files render into `#files-<id>`.
- Modals: open with `openModal(id)` / close with `closeModal(id)` — they handle focus, Esc and focus return.
- Visitor-facing errors use `visitorErrorPageHTML` (never plain `http.Error` for site visitors).
- Single-file delete: `DELETE /api/sites/{id}/files?path=<relative>`; without `path` the whole site
  directory is cleared. Paths are resolved by `storage.DeleteEntry`, which refuses directories and
  paths outside the site root.
- File list styles (`.file-list/.file-row/.file-del`) live in `uiComponentsCSS` because both the
  dashboard detail and the admin site detail render them.
- Nav height is the `--nav-h` token (52px). Anything sticky (sidebar, admin nav) must derive from it.
- Brand/ GitHub references stay quiet: text footers and the admin system tab only. Repo URL is
  `projectRepoURL`.

## Common Pitfalls

- **i18n braces**: When adding keys, ensure `en:{...}` closes with single `}` before `,zh:{...}`. Extra `}` causes JS syntax error → blank page.
- **String escaping in Go template literals**: Avoid complex regex/quote escaping in inline onclick handlers. Pass IDs only, fetch data via API instead.
- **CSS edits**: styles are in `ui_css.go`; long single-line rules mean you often need unique context or `replace_all`.
- **Stale binary**: after editing Go templates, rebuild (`go build -o bin/vibecast ./cmd/server`) and restart before screenshotting — served HTML comes from the compiled binary.
- **DB path**: Use `/home/tar/` or project directory, NOT `/tmp` (tmpfs causes SQLite I/O errors).
