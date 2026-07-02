# Vibecast

> Build with vibe. Cast instantly.

A self-hosted, multi-user static site hosting platform built in pure Go.  
No Nginx, no external web server — the application server handles everything: 
authentication, site management, ZIP deployment, and static file serving.

[中文说明](#vibecast-中文说明)

---

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
./bin/vibecast

# Or with custom config
./bin/vibecast --addr :3000 --storage ./data/sites --db ./data/vibecast.db
```

Then open `http://localhost:8080/dashboard` to register and start deploying.

## Configuration

| Flag | Env Var | Default | Description |
|------|---------|---------|-------------|
| `--addr` | `VIBECAST_ADDR` | `:8080` | Listen address |
| `--storage` | `VIBECAST_STORAGE` | `./data/sites` | Site files storage directory |
| `--db` | `VIBECAST_DB` | `./data/vibecast.db` | SQLite database path |

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

---

# Vibecast 中文说明

> Build with vibe. Cast instantly.

一个自托管的纯 Go 多用户静态站点托管平台。
不依赖 Nginx 或任何外部 Web Server —— 应用服务器全权处理：
用户认证、站点管理、ZIP 部署和静态文件服务。

## 功能特性

- **用户系统** — 注册、登录、Session 认证（bcrypt 加密 + 安全 Cookie）
- **多站点** — 每个用户可创建多个站点，支持自定义 slug
- **ZIP 部署** — 上传 ZIP 文件，自动解压并即时部署上线
- **密码保护** — 可选的站点级密码门禁，通过 Session Cookie 维持访问状态
- **纯 Go 文件服务器** — 自定义静态文件 Handler，支持 MIME 检测、路径遍历防护、SPA fallback、条件请求
- **零外部依赖** — 仅需 SQLite（纯 Go 驱动）+ Go 标准库
- **单体二进制** — 编译即用，无需其他组件

## 快速开始

```bash
# 编译
make build

# 运行
./bin/vibecast

# 或指定自定义配置
./bin/vibecast --addr :3000 --storage ./data/sites --db ./data/vibecast.db
```

然后打开 `http://localhost:8080/dashboard` 注册并开始部署。

## 配置项

| 参数 | 环境变量 | 默认值 | 说明 |
|------|----------|--------|------|
| `--addr` | `VIBECAST_ADDR` | `:8080` | 监听地址 |
| `--storage` | `VIBECAST_STORAGE` | `./data/sites` | 站点文件存储目录 |
| `--db` | `VIBECAST_DB` | `./data/vibecast.db` | SQLite 数据库路径 |

## API 接口

### 认证

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/auth/register` | 注册（`{email, password}`） |
| POST | `/api/auth/login` | 登录（`{email, password}`） |
| POST | `/api/auth/logout` | 退出登录 |
| GET | `/api/auth/me` | 获取当前用户 |

### 站点管理（需认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/sites` | 列出用户站点 |
| POST | `/api/sites` | 创建站点（`{name, slug, password}`） |
| GET | `/api/sites/{id}` | 获取站点详情 |
| PUT | `/api/sites/{id}` | 更新站点（`{name, password}`） |
| DELETE | `/api/sites/{id}` | 删除站点 |
| POST | `/api/sites/{id}/deploy` | 部署 ZIP（multipart `file` 字段） |

### 静态站点访问

- `/s/{slug}/` — 访问已部署的站点
- `/p/{slug}` — 密码验证页（站点设置了密码保护时）

## 架构

```
cmd/server/main.go          — 入口，CLI 参数，优雅关闭
internal/db/                 — SQLite schema、数据库迁移、数据模型
internal/auth/               — bcrypt 哈希、Session Token、认证中间件
internal/storage/            — ZIP 解压（含路径遍历防护）
internal/server/             — HTTP Handler、路由、静态文件服务、前端页面
```

## 技术栈

- **语言**：Go 1.23+
- **数据库**：SQLite（modernc.org/sqlite 纯 Go 驱动，无需 CGO）
- **认证**：bcrypt 密码哈希 + 随机 Session Token
- **文件服务**：自定义 Handler（非 http.FileServer），完全掌控 MIME、缓存和安全
- **前端**：单页 Dashboard，原生 JS，无构建步骤

## 开源协议

MIT
