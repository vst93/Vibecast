# DESIGN.md · Vibecast 设计锚点

本文件是 Vibecast 界面层唯一的决策来源。新增或修改 UI 前先读这里；实现必须
落在 `internal/server/ui_css.go` 的 token 与组件层上，不要在页面里另起一套样式。

## 1. 设计读取

> **读取为**：自托管静态发布控制台，面向独立开发者、小团队运维者与内容作者，
> 用「克制、可信、工具感」的语言，倾向 GitHub Primer × Linear 的现代工程工具美学。

| 拨盘 | 值 | 人话 |
|---|---|---|
| 视觉冒险度 | 4 | 这是每天用的功能页，一致性优先于惊喜 |
| 动效强度 | 3 | 只保留状态反馈，不做装饰性编排 |
| 信息密度 | 6 | 控制台效率优先，但每行要有呼吸 |

## 2. 功能契约

- **存在理由**：把一份文件在几秒内变成可分享的 URL，并能随时管住访问权限。
- **3 秒内必须看到**：我有几个站点、每个站点的地址与状态（公开/受保护）、有多少访问。
- **必须能完成**：创建站点 → 上传文件/ZIP → 打开或复制链接 → 设置密码/组织可见性 → 删除或清理文件。
- 任何视觉决策不得遮蔽这三条。创意只能长在它们之上。

## 3. 原则

1. **一个焦点**：每屏只有一个主行动。站点行只有「上传」是实心按钮。
2. **机器数据用等宽**：URL、slug、文件大小、访问数、版本号走 `--ff-mono` + `tabular-nums`。
3. **层级靠字号与颜色，不靠阴影堆叠**：卡片刻意只有 1px 描边 + 极轻投影。
4. **危险操作远离主路径**：删除进溢出菜单、二次确认、默认焦点在「取消」。
5. **每次删除都问一次，每次保存都不问**：只有不可逆操作需要确认。

## 4. Token

色板（`uiTokensCSS`）：

| 角色 | 亮色 | 暗色 | 用途 |
|---|---|---|---|
| `--bg` | `#f2f5f3` | `#0b1210` | 页面底 |
| `--surface` | `#ffffff` | `#121b17` | 卡片 |
| `--surface-2/3` | `#eef2f0` / `#e3eae6` | `#182521` / `#20302a` | 次级面、轨道、表头 |
| `--line` / `--line-strong` | `#d5ded9` / `#bccbc3` | `#2b3d35` / `#3b5148` | 描边 |
| `--text` / `--text-soft` / `--dim` / `--dim-2` | `#16201b` / `#31413a` / `#596b62` / `#61746c` | `#e8f0eb` / `#d2ded7` / `#9db1a7` / `#849790` | 正文层级 |
| `--accent` / `--accent-ink` | `#0a7d59` / `#076347` | `#4bd6a1` / `#4bd6a1` | 品牌绿；`ink` 专用于浅底文字 |
| `--warn` / `--danger` | `#8a5a00` / `#b3242c` | `#e6bd6a` / `#ff8a91` | 受保护 / 破坏性 |
| `--field` | `#ffffff` | `#0b1512` | 输入框底色（暗色比卡片更深） |

对比度实测（`WCAG AA`，全部 ≥4.5:1）：正文 16.7、次要文字 5.2–5.7、允许/品牌文字 5.1、
品牌按钮文字 5.1、危险文字 6.6、暗色正文 15.1、暗色次要文字 7.8、暗色品牌按钮 9.2。
徽标配色全部 ≥4.9。

字阶：`12 / 13 / 14 / 15 / 17 / 22px`；正文行高 1.55，中文 ≥1.5。
圆角：控件 8、卡片 12、小元素 6；间距以 4 为步进（8/12/16/20/24）。
动效：`--dur:150ms`，`--ease:cubic-bezier(.22,.61,.36,1)`，只动 transform/opacity/颜色，
`prefers-reduced-motion` 兜底。

### 字体（刻意不用 webfont）

- `--ff-sans` 以具体系统字族开头（`-apple-system` → `Segoe UI` → `Roboto` → `Noto Sans`
  → `Arial` → 中文族 → `Liberation/DejaVu` → `sans-serif`）。
- **不要**把 `system-ui` / `ui-sans-serif` 放在首位：在本机 fontconfig 上，
  这两个泛型族在 `font-weight:600` 时会被解析成 Noto **Serif** CJK，整站悄悄变衬线。
  这是实测缺陷，改动字体栈前先复测 `CSS.getPlatformFontsForNode`。
- 选择系统栈而非 webfont：自托管产品、零外部请求、中文 webfont 体积不可接受；
  个性靠字阶、留白、等宽数据和绿色单点强调来实现。

## 5. 组件规则

- **站点行**：`chevron | 状态点 + 名称 + 徽标 | 地址 + 访问数 | 上传 + 溢出菜单`。
  整行可点开展开，`role=button` + `aria-expanded` + Enter/Space；行内链接与按钮
  `stopPropagation`。地址是真 `<a target=_blank>`，不再单设「访问」按钮
  （同一意图只有一个入口）。
- **溢出菜单**：`<details class="row-menu">` + `<summary aria-label>`；复制链接、复制密码、
  分享、编辑、下载 ZIP、清空文件、删除。点击外部或 Esc 关闭。删除永远在菜单最底部。
- **详情面板**：`grid-template-rows 0fr→1fr` 展开（padding 放在 `.detail-inner`，
  外层 `.detail-clip` 负责裁剪，否则收起时漏出 32px）。内容顺序：规格网格 → 两个开关 → 文件列表。
- **文件列表**：等宽、行高 32px、大小右对齐；空态用虚线框 + 说明。
- **表格**（管理端）：11px 大写字距表头、行 hover、数字列 `tabular-nums` 右对齐、
  操作列只放文字幽灵按钮；≤720px 隐藏 ID/时间/访问量等次要列，把徽标下沉到主单元格。
- **开关**：`button[role=switch][aria-checked]`（管理端）或 `label>input[type=checkbox]`
  + `.switch-track`（表单）。只有 400/700 两档字重可用。
- **弹窗**：`role=dialog` + `aria-modal` + `aria-labelledby`；打开聚焦第一个输入/按钮，
  Esc 关闭并把焦点还给触发元素，遮罩 `overscroll-behavior:contain`。
- **Toast**：`#toast-host[aria-live=polite]`，右下角堆叠，3 秒消失。
- **错误页**：访客侧 403/404 用 `visitorErrorPageHTML`，含错误码、原因、下一步行动，
  不再返回裸文本。

## 6. 布局

| 断点 | 行为 |
|---|---|
| ≥1025px | 320px 侧栏（sticky，自带滚动）+ 主列表 |
| ≤1024px | 单列；创建面板折叠为「展开」 |
| ≤720px | 导航收进汉堡菜单；站点行动作独占一行；表格隐藏次要列 |
| ≤420px | 详情内边距收紧，访问数换行 |

并列滚动区（侧栏、文件列表、成员列表）各自 `overscroll-behavior: contain`，
侧栏固定视口高度，滚到底不会带动页面。

## 7. 结构位置

- `internal/server/ui_css.go`：`uiTokensCSS`（token/reset/a11y）、`uiComponentsCSS`
  （导航、按钮、表单、卡片、表格、弹窗、Toast、菜单）、`uiDashboardCSS`、`uiAdminCSS`、
  `uiLandingCSS`、`uiGateCSS`。
- `internal/server/pages.go`：页面骨架与内联 JS，只允许出现页面特有样式；
  通用类名复用上面的层。
- 每个页面用 `<style>` + token + 组件 + 页面层拼接；**禁止**再追加"覆盖层"式补丁 CSS。

## 8. Do / Don't

- Do：新增 UI 先找现成 token 与组件；文案给下一步；数字用 `tabular-nums`。
- Do：新增字符串同时补 `en` 与 `zh`（`Object.assign(i18n.en/zh, {...})` 更安全）。
- Don't：`font-style: italic`、`transition: all`、`scale(0)` 入场、装饰性 glow、
  左侧竖线强调、卡片包卡片、无 label 的图标按钮。
- Don't：在主页堆假数据；截图资产必须来自真实运行实例。

## 9. 本次改版记录

- 修复字体栈把整站渲染成衬线（Linux/fontconfig）。
- 合并三套互相覆盖的 CSS 补丁为 token + 组件 + 页面四层；管理端与用户端共用组件。
- 站点行重构：可展开行 + 溢出菜单，销毁按钮从 3 个实心按钮减到 1 个主行动。
- 详情面板重排：规格 → 开关 → 文件，修复收起状态漏出内容。
- 管理端：概览加"最新站点"、统计卡重做、表格重做、设置开关语义化、版本号移到侧栏底部。
- 落地页：产品预览改为真实 DOM 迷你面板（不再模糊截图）、图标换成统一 SVG、
  修复品牌字样间距、补齐主题切换。
- 访客侧 403/404 改成品牌化错误页；密码门改用共享 token。
- 无障碍：焦点环、Esc 关闭弹窗/菜单、`aria-expanded`、`aria-live`、
  图标按钮 `aria-label`、对比度全部达标。
