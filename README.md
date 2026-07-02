# Vibecast

> Build with vibe. Cast instantly.

A self-hosted, multi-user static site hosting platform built in pure Go.  
No Nginx, no external web server — the application server handles everything:
authentication, site management, ZIP deployment, and static file serving.

[中文说明](#vibecast-中文说明)

---

## Features

- **User System** — Register, login, captcha-protected auth (bcrypt + secure cookies)
- **Admin Panel** — Full admin dashboard with user management, site oversight, and system settings
- **Multi-Site** — Each user can create multiple sites with auto-generated random slugs
- **ZIP Deploy** — Upload a ZIP file, auto-extract and deploy instantly
- **Password Protection** — Optional per-site password gate with 7-day session cookies
- **Directory Listing** — nginx-style auto-index when no `index.html` is present
- **File Tree Viewer** — Click-to-expand file listing in both dashboard and admin panel
- **Dark / Light Theme** — System-wide theme toggle with CSS variables, persisted per user
- **Bilingual UI (EN / 中文)** — Full i18n across dashboard, admin panel, and API error messages
- **Captcha** — SVG-based math captcha for login and registration
- **Storage Cleanup** — Admin tool to scan and remove orphaned site directories
- **Settings Control** — Toggle open registration, public site access, and email domain restriction
- **Pure Go File Server** — Custom static file handler with MIME detection, path traversal protection, and conditional requests
- **Zero External Dependencies** — Only SQLite (pure-Go driver) + Go standard library
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
The first registered user automatically becomes an admin.

## Configuration

| Flag | Env Var | Default | Description |
|------|---------|---------|-------------|
| `--addr` | `VIBECAST_ADDR` | `:8080` | Listen address |
| `--storage` | `VIBECAST_STORAGE` | `./data/sites` | Site files storage directory |
| `--db` | `VIBECAST_DB` | `./data/vibecast.db` | SQLite database path |

## Admin Settings

Configurable from the admin panel at `/admin`:

| Setting | Description |
|---------|-------------|
| **Open Registration** | When disabled, new registration is blocked (the register link is hidden on the login page) |
| **Public Site Access** | When disabled, all deployed sites return 403 to visitors unless password-protected |
| **Email Domain Restriction** | Restrict registration to specified email domains (one per line) |

## API

### Auth

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/auth/register` | Register (`{email, password, confirm, captchaId, captchaCode}`) |
| POST | `/api/auth/login` | Login (`{email, password, captchaId, captchaCode}`) |
| POST | `/api/auth/logout` | Logout |
| GET | `/api/auth/me` | Current user info |
| GET | `/api/auth/captcha` | Get captcha image (`{id, image}`) |
| POST | `/api/auth/change-password` | Change password (`{currentPassword, newPassword}`) |

### Public

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/settings` | Public settings (`{openRegistration}`) |

### Sites (auth required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/sites?page=&perPage=&q=` | List user's sites (paginated, searchable) |
| POST | `/api/sites` | Create site (`{name, password}`) — slug auto-generated |
| GET | `/api/sites/{id}` | Get site detail |
| PUT | `/api/sites/{id}` | Update site (`{name, password}`) |
| DELETE | `/api/sites/{id}` | Delete site |
| POST | `/api/sites/{id}/deploy` | Deploy ZIP (multipart `file` field) |
| GET | `/api/sites/{id}/files` | List files in site directory |

### Admin (admin only)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/admin/stats` | Dashboard statistics |
| GET | `/api/admin/users?page=&perPage=&q=` | List all users (paginated) |
| PUT | `/api/admin/users/{id}` | Toggle admin role |
| DELETE | `/api/admin/users/{id}` | Delete user |
| GET | `/api/admin/sites?page=&perPage=&q=` | List all sites (paginated) |
| DELETE | `/api/admin/sites/{id}` | Delete any site |
| GET | `/api/admin/sites/{id}/files` | List files in any site |
| GET | `/api/admin/settings` | Get system settings |
| PUT | `/api/admin/settings` | Update settings |
| GET | `/api/admin/cleanup` | Scan for orphaned directories |
| POST | `/api/admin/cleanup` | Delete orphaned directories |

### Static Site Access

| URL | Description |
|-----|-------------|
| `/s/{slug}/` | Access deployed site (with directory listing fallback) |
| `/p/{slug}` | Password gate page (if site is protected) |

## Architecture

```
cmd/server/main.go              — Entry point, CLI flags, graceful shutdown
internal/db/                     — SQLite schema, migrations, data models
internal/auth/                  — bcrypt hashing, session tokens, middleware
internal/storage/               — ZIP extraction with path traversal protection
internal/server/
  ├── server.go                 — Router, config, slug generation
  ├── api.go                    — Auth handlers, sites API, deploy, file tree
  ├── admin.go                  — Admin API (users, sites, settings, cleanup)
  ├── static.go                 — Static file serving, directory listing
  ├── password.go               — Password gate handler
  ├── captcha.go                — SVG captcha generation
  ├── messages.go               — i18n message map (EN/ZH)
  └── pages.go                  — All HTML/CSS/JS (dashboard, admin, auth, password pages)
```

## Tech Stack

- **Language**: Go 1.23+
- **Database**: SQLite (via modernc.org/sqlite — pure Go, no CGO)
- **Auth**: bcrypt password hashing, random session tokens, 7-day site access sessions
- **Captcha**: SVG math captcha (no external image library)
- **File Serving**: Custom handler (no http.FileServer) for full control over MIME, caching, and security
- **Frontend**: Single-page apps (dashboard + admin), vanilla JS, CSS variables for theming, no build step
- **i18n**: Bilingual EN/ZH, language preference via `Accept-Language` header and localStorage

## License

MIT

---

# Vibecast 中文说明

> Build with vibe. Cast instantly.

一个自托管的纯 Go 多用户静态站点托管平台。
不依赖 Nginx 或任何外部 Web Server —— 应用服务器全权处理：
用户认证、站点管理、ZIP 部署和静态文件服务。

## 功能特性

- **用户系统** — 注册、登录、验证码保护（bcrypt 加密 + 安全 Cookie）
- **管理后台** — 完整管理面板，支持用户管理、站点总览和系统设置
- **多站点** — 每个用户可创建多个站点，自动生成随机 Slug
- **ZIP 部署** — 上传 ZIP 文件，自动解压并即时部署上线
- **密码保护** — 可选的站点级密码门禁，7 天有效期的 Session Cookie
- **目录列表** — 无 index.html 时自动展示 nginx 风格的目录列表
- **文件树查看** — Dashboard 和管理后台均支持点击展开查看站点文件列表
- **深色/浅色主题** — 基于 CSS 变量的全局主题切换，用户偏好持久化
- **中英文双语** — Dashboard、管理后台、API 错误提示全面支持中英文
- **验证码** — SVG 数学验证码，用于登录和注册
- **存储清理** — 管理员工具，扫描并清理无对应站点的孤立目录
- **设置控制** — 开关注册、公开访问、邮箱域名限制
- **纯 Go 文件服务器** — 自定义静态文件 Handler，支持 MIME 检测、路径遍历防护、条件请求
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
首个注册的用户自动成为管理员。

## 配置项

| 参数 | 环境变量 | 默认值 | 说明 |
|------|----------|--------|------|
| `--addr` | `VIBECAST_ADDR` | `:8080` | 监听地址 |
| `--storage` | `VIBECAST_STORAGE` | `./data/sites` | 站点文件存储目录 |
| `--db` | `VIBECAST_DB` | `./data/vibecast.db` | SQLite 数据库路径 |

## 管理员设置

在 `/admin` 管理后台可配置：

| 设置 | 说明 |
|------|------|
| **开放注册** | 关闭后禁止新用户注册（登录页的注册入口也会隐藏） |
| **公开站点访问** | 关闭后所有站点对访问者返回 403，除非设置了密码保护 |
| **邮箱域名限制** | 限制只能用指定域名的邮箱注册（每行一个） |

## API 接口

### 认证

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/auth/register` | 注册（`{email, password, confirm, captchaId, captchaCode}`） |
| POST | `/api/auth/login` | 登录（`{email, password, captchaId, captchaCode}`） |
| POST | `/api/auth/logout` | 退出登录 |
| GET | `/api/auth/me` | 获取当前用户 |
| GET | `/api/auth/captcha` | 获取验证码图片（`{id, image}`） |
| POST | `/api/auth/change-password` | 修改密码（`{currentPassword, newPassword}`） |

### 公开接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/settings` | 公开设置（`{openRegistration}`） |

### 站点管理（需认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/sites?page=&perPage=&q=` | 列出用户站点（分页、搜索） |
| POST | `/api/sites` | 创建站点（`{name, password}`）— Slug 自动生成 |
| GET | `/api/sites/{id}` | 获取站点详情 |
| PUT | `/api/sites/{id}` | 更新站点（`{name, password}`） |
| DELETE | `/api/sites/{id}` | 删除站点 |
| POST | `/api/sites/{id}/deploy` | 部署 ZIP（multipart `file` 字段） |
| GET | `/api/sites/{id}/files` | 查看站点文件列表 |

### 管理接口（需管理员）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/admin/stats` | 仪表盘统计 |
| GET | `/api/admin/users?page=&perPage=&q=` | 列出所有用户（分页） |
| PUT | `/api/admin/users/{id}` | 切换管理员角色 |
| DELETE | `/api/admin/users/{id}` | 删除用户 |
| GET | `/api/admin/sites?page=&perPage=&q=` | 列出所有站点（分页） |
| DELETE | `/api/admin/sites/{id}` | 删除任意站点 |
| GET | `/api/admin/sites/{id}/files` | 查看任意站点文件 |
| GET | `/api/admin/settings` | 获取系统设置 |
| PUT | `/api/admin/settings` | 更新设置 |
| GET | `/api/admin/cleanup` | 扫描孤立目录 |
| POST | `/api/admin/cleanup` | 删除孤立目录 |

### 静态站点访问

| URL | 说明 |
|-----|------|
| `/s/{slug}/` | 访问已部署站点（无 index.html 时显示目录列表） |
| `/p/{slug}` | 密码验证页（站点设置了密码保护时） |

## 架构

```
cmd/server/main.go              — 入口，CLI 参数，优雅关闭
internal/db/                     — SQLite schema、数据库迁移、数据模型
internal/auth/                   — bcrypt 哈希、Session Token、认证中间件
internal/storage/                — ZIP 解压（含路径遍历防护）
internal/server/
  ├── server.go                 — 路由、配置、Slug 生成
  ├── api.go                    — 认证、站点 API、部署、文件树
  ├── admin.go                  — 管理 API（用户、站点、设置、清理）
  ├── static.go                 — 静态文件服务、目录列表
  ├── password.go               — 密码门禁
  ├── captcha.go                — SVG 验证码生成
  ├── messages.go               — i18n 消息映射（中/英）
  └── pages.go                  — 所有 HTML/CSS/JS（Dashboard、管理后台、认证、密码页）
```

## 技术栈

- **语言**：Go 1.23+
- **数据库**：SQLite（modernc.org/sqlite 纯 Go 驱动，无需 CGO）
- **认证**：bcrypt 密码哈希 + 随机 Session Token，站点访问会话有效期 7 天
- **验证码**：SVG 数学验证码（无外部图片库）
- **文件服务**：自定义 Handler（非 http.FileServer），完全掌控 MIME、缓存和安全
- **前端**：单页应用（Dashboard + 管理后台），原生 JS，CSS 变量主题切换，无构建步骤
- **国际化**：中英文双语，通过 `Accept-Language` 请求头和 localStorage 传递语言偏好

## 开源协议

MIT
