# Vibecast

> Build with vibe. Cast instantly.

A self-hosted, multi-user static site hosting platform built in pure Go.
No Nginx, no external web server — one binary handles everything: auth, site management, ZIP deployment, and static file serving.

[中文说明](#vibecast-中文说明)

---

## Features

- **ZIP Deploy** — Upload a ZIP, auto-extract and deploy instantly
- **Password Protection** — Optional per-site password gate (7-day session cookie)
- **Admin Panel** — User management, site oversight, storage cleanup, system settings
- **Random Slugs** — Auto-generated unguessable URLs, no need to pick a slug
- **Directory Listing** — nginx-style auto-index when no `index.html` exists
- **File Tree** — Click any site to expand and browse its files
- **Dark / Light Theme** — Toggle with CSS variables, persisted per user
- **Bilingual EN / 中文** — Full i18n across UI and API errors
- **Captcha** — SVG math captcha on login and registration
- **Settings Control** — Toggle open registration, public access, email domain restriction
- **Version Info** — `vibecast --version` prints the version; admin panel shows it in the navbar
- **Zero external dependencies** — Pure Go + SQLite, no CGO

## Installation

### One-liner (Linux / macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/vst93/Vibecast/main/install.sh | bash
```

Or with wget:

```bash
wget -qO- https://raw.githubusercontent.com/vst93/Vibecast/main/install.sh | bash
```

Options:

```bash
# Install a specific version
curl -fsSL https://raw.githubusercontent.com/vst93/Vibecast/main/install.sh | bash -s -- --version 20260702

# Install to a custom directory
curl -fsSL https://raw.githubusercontent.com/vst93/Vibecast/main/install.sh | bash -s -- --dir /opt/vibecast
```

The script auto-detects your OS and architecture (linux/darwin, amd64/arm64), downloads the matching binary from [Releases](../../releases), and installs it to `/usr/local/bin/vibecast`.

### Build from source

```bash
git clone https://github.com/vst93/Vibecast.git
cd Vibecast
make build
./bin/vibecast
```

### Manual download

Grab the binary for your platform from the [Releases](../../releases) page, make it executable, and run.

## Quick Start

```bash
vibecast

# or with custom config
vibecast --addr :3000 --storage ./data/sites --db ./data/vibecast.db
```

Open `http://localhost:8080/dashboard` — the first registered user becomes admin.

## Usage

1. **Register** at `/dashboard`, first user is auto-promoted to admin
2. **Create a site** — just give it a name, optionally set an access password
3. **Deploy** — click "Deploy ZIP" and upload your site bundle
4. **Visit** — your site goes live at `/s/{slug}/`
5. **Manage** — expand any site to view its file tree; admin panel at `/admin`

Admins can toggle open registration, disable public access, restrict email domains, clean up orphaned directories, and manage all users and sites from `/admin`.

## Configuration

| Flag | Env Var | Default | Description |
|------|---------|---------|-------------|
| `--addr` | `VIBECAST_ADDR` | `:8080` | Listen address |
| `--storage` | `VIBECAST_STORAGE` | `./data/sites` | Site files storage directory |
| `--db` | `VIBECAST_DB` | `./data/vibecast.db` | SQLite database path |

## Architecture

```
cmd/server/main.go        — Entry point, CLI flags, graceful shutdown
internal/db/              — SQLite schema, migrations, data models
internal/auth/            — bcrypt, session tokens, middleware
internal/storage/         — ZIP extraction with path traversal protection
internal/server/          — HTTP handlers, routing, static serving, captcha, i18n, all UI
```

## Tech Stack

Go 1.23+ · SQLite (pure Go driver) · bcrypt · vanilla JS SPA · no build step

## License

MIT

---

# Vibecast 中文说明

> 随心构建，即刻发布。

一个自托管的纯 Go 多用户静态站点托管平台。
不依赖 Nginx 或任何外部 Web Server —— 一个二进制搞定一切：认证、站点管理、ZIP 部署、静态文件服务。

## 功能特性

- **ZIP 部署** — 上传 ZIP，自动解压即时上线
- **密码保护** — 可选的站点级密码门禁（7 天有效 Cookie）
- **管理后台** — 用户管理、站点总览、存储清理、系统设置
- **随机 Slug** — 自动生成不可猜测的 URL，无需手动填写
- **目录列表** — 无 index.html 时自动展示 nginx 风格目录列表
- **文件树** — 点击展开任意站点查看文件列表
- **深色 / 浅色主题** — CSS 变量切换，用户偏好持久化
- **中英文双语** — UI 和 API 错误提示全面支持
- **验证码** — SVG 数学验证码，登录注册保护
- **设置控制** — 开关注册、公开访问、邮箱域名限制
- **版本信息** — `vibecast --version` 查看版本号，管理后台导航栏显示版本
- **零外部依赖** — 纯 Go + SQLite，无需 CGO

## 安装

### 一键安装（Linux / macOS）

```bash
curl -fsSL https://raw.githubusercontent.com/vst93/Vibecast/main/install.sh | bash
```

或使用 wget：

```bash
wget -qO- https://raw.githubusercontent.com/vst93/Vibecast/main/install.sh | bash
```

可选参数：

```bash
# 安装指定版本
curl -fsSL https://raw.githubusercontent.com/vst93/Vibecast/main/install.sh | bash -s -- --version 20260702

# 安装到自定义目录
curl -fsSL https://raw.githubusercontent.com/vst93/Vibecast/main/install.sh | bash -s -- --dir /opt/vibecast
```

脚本自动检测操作系统和架构（linux/darwin，amd64/arm64），从 [Releases](../../releases) 下载对应二进制并安装到 `/usr/local/bin/vibecast`。

### 从源码编译

```bash
git clone https://github.com/vst93/Vibecast.git
cd Vibecast
make build
./bin/vibecast
```

### 手动下载

从 [Releases](../../releases) 页面下载对应平台的二进制，赋予执行权限后运行。

## 快速开始

```bash
vibecast

# 或指定配置
vibecast --addr :3000 --storage ./data/sites --db ./data/vibecast.db
```

打开 `http://localhost:8080/dashboard`，首个注册用户自动成为管理员。

## 使用方式

1. **注册** — 在 `/dashboard` 注册，首用户自动成为管理员
2. **创建站点** — 填个名字即可，可选设访问密码
3. **部署** — 点击 "Deploy ZIP" 上传站点压缩包
4. **访问** — 站点上线地址为 `/s/{slug}/`
5. **管理** — 点击站点展开查看文件树；管理后台在 `/admin`

管理员可在 `/admin` 开关注册、禁用公开访问、限制邮箱域名、清理孤立目录，以及管理所有用户和站点。

## 配置项

| 参数 | 环境变量 | 默认值 | 说明 |
|------|----------|--------|------|
| `--addr` | `VIBECAST_ADDR` | `:8080` | 监听地址 |
| `--storage` | `VIBECAST_STORAGE` | `./data/sites` | 站点文件存储目录 |
| `--db` | `VIBECAST_DB` | `./data/vibecast.db` | SQLite 数据库路径 |

## 架构

```
cmd/server/main.go        — 入口，CLI 参数，优雅关闭
internal/db/              — SQLite schema、迁移、数据模型
internal/auth/            — bcrypt、Session Token、认证中间件
internal/storage/         — ZIP 解压（含路径遍历防护）
internal/server/          — HTTP Handler、路由、静态服务、验证码、i18n、全部前端页面
```

## 技术栈

Go 1.23+ · SQLite（纯 Go 驱动）· bcrypt · 原生 JS 单页应用 · 无构建步骤

## 开源协议

MIT
