package server

import (
	"net/http"
	"strings"
	"time"
)

var messageMap = map[string]map[string]string{
	"method_not_allowed":     {"en": "Method not allowed", "zh": "不支持的请求方法"},
	"invalid_json":           {"en": "Invalid JSON", "zh": "无效的 JSON 格式"},
	"email_password_invalid": {"en": "Email required and password must be at least 6 characters", "zh": "邮箱必填且密码至少 6 位"},
	"password_too_long":      {"en": "Password too long (max 72 chars)", "zh": "密码过长（最多 72 字符）"},
	"password_mismatch":     {"en": "Passwords do not match", "zh": "两次输入的密码不一致"},
	"captcha_incorrect":     {"en": "Captcha incorrect", "zh": "验证码错误"},
	"registration_closed":   {"en": "Registration is closed", "zh": "注册已关闭"},
	"domain_not_allowed":    {"en": "Email domain not allowed", "zh": "邮箱域名不被允许"},
	"hash_failed":           {"en": "Failed to hash password", "zh": "密码加密失败"},
	"email_taken":           {"en": "Email already registered", "zh": "邮箱已被注册"},
	"create_user_failed":    {"en": "Failed to create user", "zh": "创建用户失败"},
	"create_session_failed": {"en": "Failed to create session", "zh": "创建会话失败"},
	"internal_error":        {"en": "Internal error", "zh": "内部错误"},
	"invalid_credentials":   {"en": "Invalid email or password", "zh": "邮箱或密码错误"},
	"unauthorized":          {"en": "Unauthorized", "zh": "未登录"},
	"new_password_too_short": {"en": "New password must be at least 6 characters", "zh": "新密码至少 6 位"},
	"current_password_wrong": {"en": "Current password is incorrect", "zh": "当前密码不正确"},
	"update_password_failed": {"en": "Failed to update password", "zh": "更新密码失败"},
	"name_required":         {"en": "Name required", "zh": "站点名称必填"},
	"invalid_slug":          {"en": "Invalid slug (2-63 chars, a-z0-9- only)", "zh": "Slug 无效（2-63 字符，仅限 a-z0-9-）"},
	"check_slug_failed":     {"en": "Failed to check slug", "zh": "检查 Slug 失败"},
	"slug_taken":            {"en": "Slug already taken", "zh": "Slug 已被占用"},
	"generate_slug_failed":  {"en": "Failed to generate slug", "zh": "生成 Slug 失败"},
	"site_password_too_short": {"en": "Site password must be at least 4 characters", "zh": "站点密码至少 4 位"},
	"public_access_disabled": {"en": "Public access is disabled: site must have a password", "zh": "已禁用公开访问：站点必须设置密码"},
	"site_limit_reached":    {"en": "Site limit reached. Delete an existing site or contact admin.", "zh": "已达站点数量上限，请删除已有站点或联系管理员"},
	"create_site_failed":    {"en": "Failed to create site", "zh": "创建站点失败"},
	"list_sites_failed":     {"en": "Failed to list sites", "zh": "获取站点列表失败"},
	"count_sites_failed":    {"en": "Failed to count sites", "zh": "统计站点失败"},
	"site_id_required":      {"en": "Site ID required", "zh": "需要站点 ID"},
	"site_not_found":        {"en": "Site not found", "zh": "站点不存在"},
	"forbidden":             {"en": "Forbidden", "zh": "无权限"},
	"file_required":         {"en": "File upload required (field name: 'file')", "zh": "需要上传文件（字段名：'file'）"},
	"zip_only":              {"en": "Only .zip files are accepted", "zh": "仅接受 .zip 文件"},
	"file_type_blocked":     {"en": "This file type is not allowed (executable/script files are blocked)", "zh": "不允许此文件类型（可执行文件/脚本文件已被拦截）"},
	"file_too_large":        {"en": "File too large (max 100 MB)", "zh": "文件过大（最大 100 MB）"},
	"file_too_large_single": {"en": "File too large: %s exceeds 100 MB limit", "zh": "文件过大：%s 超过 100 MB 限制"},
	"zip_file_too_large":    {"en": "File in zip too large: %s exceeds 100 MB limit", "zh": "ZIP 内文件过大：%s 超过 100 MB 限制"},
	"zip_bomb":              {"en": "Zip bomb detected: total uncompressed size exceeds 500 MB limit", "zh": "检测到 ZIP 炸弹：解压后总大小超过 500 MB 限制"},
	"save_file_failed":      {"en": "Failed to save file", "zh": "保存文件失败"},
	"read_upload_failed":    {"en": "Failed to read upload", "zh": "读取上传文件失败"},
	"extract_zip_failed":    {"en": "Failed to extract zip", "zh": "解压 zip 文件失败"},
	"update_failed":         {"en": "Failed to update", "zh": "更新失败"},
	"update_site_failed":    {"en": "Failed to update site", "zh": "更新站点失败"},
	"delete_site_failed":    {"en": "Failed to delete site", "zh": "删除站点失败"},
	"delete_user_failed":    {"en": "Failed to delete user", "zh": "删除用户失败"},
	// admin.go
	"invalid_user_id":         {"en": "Invalid user ID", "zh": "用户 ID 无效"},
	"user_not_found":          {"en": "User not found", "zh": "用户不存在"},
	"cannot_modify_self_admin": {"en": "Cannot modify your own admin status", "zh": "不能修改自己的管理员状态"},
	"cannot_delete_self":      {"en": "Cannot delete yourself", "zh": "不能删除自己"},
	"invalid_site_id":         {"en": "Invalid site ID", "zh": "站点 ID 无效"},
	"get_settings_failed":     {"en": "Failed to get settings", "zh": "获取设置失败"},
	"update_settings_failed":  {"en": "Failed to update settings", "zh": "更新设置失败"},
	"list_users_failed":       {"en": "Failed to list users", "zh": "获取用户列表失败"},
	"count_users_failed":      {"en": "Failed to count users", "zh": "统计用户失败"},
	"get_stats_failed":        {"en": "Failed to get stats", "zh": "获取统计数据失败"},
	// password.go
	"incorrect_password": {"en": "Incorrect password", "zh": "密码错误，请重试"},
	// update.go
	"update_fetch_failed":     {"en": "Failed to fetch update info", "zh": "获取更新信息失败"},
	"update_asset_not_found":  {"en": "No matching binary found for your platform", "zh": "未找到匹配当前平台的安装包"},
	"update_download_failed":  {"en": "Failed to download update", "zh": "下载更新失败"},
	"update_install_failed":   {"en": "Failed to install update", "zh": "安装更新失败"},
	"update_permission_denied": {"en": "Permission denied — the binary is not writable by the current user. Try: sudo vibecast update (CLI) or restart the service as the owning user.", "zh": "权限不足 — 当前用户无法写入二进制文件。请尝试：sudo vibecast update（命令行）或以拥有该文件的用户重启服务。"},
	"update_windows_locked":    {"en": "Cannot replace the running binary on Windows — stop the service first, then retry.", "zh": "Windows 下无法替换正在运行的程序 — 请先停止服务再重试。"},
	"updateRestart":           {"en": "Restart Now", "zh": "立即重启"},
	"updateRestarting":        {"en": "Restarting...", "zh": "正在重启..."},
	"updateRestartSuccess":    {"en": "Server is restarting. Page will reload automatically.", "zh": "服务器正在重启，页面将自动刷新。"},
	"updateInProgress":        {"en": "Update already in progress", "zh": "更新正在进行中"},
	"updateVerifyFailed":      {"en": "Checksum verification failed", "zh": "校验和验证失败"},
	"updateNoChecksum":        {"en": "No checksum available (skipped verification)", "zh": "无校验和（已跳过验证）"},
	"updateDownloadProgress":  {"en": "Downloading", "zh": "下载中"},
	// organizations
	"already_in_org":          {"en": "You are already in an organization. Leave or delete it first.", "zh": "你已在组织中，请先退出或删除当前组织"},
	"create_org_failed":       {"en": "Failed to create organization", "zh": "创建组织失败"},
	"invite_code_required":    {"en": "Invite code is required", "zh": "请输入邀请码"},
	"org_not_found":           {"en": "Organization not found", "zh": "组织不存在"},
	"join_org_failed":         {"en": "Failed to join organization", "zh": "加入组织失败"},
	"not_in_org":              {"en": "You are not in an organization", "zh": "你不在任何组织中"},
	"owner_cannot_leave":      {"en": "Organization owner cannot leave. Delete the organization instead.", "zh": "组织创建者不能退出，请删除组织"},
	"leave_org_failed":        {"en": "Failed to leave organization", "zh": "退出组织失败"},
	"not_org_owner":           {"en": "Only the organization owner can do this", "zh": "只有组织创建者可以执行此操作"},
	"org_has_members":         {"en": "Cannot delete: organization still has other members", "zh": "无法删除：组织中还有其他成员"},
	"delete_org_failed":       {"en": "Failed to delete organization", "zh": "删除组织失败"},
	"cannot_remove_self":      {"en": "Cannot remove yourself", "zh": "不能移除自己"},
	"org_open_requires_org":   {"en": "You must be in an organization to enable org access", "zh": "需要先加入组织才能开启组织内访问"},
}

func tMsg(r *http.Request, key string) string {
	lang := "en"
	if al := r.Header.Get("Accept-Language"); strings.HasPrefix(al, "zh") {
		lang = "zh"
	}
	if msgs, ok := messageMap[key]; ok {
		if msg, ok := msgs[lang]; ok {
			return msg
		}
	}
	return key
}

// tStatic returns the English message for a key, for use outside of HTTP
// request context (e.g. CLI output). Falls back to the key itself.
func tStatic(key string) string {
	if msgs, ok := messageMap[key]; ok {
		if msg, ok := msgs["en"]; ok {
			return msg
		}
	}
	return key
}

// isCST checks whether the system local timezone is China Standard Time (UTC+8).
// Used by tCLI() to decide CLI output language.
func isCST() bool {
	loc, err := time.LoadLocation("Local")
	if err != nil {
		return false
	}
	_, offset := time.Now().In(loc).Zone()
	// UTC+8 = 28800 seconds; also accept UTC+8 name
	return offset == 8*3600
}

// tCLI returns a localized message for CLI output based on the system timezone.
// UTC+8 (CST) → Chinese, otherwise → English.
func tCLI(key string) string {
	lang := "en"
	if isCST() {
		lang = "zh"
	}
	if msgs, ok := messageMap[key]; ok {
		if msg, ok := msgs[lang]; ok {
			return msg
		}
	}
	return key
}

// CLI message keys (not used by HTTP handlers)
var cliMessages = map[string]map[string]string{
	"cli_usage":          {"en": "Usage: vibecast [options] [command]", "zh": "用法: vibecast [选项] [命令]"},
	"cli_options":        {"en": "Options:", "zh": "选项:"},
	"cli_commands":       {"en": "Commands:", "zh": "命令:"},
	"cli_addr":           {"en": "listen address", "zh": "监听地址"},
	"cli_storage":        {"en": "site files storage directory", "zh": "站点文件存储目录"},
	"cli_db":             {"en": "SQLite database path", "zh": "SQLite 数据库路径"},
	"cli_version_cmd":    {"en": "print version and exit", "zh": "打印版本号并退出"},
	"cli_update_cmd":     {"en": "check for updates and self-update", "zh": "检查更新并自更新"},
	"cli_help_cmd":       {"en": "show this help message", "zh": "显示帮助信息"},
	"cli_unknown_cmd":   {"en": "unknown command", "zh": "未知命令"},
	"cli_listening":     {"en": "Listening:", "zh": "监听地址:"},
	"cli_storage_label": {"en": "Storage:", "zh": "存储路径:"},
	"cli_database":      {"en": "Database:", "zh": "数据库:"},
	"cli_dashboard":     {"en": "Dashboard:", "zh": "控制面板:"},
	"cli_checking":      {"en": "Checking for updates...", "zh": "正在检查更新..."},
	"cli_latest_rel":    {"en": "Latest release:", "zh": "最新版本:"},
	"cli_up_to_date":    {"en": "You are already running the latest version.", "zh": "当前已是最新版本。"},
	"cli_dev_version":   {"en": "Current version: dev (always allows update)", "zh": "当前版本: dev（始终允许更新）"},
	"cli_update_avail": {"en": "Update available!", "zh": "有可用更新！"},
	"cli_release":       {"en": "Release:", "zh": "版本:"},
	"cli_downloading":   {"en": "Downloading", "zh": "下载中"},
	"cli_downloaded":    {"en": "Downloaded", "zh": "下载完成"},
	"cli_checksum_ok":   {"en": "Checksum verified", "zh": "校验和验证通过"},
	"cli_installing":    {"en": "Installing...", "zh": "安装中..."},
	"cli_updated":       {"en": "Updated to", "zh": "已更新至"},
	"cli_restart_hint":  {"en": "Please restart vibecast to apply the update.", "zh": "请重启 vibecast 以应用更新。"},
	"cli_update_failed": {"en": "Update failed", "zh": "更新失败"},
	"cli_fetch_failed":  {"en": "failed to check for updates", "zh": "检查更新失败"},
	"cli_no_binary":     {"en": "no matching binary found for", "zh": "未找到匹配的二进制文件"},
	"cli_dl_failed":     {"en": "download failed", "zh": "下载失败"},
	"cli_install_failed": {"en": "installation failed", "zh": "安装失败"},
	"cli_empty_file":    {"en": "downloaded file is empty or invalid", "zh": "下载的文件为空或无效"},
	// service subcommand
	"svc_windows_unsupported": {"en": "Service management is not supported on Windows. Please register Vibecast as a Windows Service manually or use Task Scheduler.", "zh": "Windows 不支持服务管理。请手动注册为 Windows 服务或使用任务计划程序。"},
	"svc_windows_hint":        {"en": "Hint: You can use nssm (https://nssm.cc) to register vibecast as a Windows service.", "zh": "提示：可以使用 nssm (https://nssm.cc) 将 vibecast 注册为 Windows 服务。"},
	"svc_unsupported":         {"en": "Service management is not supported on this platform", "zh": "此平台不支持服务管理"},
	"svc_installing":          {"en": "Installing service", "zh": "正在安装服务"},
	"svc_install_failed":      {"en": "Service installation failed", "zh": "服务安装失败"},
	"svc_installed":           {"en": "Service installed and started", "zh": "服务已安装并启动"},
	"svc_uninstalled":         {"en": "Service uninstalled", "zh": "服务已卸载"},
	"svc_uninstall_failed":   {"en": "Service uninstall failed", "zh": "服务卸载失败"},
	"svc_uninstall_cmd":      {"en": "uninstall service", "zh": "卸载服务"},
	"cli_service_cmd":        {"en": "manage system service (install/status/stop/restart/uninstall)", "zh": "管理系统服务（安装/状态/停止/重启/卸载）"},
	"cli_setup_cmd":          {"en": "register vibecast as a system service", "zh": "注册为系统服务"},
	"cli_uninstall_cmd":      {"en": "uninstall the system service", "zh": "卸载系统服务"},
	"svc_manage_hint_linux":  {"en": "Manage the service with standard commands:", "zh": "使用标准命令管理服务："},
	"svc_manage_hint_macos":  {"en": "Manage the service with standard commands:", "zh": "使用标准命令管理服务："},
}

// TCLIMsg returns a CLI-specific message (from cliMessages) based on timezone.
func TCLIMsg(key string) string {
	lang := "en"
	if isCST() {
		lang = "zh"
	}
	if msgs, ok := cliMessages[key]; ok {
		if msg, ok := msgs[lang]; ok {
			return msg
		}
	}
	return key
}
