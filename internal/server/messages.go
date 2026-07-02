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
	"create_site_failed":    {"en": "Failed to create site", "zh": "创建站点失败"},
	"list_sites_failed":     {"en": "Failed to list sites", "zh": "获取站点列表失败"},
	"count_sites_failed":    {"en": "Failed to count sites", "zh": "统计站点失败"},
	"site_id_required":      {"en": "Site ID required", "zh": "需要站点 ID"},
	"site_not_found":        {"en": "Site not found", "zh": "站点不存在"},
	"forbidden":             {"en": "Forbidden", "zh": "无权限"},
	"file_required":         {"en": "File upload required (field name: 'file')", "zh": "需要上传文件（字段名：'file'）"},
	"zip_only":              {"en": "Only .zip files are accepted", "zh": "仅接受 .zip 文件"},
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
