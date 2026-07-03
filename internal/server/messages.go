package server

import (
	"net/http"
	"strings"
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
	"updateRestart":           {"en": "Restart Now", "zh": "立即重启"},
	"updateRestarting":        {"en": "Restarting...", "zh": "正在重启..."},
	"updateRestartSuccess":    {"en": "Server is restarting. Page will reload automatically.", "zh": "服务器正在重启，页面将自动刷新。"},
	"updateInProgress":        {"en": "Update already in progress", "zh": "更新正在进行中"},
	"updateVerifyFailed":      {"en": "Checksum verification failed", "zh": "校验和验证失败"},
	"updateNoChecksum":        {"en": "No checksum available (skipped verification)", "zh": "无校验和（已跳过验证）"},
	"updateDownloadProgress":  {"en": "Downloading", "zh": "下载中"},
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
