package server

import (
	"fmt"
	"strings"
)

// landingPageHTML is the public landing page.
var landingPageHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Vibecast — Build with vibe. Cast instantly.</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--ink:#f6f8fa;--surface:#fff;--surface-2:#eef1f4;--line:#d0d7de;--text:#1f2328;--dim:#636c76;--accent:#1a7f37;--accent-dim:#2da44e;--warn:#9a6700;--danger:#cf222e;--placeholder:#8c959f;--mono:'JetBrains Mono','Fira Code',monospace;--sans:system-ui,-apple-system,sans-serif}
[data-theme="dark"]{--ink:#0c1117;--surface:#161b22;--surface-2:#1c2128;--line:#30363d;--text:#e6edf3;--dim:#7d8590;--accent:#39d353;--accent-dim:#238636;--warn:#d29922;--danger:#f85149;--placeholder:#484f58}
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:var(--sans);background:var(--ink);color:var(--text);min-height:100vh;-webkit-font-smoothing:antialiased}
a{color:var(--accent);text-decoration:none}
.wrap{max-width:820px;margin:0 auto;padding:2rem 1.5rem}
.hero{text-align:center;padding:4rem 0 3rem}
.hero .logo{font-family:var(--mono);font-size:3rem;font-weight:700;letter-spacing:-.03em;color:var(--text)}
.hero .logo .accent{color:var(--accent)}
.hero .tagline{font-family:var(--mono);font-size:1rem;color:var(--dim);margin:1rem 0 2.5rem}
.hero .cta{display:inline-block;padding:10px 28px;background:var(--accent);color:var(--ink);border:none;border-radius:6px;font-family:var(--mono);font-size:.9rem;font-weight:700;cursor:pointer;transition:background .15s;text-decoration:none}
.hero .cta:hover{background:var(--accent-dim)}
.features{display:grid;grid-template-columns:repeat(auto-fit,minmax(170px,1fr));gap:1rem;margin-top:2rem}
.feature{padding:1.5rem 1.25rem;background:var(--surface);border:1px solid var(--line);border-radius:8px;transition:border-color .15s}
.feature:hover{border-color:var(--accent-dim)}
.feature .icon{font-family:var(--mono);font-size:1.4rem;color:var(--accent);margin-bottom:.5rem}
.feature h3{font-family:var(--mono);font-size:.85rem;margin-bottom:.3rem;color:var(--text)}
.feature p{font-size:.8rem;color:var(--dim);line-height:1.5}
footer{text-align:center;padding:2rem 0;color:var(--dim);font-family:var(--mono);font-size:.75rem;border-top:1px solid var(--line);margin-top:3rem}
.lang-toggle{position:absolute;top:1.5rem;right:1.5rem;display:flex;gap:.25rem}
.lang-toggle a{font-family:var(--mono);font-size:.75rem;cursor:pointer;padding:4px 10px;border-radius:5px;color:var(--dim);transition:all .15s}
.lang-toggle a.active{background:var(--surface-2);color:var(--accent)}
@media(max-width:640px){.hero{padding:2rem 0}.hero .logo{font-size:2.2rem}.features{grid-template-columns:1fr}}
@media(prefers-reduced-motion:reduce){*{transition:none!important}}
</style>
<script>var _t="light";try{_t=localStorage.getItem("theme")||"light"}catch(e){}document.documentElement.setAttribute("data-theme",_t)</script>
</head>
<body>
<div class="lang-toggle"><a id="theme-toggle" onclick="toggleTheme()">🌙</a><a id="langEn" class="active" onclick="setLang('en')">EN</a><a id="langZh" onclick="setLang('zh')">中文</a></div>
<div class="wrap">
<div class="hero">
<div class="logo">Vibe<span class="accent">cast</span></div>
<p class="tagline" data-i="tagline">Build with vibe. Cast instantly.</p>
<a href="dashboard" class="cta" data-i="getStarted">Get Started</a>
<div class="features">
<div class="feature"><div class="icon">$ deploy</div><h3 data-i="feat1Title">Instant Deploy</h3><p data-i="feat1Desc">Upload a ZIP, get a live URL in seconds.</p></div>
<div class="feature"><div class="icon">****</div><h3 data-i="feat2Title">Password Protect</h3><p data-i="feat2Desc">Keep your site private with password gating.</p></div>
<div class="feature"><div class="icon">{ }</div><h3 data-i="feat3Title">Zero Nginx</h3><p data-i="feat3Desc">Pure Go application server. No dependencies.</p></div>
<div class="feature"><div class="icon">~/</div><h3 data-i="feat4Title">Self-Hosted</h3><p data-i="feat4Desc">Your data, your rules. SQLite + filesystem.</p></div>
</div>
</div>
<footer data-i="footer">Vibecast — Self-hosted static site hosting</footer>
</div>
<script>
var i18n={en:{tagline:"Build with vibe. Cast instantly.",getStarted:"Get Started",feat1Title:"Instant Deploy",feat1Desc:"Upload a ZIP, get a live URL in seconds.",feat2Title:"Password Protect",feat2Desc:"Keep your site private with password gating.",feat3Title:"Zero Nginx",feat3Desc:"Pure Go application server. No dependencies.",feat4Title:"Self-Hosted",feat4Desc:"Your data, your rules. SQLite + filesystem.",footer:"Vibecast — Self-hosted static site hosting"},zh:{tagline:"Build with vibe. Cast instantly.",getStarted:"开始使用",feat1Title:"即时部署",feat1Desc:"上传 ZIP，秒级获取线上 URL。",feat2Title:"密码保护",feat2Desc:"通过密码门禁保护你的站点隐私。",feat3Title:"零依赖",feat3Desc:"纯 Go 应用服务器，无需 Nginx。",feat4Title:"自托管",feat4Desc:"数据归你所有，SQLite + 文件系统。",footer:"Vibecast — 自托管静态站点托管平台"}};
function getTheme(){return document.documentElement.getAttribute("data-theme")||"light"}
function setTheme(t){document.documentElement.setAttribute("data-theme",t);try{localStorage.setItem("theme",t)}catch(e){}var b=document.getElementById("theme-toggle");if(b)b.textContent=t==="dark"?"☀":"🌙"}
function toggleTheme(){setTheme(getTheme()==="dark"?"light":"dark")}
var lang="en";
function setLang(l){lang=l;document.querySelectorAll("[data-i]").forEach(function(e){var k=e.getAttribute("data-i");if(i18n[l][k])e.textContent=i18n[l][k]});document.getElementById("langEn").className=l==="en"?"active":"";document.getElementById("langZh").className=l==="zh"?"active":"";try{localStorage.setItem("lang",l)}catch(e){}}
var saved="en";try{saved=localStorage.getItem("lang")||"en"}catch(e){}setLang(saved);
setTheme(getTheme());
</script>
</body>
</html>`

// dashboardHTML is the user dashboard SPA.
var dashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Vibecast Dashboard</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--ink:#f6f8fa;--surface:#fff;--surface-2:#eef1f4;--line:#d0d7de;--text:#1f2328;--dim:#636c76;--accent:#1a7f37;--accent-dim:#2da44e;--warn:#9a6700;--danger:#cf222e;--placeholder:#8c959f;--mono:'JetBrains Mono','Fira Code',monospace;--sans:system-ui,-apple-system,sans-serif}
[data-theme="dark"]{--ink:#0c1117;--surface:#161b22;--surface-2:#1c2128;--line:#30363d;--text:#e6edf3;--dim:#7d8590;--accent:#39d353;--accent-dim:#238636;--warn:#d29922;--danger:#f85149;--placeholder:#484f58}
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:var(--sans);background:var(--ink);color:var(--text);min-height:100vh;-webkit-font-smoothing:antialiased}
a{color:var(--accent);text-decoration:none}
input,button,select{font-family:inherit}
.navbar{display:flex;justify-content:space-between;align-items:center;padding:0 1rem;height:52px;background:var(--surface);border-bottom:1px solid var(--line);position:sticky;top:0;z-index:100}
.navbar .logo{font-family:var(--mono);font-size:1rem;font-weight:700;color:var(--text)}
.navbar .logo .accent{color:var(--accent)}
.navbar .nav-right{display:flex;align-items:center;gap:.5rem;flex-wrap:nowrap}
.navbar .nav-right .email{font-family:var(--mono);font-size:.75rem;color:var(--dim);max-width:140px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.navbar .nav-right .btn-link{font-size:.8rem;cursor:pointer;color:var(--dim);background:none;border:none;white-space:nowrap;transition:color .15s}
.navbar .nav-right .btn-link:hover{color:var(--text)}
.navbar .admin-link{font-size:.8rem;color:var(--warn);font-weight:600;white-space:nowrap}
.navbar .btn-icon{font-size:.8rem;cursor:pointer;color:var(--dim);background:none;border:1px solid transparent;padding:4px 8px;border-radius:5px;transition:all .15s;white-space:nowrap}
.navbar .btn-icon:hover{background:var(--surface-2);border-color:var(--line);color:var(--text)}
.lang-toggle a{font-family:var(--mono);font-size:.75rem;cursor:pointer;padding:3px 7px;border-radius:4px;color:var(--dim);transition:all .15s}
.lang-toggle a.active{background:var(--surface-2);color:var(--accent)}
.container{max-width:1280px;margin:1.5rem auto;padding:0 1.5rem}
.card{background:var(--surface);border:1px solid var(--line);border-radius:8px;margin-bottom:1.5rem}
.card-header{padding:.85rem 1.1rem;border-bottom:1px solid var(--line);display:flex;justify-content:space-between;align-items:center}
.card-header h2{font-family:var(--mono);font-size:.85rem;font-weight:600;color:var(--text)}
.card-header .hint{font-size:.7rem;color:var(--dim)}
.card-body{padding:1.1rem}
.form-field{display:flex;flex-direction:column;gap:.25rem;margin-bottom:.75rem}
.form-field label{font-size:.75rem;font-weight:600;color:var(--dim)}
.form-field input{padding:8px 11px;border:1px solid var(--line);border-radius:6px;font-size:.85rem;outline:none;transition:border-color .15s;background:var(--ink);color:var(--text)}
.form-field input:focus{border-color:var(--accent)}
.form-field input::placeholder{color:var(--placeholder)}
.form-field .desc{font-size:.7rem;color:var(--dim)}
.form-actions{display:flex;justify-content:flex-end;margin-top:.5rem}
.btn{padding:8px 16px;border:none;border-radius:6px;font-size:.85rem;font-weight:600;cursor:pointer;transition:opacity .15s}
.btn-primary{background:var(--accent);color:var(--ink)}
.btn-primary:hover{background:var(--accent-dim)}
.btn-sm{padding:5px 10px;font-size:.75rem;border-radius:5px}
.btn-danger{background:var(--danger);color:#fff}
.btn-danger:hover{opacity:.85}
.btn-outline{background:transparent;border:1px solid var(--line);color:var(--text)}
.btn-outline:hover{background:var(--surface-2)}
.upload-btn{display:inline-block;padding:5px 10px;background:var(--accent-dim);color:var(--text);border-radius:5px;font-size:.75rem;cursor:pointer;font-weight:600;position:relative;overflow:hidden;white-space:nowrap}
.upload-btn:hover{background:var(--accent);color:var(--ink)}
.copy-btn{background:none;border:1px solid var(--line);border-radius:4px;padding:2px 6px;font-size:.7rem;cursor:pointer;color:var(--dim);line-height:1;vertical-align:middle}
.copy-btn:hover{color:var(--accent);border-color:var(--accent)}
.upload-btn input[type=file]{position:absolute;top:0;left:0;width:100%;height:100%;opacity:0;cursor:pointer}
.site-list{list-style:none}
.site-item{border:1px solid var(--line);border-radius:6px;margin-bottom:.6rem;overflow:hidden;transition:border-color .15s}
.site-item:hover{border-color:#424a53}
.site-item .site-head{display:flex;justify-content:space-between;align-items:center;padding:.7rem .9rem;cursor:pointer}
.site-item .site-head .info{flex:1;min-width:0}
.site-item .site-head .name{font-weight:600;font-size:.85rem;display:flex;align-items:center;gap:.4rem}
.site-item .site-head .url{font-family:var(--mono);font-size:.75rem;color:var(--accent);margin-top:.15rem;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.site-item .site-head .actions{display:flex;gap:.35rem;align-items:center;flex-shrink:0}
.site-item .site-detail{padding:.65rem .9rem;border-top:1px solid var(--line);background:var(--ink);font-size:.75rem;color:var(--dim);display:none}
.site-item .site-detail.show{display:block}
.site-item .site-detail .detail-row{display:flex;gap:.5rem;padding:.15rem 0}
.site-item .site-detail .detail-row .label{font-weight:600;min-width:80px;color:var(--dim)}
.site-item .site-detail .detail-row .value{word-break:break-all;color:var(--text)}
.status-dot{display:inline-block;width:7px;height:7px;border-radius:50%;flex-shrink:0}
.status-dot.public{background:var(--accent)}
.status-dot.protected{background:var(--warn)}
.status-dot.disabled{background:var(--danger)}
.badge{font-family:var(--mono);font-size:.65rem;padding:1px 6px;border-radius:3px;font-weight:500}
.badge-protected{background:rgba(210,153,34,.15);color:var(--warn)}
.badge-public{background:rgba(57,211,83,.15);color:var(--accent)}
.badge-disabled{background:rgba(248,81,73,.15);color:var(--danger)}
.empty{text-align:center;color:var(--dim);padding:2rem;font-size:.85rem}
.toast{position:fixed;bottom:1.5rem;right:1.5rem;padding:10px 18px;border-radius:6px;font-size:.85rem;z-index:999;opacity:0;transform:translateY(8px);transition:all .2s;font-family:var(--mono)}
.toast.show{opacity:1;transform:translateY(0)}
.toast.success{background:var(--accent-dim);color:var(--text);border:1px solid var(--accent)}
.toast.error{background:rgba(248,81,73,.15);color:var(--danger);border:1px solid var(--danger)}
.auth-screen{display:flex;align-items:center;justify-content:center;min-height:100vh;padding:1rem}
.auth-card{background:var(--surface);border:1px solid var(--line);border-radius:10px;padding:2.5rem;width:100%;max-width:360px}
.auth-card h1{font-family:var(--mono);font-size:1.4rem;font-weight:700;text-align:center;margin-bottom:.3rem}
.auth-card h1 .accent{color:var(--accent)}
.auth-card .subtitle{font-family:var(--mono);color:var(--dim);text-align:center;margin-bottom:1.5rem;font-size:.75rem}
.auth-field{margin-bottom:.75rem}
.auth-field label{display:block;font-size:.75rem;font-weight:600;margin-bottom:.2rem;color:var(--dim)}
.auth-field input{width:100%;padding:10px 12px;border:1px solid var(--line);border-radius:6px;font-size:.85rem;outline:none;transition:border-color .15s;background:var(--ink);color:var(--text)}
.auth-field input:focus{border-color:var(--accent)}
.auth-field input::placeholder{color:var(--placeholder)}
.auth-field .field-hint{font-size:.7rem;color:var(--dim);margin-top:.15rem}
.auth-card .btn{width:100%;margin-top:.3rem;font-family:var(--mono)}
.auth-card .switch{text-align:center;margin-top:1rem;font-size:.8rem;color:var(--dim)}
.auth-card .switch a{color:var(--accent);cursor:pointer;font-weight:500}
.captcha-label{font-size:.75rem;font-weight:600;margin-bottom:.2rem;color:var(--dim)}
.captcha-row{display:flex;align-items:stretch;gap:.5rem;margin-bottom:.75rem}
.captcha-row .captcha-img{flex-shrink:0;border:1px solid var(--line);border-radius:6px;background:var(--ink);cursor:pointer;height:40px;width:150px;object-fit:contain;transition:border-color .15s}
.captcha-row .captcha-img:hover{border-color:var(--accent)}
.captcha-row input{flex:1;min-width:0;padding:10px 12px;border:1px solid var(--line);border-radius:6px;font-size:.85rem;outline:none;background:var(--ink);color:var(--text);transition:border-color .15s}
.captcha-row input:focus{border-color:var(--accent)}
.captcha-row input::placeholder{color:var(--placeholder)}
.list-toolbar{display:flex;align-items:center;gap:.5rem;margin-bottom:.6rem}
.list-toolbar input[type=text]{flex:1;min-width:100px;padding:7px 11px;border:1px solid var(--line);border-radius:6px;font-size:.8rem;outline:none;background:var(--ink);color:var(--text)}
.list-toolbar input[type=text]:focus{border-color:var(--accent)}
.list-toolbar input[type=text]::placeholder{color:var(--placeholder)}
.pagination{display:flex;align-items:center;gap:.3rem;justify-content:center;padding:.6rem 0 0;flex-wrap:wrap}
.pagination .page-info{font-family:var(--mono);font-size:.7rem;color:var(--dim);margin:0 .3rem}
.pagination button{padding:3px 9px;border:1px solid var(--line);background:var(--surface);color:var(--text);border-radius:4px;font-size:.7rem;cursor:pointer;font-weight:500;font-family:var(--mono)}
.pagination button:hover:not(:disabled){background:var(--surface-2)}
.pagination button.active{background:var(--accent);color:var(--ink);border-color:var(--accent)}
.pagination button:disabled{opacity:.4;cursor:default}
.dashboard-grid{display:grid;grid-template-columns:320px 1fr;gap:1.5rem;align-items:start}
.sidebar{position:sticky;top:64px}
.main-content{min-width:0}
.modal-overlay{position:fixed;inset:0;background:rgba(0,0,0,.5);z-index:200;display:flex;align-items:center;justify-content:center;opacity:0;visibility:hidden;transition:opacity .2s,visibility .2s}
.modal-overlay.show{opacity:1;visibility:visible}
.modal{background:var(--surface);border:1px solid var(--line);border-radius:10px;padding:1.5rem;width:100%;max-width:380px;transform:scale(.95);transition:transform .2s}
.modal-overlay.show .modal{transform:scale(1)}
.modal h3{font-family:var(--mono);font-size:.95rem;font-weight:700;margin-bottom:1rem}
.modal p{font-size:.85rem;color:var(--dim);margin-bottom:1rem}
.modal .modal-field{margin-bottom:.65rem}
.modal .modal-field label{display:block;font-size:.75rem;font-weight:600;margin-bottom:.2rem;color:var(--dim)}
.modal .modal-field input{width:100%;padding:9px 11px;border:1px solid var(--line);border-radius:6px;font-size:.85rem;outline:none;background:var(--ink);color:var(--text)}
.modal .modal-field input:focus{border-color:var(--accent)}
.modal .modal-actions{display:flex;gap:.5rem;justify-content:flex-end;margin-top:1rem}
@media(max-width:768px){.dashboard-grid{grid-template-columns:1fr}.sidebar{position:static}}
@media(max-width:640px){.navbar{padding:0 .75rem}.navbar .nav-right .email{display:none}.navbar .btn-icon span{display:none}.lang-toggle a{padding:2px 5px;font-size:.7rem}}
@media(prefers-reduced-motion:reduce){*{transition:none!important}}
</style>
<script>var _t="light";try{_t=localStorage.getItem("theme")||"light"}catch(e){}document.documentElement.setAttribute("data-theme",_t)</script>
</head>
<body>
<div id="app"></div>
<script>
var BASE=location.pathname.replace(/\/(?:dashboard|admin)\/?$/,"/")||"/";var API=BASE+"api";
var currentUser=null;
function getTheme(){return document.documentElement.getAttribute("data-theme")||"light"}
function setTheme(t){document.documentElement.setAttribute("data-theme",t);try{localStorage.setItem("theme",t)}catch(e){}var b=document.getElementById("theme-toggle");if(b)b.textContent=t==="dark"?"☀":"🌙"}
function toggleTheme(){setTheme(getTheme()==="dark"?"light":"dark")}
var lang="en";
var sitePage=1,sitePerPage=10,siteTotal=0,siteSearch="";
var i18n={en:{siteName:"Site Name",siteNamePh:"e.g. My Portfolio",slug:"URL Slug",slugPh:"my-portfolio",sitePwd:"Access Password",sitePwdPh:"Leave empty for public",create:"Create Site",yourSites:"Your Sites",noSites:"No sites yet. Create one above.",deployZip:"Deploy ZIP",delete:"Delete",deleteConfirm:"Delete this site? This removes all files.",protected:"Protected",public:"Public",login:"Login",register:"Register",email:"Email",emailPh:"you@example.com",password:"Password",pwdHint:"At least 6 characters",noAccount:"No account?",haveAccount:"Have an account?",logout:"Logout",adminPanel:"Admin",deployed:"Deployed!",siteCreated:"Site created",deleted:"Deleted",loginFailed:"Login failed",registerFailed:"Registration failed",slugDesc:"Auto-generated from name if blank. a-z, 0-9, hyphens only.",pwdDesc:"If set, visitors need this password.",sitesHint:"Click to expand",deployHint:"Upload .zip to deploy or update",storagePath:"Storage",accessPassword:"Password",none:"None",accessDisabled:"Disabled",pwdRequired:"Public access disabled — password required",pwdOptional:"Optional password protection",search:"Search",searchSitesPh:"Search sites...",prev:"Prev",next:"Next",page:"Page",of:"of",captcha:"Captcha",captchaPh:"Answer",captchaLabel:"Verification",confirmPassword:"Confirm Password",pwdMismatch:"Passwords do not match",changePassword:"Change Password",currentPassword:"Current Password",newPassword:"New Password",newPasswordPh:"New password (min 6)",passwordChanged:"Password changed",cancel:"Cancel",save:"Save",confirm:"Confirm",emailRequired:"Please enter your email",pwdRequired:"Please enter your password",confirmRequired:"Please confirm your password",captchaRequired:"Please solve the captcha",copyUrl:"Copy URL",copied:"Copied!",copy:"Copy",copyFailed:"Copy failed",visit:"Visit",files:"Files",noFiles:"No files"},zh:{siteName:"站点名称",siteNamePh:"例如：我的作品集",slug:"URL Slug",slugPh:"my-portfolio",sitePwd:"访问密码",sitePwdPh:"留空则公开访问",create:"创建站点",yourSites:"我的站点",noSites:"还没有站点，在左侧创建一个。",deployZip:"部署 ZIP",delete:"删除",deleteConfirm:"确定删除此站点？所有文件将被移除。",protected:"已保护",public:"公开",login:"登录",register:"注册",email:"邮箱",emailPh:"you@example.com",password:"密码",pwdHint:"至少 6 个字符",noAccount:"没有账号？",haveAccount:"已有账号？",logout:"退出",adminPanel:"管理",deployed:"部署成功！",siteCreated:"站点已创建",deleted:"已删除",loginFailed:"登录失败",registerFailed:"注册失败",slugDesc:"留空则自动生成。仅限 a-z、0-9、连字符。",pwdDesc:"设置后，访问者需要输入此密码。",sitesHint:"点击展开详情",deployHint:"上传 .zip 文件来部署或更新站点",storagePath:"存储路径",accessPassword:"访问密码",none:"无",accessDisabled:"已禁用",pwdRequired:"公开访问已关闭 — 必须设置密码",pwdOptional:"可选的密码保护",search:"搜索",searchSitesPh:"搜索站点...",prev:"上一页",next:"下一页",page:"第",of:"/ 共",captcha:"验证码",captchaPh:"输入答案",captchaLabel:"验证",confirmPassword:"确认密码",pwdMismatch:"两次密码不一致",changePassword:"修改密码",currentPassword:"当前密码",newPassword:"新密码",newPasswordPh:"新密码（至少 6 位）",passwordChanged:"密码已修改",cancel:"取消",save:"保存",confirm:"确认",emailRequired:"请输入邮箱",pwdRequired:"请输入密码",confirmRequired:"请确认密码",captchaRequired:"请输入验证码",copyUrl:"复制链接",copied:"已复制！",copy:"复制",copyFailed:"复制失败",visit:"访问",files:"文件",noFiles:"暂无文件"}};
function t(k){return(i18n[lang]||i18n.en)[k]||(i18n.en[k]||k)}
function setLang(l){lang=l;try{localStorage.setItem("lang",l)}catch(e){}if(currentUser)renderDashboard();else renderAuth();}
function esc(s){return String(s||"").replace(/[&<>"']/g,function(c){return{"&":"&amp;","<":"&lt;",">":"&gt;",'"':"&quot;","'":"&#39;"}[c]})}
function siteUrl(u){return BASE+(u||"").replace(/^\//,"")}
function toast(msg,type){type=type||"success";var el=document.createElement("div");el.className="toast "+type+" show";el.textContent=msg;document.body.appendChild(el);var ms=type==="error"?5000:2500;setTimeout(function(){el.classList.remove("show");setTimeout(function(){el.remove()},200)},ms)}
function getToken(){try{return localStorage.getItem("vibecast_token")}catch(e){return""}}
function setToken(tk){try{localStorage.setItem("vibecast_token",tk)}catch(e){}}
function clearToken(){try{localStorage.removeItem("vibecast_token")}catch(e){}}
function customConfirm(msg,onConfirm){
var o=document.createElement("div");o.className="modal-overlay";o.style.zIndex="300";
o.innerHTML='<div class="modal"><h3>'+t("confirm")+'</h3><p>'+esc(msg)+'</p><div class="modal-actions"><button class="btn btn-outline" id="cf-no">'+t("cancel")+'</button><button class="btn btn-danger" id="cf-yes">'+t("confirm")+'</button></div></div>';
document.body.appendChild(o);setTimeout(function(){o.classList.add("show")},10);
var cl=function(){o.classList.remove("show");setTimeout(function(){o.remove()},200)};
o.addEventListener("click",function(e){if(e.target===o)cl()});
o.querySelector("#cf-no").addEventListener("click",cl);
o.querySelector("#cf-yes").addEventListener("click",function(){cl();onConfirm()})
}
function api(path,opts){
opts=opts||{};
var token=localStorage.getItem("vibecast_token")||"";
var headers={"Content-Type":"application/json","Accept-Language":lang||"en"};
if(token)headers["Authorization"]="Bearer "+token;
if(opts.headers)Object.assign(headers,opts.headers);
return fetch(API+path,Object.assign({},opts,{headers:headers,credentials:"same-origin"})).then(function(r){
if(r.status===401){try{localStorage.removeItem("vibecast_token")}catch(e){}var onAuth=/\/(dashboard|admin)\/?$/.test(location.pathname);if(!onAuth||getToken()){location.reload()}}
return r.json().catch(function(){return{error:"network error"}}).then(function(data){if(!r.ok)throw new Error(data.error||"request failed");return data})
})
}
function checkAuth(){if(!getToken())return Promise.resolve(false);return api("/auth/me").then(function(d){currentUser=d.data;return true}).catch(function(){clearToken();return false})}
var loginCaptchaId="",registerCaptchaId="";var regOpen=true;
function loadCaptcha(v){
return fetch(BASE+"api/auth/captcha").then(function(r){return r.json()}).then(function(d){
if(v==="login"){loginCaptchaId=d.data.id;var img=document.getElementById("login-captcha-img");if(img)img.src=d.data.image}
else{registerCaptchaId=d.data.id;var img=document.getElementById("register-captcha-img");if(img)img.src=d.data.image}
}).catch(function(){})
}
function renderAuth(){
var lh='<div class="lang-toggle" style="position:absolute;top:1rem;right:1rem"><a id="theme-toggle" onclick="toggleTheme()">'+(getTheme()==="dark"?"☀":"🌙")+'</a><a class="'+(lang==="en"?"active":"")+'" onclick="setLang(\'en\')">EN</a><a class="'+(lang==="zh"?"active":"")+'" onclick="setLang(\'zh\')">中文</a></div>';
document.getElementById("app").innerHTML=lh+'<div class="auth-screen"><div class="auth-card"><h1>Vibe<span class="accent">cast</span></h1><p class="subtitle">Build with vibe. Cast instantly.</p><div id="auth-form"></div></div></div>';
fetch(BASE+"api/settings").then(function(r){return r.json()}).then(function(d){regOpen=d.data&&d.data.openRegistration!==false;showLogin()}).catch(function(){regOpen=true;showLogin()});
}
function showLogin(){
document.getElementById("auth-form").innerHTML='<form onsubmit="doLogin();return false"><div class="auth-field"><label>'+t("email")+'</label><input id="email" type="email" placeholder="'+t("emailPh")+'" autocomplete="email"></div><div class="auth-field"><label>'+t("password")+'</label><input id="password" type="password" placeholder="'+t("password")+'" autocomplete="current-password"></div><div class="auth-field"><label class="captcha-label">'+t("captchaLabel")+'</label><div class="captcha-row"><img class="captcha-img" id="login-captcha-img" src="" alt="captcha" width="150" height="40" onclick="loadCaptcha(\'login\')" title="'+t("captcha")+'"><input id="login-captcha" type="text" placeholder="'+t("captchaPh")+'"></div></div><button class="btn btn-primary" type="submit">'+t("login")+'</button></form>'+(regOpen?'<p class="switch">'+t("noAccount")+' <a onclick="showRegister()">'+t("register")+'</a></p>':'');
loadCaptcha("login");
}
function showRegister(){
document.getElementById("auth-form").innerHTML='<form onsubmit="doRegister();return false"><div class="auth-field"><label>'+t("email")+'</label><input id="email" type="email" placeholder="'+t("emailPh")+'" autocomplete="email"></div><div class="auth-field"><label>'+t("password")+'</label><input id="password" type="password" placeholder="'+t("password")+'" autocomplete="new-password"><div class="field-hint">'+t("pwdHint")+'</div></div><div class="auth-field"><label>'+t("confirmPassword")+'</label><input id="confirm-pwd" type="password" placeholder="'+t("confirmPassword")+'" autocomplete="new-password"></div><div class="auth-field"><label class="captcha-label">'+t("captchaLabel")+'</label><div class="captcha-row"><img class="captcha-img" id="register-captcha-img" src="" alt="captcha" width="150" height="40" onclick="loadCaptcha(\'register\')" title="'+t("captcha")+'"><input id="register-captcha" type="text" placeholder="'+t("captchaPh")+'"></div></div><button class="btn btn-primary" type="submit">'+t("register")+'</button></form><p class="switch">'+t("haveAccount")+' <a onclick="showLogin()">'+t("login")+'</a></p>';
loadCaptcha("register");
}
function doLogin(){
var em=document.getElementById("email").value.trim(),pw=document.getElementById("password").value,ca=document.getElementById("login-captcha").value;
if(!em){toast(t("emailRequired"),"error");return}if(!pw){toast(t("pwdRequired"),"error");return}if(!ca){toast(t("captchaRequired"),"error");return}
api("/auth/login",{method:"POST",body:JSON.stringify({email:em,password:pw,captchaId:loginCaptchaId,captchaCode:ca})}).then(function(d){if(d.data&&d.data.token){setToken(d.data.token);location.reload()}else{toast(t("loginFailed"),"error")}}).catch(function(e){toast(e.message,"error");loadCaptcha("login")})
}
function doRegister(){
var em=document.getElementById("email").value.trim(),pw=document.getElementById("password").value,cf=document.getElementById("confirm-pwd").value,ca=document.getElementById("register-captcha").value;
if(!em){toast(t("emailRequired"),"error");return}if(!pw){toast(t("pwdRequired"),"error");return}if(!cf){toast(t("confirmRequired"),"error");return}if(!ca){toast(t("captchaRequired"),"error");return}if(pw!==cf){toast(t("pwdMismatch"),"error");return}
api("/auth/register",{method:"POST",body:JSON.stringify({email:em,password:pw,confirm:cf,captchaId:registerCaptchaId,captchaCode:ca})}).then(function(d){if(d.data&&d.data.token){setToken(d.data.token);location.reload()}else{toast(t("registerFailed"),"error")}}).catch(function(e){toast(e.message,"error");loadCaptcha("register")})
}
function doLogout(){api("/auth/logout",{method:"POST"}).then(function(){clearToken();location.href=BASE}).catch(function(){clearToken();location.href=BASE})}
function renderDashboard(){
var al=currentUser.isAdmin?'<a class="admin-link" href="'+BASE+'admin">'+t("adminPanel")+'</a>':'';
var th='<a id="theme-toggle" class="btn-icon" onclick="toggleTheme()">'+(getTheme()==="dark"?"☀":"🌙")+'</a>';
var lh='<div class="lang-toggle"><a class="'+(lang==="en"?"active":"")+'" onclick="setLang(\'en\')">EN</a><a class="'+(lang==="zh"?"active":"")+'" onclick="setLang(\'zh\')">中文</a></div>';
document.getElementById("app").innerHTML='<nav class="navbar"><div class="logo">Vibe<span class="accent">cast</span></div><div class="nav-right">'+al+'<button class="btn-icon" onclick="openChangePwdModal()">🔒<span> '+t("changePassword")+'</span></button>'+th+lh+'<span class="email">'+esc(currentUser.email)+'</span><button class="btn-link" onclick="doLogout()">'+t("logout")+'</button></div></nav><div class="container"><div class="dashboard-grid"><div class="sidebar"><div class="card"><div class="card-header"><h2>'+t("create")+'</h2></div><div class="card-body"><div class="form-field"><label>'+t("siteName")+'</label><input id="site-name" placeholder="'+t("siteNamePh")+'"></div><div class="form-field"><label>'+t("sitePwd")+'</label><input id="site-pwd" type="password" placeholder="'+t("sitePwdPh")+'"><div class="desc" id="pwd-desc">'+t("pwdDesc")+'</div></div><div class="form-actions"><button class="btn btn-primary" style="width:100%" onclick="createSite()">'+t("create")+'</button></div></div></div></div><div class="main-content"><div class="card"><div class="card-header"><h2>'+t("yourSites")+'</h2><span class="hint">'+t("sitesHint")+'</span></div><div class="card-body"><div class="list-toolbar"><input type="text" id="site-search" placeholder="'+t("searchSitesPh")+'" onkeydown="if(event.key===\'Enter\')searchSites()" value="'+esc(siteSearch)+'"></div><div id="site-list"></div></div></div></div></div></div><div class="modal-overlay" id="pwd-modal" onclick="if(event.target===this)closeChangePwdModal()"><div class="modal"><h3>'+t("changePassword")+'</h3><div class="modal-field"><label>'+t("currentPassword")+'</label><input id="old-pwd" type="password" placeholder="'+t("currentPassword")+'"></div><div class="modal-field"><label>'+t("newPassword")+'</label><input id="new-pwd" type="password" placeholder="'+t("newPasswordPh")+'"></div><div class="modal-field"><label>'+t("confirmPassword")+'</label><input id="confirm-new-pwd" type="password" placeholder="'+t("confirmPassword")+'"></div><div class="modal-actions"><button class="btn btn-outline" onclick="closeChangePwdModal()">'+t("cancel")+'</button><button class="btn btn-primary" onclick="changePassword()">'+t("save")+'</button></div></div></div>';
loadSites();
}
function openChangePwdModal(){document.getElementById("pwd-modal").classList.add("show")}
function closeChangePwdModal(){var m=document.getElementById("pwd-modal");m.classList.remove("show");m.querySelectorAll("input").forEach(function(i){i.value=""})}
function searchSites(){siteSearch=document.getElementById("site-search").value;sitePage=1;loadSites()}
function sitePageGo(p){sitePage=p;loadSites()}
function loadSites(){
var q=siteSearch?"&q="+encodeURIComponent(siteSearch):"";
api("/sites?page="+sitePage+"&perPage="+sitePerPage+q).then(function(d){
var r=d.data||{},sites=r.items||[];siteTotal=r.total||0;
var el=document.getElementById("site-list");
var pd=document.getElementById("pwd-desc");
if(sites.length>0&&sites[0].publicAccessDisabled&&pd){pd.textContent=t("pwdRequired");pd.style.color="var(--danger)"}
var tp=Math.ceil(siteTotal/sitePerPage)||1,pg=paginationHtml(sitePage,tp,"sitePageGo");
if(!sites.length){el.innerHTML='<div class="empty">'+t("noSites")+'</div>'+pg;return}
var h='<ul class="site-list">';
for(var i=0;i<sites.length;i++){var s=sites[i],dot="",badge="";
if(s.publicAccessDisabled&&!s.protected){dot="disabled";badge='<span class="badge badge-disabled">'+t("accessDisabled")+'</span>'}
else if(s.protected){dot="protected";badge='<span class="badge badge-protected">'+t("protected")+'</span>'}
else{dot="public";badge='<span class="badge badge-public">'+t("public")+'</span>'}
var pwd=s.protected?'<code style="font-family:var(--mono);font-size:.75rem;color:var(--text)">'+esc(s.password)+'</code> <button class="copy-btn" data-pwd="'+esc(s.password)+'" onclick="event.stopPropagation();copyText(this.getAttribute(\'data-pwd\'))" title="'+t("copy")+'">⧉</button>':'<span style="color:var(--dim)">'+t("none")+'</span>';
h+='<li class="site-item"><div class="site-head" onclick="toggleDetail('+s.id+')"><div class="info"><div class="name"><span class="status-dot '+dot+'"></span>'+esc(s.name)+' '+badge+'</div><div class="url">~/sites/'+esc(s.slug)+'/</div></div><div class="actions"><a class="btn btn-sm btn-outline" href="'+siteUrl(s.url)+'" target="_blank" onclick="event.stopPropagation()">'+t("visit")+'</a><label class="upload-btn" onclick="event.stopPropagation()">'+t("deployZip")+'<input type="file" accept=".zip" onchange="deploy('+s.id+',this.files[0])"></label><button class="btn btn-sm btn-danger" onclick="event.stopPropagation();delSite('+s.id+')">'+t("delete")+'</button></div></div><div class="site-detail" id="detail-'+s.id+'"><div class="detail-row"><span class="label">'+t("accessPassword")+'</span><span class="value">'+pwd+'</span></div></div></li>'}
h+='</ul><div style="font-size:.7rem;color:var(--dim);margin-top:.5rem;font-family:var(--mono)">'+t("deployHint")+'</div>'+pg;
el.innerHTML=h}).catch(function(e){toast(e.message,"error")})
}
function paginationHtml(p,tp,gf){
if(tp<=1)return"";
var h='<div class="pagination"><button '+(p<=1?"disabled":"")+' onclick="'+gf+'('+(p-1)+')">'+t("prev")+'</button>';
var s=Math.max(1,p-2),e=Math.min(tp,p+2);
if(s>1){h+='<button onclick="'+gf+'(1)">1</button>';if(s>2)h+='<span class="page-info">...</span>'}
for(var i=s;i<=e;i++){h+='<button class="'+(i===p?"active":"")+'" onclick="'+gf+'('+i+')">'+i+'</button>'}
if(e<tp){if(e<tp-1)h+='<span class="page-info">...</span>';h+='<button onclick="'+gf+'('+tp+')">'+tp+'</button>'}
h+='<button '+(p>=tp?"disabled":"")+' onclick="'+gf+'('+(p+1)+')">'+t("next")+'</button><span class="page-info">'+t("page")+' '+p+' '+t("of")+' '+tp+'</span></div>';
return h
}
function toggleDetail(id){var e=document.getElementById("detail-"+id);if(e){var wasShown=e.classList.contains("show");e.classList.toggle("show");if(!wasShown)loadFileTree(id)}}
function loadFileTree(id){
api("/sites/"+id+"/files").then(function(d){var files=d.data||[];var el=document.getElementById("detail-"+id);if(!el)return;
var h=el.innerHTML;
var ft='<div class="detail-row"><span class="label">'+t("files")+'</span><span class="value">';
if(!files.length){ft+=t("noFiles")+'</span></div>';el.innerHTML=h+ft;return}
ft+='<div style="margin-top:.3rem">';
for(var i=0;i<files.length;i++){var f=files[i],icon=f.dir?"📁":"📄",sz=f.dir?"-":formatSize(f.size);
ft+='<div style="display:flex;justify-content:space-between;align-items:baseline;padding:1px 0;font-family:var(--mono);font-size:.75rem"><span style="overflow:hidden;text-overflow:ellipsis;white-space:nowrap">'+icon+' '+esc(f.name)+(f.dir?"/":"")+'</span><span style="color:var(--dim);flex-shrink:0;margin-left:1rem;text-align:right;min-width:70px">'+sz+'</span></div>'}
ft+='</div></div>';
el.innerHTML=h+ft}).catch(function(){})
}
function formatSize(n){if(n<1024)return n+" B";if(n<1048576)return(n/1024).toFixed(1)+" KB";if(n<1073741824)return(n/1048576).toFixed(1)+" MB";return(n/1073741824).toFixed(1)+" GB"}
function copyText(text){if(navigator.clipboard){navigator.clipboard.writeText(text).then(function(){toast(t("copied"))}).catch(function(){fallbackCopy(text)})}else{fallbackCopy(text)}}
function fallbackCopy(text){var ta=document.createElement("textarea");ta.value=text;ta.style.position="fixed";ta.style.opacity="0";document.body.appendChild(ta);ta.select();try{document.execCommand("copy");toast(t("copied"))}catch(e){toast(t("copyFailed"),"error")}document.body.removeChild(ta)}
function createSite(){
api("/sites",{method:"POST",body:JSON.stringify({name:document.getElementById("site-name").value,password:document.getElementById("site-pwd").value})}).then(function(){document.getElementById("site-name").value="";document.getElementById("site-pwd").value="";toast(t("siteCreated"));loadSites()}).catch(function(e){toast(e.message,"error")})
}
function deploy(id,file){
if(!file)return;var fd=new FormData();fd.append("file",file);
var token=localStorage.getItem("vibecast_token")||"";
fetch(API+"/sites/"+id+"/deploy",{method:"POST",body:fd,headers:{"Authorization":"Bearer "+token,"Accept-Language":lang||"en"},credentials:"same-origin"}).then(function(r){return r.json().catch(function(){return{error:"network error"}}).then(function(data){if(!r.ok)throw new Error(data.error||"request failed");return data})}).then(function(){toast(t("deployed"));loadSites()}).catch(function(e){toast(e.message,"error")})
}
function delSite(id){
customConfirm(t("deleteConfirm"),function(){api("/sites/"+id,{method:"DELETE"}).then(function(){toast(t("deleted"));loadSites()}).catch(function(e){toast(e.message,"error")})})
}
function changePassword(){
var o=document.getElementById("old-pwd").value,n=document.getElementById("new-pwd").value,c=document.getElementById("confirm-new-pwd").value;
if(!o||!n){toast(t("currentPassword")+" & "+t("newPassword"),"error");return}
if(n.length<6){toast(t("pwdHint"),"error");return}
if(n!==c){toast(t("pwdMismatch"),"error");return}
api("/auth/change-password",{method:"PUT",body:JSON.stringify({oldPassword:o,newPassword:n})}).then(function(){toast(t("passwordChanged"));closeChangePwdModal()}).catch(function(e){toast(e.message,"error")})
}
var sl="en";try{sl=localStorage.getItem("lang")||"en"}catch(e){}lang=sl;
checkAuth().then(function(ok){if(ok)renderDashboard();else renderAuth()})
</script>
</body>
</html>`

// adminPageHTML is the admin dashboard SPA.
var adminPageHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Vibecast Admin</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--ink:#f6f8fa;--surface:#fff;--surface-2:#eef1f4;--line:#d0d7de;--text:#1f2328;--dim:#636c76;--accent:#1a7f37;--accent-dim:#2da44e;--warn:#9a6700;--danger:#cf222e;--placeholder:#8c959f;--mono:'JetBrains Mono','Fira Code',monospace;--sans:system-ui,-apple-system,sans-serif}
[data-theme="dark"]{--ink:#0c1117;--surface:#161b22;--surface-2:#1c2128;--line:#30363d;--text:#e6edf3;--dim:#7d8590;--accent:#39d353;--accent-dim:#238636;--warn:#d29922;--danger:#f85149;--placeholder:#484f58}
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:var(--sans);background:var(--ink);color:var(--text);min-height:100vh;-webkit-font-smoothing:antialiased}
a{color:var(--accent);text-decoration:none}
button{font-family:inherit;cursor:pointer}
.navbar{display:flex;justify-content:space-between;align-items:center;padding:0 1rem;height:52px;background:var(--surface);border-bottom:1px solid var(--line);position:sticky;top:0;z-index:100}
.navbar .logo{font-family:var(--mono);font-size:1rem;font-weight:700;color:var(--warn)}
.navbar .nav-right{display:flex;align-items:center;gap:.75rem}
.navbar .nav-right .btn-link{font-size:.8rem;cursor:pointer;color:var(--dim);background:none;border:none;transition:color .15s}
.navbar .nav-right .btn-link:hover{color:var(--text)}
.lang-toggle a{font-family:var(--mono);font-size:.75rem;cursor:pointer;padding:3px 7px;border-radius:4px;color:var(--dim);transition:all .15s}
.lang-toggle a.active{background:var(--surface-2);color:var(--accent)}
.container{max-width:1280px;margin:1.5rem auto;padding:0 1.5rem}
.card{background:var(--surface);border:1px solid var(--line);border-radius:8px;margin-bottom:1.5rem}
.card-header{padding:.85rem 1.1rem;border-bottom:1px solid var(--line)}
.card-header h2{font-family:var(--mono);font-size:.85rem;font-weight:600}
.card-body{padding:1.1rem}
.stats-grid{display:grid;grid-template-columns:repeat(3,1fr);gap:1rem}
.stat-card{background:var(--ink);border:1px solid var(--line);border-radius:6px;padding:1.1rem;text-align:center}
.stat-card .num{font-family:var(--mono);font-size:1.6rem;font-weight:700;color:var(--accent)}
.stat-card .label{font-size:.7rem;color:var(--dim);margin-top:.15rem;font-family:var(--mono)}
table{width:100%;border-collapse:collapse}
th,td{padding:.5rem .7rem;text-align:left;border-bottom:1px solid var(--line)}
th{color:var(--dim);font-family:var(--mono);font-size:.7rem;font-weight:600;text-transform:uppercase}
td{font-size:.8rem}
.badge{font-family:var(--mono);font-size:.65rem;padding:1px 6px;border-radius:3px;font-weight:500}
.badge-admin{background:rgba(210,153,34,.15);color:var(--warn)}
.badge-user{background:var(--surface-2);color:var(--dim)}
.badge-protected{background:rgba(210,153,34,.15);color:var(--warn)}
.badge-public{background:rgba(57,211,83,.15);color:var(--accent)}
.badge-disabled{background:rgba(248,81,73,.15);color:var(--danger)}
.btn{padding:5px 11px;border:none;border-radius:5px;font-size:.75rem;font-weight:600}
.btn-promote{background:var(--accent-dim);color:var(--text)}
.btn-demote{background:var(--warn);color:var(--ink)}
.btn-danger{background:var(--danger);color:#fff}
.btn:hover{opacity:.85}
.toggle-row{display:flex;align-items:center;justify-content:space-between;padding:.6rem 0;border-bottom:1px solid var(--line)}
.toggle-row:last-child{border-bottom:none}
.toggle-row .toggle-info{flex:1}
.toggle-row .toggle-label{font-size:.85rem;font-weight:500}
.toggle-row .toggle-desc{font-size:.7rem;color:var(--dim);margin-top:.1rem}
.toggle-switch{position:relative;width:38px;height:20px;background:var(--line);border-radius:10px;cursor:pointer;transition:background .2s;flex-shrink:0}
.toggle-switch.on{background:var(--accent)}
.toggle-switch::after{content:"";position:absolute;top:2px;left:2px;width:16px;height:16px;background:#fff;border-radius:50%;transition:transform .2s}
.toggle-switch.on::after{transform:translateX(18px)}
.empty{text-align:center;color:var(--dim);padding:1.5rem;font-size:.8rem}
.toast{position:fixed;bottom:1.5rem;right:1.5rem;padding:10px 18px;border-radius:6px;font-size:.85rem;z-index:999;opacity:0;transform:translateY(8px);transition:all .2s;font-family:var(--mono)}
.toast.show{opacity:1;transform:translateY(0)}
.toast.success{background:var(--accent-dim);color:var(--text);border:1px solid var(--accent)}
.toast.error{background:rgba(248,81,73,.15);color:var(--danger);border:1px solid var(--danger)}
.list-toolbar{display:flex;align-items:center;gap:.5rem;margin-bottom:.6rem}
.list-toolbar input[type=text]{flex:1;min-width:100px;padding:7px 11px;border:1px solid var(--line);border-radius:6px;font-size:.8rem;outline:none;background:var(--ink);color:var(--text)}
.list-toolbar input[type=text]:focus{border-color:var(--accent)}
.list-toolbar input[type=text]::placeholder{color:var(--placeholder)}
.pagination{display:flex;align-items:center;gap:.3rem;justify-content:center;padding:.6rem 0 0;flex-wrap:wrap}
.pagination .page-info{font-family:var(--mono);font-size:.7rem;color:var(--dim);margin:0 .3rem}
.pagination button{padding:3px 9px;border:1px solid var(--line);background:var(--surface);color:var(--text);border-radius:4px;font-size:.7rem;cursor:pointer;font-weight:500;font-family:var(--mono)}
.pagination button:hover:not(:disabled){background:var(--surface-2)}
.pagination button.active{background:var(--accent);color:var(--ink);border-color:var(--accent)}
.pagination button:disabled{opacity:.4;cursor:default}
textarea{font-family:var(--mono);background:var(--ink);color:var(--text);border:1px solid var(--line);border-radius:6px;font-size:.8rem}
.site-row{cursor:pointer;transition:background .15s}
.site-row:hover{background:var(--surface-2)}
.site-row.expanded{background:var(--surface-2)}
.admin-detail-row{display:none}
.admin-detail-row.show{display:table-row}
.admin-detail-row>td{padding:0!important;background:var(--ink);border-bottom:1px solid var(--line)}
.modal-overlay{position:fixed;inset:0;background:rgba(0,0,0,.5);z-index:200;display:flex;align-items:center;justify-content:center;opacity:0;visibility:hidden;transition:opacity .2s,visibility .2s}
.modal-overlay.show{opacity:1;visibility:visible}
.modal{background:var(--surface);border:1px solid var(--line);border-radius:10px;padding:1.5rem;width:100%;max-width:380px;transform:scale(.95);transition:transform .2s}
.modal-overlay.show .modal{transform:scale(1)}
.modal h3{font-family:var(--mono);font-size:.95rem;font-weight:700;margin-bottom:1rem}
.modal p{font-size:.85rem;color:var(--dim);margin-bottom:1rem}
.modal .modal-actions{display:flex;gap:.5rem;justify-content:flex-end;margin-top:1rem}
@media(max-width:640px){.navbar{padding:0 .75rem}.stats-grid{grid-template-columns:1fr}.lang-toggle a{padding:2px 5px;font-size:.7rem}}
@media(prefers-reduced-motion:reduce){*{transition:none!important}}
</style>
<script>var _t="light";try{_t=localStorage.getItem("theme")||"light"}catch(e){}document.documentElement.setAttribute("data-theme",_t)</script>
</head>
<body>
<div id="app"></div>
<script>
var BASE=location.pathname.replace(/\/(?:dashboard|admin)\/?$/,"/")||"/";var API=BASE+"api";
function getTheme(){return document.documentElement.getAttribute("data-theme")||"light"}
function setTheme(t){document.documentElement.setAttribute("data-theme",t);try{localStorage.setItem("theme",t)}catch(e){}var b=document.getElementById("theme-toggle");if(b)b.textContent=t==="dark"?"☀":"🌙"}
function toggleTheme(){setTheme(getTheme()==="dark"?"light":"dark")}
var lang="en";
var userPage=1,userPerPage=10,userTotal=0,userSearch="";
var adminSitePage=1,adminSitePerPage=10,adminSiteTotal=0,adminSiteSearch="";
var i18n={en:{overview:"Overview",settings:"Settings",users:"Users",allSites:"All Sites",userCount:"Users",siteCount:"Sites",adminCount:"Admins",openReg:"Open Registration",openRegDesc:"When enabled, anyone can register a new account.",publicAccess:"Public Site Access",publicAccessDesc:"When disabled, all deployed sites return 403 to visitors.",dashboard:"Dashboard",logout:"Logout",role:"Role",created:"Created",actions:"Actions",promote:"Promote",demote:"Demote",delete:"Delete",deleteUserConfirm:"Delete this user and all their sites?",deleteSiteConfirm:"Delete this site?",name:"Name",slug:"Slug",owner:"Owner",protected:"Protected",public:"Public",password:"Password",storagePath:"Storage",url:"URL",noUsers:"No users",noSites:"No sites",updated:"Updated",deleted:"Deleted",roleUpdated:"Role updated",userDeleted:"User deleted",siteDeleted:"Site deleted",settingsSaved:"Settings saved",none:"None",accessDisabled:"Disabled",search:"Search",searchUsersPh:"Search users...",searchSitesPh:"Search sites...",prev:"Prev",next:"Next",page:"Page",of:"of",ownerEmail:"Owner",domainRestriction:"Email Domain Restriction",domainRestrictionDesc:"Only allow registration from specified domains (one per line).",save:"Save",cancel:"Cancel",confirm:"Confirm",cleanup:"Storage Cleanup",cleanupDesc:"Remove orphaned directories that exist on disk but have no corresponding site in the database.",checkOrphans:"Check for Orphans",noOrphans:"No orphaned directories found",orphanFound:"Orphaned directories found",confirmDelete:"Delete All",cleanupConfirmWarning:"Delete all orphaned directories? This cannot be undone.",cleanupDone:"Cleanup complete",scanning:"Scanning...",files:"Files",noFiles:"No files",systemUpdate:"System Update",updateCheck:"Check for Updates",updateAvailable:"Update Available",updateToDate:"Up to Date",currentVersion:"Current Version",latestVersion:"Latest Version",updateNow:"Update Now",updating:"Updating...",updateSuccess:"Update successful! Please restart vibecast.",updateFailed:"Update failed",noUpdate:"You are running the latest version.",checkingUpdate:"Checking...",releaseNotes:"Release Notes",viewRelease:"View on GitHub"},zh:{overview:"概览",settings:"设置",users:"用户",allSites:"全部站点",userCount:"用户",siteCount:"站点",adminCount:"管理员",openReg:"开放注册",openRegDesc:"开启后，任何人都可以注册新账号。",publicAccess:"公开站点访问",publicAccessDesc:"关闭后，所有已部署站点对访问者返回 403。",dashboard:"Dashboard",logout:"退出",role:"角色",created:"创建时间",actions:"操作",promote:"提升",demote:"降级",delete:"删除",deleteUserConfirm:"删除此用户及其所有站点？",deleteSiteConfirm:"删除此站点？",name:"名称",slug:"Slug",owner:"所有者",protected:"已保护",public:"公开",password:"密码",storagePath:"存储路径",url:"URL",noUsers:"暂无用户",noSites:"暂无站点",updated:"已更新",deleted:"已删除",roleUpdated:"角色已更新",userDeleted:"用户已删除",siteDeleted:"站点已删除",settingsSaved:"设置已保存",none:"无",accessDisabled:"已禁用",search:"搜索",searchUsersPh:"搜索用户...",searchSitesPh:"搜索站点...",prev:"上一页",next:"下一页",page:"第",of:"/ 共",ownerEmail:"所有者",domainRestriction:"邮箱域名限制",domainRestrictionDesc:"仅允许指定域名的邮箱注册（每行一个）。",save:"保存",cancel:"取消",confirm:"确认",cleanup:"存储清理",cleanupDesc:"删除磁盘上存在但数据库中无对应站点的孤立目录。",checkOrphans:"检查孤立目录",noOrphans:"未发现孤立目录",orphanFound:"发现孤立目录",confirmDelete:"全部删除",cleanupConfirmWarning:"确认删除所有孤立目录？此操作不可撤销。",cleanupDone:"清理完成",scanning:"扫描中...",files:"文件",noFiles:"暂无文件",systemUpdate:"系统更新",updateCheck:"检查更新",updateAvailable:"有可用更新",updateToDate:"已是最新版本",currentVersion:"当前版本",latestVersion:"最新版本",updateNow:"立即更新",updating:"更新中...",updateSuccess:"更新成功！请重启 vibecast。",updateFailed:"更新失败",noUpdate:"当前已是最新版本。",checkingUpdate:"检查中...",releaseNotes:"更新日志",viewRelease:"在 GitHub 查看"}};
function t(k){return(i18n[lang]||i18n.en)[k]||(i18n.en[k]||k)}
function setLang(l){lang=l;try{localStorage.setItem("lang",l)}catch(e){}renderAdmin()}
function esc(s){return String(s||"").replace(/[&<>"']/g,function(c){return{"&":"&amp;","<":"&lt;",">":"&gt;",'"':"&quot;","'":"&#39;"}[c]})}
function siteUrl(u){return BASE+(u||"").replace(/^\//,"")}
function fmtDate(s){if(!s)return"-";var d=new Date(s);return d.toLocaleDateString()+" "+d.toLocaleTimeString([],{hour:"2-digit",minute:"2-digit"})}
function toast(msg,type){type=type||"success";var el=document.createElement("div");el.className="toast "+type+" show";el.textContent=msg;document.body.appendChild(el);var ms=type==="error"?5000:2500;setTimeout(function(){el.classList.remove("show");setTimeout(function(){el.remove()},200)},ms)}
function getToken(){try{return localStorage.getItem("vibecast_token")}catch(e){return""}}
function clearToken(){try{localStorage.removeItem("vibecast_token")}catch(e){}}
function customConfirm(msg,onConfirm){
var o=document.createElement("div");o.className="modal-overlay";o.style.zIndex="300";
o.innerHTML='<div class="modal"><h3>'+t("confirm")+'</h3><p>'+esc(msg)+'</p><div class="modal-actions"><button class="btn btn-outline" style="background:transparent;border:1px solid var(--line);color:var(--text)" id="cf-no">'+t("cancel")+'</button><button class="btn btn-danger" id="cf-yes">'+t("confirm")+'</button></div></div>';
document.body.appendChild(o);setTimeout(function(){o.classList.add("show")},10);
var cl=function(){o.classList.remove("show");setTimeout(function(){o.remove()},200)};
o.addEventListener("click",function(e){if(e.target===o)cl()});
o.querySelector("#cf-no").addEventListener("click",cl);
o.querySelector("#cf-yes").addEventListener("click",function(){cl();onConfirm()})
}
function api(path,opts){
opts=opts||{};
var token=localStorage.getItem("vibecast_token")||"";
var headers={"Content-Type":"application/json","Accept-Language":lang||"en"};
if(token)headers["Authorization"]="Bearer "+token;
if(opts.headers)Object.assign(headers,opts.headers);
return fetch(API+path,Object.assign({},opts,{headers:headers,credentials:"same-origin"})).then(function(r){
if(r.status===401){try{localStorage.removeItem("vibecast_token")}catch(e){}location.href=BASE+"dashboard"}
return r.json().catch(function(){return{error:"network error"}}).then(function(data){if(!r.ok)throw new Error(data.error||"request failed");return data})
})
}
function checkAuth(){if(!getToken())return Promise.reject(new Error("no token"));return api("/auth/me").then(function(d){if(!d.data||!d.data.isAdmin)throw new Error("not admin");return d.data}).catch(function(){clearToken();location.href=BASE+"dashboard"})}
function doLogout(){api("/auth/logout",{method:"POST"}).then(function(){clearToken();location.href=BASE}).catch(function(){clearToken();location.href=BASE})}
function renderAdmin(){
var lh='<div class="lang-toggle"><a id="theme-toggle" class="btn-icon" onclick="toggleTheme()">'+(getTheme()==="dark"?"☀":"🌙")+'</a><a class="'+(lang==="en"?"active":"")+'" onclick="setLang(\'en\')">EN</a><a class="'+(lang==="zh"?"active":"")+'" onclick="setLang(\'zh\')">中文</a></div>';
document.getElementById("app").innerHTML='<nav class="navbar"><div class="logo">Vibecast Admin <span id="ver" style="font-size:.6rem;color:var(--dim);font-weight:400"></span></div><div class="nav-right">'+lh+'<a href="'+BASE+'dashboard" style="font-size:.8rem;color:var(--dim)">'+t("dashboard")+'</a><button class="btn-link" onclick="doLogout()">'+t("logout")+'</button></div></nav><div class="container"><div class="card"><div class="card-header"><h2>'+t("overview")+'</h2></div><div class="card-body"><div id="stats" class="stats-grid"></div></div></div><div class="card"><div class="card-header"><h2>'+t("settings")+'</h2></div><div class="card-body"><div id="settings"></div></div></div><div class="card"><div class="card-header"><h2>'+t("cleanup")+'</h2></div><div class="card-body"><div id="cleanup-section"><p style="font-size:.8rem;color:var(--dim);margin-bottom:.8rem">'+t("cleanupDesc")+'</p><button class="btn btn-outline" onclick="checkCleanup()">'+t("checkOrphans")+'</button><div id="cleanup-result" style="margin-top:.8rem"></div></div></div></div><div class="card"><div class="card-header"><h2>'+t("systemUpdate")+'</h2></div><div class="card-body"><div id="update-section"><p style="font-size:.8rem;color:var(--dim);margin-bottom:.8rem">'+t("currentVersion")+': <span id="update-current-ver" style="font-family:var(--mono);color:var(--text)">—</span></p><button class="btn btn-outline" id="update-check-btn" onclick="checkUpdate()">'+t("updateCheck")+'</button><div id="update-result" style="margin-top:.8rem"></div></div></div></div><div class="card"><div class="card-header"><h2>'+t("users")+'</h2></div><div class="card-body"><div class="list-toolbar"><input type="text" id="user-search" placeholder="'+t("searchUsersPh")+'" onkeydown="if(event.key===\'Enter\')searchUsers()" value="'+esc(userSearch)+'"></div><div id="users"></div></div></div><div class="card"><div class="card-header"><h2>'+t("allSites")+'</h2></div><div class="card-body"><div class="list-toolbar"><input type="text" id="admin-site-search" placeholder="'+t("searchSitesPh")+'" onkeydown="if(event.key===\'Enter\')searchAdminSites()" value="'+esc(adminSiteSearch)+'"></div><div id="sites"></div></div></div></div>';
loadStats();loadSettings();loadUsers();loadSites();
fetch(BASE+"api/version").then(function(r){return r.json()}).then(function(d){var ver=(d.data?d.data.version:"dev");var v=document.getElementById("ver");if(v)v.textContent="v"+ver;var uv=document.getElementById("update-current-ver");if(uv)uv.textContent="v"+ver}).catch(function(){})
}
function loadStats(){
api("/admin/stats").then(function(d){var s=d.data;document.getElementById("stats").innerHTML='<div class="stat-card"><div class="num">'+s.users+'</div><div class="label">'+t("userCount")+'</div></div><div class="stat-card"><div class="num">'+s.sites+'</div><div class="label">'+t("siteCount")+'</div></div><div class="stat-card"><div class="num">'+s.admins+'</div><div class="label">'+t("adminCount")+'</div></div>'}).catch(function(e){toast(e.message,"error")})
}
function loadSettings(){
api("/admin/settings").then(function(d){var s=d.data;
var regOn=s.openRegistration,pubOn=s.allowPublicAccess!==false,drOn=s.domainRestriction;
var h='<div class="toggle-row"><div class="toggle-info"><div class="toggle-label">'+t("openReg")+'</div><div class="toggle-desc">'+t("openRegDesc")+'</div></div><div class="toggle-switch '+(regOn?"on":"")+'" onclick="toggleReg()"></div></div><div class="toggle-row"><div class="toggle-info"><div class="toggle-label">'+t("publicAccess")+'</div><div class="toggle-desc">'+t("publicAccessDesc")+'</div></div><div class="toggle-switch '+(pubOn?"on":"")+'" onclick="togglePub()"></div></div><div class="toggle-row"><div class="toggle-info"><div class="toggle-label">'+t("domainRestriction")+'</div><div class="toggle-desc">'+t("domainRestrictionDesc")+'</div></div><div class="toggle-switch '+(drOn?"on":"")+'" onclick="toggleDomainRestriction()"></div></div>';
if(drOn){h+='<div style="padding:.6rem 0"><textarea id="allowed-domains" rows="4" style="width:100%;padding:8px" placeholder="example.com&#10;gmail.com">'+esc(s.allowedDomains||"")+'</textarea><div style="margin-top:.4rem"><button class="btn btn-promote" onclick="saveDomains()">'+t("save")+'</button></div></div>'}
document.getElementById("settings").innerHTML=h}).catch(function(e){toast(e.message,"error")})
}
function toggleReg(){api("/admin/settings").then(function(d){var c=d.data;return api("/admin/settings",{method:"PUT",body:JSON.stringify({openRegistration:!c.openRegistration,allowPublicAccess:c.allowPublicAccess!==false,domainRestriction:c.domainRestriction,allowedDomains:c.allowedDomains||""})}).then(function(){toast(t("settingsSaved"));loadSettings()})}).catch(function(e){toast(e.message,"error")})}
function togglePub(){api("/admin/settings").then(function(d){var c=d.data;return api("/admin/settings",{method:"PUT",body:JSON.stringify({openRegistration:c.openRegistration,allowPublicAccess:!(c.allowPublicAccess!==false),domainRestriction:c.domainRestriction,allowedDomains:c.allowedDomains||""})}).then(function(){toast(t("settingsSaved"));loadSettings()})}).catch(function(e){toast(e.message,"error")})}
function toggleDomainRestriction(){api("/admin/settings").then(function(d){var c=d.data;return api("/admin/settings",{method:"PUT",body:JSON.stringify({openRegistration:c.openRegistration,allowPublicAccess:c.allowPublicAccess!==false,domainRestriction:!c.domainRestriction,allowedDomains:c.allowedDomains||""})}).then(function(){toast(t("settingsSaved"));loadSettings()})}).catch(function(e){toast(e.message,"error")})}
function saveDomains(){var v=document.getElementById("allowed-domains").value;api("/admin/settings").then(function(d){return api("/admin/settings",{method:"PUT",body:JSON.stringify({openRegistration:d.data.openRegistration,allowPublicAccess:d.data.allowPublicAccess!==false,domainRestriction:d.data.domainRestriction,allowedDomains:v})}).then(function(){toast(t("settingsSaved"));loadSettings()})}).catch(function(e){toast(e.message,"error")})}
function checkUpdate(){
var el=document.getElementById("update-result");var btn=document.getElementById("update-check-btn");
if(btn){btn.disabled=true;btn.textContent=t("checkingUpdate")}
el.innerHTML='<span style="color:var(--dim);font-size:.8rem">'+t("checkingUpdate")+'</span>';
api("/admin/update/check").then(function(d){
var r=d.data||{};
if(btn){btn.disabled=false;btn.textContent=t("updateCheck")}
if(r.error){el.innerHTML='<span style="color:var(--danger);font-size:.8rem">✗ '+esc(r.error)+'</span>';return}
if(r.updateAvailable){
var h='<div style="margin-bottom:.6rem"><span style="color:var(--accent);font-size:.8rem;font-family:var(--mono)">✓ '+t("updateAvailable")+'</span></div>';
h+='<div style="font-size:.8rem;color:var(--dim);margin-bottom:.4rem">'+t("currentVersion")+': <span style="font-family:var(--mono);color:var(--text)">v'+esc(r.currentVersion)+'</span> → '+t("latestVersion")+': <span style="font-family:var(--mono);color:var(--accent)">v'+esc(r.latestVersion)+'</span></div>';
if(r.releaseNotes){h+='<details style="margin-bottom:.6rem"><summary style="font-size:.75rem;cursor:pointer;color:var(--dim)">'+t("releaseNotes")+'</summary><div style="font-size:.75rem;color:var(--dim);margin-top:.4rem;white-space:pre-wrap;max-height:200px;overflow-y:auto">'+esc(r.releaseNotes)+'</div></details>'}
h+='<button class="btn btn-promote" id="update-apply-btn" onclick="applyUpdate()">'+t("updateNow")+'</button>';
el.innerHTML=h}else{
var h2='<span style="color:var(--dim);font-size:.8rem">✓ '+t("noUpdate")+'</span>';
if(r.latestVersion){h2+=' <span style="color:var(--dim);font-size:.75rem;font-family:var(--mono)">v'+esc(r.latestVersion)+'</span>'}
el.innerHTML=h2}
}).catch(function(e){if(btn){btn.disabled=false;btn.textContent=t("updateCheck")}el.innerHTML='<span style="color:var(--danger);font-size:.8rem">'+esc(e.message)+'</span>'})
}
function applyUpdate(){
var btn=document.getElementById("update-apply-btn");var el=document.getElementById("update-result");
if(btn){btn.disabled=true;btn.textContent=t("updating")}
el.innerHTML+='<div style="margin-top:.6rem"><span style="color:var(--dim);font-size:.8rem">'+t("updating")+'</span></div>';
customConfirm(t("updateNow")+"?",function(){
api("/admin/update/apply",{method:"POST"}).then(function(d){
var r=d.data||{};
var h='<div style="color:var(--accent);font-size:.8rem;font-family:var(--mono)">✓ '+t("updateSuccess")+'</div>';
if(r.previousVersion&&r.newVersion){h+='<div style="font-size:.75rem;color:var(--dim);margin-top:.3rem">v'+esc(r.previousVersion)+' → v'+esc(r.newVersion)+'</div>'}
el.innerHTML=h;
toast(t("updateSuccess"),"success")
}).catch(function(e){
var h='<div style="color:var(--danger);font-size:.8rem">✗ '+t("updateFailed")+': '+esc(e.message)+'</div><button class="btn btn-outline" style="margin-top:.4rem" onclick="checkUpdate()">'+t("updateCheck")+'</button>';
el.innerHTML=h;toast(e.message,"error")
})
})
}
function checkCleanup(){var el=document.getElementById("cleanup-result");el.innerHTML='<span style="color:var(--dim);font-size:.8rem">'+t("scanning")+'</span>';api("/admin/cleanup").then(function(d){var r=d.data||{},orphans=r.orphans||[];if(!orphans.length){el.innerHTML='<span style="color:var(--accent);font-size:.8rem">✓ '+t("noOrphans")+'</span>';return}var h='<div style="margin-bottom:.6rem"><span style="color:var(--warn);font-size:.8rem;font-family:var(--mono)">'+t("orphanFound")+' ('+orphans.length+')</span></div>';h+='<div style="max-height:200px;overflow-y:auto;margin-bottom:.6rem">';for(var i=0;i<orphans.length;i++){h+='<div style="font-family:var(--mono);font-size:.75rem;padding:2px 0;color:var(--dim)">📁 '+esc(orphans[i])+'</div>'}h+='</div>';h+='<button class="btn btn-danger" onclick="confirmCleanup()">'+t("confirmDelete")+'</button>';el.innerHTML=h}).catch(function(e){el.innerHTML='<span style="color:var(--danger);font-size:.8rem">'+esc(e.message)+'</span>'})}
function confirmCleanup(){customConfirm(t("cleanupConfirmWarning"),function(){api("/admin/cleanup",{method:"POST",body:JSON.stringify({confirm:true})}).then(function(d){var r=d.data||{};toast(t("cleanupDone")+" ("+r.deleted+"/"+r.total+")");document.getElementById("cleanup-result").innerHTML='<span style="color:var(--accent);font-size:.8rem">✓ '+t("cleanupDone")+' — '+r.deleted+'/'+r.total+'</span>';loadStats()}).catch(function(e){toast(e.message,"error")})})}
function searchUsers(){userSearch=document.getElementById("user-search").value;userPage=1;loadUsers()}
function userPageGo(p){userPage=p;loadUsers()}
function paginationHtml(p,tp,gf){
if(tp<=1)return"";
var h='<div class="pagination"><button '+(p<=1?"disabled":"")+' onclick="'+gf+'('+(p-1)+')">'+t("prev")+'</button>';
var s=Math.max(1,p-2),e=Math.min(tp,p+2);
if(s>1){h+='<button onclick="'+gf+'(1)">1</button>';if(s>2)h+='<span class="page-info">...</span>'}
for(var i=s;i<=e;i++){h+='<button class="'+(i===p?"active":"")+'" onclick="'+gf+'('+i+')">'+i+'</button>'}
if(e<tp){if(e<tp-1)h+='<span class="page-info">...</span>';h+='<button onclick="'+gf+'('+tp+')">'+tp+'</button>'}
h+='<button '+(p>=tp?"disabled":"")+' onclick="'+gf+'('+(p+1)+')">'+t("next")+'</button><span class="page-info">'+t("page")+' '+p+' '+t("of")+' '+tp+'</span></div>';
return h
}
function loadUsers(){
var q=userSearch?"&q="+encodeURIComponent(userSearch):"";
api("/admin/users?page="+userPage+"&perPage="+userPerPage+q).then(function(d){
var r=d.data||{},users=r.items||[];userTotal=r.total||0;
var el=document.getElementById("users"),tp=Math.ceil(userTotal/userPerPage)||1,pg=paginationHtml(userPage,tp,"userPageGo");
if(!users.length){el.innerHTML='<div class="empty">'+t("noUsers")+'</div>'+pg;return}
var h='<table><thead><tr><th>ID</th><th>Email</th><th>'+t("role")+'</th><th>'+t("created")+'</th><th>'+t("actions")+'</th></tr></thead><tbody>';
for(var i=0;i<users.length;i++){var u=users[i],b=u.isAdmin?'<span class="badge badge-admin">Admin</span>':'<span class="badge badge-user">User</span>',btn=u.isAdmin?'<button class="btn btn-demote" onclick="toggleAdmin('+u.id+')">'+t("demote")+'</button>':'<button class="btn btn-promote" onclick="toggleAdmin('+u.id+')">'+t("promote")+'</button>';
h+='<tr><td style="font-family:var(--mono);color:var(--dim)">'+u.id+'</td><td>'+esc(u.email)+'</td><td>'+b+'</td><td style="font-family:var(--mono);font-size:.7rem;color:var(--dim)">'+fmtDate(u.createdAt)+'</td><td>'+btn+' <button class="btn btn-danger" onclick="delUser('+u.id+')">'+t("delete")+'</button></td></tr>'}
h+='</tbody></table>'+pg;el.innerHTML=h}).catch(function(e){toast(e.message,"error")})
}
function toggleAdmin(id){api("/admin/users/"+id,{method:"PUT"}).then(function(){toast(t("roleUpdated"));loadUsers();loadStats()}).catch(function(e){toast(e.message,"error")})}
function delUser(id){customConfirm(t("deleteUserConfirm"),function(){api("/admin/users/"+id,{method:"DELETE"}).then(function(){toast(t("userDeleted"));loadUsers();loadStats();loadSites()}).catch(function(e){toast(e.message,"error")})})}
function searchAdminSites(){adminSiteSearch=document.getElementById("admin-site-search").value;adminSitePage=1;loadSites()}
function adminSitePageGo(p){adminSitePage=p;loadSites()}
function loadSites(){
var q=adminSiteSearch?"&q="+encodeURIComponent(adminSiteSearch):"";
api("/admin/sites?page="+adminSitePage+"&perPage="+adminSitePerPage+q).then(function(d){
var r=d.data||{},sites=r.items||[];adminSiteTotal=r.total||0;
var el=document.getElementById("sites"),tp=Math.ceil(adminSiteTotal/adminSitePerPage)||1,pg=paginationHtml(adminSitePage,tp,"adminSitePageGo");
if(!sites.length){el.innerHTML='<div class="empty">'+t("noSites")+'</div>'+pg;return}
var h='<table><thead><tr><th>ID</th><th>'+t("name")+'</th><th>'+t("slug")+'</th><th>'+t("owner")+'</th><th>'+t("protected")+'</th><th>'+t("password")+'</th><th>'+t("url")+'</th><th>'+t("actions")+'</th></tr></thead><tbody>';
for(var i=0;i<sites.length;i++){var s=sites[i],b="";
if(s.publicAccessDisabled&&!s.protected){b='<span class="badge badge-disabled">'+t("accessDisabled")+'</span>'}
else if(s.protected){b='<span class="badge badge-protected">'+t("protected")+'</span>'}
else{b='<span class="badge badge-public">'+t("public")+'</span>'}
var pwd=s.protected?'<code style="font-family:var(--mono);font-size:.7rem">'+esc(s.password)+'</code>':'<span style="color:var(--dim)">'+t("none")+'</span>';
h+='<tr class="site-row" id="row-'+s.id+'" onclick="adminToggleDetail('+s.id+')"><td style="font-family:var(--mono);color:var(--dim)">'+s.id+'</td><td>'+esc(s.name)+'</td><td style="font-family:var(--mono);font-size:.75rem">'+esc(s.slug)+'</td><td style="font-size:.75rem">'+esc(s.ownerEmail||"-")+'</td><td>'+b+'</td><td>'+pwd+'</td><td><a href="'+siteUrl(s.url)+'" target="_blank" style="font-family:var(--mono);font-size:.75rem" onclick="event.stopPropagation()">'+s.url+'</a></td><td onclick="event.stopPropagation()"><button class="btn btn-danger" onclick="delSite('+s.id+')">'+t("delete")+'</button></td></tr>';
h+='<tr class="admin-detail-row" id="detail-'+s.id+'"><td colspan="8"></td></tr>'}
h+='</tbody></table>'+pg;el.innerHTML=h}).catch(function(e){toast(e.message,"error")})
}
function delSite(id){customConfirm(t("deleteSiteConfirm"),function(){api("/admin/sites/"+id,{method:"DELETE"}).then(function(){toast(t("siteDeleted"));loadSites();loadStats()}).catch(function(e){toast(e.message,"error")})})}
function formatSize(n){if(n<1024)return n+" B";if(n<1048576)return(n/1024).toFixed(1)+" KB";if(n<1073741824)return(n/1048576).toFixed(1)+" MB";return(n/1073741824).toFixed(1)+" GB"}
function adminToggleDetail(id){var r=document.getElementById("row-"+id);var d=document.getElementById("detail-"+id);if(!r||!d)return;var wasShown=d.classList.contains("show");d.classList.toggle("show");r.classList.toggle("expanded");if(!wasShown)adminLoadFileTree(id)}
function adminLoadFileTree(id){api("/admin/sites/"+id+"/files").then(function(d){var files=d.data||[];var el=document.getElementById("detail-"+id);if(!el)return;var td=el.querySelector("td");if(!td)return;
var h='<div style="padding:.5rem .8rem"><div style="font-weight:600;margin-bottom:.3rem;font-size:.75rem;color:var(--dim)">'+t("files")+'</div>';
if(!files.length){h+='<div style="font-size:.75rem;color:var(--dim)">'+t("noFiles")+'</div></div>';td.innerHTML=h;return}
h+='<div style="max-height:300px;overflow-y:auto">';
for(var i=0;i<files.length;i++){var f=files[i],icon=f.dir?"📁":"📄",sz=f.dir?"-":formatSize(f.size);
h+='<div style="display:flex;justify-content:space-between;align-items:baseline;padding:2px 0;font-family:var(--mono);font-size:.75rem"><span style="overflow:hidden;text-overflow:ellipsis;white-space:nowrap">'+icon+' '+esc(f.name)+(f.dir?"/":"")+'</span><span style="color:var(--dim);flex-shrink:0;margin-left:1rem;text-align:right;min-width:70px">'+sz+'</span></div>'}
h+='</div></div>';td.innerHTML=h}).catch(function(){})}
var sl="en";try{sl=localStorage.getItem("lang")||"en"}catch(e){}lang=sl;
checkAuth().then(function(){renderAdmin()})
</script>
</body>
</html>`

// passwordPageHTML returns the password gate page for a protected site.
func passwordPageHTML(slug, siteName string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>%s — Password Required</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--ink:#f6f8fa;--surface:#fff;--line:#d0d7de;--text:#1f2328;--dim:#636c76;--accent:#1a7f37;--accent-dim:#2da44e;--warn:#9a6700;--danger:#cf222e;--placeholder:#8c959f;--mono:'JetBrains Mono','Fira Code',monospace;--sans:system-ui,-apple-system,sans-serif}
[data-theme="dark"]{--ink:#0c1117;--surface:#161b22;--line:#30363d;--text:#e6edf3;--dim:#7d8590;--accent:#39d353;--accent-dim:#238636;--warn:#d29922;--danger:#f85149;--placeholder:#484f58}
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:var(--sans);background:var(--ink);color:var(--text);min-height:100vh;display:flex;align-items:center;justify-content:center}
.card{background:var(--surface);border:1px solid var(--line);border-radius:10px;padding:2.5rem;width:100%%;max-width:360px}
.card h1{font-family:var(--mono);font-size:1.2rem;margin-bottom:.3rem}
.card .site-name{color:var(--accent);font-weight:600}
.card p{color:var(--dim);margin-bottom:1.2rem;font-size:.8rem;font-family:var(--mono)}
input[type=password]{width:100%%;padding:10px 12px;border:1px solid var(--line);border-radius:6px;font-size:.85rem;margin-bottom:.8rem;outline:none;background:var(--ink);color:var(--text)}
input[type=password]:focus{border-color:var(--accent)}
input[type=password]::placeholder{color:var(--placeholder)}
button{width:100%%;padding:10px;background:var(--accent);color:var(--ink);border:none;border-radius:6px;font-family:var(--mono);font-size:.85rem;font-weight:700;cursor:pointer}
button:hover{background:var(--accent-dim)}
.err{color:var(--danger);margin-bottom:.8rem;font-size:.8rem;font-family:var(--mono)}
</style>
<script>var _t="light";try{_t=localStorage.getItem("theme")||"light"}catch(e){}document.documentElement.setAttribute("data-theme",_t)</script>
</head>
<body>
<div class="card">
<h1>🔒 <span class="site-name">%s</span></h1>
<p>This site is password-protected. Enter the password to continue.</p>
<form method="POST" action="">
<input type="password" id="password" name="password" placeholder="Password" autofocus required>
<button type="submit">Enter Site</button>
</form>
</div>
</body>
</html>`, escHTML(siteName), escHTML(siteName))
}

// passwordPageHTMLWithErr returns the password gate page with an error message.
func passwordPageHTMLWithErr(slug, siteName, errMsg string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>%s — Password Required</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--ink:#f6f8fa;--surface:#fff;--line:#d0d7de;--text:#1f2328;--dim:#636c76;--accent:#1a7f37;--accent-dim:#2da44e;--warn:#9a6700;--danger:#cf222e;--placeholder:#8c959f;--mono:'JetBrains Mono','Fira Code',monospace;--sans:system-ui,-apple-system,sans-serif}
[data-theme="dark"]{--ink:#0c1117;--surface:#161b22;--line:#30363d;--text:#e6edf3;--dim:#7d8590;--accent:#39d353;--accent-dim:#238636;--warn:#d29922;--danger:#f85149;--placeholder:#484f58}
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:var(--sans);background:var(--ink);color:var(--text);min-height:100vh;display:flex;align-items:center;justify-content:center}
.card{background:var(--surface);border:1px solid var(--line);border-radius:10px;padding:2.5rem;width:100%%;max-width:360px}
.card h1{font-family:var(--mono);font-size:1.2rem;margin-bottom:.3rem}
.card .site-name{color:var(--accent);font-weight:600}
.card p{color:var(--dim);margin-bottom:1.2rem;font-size:.8rem;font-family:var(--mono)}
input[type=password]{width:100%%;padding:10px 12px;border:1px solid var(--line);border-radius:6px;font-size:.85rem;margin-bottom:.8rem;outline:none;background:var(--ink);color:var(--text)}
input[type=password]:focus{border-color:var(--accent)}
input[type=password]::placeholder{color:var(--placeholder)}
button{width:100%%;padding:10px;background:var(--accent);color:var(--ink);border:none;border-radius:6px;font-family:var(--mono);font-size:.85rem;font-weight:700;cursor:pointer}
button:hover{background:var(--accent-dim)}
.err{color:var(--danger);margin-bottom:.8rem;font-size:.8rem;font-family:var(--mono)}
</style>
<script>var _t="light";try{_t=localStorage.getItem("theme")||"light"}catch(e){}document.documentElement.setAttribute("data-theme",_t)</script>
</head>
<body>
<div class="card">
<h1>🔒 <span class="site-name">%s</span></h1>
<p>This site is password-protected. Enter the password to continue.</p>
<div class="err">%s</div>
<form method="POST" action="">
<input type="password" name="password" placeholder="Password" autofocus required>
<button type="submit">Enter Site</button>
</form>
</div>
</body>
</html>`, escHTML(siteName), escHTML(siteName), escHTML(errMsg))
}

// escHTML escapes HTML special characters to prevent XSS.
func escHTML(s string) string {
	return strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&quot;",
		"'", "&#39;",
	).Replace(s)
}
