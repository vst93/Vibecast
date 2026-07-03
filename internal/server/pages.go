package server

import (
	"fmt"
	"strings"
)

// logoFavicon is the SVG logo as a data URI for use as favicon.
const logoFavicon = `<link rel="icon" href="data:image/svg+xml,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20viewBox%3D%220%200%20128%20128%22%3E%3Crect%20x%3D%228%22%20y%3D%228%22%20width%3D%22112%22%20height%3D%22112%22%20rx%3D%2224%22%20fill%3D%22%2300e5a0%22%2F%3E%3Cpolyline%20points%3D%2238%2C43%2064%2C97%2090%2C43%22%20fill%3D%22none%22%20stroke%3D%22%230a0a0c%22%20stroke-width%3D%2210%22%20stroke-linecap%3D%22round%22%20stroke-linejoin%3D%22round%22%2F%3E%3Ccircle%20cx%3D%2264%22%20cy%3D%2231%22%20r%3D%226%22%20fill%3D%22%230a0a0c%22%2F%3E%3C%2Fsvg%3E">`

// logoIcon is the inline SVG logo for navbar display.
const logoIcon = `<svg viewBox="0 0 128 128" width="20" height="20" style="vertical-align:middle;margin-right:6px"><rect x="8" y="8" width="112" height="112" rx="24" fill="#00e5a0"/><polyline points="38,43 64,97 90,43" fill="none" stroke="#0a0a0c" stroke-width="10" stroke-linecap="round" stroke-linejoin="round"/><circle cx="64" cy="31" r="6" fill="#0a0a0c"/></svg>`

// landingPageHTML is the public landing page.
var landingPageHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
` + logoFavicon + `
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
@media(max-width:640px){.hero{padding:2rem 0}.hero .logo{font-size:2.2rem}.features{grid-template-columns:1fr}.wrap{padding:1.5rem 1rem}.lang-toggle{top:1rem;right:1rem}}
@media(prefers-reduced-motion:reduce){*{transition:none!important}}
</style>
<script>var _t="light";try{_t=localStorage.getItem("theme")||"light"}catch(e){}document.documentElement.setAttribute("data-theme",_t)</script>
</head>
<body>
<div class="lang-toggle"><a id="theme-toggle" onclick="toggleTheme()">🌙</a><a id="langEn" class="active" onclick="setLang('en')">EN</a><a id="langZh" onclick="setLang('zh')">中文</a></div>
<div class="wrap">
<div class="hero">
<div class="logo">` + logoIcon + `Vibe<span class="accent">cast</span></div>
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
` + logoFavicon + `
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
.copy-btn{background:none;border:1px solid var(--line);border-radius:4px;padding:3px 5px;cursor:pointer;color:var(--dim);line-height:1;vertical-align:middle;display:inline-flex;align-items:center;justify-content:center;transition:color .15s,border-color .15s}
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
.captcha-row .captcha-img{flex-shrink:0;border:1px solid var(--line);border-radius:6px;background:#f6f8fa;cursor:pointer;height:40px;width:150px;overflow:hidden;transition:border-color .15s}
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
.pwd-wrap{position:relative;display:flex;align-items:center}.pwd-wrap input{width:100%;padding-right:38px!important;box-sizing:border-box}.pwd-toggle{position:absolute;right:8px;top:50%;transform:translateY(-50%);background:none;border:none;cursor:pointer;display:flex;align-items:center;justify-content:center;width:28px;height:28px;border-radius:5px;color:var(--dim);padding:0;transition:color .15s,background .15s}.pwd-toggle:hover{color:var(--text);background:var(--surface-2)}.pwd-toggle svg{width:16px;height:16px;flex-shrink:0}
@media(max-width:768px){.dashboard-grid{grid-template-columns:1fr}.sidebar{position:static}}
.badge-org{background:rgba(59,130,246,.12);color:#3b82f6;border:1px solid rgba(59,130,246,.3);padding:1px 6px;border-radius:3px;font-size:.65rem;font-family:var(--mono);text-transform:uppercase;letter-spacing:.3px}
[data-theme="dark"] .badge-org{color:#60a5fa;border-color:rgba(96,165,250,.3)}
.org-toggle-field{margin-top:.6rem;padding-top:.6rem;border-top:1px solid var(--line)}
.checkbox-label{display:flex;align-items:flex-start;gap:.4rem;cursor:pointer;font-size:.8rem;color:var(--text)}
.checkbox-label input{margin-top:2px;flex-shrink:0;accent-color:var(--accent)}
.checkbox-label .desc{margin-top:.2rem}
.modal-lg{max-width:900px}
.modal-header-row{display:flex;align-items:center;justify-content:space-between;margin-bottom:1rem}
.modal-header-row h3{margin:0}
.modal-close{background:none;border:none;font-size:1.2rem;cursor:pointer;color:var(--dim);padding:4px 8px;border-radius:6px;transition:all .15s}
.modal-close:hover{color:var(--text);background:var(--surface-2)}
.org-feature-desc{font-size:.78rem;color:var(--dim);line-height:1.6;margin-bottom:1.2rem;padding:0 0 .8rem;border-bottom:1px solid var(--line)}
#org-modal .modal{max-height:85vh;overflow-y:auto}
.org-empty{display:grid;grid-template-columns:1fr 1fr;gap:0}
.org-empty-panel{padding:1.5rem;display:flex;flex-direction:column;gap:.8rem}
.org-empty-panel:first-child{border-right:1px solid var(--line)}
.org-panel-title{font-family:var(--mono);font-size:.85rem;font-weight:700;color:var(--text);display:flex;align-items:center;gap:.4rem}
.org-panel-icon{width:32px;height:32px;border-radius:8px;display:flex;align-items:center;justify-content:center;flex-shrink:0}
.org-panel-icon.create{background:rgba(26,127,55,.1);color:var(--accent)}
.org-panel-icon.join{background:rgba(59,130,246,.1);color:#3b82f6}
[data-theme="dark"] .org-panel-icon.create{background:rgba(57,211,83,.1);color:var(--accent)}
[data-theme="dark"] .org-panel-icon.join{background:rgba(96,165,250,.1);color:#60a5fa}
.org-panel-desc{font-size:.75rem;color:var(--dim);line-height:1.5}
.org-panel-input{padding:9px 12px;border:1px solid var(--line);border-radius:6px;font-size:.85rem;outline:none;transition:border-color .15s,box-shadow .15s;background:var(--ink);color:var(--text);width:100%}
.org-panel-input:focus{border-color:var(--accent);box-shadow:0 0 0 3px rgba(26,127,55,.08)}
[data-theme="dark"] .org-panel-input:focus{box-shadow:0 0 0 3px rgba(57,211,83,.08)}
.org-bound{display:grid;grid-template-columns:1fr 1.5fr;gap:0}
.org-bound-left{padding:1.5rem;border-right:1px solid var(--line);display:flex;flex-direction:column;gap:1rem}
.org-bound-right{padding:1.5rem;display:flex;flex-direction:column;gap:.8rem}
/* Mobile: fold to single column — must come AFTER base rules to win cascade */
@media(max-width:768px){.org-bound{grid-template-columns:1fr!important}.org-bound-left{border-right:none;border-bottom:1px solid var(--line)}.org-empty{grid-template-columns:1fr!important}.org-empty-panel:first-child{border-right:none;border-bottom:1px solid var(--line)}.modal-lg{max-width:calc(100vw - 2rem)!important}}
@media(max-width:640px){.org-bound-left{padding:1rem}.org-bound-right{padding:1rem}.org-empty-panel{padding:1rem}.modal-lg{max-width:calc(100vw - 1rem)!important;padding:1rem!important}#org-modal .modal{max-height:90vh}}
.org-header{display:flex;align-items:center;gap:.6rem}
.org-avatar{width:40px;height:40px;border-radius:10px;background:linear-gradient(135deg,var(--accent),var(--accent-dim));display:flex;align-items:center;justify-content:center;font-family:var(--mono);font-size:1rem;font-weight:700;color:#fff;flex-shrink:0}
.org-title{font-family:var(--mono);font-size:1rem;font-weight:700;color:var(--text);word-break:break-word}
.org-role-badge{font-size:.65rem;padding:2px 8px;border-radius:4px;font-weight:600;font-family:var(--mono);text-transform:uppercase;letter-spacing:.3px}
.org-role-badge.owner{background:rgba(210,153,34,.12);color:var(--warn);border:1px solid rgba(210,153,34,.25)}
.org-role-badge.member{background:rgba(57,211,83,.12);color:var(--accent);border:1px solid rgba(57,211,83,.25)}
.org-invite-box{background:var(--ink);border:1px solid var(--line);border-radius:8px;padding:.8rem}
.org-invite-box label{font-size:.7rem;color:var(--dim);display:block;margin-bottom:.4rem;font-weight:600}
.org-invite-row{display:flex;align-items:center;gap:.4rem}
.org-invite-code{font-family:var(--mono);font-size:.9rem;color:var(--accent);letter-spacing:1px;flex:1;word-break:break-all}
.org-copy-btn{padding:6px 10px;border:1px solid var(--line);border-radius:6px;background:var(--surface);color:var(--dim);cursor:pointer;font-size:.7rem;font-weight:600;display:flex;align-items:center;gap:.3rem;transition:all .15s;white-space:nowrap}
.org-copy-btn:hover{border-color:var(--accent);color:var(--accent);background:var(--surface-2)}
.org-actions-row{display:flex;gap:.5rem}
.org-actions-row .btn{flex:1}
.org-members-header{display:flex;align-items:center;justify-content:space-between;margin-bottom:.5rem}
.org-members-header h3{font-family:var(--mono);font-size:.8rem;font-weight:600;color:var(--text)}
.org-members-count{font-family:var(--mono);font-size:.7rem;color:var(--dim);background:var(--surface-2);padding:2px 8px;border-radius:10px}
.org-member-item{display:flex;align-items:center;gap:.6rem;padding:.55rem .7rem;border:1px solid var(--line);border-radius:6px;margin-bottom:.4rem;transition:border-color .15s}
.org-member-item:hover{border-color:var(--accent)}
.org-member-avatar{width:28px;height:28px;border-radius:50%;background:var(--surface-2);display:flex;align-items:center;justify-content:center;font-size:.7rem;font-weight:700;color:var(--dim);flex-shrink:0;font-family:var(--mono)}
.org-member-avatar.owner{background:rgba(210,153,34,.12);color:var(--warn)}
.org-member-info{flex:1;min-width:0}
.org-member-email{font-size:.8rem;font-weight:500;color:var(--text);overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.org-member-role{font-size:.65rem;color:var(--dim);font-family:var(--mono)}
.org-member-remove{padding:4px 8px;border:1px solid var(--danger);border-radius:5px;background:transparent;color:var(--danger);cursor:pointer;font-size:.7rem;font-weight:600;transition:all .15s;white-space:nowrap}
.org-member-remove:hover{background:var(--danger);color:#fff}
@media(max-width:640px){.navbar{padding:0 .5rem;height:48px}.navbar .logo{font-size:0}.navbar .logo-text{display:none}.navbar .nav-right .email{display:none}.navbar .btn-icon span{display:none}.navbar .btn-icon{padding:4px 6px;font-size:.85rem}.lang-toggle a{padding:2px 5px;font-size:.7rem}.container{padding:0 .75rem;margin-top:1rem}.auth-card{padding:1.5rem}.auth-card h1{font-size:1.1rem}.auth-field input{padding:9px 10px;font-size:.85rem}.captcha-row .captcha-img{width:120px;height:38px}.site-item .site-head{flex-direction:column;align-items:flex-start;gap:.5rem}.site-item .site-head .actions{width:100%;justify-content:flex-start;flex-wrap:wrap;gap:.3rem}.site-item .site-head .name{font-size:.8rem;flex-wrap:wrap}.site-item .site-detail .detail-row{flex-direction:column;gap:.1rem}.site-item .site-detail .detail-row .label{min-width:auto}.modal{max-width:calc(100vw - 2rem);padding:1.2rem}.toast{bottom:1rem;right:1rem;left:1rem;text-align:center}.list-toolbar input[type=text]{font-size:.85rem}.form-field input{font-size:.85rem}}
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
var maxUploadMB=50;
var maxSites=30;
var i18n={en:{siteName:"Site Name",siteNamePh:"e.g. My Portfolio",slug:"URL Slug",slugPh:"my-portfolio",sitePwd:"Access Password",sitePwdPh:"Leave empty for public",create:"Create Site",yourSites:"Your Sites",noSites:"No sites yet. Create one above.",deployBtn:"Upload",delete:"Delete",deleteConfirm:"Delete this site? This removes all files.",protected:"Protected",public:"Public",login:"Login",register:"Register",email:"Email",emailPh:"you@example.com",password:"Password",pwdHint:"At least 6 characters",noAccount:"No account?",haveAccount:"Have an account?",logout:"Logout",adminPanel:"Admin",deployed:"Deployed!",siteCreated:"Site created",deleted:"Deleted",loginFailed:"Login failed",registerFailed:"Registration failed",slugDesc:"Auto-generated from name if blank. a-z, 0-9, hyphens only.",pwdDesc:"If set, visitors need this password.",sitesHint:"Click to expand",deployHint:"Upload ZIP or single file (PDF, Word, images, etc.)",storagePath:"Storage",accessPassword:"Password",none:"None",accessDisabled:"Disabled",pwdRequired:"Public access disabled — password required",pwdOptional:"Optional password protection",search:"Search",searchSitesPh:"Search sites...",prev:"Prev",next:"Next",page:"Page",of:"of",captcha:"Captcha",captchaPh:"Answer",captchaLabel:"Verification",confirmPassword:"Confirm Password",pwdMismatch:"Passwords do not match",changePassword:"Change Password",currentPassword:"Current Password",newPassword:"New Password",newPasswordPh:"New password (min 6)",passwordChanged:"Password changed",cancel:"Cancel",save:"Save",confirm:"Confirm",emailRequired:"Please enter your email",pwdRequired:"Please enter your password",confirmRequired:"Please confirm your password",captchaRequired:"Please solve the captcha",copyUrl:"Copy URL",copied:"Copied!",copy:"Copy",copyFailed:"Copy failed",visit:"Visit",files:"Files",noFiles:"No files",visitsToday:"Today",visitsMonth:"Month",visitsTotal:"Total",visits:"Visits",uploading:"Uploading...",dragDrop:"Drag file here or click to upload",dragDropHover:"Drop to upload",fileTooLarge:"File too large",fileTypeBlocked:"This file type is not allowed",uploadFailed:"Upload failed",networkError:"Network error",allowedTypes:"ZIP, PDF, Word, Excel, PPT, images, audio, video, text",share:"Share",shareCopied:"Share text copied!",shareTemplatePwd:"Site: {name}\nURL: {url}\nPassword: {password}",shareTemplatePublic:"Site: {name}\nURL: {url}",organization:"Organization",orgManage:"Organization",orgCreate:"Create Organization",orgJoin:"Join Organization",orgName:"Organization Name",orgNamePh:"e.g. My Team",orgInviteCode:"Invite Code",orgInviteCodePh:"Enter 12-char code",orgCreateDesc:"Create a new organization. You will get an invite code to share with your team.",orgJoinDesc:"Enter an invite code to join an existing organization.",orgLeave:"Leave Organization",orgDelete:"Delete Organization",orgMembers:"Members",orgNoMembers:"No members",orgOwner:"Owner",orgMember:"Member",orgRemove:"Remove",orgLeaveConfirm:"Leave this organization?",orgDeleteConfirm:"Delete this organization? This cannot be undone.",orgCreated:"Organization created",orgJoined:"Joined organization",orgLeft:"Left organization",orgDeleted:"Organization deleted",orgMemberRemoved:"Member removed",orgInviteCodeCopied:"Invite code copied!",orgOpenLabel:"Open to org members",orgOpenDesc:"If enabled, logged-in users in the same organization can access this site without a password.",orgBadge:"Org",orgFeatureDesc:"Organizations let you group users together. When a site is set to open to org members, anyone in the same org can access it without a password. Each user can only belong to one organization.",noOrg:"Not in an organization"},zh:{siteName:"站点名称",siteNamePh:"例如：我的作品集",slug:"URL Slug",slugPh:"my-portfolio",sitePwd:"访问密码",sitePwdPh:"留空则公开访问",create:"创建站点",yourSites:"我的站点",noSites:"还没有站点，在左侧创建一个。",deployBtn:"上传",delete:"删除",deleteConfirm:"确定删除此站点？所有文件将被移除。",protected:"已保护",public:"公开",login:"登录",register:"注册",email:"邮箱",emailPh:"you@example.com",password:"密码",pwdHint:"至少 6 个字符",noAccount:"没有账号？",haveAccount:"已有账号？",logout:"退出",adminPanel:"管理",deployed:"部署成功！",siteCreated:"站点已创建",deleted:"已删除",loginFailed:"登录失败",registerFailed:"注册失败",slugDesc:"留空则自动生成。仅限 a-z、0-9、连字符。",pwdDesc:"设置后，访问者需要输入此密码。",sitesHint:"点击展开详情",deployHint:"上传 ZIP 压缩包或单个文件（PDF、Word、图片等）",storagePath:"存储路径",accessPassword:"访问密码",none:"无",accessDisabled:"已禁用",pwdRequired:"公开访问已关闭 — 必须设置密码",pwdOptional:"可选的密码保护",search:"搜索",searchSitesPh:"搜索站点...",prev:"上一页",next:"下一页",page:"第",of:"/ 共",captcha:"验证码",captchaPh:"输入答案",captchaLabel:"验证",confirmPassword:"确认密码",pwdMismatch:"两次密码不一致",changePassword:"修改密码",currentPassword:"当前密码",newPassword:"新密码",newPasswordPh:"新密码（至少 6 位）",passwordChanged:"密码已修改",cancel:"取消",save:"保存",confirm:"确认",emailRequired:"请输入邮箱",pwdRequired:"请输入密码",confirmRequired:"请确认密码",captchaRequired:"请输入验证码",copyUrl:"复制链接",copied:"已复制！",copy:"复制",copyFailed:"复制失败",visit:"访问",files:"文件",noFiles:"暂无文件",visitsToday:"今日",visitsMonth:"本月",visitsTotal:"总计",visits:"访问",uploading:"上传中...",dragDrop:"拖拽文件到此处或点击上传",dragDropHover:"松开以上传",fileTooLarge:"文件过大",fileTypeBlocked:"不允许此文件类型",uploadFailed:"上传失败",networkError:"网络错误",allowedTypes:"ZIP、PDF、Word、Excel、PPT、图片、音频、视频、文本",share:"分享",shareCopied:"分享文本已复制！",shareTemplatePwd:"站点：{name}\n地址：{url}\n密码：{password}",shareTemplatePublic:"站点：{name}\n地址：{url}",organization:"组织",orgManage:"组织管理",orgCreate:"创建组织",orgJoin:"加入组织",orgName:"组织名称",orgNamePh:"例如：我的团队",orgInviteCode:"邀请码",orgInviteCodePh:"输入 12 位邀请码",orgCreateDesc:"创建一个新组织，你将获得一个邀请码来分享给团队成员。",orgJoinDesc:"输入邀请码加入已有组织。",orgLeave:"退出组织",orgDelete:"删除组织",orgMembers:"成员",orgNoMembers:"暂无成员",orgOwner:"创建者",orgMember:"成员",orgRemove:"移除",orgLeaveConfirm:"确定退出此组织？",orgDeleteConfirm:"确定删除此组织？此操作不可撤销。",orgCreated:"组织已创建",orgJoined:"已加入组织",orgLeft:"已退出组织",orgDeleted:"组织已删除",orgMemberRemoved:"成员已移除",orgInviteCodeCopied:"邀请码已复制！",orgOpenLabel:"对组织成员开放",orgOpenDesc:"开启后，同一组织内已登录的用户可以直接访问此站点，无需密码。",orgBadge:"组织",orgFeatureDesc:"组织用于将用户分组。当站点设置为对组织成员开放时，同一组织内的已登录用户无需密码即可访问。每个用户只能加入一个组织。",noOrg:"未加入任何组织"}};
function t(k){return(i18n[lang]||i18n.en)[k]||(i18n.en[k]||k)}
function setLang(l){lang=l;try{localStorage.setItem("lang",l)}catch(e){}if(currentUser)renderDashboard();else renderAuth();}
function esc(s){return String(s||"").replace(/[&<>"']/g,function(c){return{"&":"&amp;","<":"&lt;",">":"&gt;",'"':"&quot;","'":"&#39;"}[c]})}
function siteUrl(u){return BASE+(u||"").replace(/^\//,"")}
function toast(msg,type){type=type||"success";var el=document.createElement("div");el.className="toast "+type+" show";el.textContent=msg;document.body.appendChild(el);var ms=type==="error"?5000:2500;setTimeout(function(){el.classList.remove("show");setTimeout(function(){el.remove()},200)},ms)}
function togglePwd(btn){var inp=btn.parentElement.querySelector("input");if(!inp)return;var show=inp.type==="password";inp.type=show?"text":"password";btn.innerHTML=show?PWD_HIDE_ICON:PWD_SHOW_ICON}
var PWD_SHOW_ICON='<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7z"/><circle cx="12" cy="12" r="3"/></svg>';
var PWD_HIDE_ICON='<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-10-8-10-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 10 8 10 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/><line x1="1" y1="1" x2="23" y2="23"/></svg>';
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
if(v==="login"){loginCaptchaId=d.data.id;var el=document.getElementById("login-captcha-img");if(el)el.innerHTML=d.data.image}
else{registerCaptchaId=d.data.id;var el=document.getElementById("register-captcha-img");if(el)el.innerHTML=d.data.image}
}).catch(function(){})
}
function renderAuth(){
var lh='<div class="lang-toggle" style="position:absolute;top:1rem;right:1rem"><a id="theme-toggle" onclick="toggleTheme()">'+(getTheme()==="dark"?"☀":"🌙")+'</a><a class="'+(lang==="en"?"active":"")+'" onclick="setLang(\'en\')">EN</a><a class="'+(lang==="zh"?"active":"")+'" onclick="setLang(\'zh\')">中文</a></div>';
document.getElementById("app").innerHTML=lh+'<div class="auth-screen"><div class="auth-card"><h1>` + logoIcon + `Vibe<span class="accent">cast</span></h1><p class="subtitle">Build with vibe. Cast instantly.</p><div id="auth-form"></div></div></div>';
fetch(BASE+"api/settings").then(function(r){return r.json()}).then(function(d){regOpen=d.data&&d.data.openRegistration!==false;showLogin()}).catch(function(){regOpen=true;showLogin()});
}
function showLogin(){
document.getElementById("auth-form").innerHTML='<form onsubmit="doLogin();return false"><div class="auth-field"><label>'+t("email")+'</label><input id="email" type="email" placeholder="'+t("emailPh")+'" autocomplete="email"></div><div class="auth-field"><label>'+t("password")+'</label><div class="pwd-wrap"><input id="password" type="password" placeholder="'+t("password")+'" autocomplete="current-password"><button type="button" class="pwd-toggle" onclick="togglePwd(this)"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7z"/><circle cx="12" cy="12" r="3"/></svg></button></div></div><div class="auth-field"><label class="captcha-label">'+t("captchaLabel")+'</label><div class="captcha-row"><div class="captcha-img" id="login-captcha-img" onclick="loadCaptcha(\'login\')" title="'+t("captcha")+'"></div><input id="login-captcha" type="text" placeholder="'+t("captchaPh")+'"></div></div><button class="btn btn-primary" type="submit">'+t("login")+'</button></form>'+(regOpen?'<p class="switch">'+t("noAccount")+' <a onclick="showRegister()">'+t("register")+'</a></p>':'');
loadCaptcha("login");
}
function showRegister(){
document.getElementById("auth-form").innerHTML='<form onsubmit="doRegister();return false"><div class="auth-field"><label>'+t("email")+'</label><input id="email" type="email" placeholder="'+t("emailPh")+'" autocomplete="email"></div><div class="auth-field"><label>'+t("password")+'</label><div class="pwd-wrap"><input id="password" type="password" placeholder="'+t("password")+'" autocomplete="new-password"><button type="button" class="pwd-toggle" onclick="togglePwd(this)"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7z"/><circle cx="12" cy="12" r="3"/></svg></button></div><div class="field-hint">'+t("pwdHint")+'</div></div><div class="auth-field"><label>'+t("confirmPassword")+'</label><div class="pwd-wrap"><input id="confirm-pwd" type="password" placeholder="'+t("confirmPassword")+'" autocomplete="new-password"><button type="button" class="pwd-toggle" onclick="togglePwd(this)"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7z"/><circle cx="12" cy="12" r="3"/></svg></button></div></div><div class="auth-field"><label class="captcha-label">'+t("captchaLabel")+'</label><div class="captcha-row"><div class="captcha-img" id="register-captcha-img" onclick="loadCaptcha(\'register\')" title="'+t("captcha")+'"></div><input id="register-captcha" type="text" placeholder="'+t("captchaPh")+'"></div></div><button class="btn btn-primary" type="submit">'+t("register")+'</button></form><p class="switch">'+t("haveAccount")+' <a onclick="showLogin()">'+t("login")+'</a></p>';
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
document.getElementById("app").innerHTML='<nav class="navbar"><div class="logo">` + logoIcon + `<span class="logo-text">Vibe<span class="accent">cast</span></span></div><div class="nav-right">'+al+'<button class="btn-icon" onclick="openOrgModal()" id="nav-org-btn">🏢<span> '+t("orgManage")+'</span></button><button class="btn-icon" onclick="openChangePwdModal()">🔒<span> '+t("changePassword")+'</span></button>'+th+lh+'<span class="email">'+esc(currentUser.email)+'</span><button class="btn-link" onclick="doLogout()">'+t("logout")+'</button></div></nav><div class="container"><div class="dashboard-grid"><div class="sidebar"><div class="card"><div class="card-header"><h2>'+t("create")+'</h2></div><div class="card-body"><div class="form-field"><label>'+t("siteName")+'</label><input id="site-name" placeholder="'+t("siteNamePh")+'"></div><div class="form-field"><label>'+t("sitePwd")+'</label><div class="pwd-wrap"><input id="site-pwd" type="password" placeholder="'+t("sitePwdPh")+'"><button type="button" class="pwd-toggle" onclick="togglePwd(this)"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7z"/><circle cx="12" cy="12" r="3"/></svg></button></div><div class="desc" id="pwd-desc">'+t("pwdDesc")+'</div></div><div class="form-field org-toggle-field"><label class="checkbox-label"><input type="checkbox" id="site-org-open"> <span>'+t("orgOpenLabel")+'</span></label><div class="desc">'+t("orgOpenDesc")+'</div></div><div class="form-actions"><button class="btn btn-primary" style="width:100%" onclick="createSite()">'+t("create")+'</button></div></div></div></div><div class="main-content"><div class="card"><div class="card-header"><h2>'+t("yourSites")+'</h2><span class="hint"><span id="site-limit-badge" style="margin-right:.5rem"></span>'+t("sitesHint")+'</span></div><div class="card-body"><div class="list-toolbar"><input type="text" id="site-search" placeholder="'+t("searchSitesPh")+'" onkeydown="if(event.key===\'Enter\')searchSites()" value="'+esc(siteSearch)+'"></div><div id="site-list"></div></div></div></div></div></div><div class="modal-overlay" id="pwd-modal" onclick="if(event.target===this)closeChangePwdModal()"><div class="modal"><h3>'+t("changePassword")+'</h3><div class="modal-field"><label>'+t("currentPassword")+'</label><div class="pwd-wrap"><input id="old-pwd" type="password" placeholder="'+t("currentPassword")+'"><button type="button" class="pwd-toggle" onclick="togglePwd(this)"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7z"/><circle cx="12" cy="12" r="3"/></svg></button></div></div><div class="modal-field"><label>'+t("newPassword")+'</label><div class="pwd-wrap"><input id="new-pwd" type="password" placeholder="'+t("newPasswordPh")+'"><button type="button" class="pwd-toggle" onclick="togglePwd(this)"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7z"/><circle cx="12" cy="12" r="3"/></svg></button></div></div><div class="modal-field"><label>'+t("confirmPassword")+'</label><div class="pwd-wrap"><input id="confirm-new-pwd" type="password" placeholder="'+t("confirmPassword")+'"><button type="button" class="pwd-toggle" onclick="togglePwd(this)"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7z"/><circle cx="12" cy="12" r="3"/></svg></button></div></div><div class="modal-actions"><button class="btn btn-outline" onclick="closeChangePwdModal()">'+t("cancel")+'</button><button class="btn btn-primary" onclick="changePassword()">'+t("save")+'</button></div></div></div><div class="modal-overlay" id="org-modal" onclick="if(event.target===this)closeOrgModal()"><div class="modal modal-lg"><div class="modal-header-row"><h3>'+t("orgManage")+'</h3><button class="modal-close" onclick="closeOrgModal()">✕</button></div><p class="org-feature-desc">'+t("orgFeatureDesc")+'</p><div id="org-section"></div></div></div>';
loadSites();
loadOrg();
fetch(BASE+"api/settings").then(function(r){return r.json()}).then(function(d){if(d.data&&d.data.maxUploadSize)maxUploadMB=d.data.maxUploadSize;if(d.data&&d.data.maxSitesPerUser)maxSites=d.data.maxSitesPerUser}).catch(function(){});
}
function openChangePwdModal(){document.getElementById("pwd-modal").classList.add("show")}
function closeChangePwdModal(){var m=document.getElementById("pwd-modal");m.classList.remove("show");m.querySelectorAll("input").forEach(function(i){i.value=""})}
function openOrgModal(){document.getElementById("org-modal").classList.add("show");loadOrg()}
function closeOrgModal(){document.getElementById("org-modal").classList.remove("show")}
function updateOrgNavBtn(){var btn=document.getElementById("nav-org-btn");if(!btn)return;var label=orgInfo&&orgInfo.name?esc(orgInfo.name):t("noOrg");btn.innerHTML='🏢<span> '+label+'</span>'}
function searchSites(){siteSearch=document.getElementById("site-search").value;sitePage=1;loadSites()}
function sitePageGo(p){sitePage=p;loadSites()}
function loadSites(){
var q=siteSearch?"&q="+encodeURIComponent(siteSearch):"";
api("/sites?page="+sitePage+"&perPage="+sitePerPage+q).then(function(d){
var r=d.data||{},sites=r.items||[];siteTotal=r.total||0;
var sb=document.getElementById("site-limit-badge");if(sb&&maxSites>0)sb.textContent=siteTotal+"/"+maxSites+" · ";
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
if(s.orgOpen){badge+=' <span class="badge badge-org">'+t("orgBadge")+'</span>'}
var pwd=s.protected?'<code style="font-family:var(--mono);font-size:.75rem;color:var(--text)">'+esc(s.password)+'</code> <span style="display:inline-flex;gap:2px;vertical-align:middle"><button class="copy-btn" data-pwd="'+esc(s.password)+'" onclick="event.stopPropagation();copyText(this.getAttribute(\'data-pwd\'))" title="'+t("copy")+'"><svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg></button><button class="copy-btn" onclick="event.stopPropagation();shareSite(\''+esc(s.name)+'\',\''+siteUrl(s.url)+'\','+(s.protected?'\''+esc(s.password)+'\'':'null')+')" title="'+t("share")+'"><svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="18" cy="5" r="3"/><circle cx="6" cy="12" r="3"/><circle cx="18" cy="19" r="3"/><line x1="8.59" y1="13.51" x2="15.42" y2="17.49"/><line x1="15.41" y1="6.51" x2="8.59" y2="10.49"/></svg></button></span>':'<span style="color:var(--dim)">'+t("none")+'</span> <button class="copy-btn" onclick="event.stopPropagation();shareSite(\''+esc(s.name)+'\',\''+siteUrl(s.url)+'\',null)" title="'+t("share")+'"><svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="18" cy="5" r="3"/><circle cx="6" cy="12" r="3"/><circle cx="18" cy="19" r="3"/><line x1="8.59" y1="13.51" x2="15.42" y2="17.49"/><line x1="15.41" y1="6.51" x2="8.59" y2="10.49"/></svg></button>';
var vs=s.visits||{today:0,month:0,total:0};
var vis='<span style="font-family:var(--mono);font-size:.7rem;color:var(--dim)" title="'+t("visits")+'">👁 '+vs.today+' / '+vs.month+' / '+vs.total+'</span>';
h+='<li class="site-item"><div class="site-head" onclick="toggleDetail('+s.id+')"><div class="info"><div class="name"><span class="status-dot '+dot+'"></span>'+esc(s.name)+' '+badge+' '+vis+'</div><div class="url">~/sites/'+esc(s.slug)+'/</div></div><div class="actions"><a class="btn btn-sm btn-outline" href="'+siteUrl(s.url)+'" target="_blank" onclick="event.stopPropagation()">'+t("visit")+'</a><label class="upload-btn" id="upload-btn-'+s.id+'" onclick="event.stopPropagation()">'+t("deployBtn")+'<input type="file" accept=".zip,.pdf,.doc,.docx,.xls,.xlsx,.ppt,.pptx,.odt,.ods,.odp,.rtf,.txt,.csv,.json,.xml,.html,.htm,.css,.js,.mjs,.svg,.png,.jpg,.jpeg,.gif,.webp,.avif,.ico,.mp4,.webm,.mp3,.ogg,.wav,.flac,.wasm,.woff,.woff2,.ttf,.otf,.eot,.webmanifest" onchange="deploy('+s.id+',this.files[0])"></label><button class="btn btn-sm btn-danger" onclick="event.stopPropagation();delSite('+s.id+')">'+t("delete")+'</button></div></div><div class="site-detail" id="detail-'+s.id+'"><div class="detail-row"><span class="label">'+t("accessPassword")+'</span><span class="value">'+pwd+'</span></div><div class="detail-row"><span class="label">'+t("orgOpenLabel")+'</span><span class="value"><label class="checkbox-label" style="display:inline-flex;align-items:center;gap:.4rem"><input type="checkbox" id="org-open-'+s.id+'" '+(s.orgOpen?'checked':'')+' onchange="toggleOrgOpen('+s.id+')"> <span>'+t("orgOpenDesc")+'</span></label></span></div><div class="detail-row"><span class="label">'+t("visits")+'</span><span class="value" style="display:flex;gap:1rem;font-family:var(--mono);font-size:.75rem"><span>'+t("visitsToday")+': <b style="color:var(--accent)">'+vs.today+'</b></span><span>'+t("visitsMonth")+': <b style="color:var(--accent)">'+vs.month+'</b></span><span>'+t("visitsTotal")+': <b style="color:var(--accent)">'+vs.total+'</b></span></span></div></div></li>'}
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
if(!el.hasAttribute('data-orig'))el.setAttribute('data-orig',el.innerHTML);
var h=el.getAttribute('data-orig');
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
function shareSite(name,url,pwd){
var fullUrl=location.origin+url;
var key=pwd?"shareTemplatePwd":"shareTemplatePublic";
var txt=t(key).replace("{name}",name).replace("{url}",fullUrl);
if(pwd)txt=txt.replace("{password}",pwd);
copyText(txt);
toast(t("shareCopied"));
}
function createSite(){
api("/sites",{method:"POST",body:JSON.stringify({name:document.getElementById("site-name").value,password:document.getElementById("site-pwd").value,orgOpen:document.getElementById("site-org-open").checked})}).then(function(){document.getElementById("site-name").value="";document.getElementById("site-pwd").value="";toast(t("siteCreated"));loadSites()}).catch(function(e){toast(e.message,"error")})
}
function deploy(id,file){
if(!file)return;
var btn=document.getElementById("upload-btn-"+id);
var origText=btn?btn.firstChild.textContent:"";
var allowedExts=[".zip",".pdf",".doc",".docx",".xls",".xlsx",".ppt",".pptx",".odt",".ods",".odp",".rtf",".txt",".csv",".json",".xml",".html",".htm",".css",".js",".mjs",".svg",".png",".jpg",".jpeg",".gif",".webp",".avif",".ico",".mp4",".webm",".mp3",".ogg",".wav",".flac",".wasm",".woff",".woff2",".ttf",".otf",".eot",".webmanifest"];
var fName=file.name.toLowerCase();
var ext="."+fName.split(".").pop();
if(allowedExts.indexOf(ext)<0){toast(t("fileTypeBlocked")+" ("+ext+")","error");return}
if(file.size>maxUploadMB*1024*1024){toast(t("fileTooLarge")+" ("+maxUploadMB+"MB)","error");return}
if(btn){btn.firstChild.textContent=t("uploading");btn.style.opacity="0.6";btn.style.pointerEvents="none"}
var fd=new FormData();fd.append("file",file);
var token=localStorage.getItem("vibecast_token")||"";
var xhr=new XMLHttpRequest();
xhr.open("POST",API+"/sites/"+id+"/deploy");
xhr.setRequestHeader("Authorization","Bearer "+token);
xhr.setRequestHeader("Accept-Language",lang||"en");
xhr.upload.onprogress=function(e){if(e.lengthComputable&&btn){var pct=Math.round(e.loaded/e.total*100);btn.firstChild.textContent=t("uploading")+" "+pct+"%"}};
xhr.onload=function(){
if(btn){btn.firstChild.textContent=origText;btn.style.opacity="";btn.style.pointerEvents=""}
try{var data=JSON.parse(xhr.responseText);if(xhr.status>=200&&xhr.status<300){toast(t("deployed"));loadSites()}else{toast(data.error||t("uploadFailed"),"error")}}catch(e){toast(t("networkError"),"error")}
};
xhr.onerror=function(){if(btn){btn.firstChild.textContent=origText;btn.style.opacity="";btn.style.pointerEvents=""}toast(t("fileTooLarge")+" / "+t("networkError"),"error")};
xhr.send(fd);
}
function delSite(id){
customConfirm(t("deleteConfirm"),function(){api("/sites/"+id,{method:"DELETE"}).then(function(){toast(t("deleted"));loadSites()}).catch(function(e){toast(e.message,"error")})})
}
function toggleOrgOpen(id){
var cb=document.getElementById("org-open-"+id);if(!cb)return;
api("/sites/"+id,{method:"PUT",body:JSON.stringify({orgOpen:cb.checked})}).then(function(){toast(t("save"));loadSites()}).catch(function(e){toast(e.message,"error");cb.checked=!cb.checked})
}
function changePassword(){
var o=document.getElementById("old-pwd").value,n=document.getElementById("new-pwd").value,c=document.getElementById("confirm-new-pwd").value;
if(!o||!n){toast(t("currentPassword")+" & "+t("newPassword"),"error");return}
if(n.length<6){toast(t("pwdHint"),"error");return}
if(n!==c){toast(t("pwdMismatch"),"error");return}
api("/auth/change-password",{method:"PUT",body:JSON.stringify({oldPassword:o,newPassword:n})}).then(function(){toast(t("passwordChanged"));closeChangePwdModal()}).catch(function(e){toast(e.message,"error")})
}

var orgInfo=null,orgMembers=[],orgMemberPage=1,orgMemberPerPage=10,orgMemberTotal=0,orgMemberSearch="";
function loadOrg(){
api("/org").then(function(d){
var r=d.data||{};orgInfo=r.hasOrg?r:null;
updateOrgNavBtn();
renderOrgSection();
}).catch(function(){})
}
function renderOrgSection(){
var el=document.getElementById("org-section");if(!el)return;
if(!orgInfo){
el.innerHTML='<div class="org-empty"><div class="org-empty-panel"><div class="org-panel-title"><span class="org-panel-icon create"><svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg></span>'+t("orgCreate")+'</div><p class="org-panel-desc">'+t("orgCreateDesc")+'</p><input class="org-panel-input" id="org-create-name" placeholder="'+t("orgNamePh")+'"><button class="btn btn-primary" onclick="createOrg()">'+t("orgCreate")+'</button></div><div class="org-empty-panel"><div class="org-panel-title"><span class="org-panel-icon join"><svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M15 22v-4a4.8 4.8 0 0 0-1-3.5c3-.3 6-1.5 6-7a5.5 5.5 0 0 0-1.5-3.8 5.1 5.1 0 0 0-.1-3.8s-1.2-.3-3.8 1.4a13 13 0 0 0-7 0C5.2 3.3 4 3.6 4 3.6a5.1 5.1 0 0 0-.1 3.8A5.5 5.5 0 0 0 2.4 11c0 5.5 3 6.7 6 7a4.8 4.8 0 0 0-1 3.5V22"/></svg></span>'+t("orgJoin")+'</div><p class="org-panel-desc">'+t("orgJoinDesc")+'</p><input class="org-panel-input" id="org-join-code" placeholder="'+t("orgInviteCodePh")+'"><button class="btn btn-outline" style="border-color:#3b82f6;color:#3b82f6" onclick="joinOrg()">'+t("orgJoin")+'</button></div></div>';
return
}
var inviteCode=orgInfo.inviteCode||"";
var name=orgInfo.name||"";
var isOwner=orgInfo.isOwner||false;
var initials=name.substring(0,2).toUpperCase();
var roleBadge=isOwner?'<span class="org-role-badge owner">'+t("orgOwner")+'</span>':'<span class="org-role-badge member">'+t("orgMember")+'</span>';
var h='<div class="org-bound"><div class="org-bound-left"><div class="org-header"><div class="org-avatar">'+esc(initials)+'</div><div><div class="org-title">'+esc(name)+'</div>'+roleBadge+'</div></div>';
if(isOwner){h+='<div class="org-invite-box"><label>'+t("orgInviteCode")+'</label><div class="org-invite-row"><code class="org-invite-code">'+esc(inviteCode)+'</code><button class="org-copy-btn" onclick="copyText(\''+esc(inviteCode)+'\')"><svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg>'+t("copy")+'</button></div></div>'}
h+='<div class="org-actions-row">';
if(isOwner){h+='<button class="btn btn-danger btn-sm" onclick="deleteOrg()">'+t("orgDelete")+'</button>'}
else{h+='<button class="btn btn-outline btn-sm" onclick="leaveOrg()">'+t("orgLeave")+'</button>'}
h+='</div></div><div class="org-bound-right"><div class="org-members-header"><h3>'+t("orgMembers")+'</h3><span class="org-members-count" id="org-members-count">—</span></div><div class="list-toolbar"><input type="text" id="org-member-search" placeholder="'+t("search")+'..." onkeydown="if(event.key===\'Enter\')searchOrgMembers()" value="'+esc(orgMemberSearch)+'"></div><div id="org-member-list"></div></div></div>';
el.innerHTML=h;
loadOrgMembers()
}
function orgMemberPageGo(p){orgMemberPage=p;loadOrgMembers()}
function searchOrgMembers(){orgMemberSearch=document.getElementById("org-member-search").value;orgMemberPage=1;loadOrgMembers()}
function loadOrgMembers(){
var q=orgMemberSearch?"&q="+encodeURIComponent(orgMemberSearch):"";
api("/org/members?page="+orgMemberPage+"&perPage="+orgMemberPerPage+q).then(function(d){
var r=d.data||{},members=r.items||[];orgMemberTotal=r.total||0;
var el=document.getElementById("org-member-list");if(!el)return;
var cnt=document.getElementById("org-members-count");if(cnt)cnt.textContent=orgMemberTotal;
var tp=Math.ceil(orgMemberTotal/orgMemberPerPage)||1,pg=paginationHtml(orgMemberPage,tp,"orgMemberPageGo");
if(!members.length){el.innerHTML='<div class="empty">'+t("orgNoMembers")+'</div>'+pg;return}
var h='';
for(var i=0;i<members.length;i++){var m=members[i];
var initials=(m.email||"?").substring(0,2).toUpperCase();
var avClass=m.isOwner?'org-member-avatar owner':'org-member-avatar';
var role=m.isOwner?t("orgOwner"):t("orgMember");
var rmv=m.isOwner?'':'<button class="org-member-remove" onclick="removeOrgMember('+m.userId+')">'+t("orgRemove")+'</button>';
h+='<div class="org-member-item"><div class="'+avClass+'">'+esc(initials)+'</div><div class="org-member-info"><div class="org-member-email">'+esc(m.email)+'</div><div class="org-member-role">'+role+'</div></div>'+rmv+'</div>'}
h+=pg;
el.innerHTML=h
}).catch(function(){})
}
function createOrg(){
var name=document.getElementById("org-create-name").value.trim();
api("/org",{method:"POST",body:JSON.stringify({name:name})}).then(function(){toast(t("orgCreated"));loadOrg()}).catch(function(e){toast(e.message,"error")})
}
function joinOrg(){
var code=document.getElementById("org-join-code").value.trim();
if(!code){toast(t("orgInviteCodePh"),"error");return}
api("/org/join",{method:"POST",body:JSON.stringify({inviteCode:code})}).then(function(){toast(t("orgJoined"));loadOrg()}).catch(function(e){toast(e.message,"error")})
}
function leaveOrg(){
customConfirm(t("orgLeaveConfirm"),function(){api("/org/leave",{method:"POST"}).then(function(){toast(t("orgLeft"));loadOrg()}).catch(function(e){toast(e.message,"error")})})
}
function deleteOrg(){
customConfirm(t("orgDeleteConfirm"),function(){api("/org",{method:"DELETE"}).then(function(){toast(t("orgDeleted"));loadOrg()}).catch(function(e){toast(e.message,"error")})})
}
function removeOrgMember(uid){
api("/org/members/"+uid,{method:"DELETE"}).then(function(){toast(t("orgMemberRemoved"));loadOrgMembers()}).catch(function(e){toast(e.message,"error")})
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
` + logoFavicon + `
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
.pwd-wrap{position:relative;display:flex;align-items:center}.pwd-wrap input{width:100%;padding-right:38px!important;box-sizing:border-box}.pwd-toggle{position:absolute;right:8px;top:50%;transform:translateY(-50%);background:none;border:none;cursor:pointer;display:flex;align-items:center;justify-content:center;width:28px;height:28px;border-radius:5px;color:var(--dim);padding:0;transition:color .15s,background .15s}.pwd-toggle:hover{color:var(--text);background:var(--surface-2)}.pwd-toggle svg{width:16px;height:16px;flex-shrink:0}
.num-setting{display:flex;align-items:center;justify-content:space-between;padding:.6rem 0;border-bottom:1px solid var(--line)}.num-setting:last-child{border-bottom:none}.num-setting .num-info{flex:1;min-width:0}.num-setting .num-label{font-size:.85rem;font-weight:500}.num-setting .num-desc{font-size:.7rem;color:var(--dim);margin-top:.1rem}.num-setting .num-control{display:flex;align-items:center;gap:.4rem;flex-shrink:0}.num-setting .num-control input{width:72px;padding:6px 8px;border:1px solid var(--line);border-radius:5px;font-size:.85rem;background:var(--ink);color:var(--text);text-align:center;outline:none;transition:border-color .15s;font-family:var(--mono)}.num-setting .num-control input:focus{border-color:var(--accent)}.num-setting .num-control .unit{font-size:.75rem;color:var(--dim);font-family:var(--mono);white-space:nowrap}
@media(max-width:640px){.navbar{padding:0 .5rem;height:48px}.navbar .logo{font-size:0}.navbar .logo-text{display:none}.navbar .nav-right{gap:.4rem}.stats-grid{grid-template-columns:1fr}.lang-toggle a{padding:2px 5px;font-size:.7rem}.admin-layout{display:flex!important;flex-direction:column!important;min-height:auto!important;grid-template-columns:none!important}.admin-sidebar{position:sticky;top:48px;flex-direction:row!important;overflow-x:auto;overflow-y:hidden!important;height:44px!important;flex-shrink:0;border-right:none;border-bottom:1px solid var(--line);padding:0!important;z-index:99;background:var(--surface);-webkit-overflow-scrolling:touch;scrollbar-width:none;gap:0!important}.admin-sidebar::-webkit-scrollbar{display:none}.admin-nav-item{flex:0 0 auto;white-space:nowrap;padding:0 .65rem;height:44px;font-size:.72rem;flex-direction:row!important;gap:.25rem;border-left:none;border-bottom:2px solid transparent;border-top:none;justify-content:center;align-items:center}.admin-nav-item.active{border-left:none;border-bottom-color:var(--accent)}.admin-nav-item span{display:inline;font-size:.72rem}.admin-nav-item svg{width:13px;height:13px;flex-shrink:0}.admin-main{padding:.75rem!important;max-width:100%!important}.container{padding:0 .75rem;margin-top:1rem}.modal{max-width:calc(100vw - 2rem);padding:1.2rem}.toast{bottom:1rem;right:1rem;left:1rem;text-align:center}table{font-size:.7rem;display:block;overflow-x:auto;-webkit-overflow-scrolling:touch}th,td{padding:.3rem .4rem}.info-grid{grid-template-columns:1fr}}
@media(prefers-reduced-motion:reduce){*{transition:none!important}}
.admin-layout{display:grid;grid-template-columns:180px 1fr;gap:0;min-height:calc(100vh - 52px)}
.admin-sidebar{background:var(--surface);border-right:1px solid var(--line);padding:.75rem 0;display:flex;flex-direction:column;gap:2px;position:sticky;top:52px;height:calc(100vh - 52px);overflow-y:auto}
.admin-nav-item{display:flex;align-items:center;gap:.5rem;padding:.6rem 1rem;font-size:.8rem;color:var(--dim);cursor:pointer;border-left:2px solid transparent;transition:all .15s;font-weight:500}
.admin-nav-item:hover{color:var(--text);background:var(--surface-2)}
.admin-nav-item.active{color:var(--accent);border-left-color:var(--accent);background:var(--surface-2)}
.admin-nav-item svg{width:14px;height:14px;flex-shrink:0}
.admin-main{padding:1.5rem;max-width:1100px}
.admin-tab-content{display:none}
.admin-tab-content.active{display:block}
.info-grid{display:grid;grid-template-columns:1fr 1fr;gap:.6rem}
.info-item{background:var(--ink);border:1px solid var(--line);border-radius:6px;padding:.7rem .9rem}
.info-item .info-label{font-size:.65rem;color:var(--dim);font-family:var(--mono);text-transform:uppercase;letter-spacing:.5px;margin-bottom:.2rem}
.info-item .info-value{font-size:.85rem;font-family:var(--mono);color:var(--text);word-break:break-all}
.info-item .info-value.accent{color:var(--accent)}
.update-inline{margin-top:1rem;padding-top:1rem;border-top:1px solid var(--line)}
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
var i18n={en:{overview:"Overview",settings:"Settings",users:"Users",allSites:"All Sites",userCount:"Users",siteCount:"Sites",adminCount:"Admins",openReg:"Open Registration",openRegDesc:"When enabled, anyone can register a new account.",publicAccess:"Public Site Access",publicAccessDesc:"When disabled, all deployed sites return 403 to visitors.",dashboard:"Dashboard",logout:"Logout",role:"Role",created:"Created",actions:"Actions",promote:"Promote",demote:"Demote",delete:"Delete",deleteUserConfirm:"Delete this user and all their sites?",deleteSiteConfirm:"Delete this site?",name:"Name",slug:"Slug",owner:"Owner",protected:"Protected",public:"Public",password:"Password",storagePath:"Storage",url:"URL",noUsers:"No users",noSites:"No sites",updated:"Updated",deleted:"Deleted",roleUpdated:"Role updated",userDeleted:"User deleted",siteDeleted:"Site deleted",settingsSaved:"Settings saved",none:"None",accessDisabled:"Disabled",search:"Search",searchUsersPh:"Search users...",searchSitesPh:"Search sites...",prev:"Prev",next:"Next",page:"Page",of:"of",ownerEmail:"Owner",domainRestriction:"Email Domain Restriction",domainRestrictionDesc:"Only allow registration from specified domains (one per line).",save:"Save",cancel:"Cancel",confirm:"Confirm",cleanup:"Storage Cleanup",cleanupDesc:"Remove orphaned directories that exist on disk but have no corresponding site in the database.",checkOrphans:"Check for Orphans",noOrphans:"No orphaned directories found",orphanFound:"Orphaned directories found",confirmDelete:"Delete All",cleanupConfirmWarning:"Delete all orphaned directories? This cannot be undone.",cleanupDone:"Cleanup complete",scanning:"Scanning...",files:"Files",noFiles:"No files",visits:"Visits",visitsToday:"Today",visitsMonth:"Month",visitsTotal:"Total",systemUpdate:"System Update",updateCheck:"Check for Updates",updateAvailable:"Update Available",updateToDate:"Up to Date",currentVersion:"Current Version",latestVersion:"Latest Version",updateNow:"Update Now",updating:"Updating...",updateSuccess:"Update successful! Please restart vibecast.",updateFailed:"Update failed",noUpdate:"You are running the latest version.",checkingUpdate:"Checking...",releaseNotes:"Release Notes",viewRelease:"View on GitHub",updateRestart:"Restart Now",updateRestarting:"Restarting...",updateRestartSuccess:"Server is restarting. Page will reload automatically.",updateInProgress:"Update already in progress",updateVerifyFailed:"Checksum verification failed",updateNoChecksum:"No checksum available (skipped verification)",updateDownloadProgress:"Downloading",systemInfo:"System Info",systemInfoDesc:"Runtime environment and resource paths",tabOverview:"Overview",tabUsers:"Users",tabSites:"Sites",tabSettings:"Settings",tabMaintenance:"Maintenance",tabSystem:"System",siVersion:"Version",siStoragePath:"Storage Path",siDbPath:"Database Path",siListenAddr:"Listen Address",siGoVersion:"Go Runtime",siPlatform:"Platform",siStartTime:"Started At",siUptime:"Uptime",siLanURLs:"LAN Access URLs",maxUploadSize:"Max Upload Size (MB)",maxUploadSizeDesc:"Maximum file size for uploads. Default: 50 MB.",maxSitesPerUser:"Max Sites per User",maxSitesPerUserDesc:"Maximum number of sites each user can create. Set to 0 for unlimited."},zh:{overview:"概览",settings:"设置",users:"用户",allSites:"所有站点",userCount:"用户数",siteCount:"站点数",adminCount:"管理员数",openReg:"开放注册",openRegDesc:"启用后，任何人都可以注册新账号。",publicAccess:"公开访问",publicAccessDesc:"禁用后，所有已部署站点对访客返回 403。",dashboard:"仪表盘",logout:"退出",role:"角色",created:"创建时间",actions:"操作",promote:"提升",demote:"降级",delete:"删除",deleteUserConfirm:"删除该用户及其所有站点？",deleteSiteConfirm:"删除此站点？",name:"名称",slug:"Slug",owner:"所有者",protected:"已加密",public:"公开",password:"密码",storagePath:"存储",url:"URL",noUsers:"暂无用户",noSites:"暂无站点",updated:"已更新",deleted:"已删除",roleUpdated:"角色已更新",userDeleted:"用户已删除",siteDeleted:"站点已删除",settingsSaved:"设置已保存",none:"无",accessDisabled:"已禁用",search:"搜索",searchUsersPh:"搜索用户...",searchSitesPh:"搜索站点...",prev:"上一页",next:"下一页",page:"第",of:"页 共",ownerEmail:"所有者",domainRestriction:"邮箱域名限制",domainRestrictionDesc:"仅允许指定域名的邮箱注册（每行一个）。",save:"保存",cancel:"取消",confirm:"确认",cleanup:"存储清理",cleanupDesc:"删除磁盘上存在但数据库中无对应站点的孤立目录。",checkOrphans:"检查孤立目录",noOrphans:"未发现孤立目录",orphanFound:"发现孤立目录",confirmDelete:"全部删除",cleanupConfirmWarning:"确认删除所有孤立目录？此操作不可撤销。",cleanupDone:"清理完成",scanning:"扫描中...",files:"文件",noFiles:"暂无文件",visits:"访问",visitsToday:"今日",visitsMonth:"本月",visitsTotal:"总计",systemUpdate:"系统更新",updateCheck:"检查更新",updateAvailable:"有可用更新",updateToDate:"已是最新版本",currentVersion:"当前版本",latestVersion:"最新版本",updateNow:"立即更新",updating:"更新中...",updateSuccess:"更新成功！请重启 vibecast。",updateFailed:"更新失败",noUpdate:"当前已是最新版本。",checkingUpdate:"检查中...",releaseNotes:"更新日志",viewRelease:"在 GitHub 查看",updateRestart:"立即重启",updateRestarting:"正在重启...",updateRestartSuccess:"服务器正在重启，页面将自动刷新。",updateInProgress:"更新正在进行中",updateVerifyFailed:"校验和验证失败",updateNoChecksum:"无校验和（已跳过验证）",updateDownloadProgress:"下载中",systemInfo:"系统信息",systemInfoDesc:"运行环境与资源路径",tabOverview:"概览",tabUsers:"用户",tabSites:"站点",tabSettings:"设置",tabMaintenance:"维护",tabSystem:"系统",siVersion:"版本",siStoragePath:"存储路径",siDbPath:"数据库路径",siListenAddr:"监听地址",siGoVersion:"Go 运行时",siPlatform:"平台",siStartTime:"启动时间",siUptime:"运行时长",siLanURLs:"内网访问地址",maxUploadSize:"最大上传大小 (MB)",maxUploadSizeDesc:"上传文件的最大大小。默认：50 MB。",maxSitesPerUser:"每用户站点上限",maxSitesPerUserDesc:"每个用户可创建的最大站点数。设为 0 则不限制。"}};
function t(k){return(i18n[lang]||i18n.en)[k]||(i18n.en[k]||k)}
function setLang(l){lang=l;try{localStorage.setItem("lang",l)}catch(e){}renderAdmin()}
function esc(s){return String(s||"").replace(/[&<>"']/g,function(c){return{"&":"&amp;","<":"&lt;",">":"&gt;",'"':"&quot;","'":"&#39;"}[c]})}
function siteUrl(u){return BASE+(u||"").replace(/^\//,"")}
function fmtDate(s){if(!s)return"-";var d=new Date(s);return d.toLocaleDateString()+" "+d.toLocaleTimeString([],{hour:"2-digit",minute:"2-digit"})}
function toast(msg,type){type=type||"success";var el=document.createElement("div");el.className="toast "+type+" show";el.textContent=msg;document.body.appendChild(el);var ms=type==="error"?5000:2500;setTimeout(function(){el.classList.remove("show");setTimeout(function(){el.remove()},200)},ms)}
function togglePwd(btn){var inp=btn.parentElement.querySelector("input");if(!inp)return;var show=inp.type==="password";inp.type=show?"text":"password";btn.innerHTML=show?PWD_HIDE_ICON:PWD_SHOW_ICON}
var PWD_SHOW_ICON='<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7z"/><circle cx="12" cy="12" r="3"/></svg>';
var PWD_HIDE_ICON='<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-10-8-10-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 10 8 10 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/><line x1="1" y1="1" x2="23" y2="23"/></svg>';
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
var adminTab="overview";
function renderAdmin(){
var lh='<div class="lang-toggle"><a id="theme-toggle" class="btn-icon" onclick="toggleTheme()">'+(getTheme()==="dark"?"☀":"🌙")+'</a><a class="'+(lang==="en"?"active":"")+'" onclick="setLang(\'en\')">EN</a><a class="'+(lang==="zh"?"active":"")+'" onclick="setLang(\'zh\')">中文</a></div>';
var icons={overview:'<svg viewBox="0 0 16 16" fill="currentColor"><path d="M1 2.5A1.5 1.5 0 0 1 2.5 1h3A1.5 1.5 0 0 1 7 2.5v3A1.5 1.5 0 0 1 5.5 7h-3A1.5 1.5 0 0 1 1 5.5v-3zm8 0A1.5 1.5 0 0 1 10.5 1h3A1.5 1.5 0 0 1 15 2.5v3A1.5 1.5 0 0 1 13.5 7h-3A1.5 1.5 0 0 1 9 5.5v-3zm-8 8A1.5 1.5 0 0 1 2.5 9h3A1.5 1.5 0 0 1 7 10.5v3A1.5 1.5 0 0 1 5.5 15h-3A1.5 1.5 0 0 1 1 13.5v-3zm8 0A1.5 1.5 0 0 1 10.5 9h3a1.5 1.5 0 0 1 1.5 1.5v3a1.5 1.5 0 0 1-1.5 1.5h-3A1.5 1.5 0 0 1 9 13.5v-3z"/></svg>',users:'<svg viewBox="0 0 16 16" fill="currentColor"><path d="M5.5 3.5a2 2 0 1 1-4 0 2 2 0 0 1 4 0zm7 0a2 2 0 1 1-4 0 2 2 0 0 1 4 0zM1.5 8a.5.5 0 0 1 .5-.5h3a.5.5 0 0 1 .5.5v4.5a.5.5 0 0 1-.5.5H2a.5.5 0 0 1-.5-.5V8zm7 0a.5.5 0 0 1 .5-.5h3a.5.5 0 0 1 .5.5v4.5a.5.5 0 0 1-.5.5H9a.5.5 0 0 1-.5-.5V8z"/></svg>',sites:'<svg viewBox="0 0 16 16" fill="currentColor"><path d="M1 3.5A1.5 1.5 0 0 1 2.5 2h11A1.5 1.5 0 0 1 15 3.5v9a1.5 1.5 0 0 1-1.5 1.5h-11A1.5 1.5 0 0 1 1 12.5v-9zM2.5 5v7.5h11V5h-11z"/></svg>',settings:'<svg viewBox="0 0 16 16" fill="currentColor"><path d="M8 4.754a3.246 3.246 0 1 0 0 6.492 3.246 3.246 0 0 0 0-6.492zM5.754 8a2.246 2.246 0 1 1 4.492 0 2.246 2.246 0 0 1-4.492 0z"/><path d="M9.796 1.343c-.527-1.79-3.065-1.79-3.592 0l-.094.319a.873.873 0 0 1-1.255.52l-.292-.16c-1.64-.892-3.433.902-2.54 2.541l.159.292a.873.873 0 0 1-.52 1.255l-.319.094c-1.79.527-1.79 3.065 0 3.592l.319.094a.873.873 0 0 1 .52 1.255l-.16.292c-.892 1.64.901 3.434 2.541 2.54l.292-.159a.873.873 0 0 1 1.255.52l.094.319c.527 1.79 3.065 1.79 3.592 0l.094-.319a.873.873 0 0 1 1.255-.52l.292.16c1.64.893 3.434-.902 2.54-2.541l-.159-.292a.873.873 0 0 1 .52-1.255l.319-.094c1.79-.527 1.79-3.065 0-3.592l-.319-.094a.873.873 0 0 1-.52-1.255l.16-.292c.893-1.64-.902-3.433-2.541-2.54l-.292.159a.873.873 0 0 1-1.255-.52l-.094-.319z"/></svg>',maintenance:'<svg viewBox="0 0 16 16" fill="currentColor"><path d="M2 1.5A1.5 1.5 0 0 1 3.5 0h9A1.5 1.5 0 0 1 14 1.5v13a.5.5 0 0 1-.5.5h-11a.5.5 0 0 1-.5-.5v-13zM3.5 1a.5.5 0 0 0-.5.5V14h10V1.5a.5.5 0 0 0-.5-.5h-9z"/><path d="M5 5.5a.5.5 0 0 1 .5-.5h5a.5.5 0 0 1 0 1h-5a.5.5 0 0 1-.5-.5zM5 8a.5.5 0 0 1 .5-.5h5a.5.5 0 0 1 0 1h-5A.5.5 0 0 1 5 8z"/></svg>',system:'<svg viewBox="0 0 16 16" fill="currentColor"><path d="M8 0a8 8 0 1 0 0 16A8 8 0 0 0 8 0zM1.5 8a6.5 6.5 0 1 1 13 0 6.5 6.5 0 0 1-13 0zM7.25 4a.75.75 0 0 1 1.5 0v3.69l2.28 2.28a.75.75 0 0 1-1.06 1.06L7.25 8.31V4z"/></svg>'};
var navItems=[["overview","tabOverview"],["users","tabUsers"],["sites","tabSites"],["settings","tabSettings"],["maintenance","tabMaintenance"],["system","tabSystem"]];
var navHtml='';
for(var i=0;i<navItems.length;i++){var k=navItems[i][0],lk=navItems[i][1];navHtml+='<div class="admin-nav-item'+(adminTab===k?" active":"")+'" onclick="switchTab(\''+k+'\')">'+(icons[k]||'')+'<span>'+t(lk)+'</span></div>'}
document.getElementById("app").innerHTML='<nav class="navbar"><div class="logo">` + logoIcon + `<span class="logo-text">Vibecast Admin <span id="ver" style="font-size:.6rem;color:var(--dim);font-weight:400"></span></span></div><div class="nav-right">'+lh+'<a href="'+BASE+'dashboard" style="font-size:.8rem;color:var(--dim)">'+t("dashboard")+'</a><button class="btn-link" onclick="doLogout()">'+t("logout")+'</button></div></nav><div class="admin-layout"><div class="admin-sidebar">'+navHtml+'</div><div class="admin-main"><div class="admin-tab-content" id="tab-overview"><div class="card"><div class="card-header"><h2>'+t("overview")+'</h2></div><div class="card-body"><div id="stats" class="stats-grid"></div></div></div></div><div class="admin-tab-content" id="tab-users"><div class="card"><div class="card-header"><h2>'+t("users")+'</h2></div><div class="card-body"><div class="list-toolbar"><input type="text" id="user-search" placeholder="'+t("searchUsersPh")+'" onkeydown="if(event.key===\'Enter\')searchUsers()" value="'+esc(userSearch)+'"></div><div id="users"></div></div></div></div><div class="admin-tab-content" id="tab-sites"><div class="card"><div class="card-header"><h2>'+t("allSites")+'</h2></div><div class="card-body"><div class="list-toolbar"><input type="text" id="admin-site-search" placeholder="'+t("searchSitesPh")+'" onkeydown="if(event.key===\'Enter\')searchAdminSites()" value="'+esc(adminSiteSearch)+'"></div><div id="sites"></div></div></div></div><div class="admin-tab-content" id="tab-settings"><div class="card"><div class="card-header"><h2>'+t("settings")+'</h2></div><div class="card-body"><div id="settings"></div></div></div></div><div class="admin-tab-content" id="tab-maintenance"><div class="card"><div class="card-header"><h2>'+t("cleanup")+'</h2></div><div class="card-body"><div id="cleanup-section"><p style="font-size:.8rem;color:var(--dim);margin-bottom:.8rem">'+t("cleanupDesc")+'</p><button class="btn btn-outline" onclick="checkCleanup()">'+t("checkOrphans")+'</button><div id="cleanup-result" style="margin-top:.8rem"></div></div></div></div></div><div class="admin-tab-content" id="tab-system"><div class="card"><div class="card-header"><h2>'+t("systemInfo")+'</h2></div><div class="card-body"><p style="font-size:.8rem;color:var(--dim);margin-bottom:1rem">'+t("systemInfoDesc")+'</p><div id="sysinfo" class="info-grid"></div><div class="update-inline"><div id="update-section"><p style="font-size:.8rem;color:var(--dim);margin-bottom:.8rem">'+t("currentVersion")+': <span id="update-current-ver" style="font-family:var(--mono);color:var(--text)">—</span></p><button class="btn btn-outline" id="update-check-btn" onclick="checkUpdate()">'+t("updateCheck")+'</button><div id="update-result" style="margin-top:.8rem"></div></div></div></div></div></div></div></div>';
switchTab(adminTab);
loadStats();loadSettings();loadUsers();loadSites();loadSystemInfo();
fetch(BASE+"api/version").then(function(r){return r.json()}).then(function(d){var ver=(d.data?d.data.version:"dev");var v=document.getElementById("ver");if(v)v.textContent="v"+ver;var uv=document.getElementById("update-current-ver");if(uv)uv.textContent="v"+ver}).catch(function(){})
}
function switchTab(tab){
adminTab=tab;
var items=document.querySelectorAll(".admin-nav-item");
for(var i=0;i<items.length;i++){items[i].classList.remove("active")}
var contents=document.querySelectorAll(".admin-tab-content");
for(var i=0;i<contents.length;i++){contents[i].classList.remove("active")}
var navMap={"overview":0,"users":1,"sites":2,"settings":3,"maintenance":4,"system":5};
if(items[navMap[tab]])items[navMap[tab]].classList.add("active");
var el=document.getElementById("tab-"+tab);
if(el)el.classList.add("active")
}
function loadSystemInfo(){
api("/admin/system-info").then(function(d){var r=d.data||{};var el=document.getElementById("sysinfo");if(!el)return;
var items=[[t("siVersion"),r.version,"accent"],[t("siPlatform"),(r.os||"?")+"/"+(r.arch||"?"),""],[t("siGoVersion"),r.goVersion||"-",""],[t("siListenAddr"),r.listenAddr||"-",""],[t("siStoragePath"),r.storagePath||"-",""],[t("siDbPath"),r.dbPath||"-",""],[t("siStartTime"),r.startTime||"-",""],[t("siUptime"),r.startTime?fmtUptime(r.startTime):"-",""]];
var h='';
for(var i=0;i<items.length;i++){h+='<div class="info-item"><div class="info-label">'+items[i][0]+'</div><div class="info-value'+(items[i][2]?" "+items[i][2]:"")+'">'+esc(items[i][1])+'</div></div>'}
var ips=r.localIPs||[];var port=r.listenAddr?r.listenAddr.split(":").pop():"8080";
if(ips.length>0){var urls="";for(var i=0;i<ips.length;i++){urls+='<a href="http://'+ips[i]+':'+port+'" target="_blank" style="display:block;font-family:var(--mono);font-size:.8rem;color:var(--accent);margin-bottom:.2rem">http://'+ips[i]+':'+port+'</a>'}h+='<div class="info-item"><div class="info-label">'+t("siLanURLs")+'</div><div class="info-value">'+urls+'</div></div>'}
el.innerHTML=h}).catch(function(e){var el=document.getElementById("sysinfo");if(el)el.innerHTML='<div style="color:var(--danger);font-size:.8rem">'+esc(e.message)+'</div>'})
}
function fmtUptime(startStr){var start=new Date(startStr);var now=new Date();var diff=Math.floor((now-start)/1000);if(diff<60)return diff+"s";var m=Math.floor(diff/60),s=diff%60;if(m<60)return m+"m "+s+"s";var h=Math.floor(m/60),mm=m%60;if(h<24)return h+"h "+mm+"m";var d=Math.floor(h/24),hh=h%24;return d+"d "+hh+"h"}
function loadStats(){
api("/admin/stats").then(function(d){var s=d.data;document.getElementById("stats").innerHTML='<div class="stat-card"><div class="num">'+s.users+'</div><div class="label">'+t("userCount")+'</div></div><div class="stat-card"><div class="num">'+s.sites+'</div><div class="label">'+t("siteCount")+'</div></div><div class="stat-card"><div class="num">'+s.admins+'</div><div class="label">'+t("adminCount")+'</div></div>'}).catch(function(e){toast(e.message,"error")})
}
function loadSettings(){
api("/admin/settings").then(function(d){var s=d.data;
var regOn=s.openRegistration,pubOn=s.allowPublicAccess!==false,drOn=s.domainRestriction;
var mus=s.maxUploadSize||50;
var h='<div class="toggle-row"><div class="toggle-info"><div class="toggle-label">'+t("openReg")+'</div><div class="toggle-desc">'+t("openRegDesc")+'</div></div><div class="toggle-switch '+(regOn?"on":"")+'" onclick="toggleReg()"></div></div><div class="toggle-row"><div class="toggle-info"><div class="toggle-label">'+t("publicAccess")+'</div><div class="toggle-desc">'+t("publicAccessDesc")+'</div></div><div class="toggle-switch '+(pubOn?"on":"")+'" onclick="togglePub()"></div></div><div class="toggle-row"><div class="toggle-info"><div class="toggle-label">'+t("domainRestriction")+'</div><div class="toggle-desc">'+t("domainRestrictionDesc")+'</div></div><div class="toggle-switch '+(drOn?"on":"")+'" onclick="toggleDomainRestriction()"></div></div>';
if(drOn){h+='<div style="padding:.6rem 0"><textarea id="allowed-domains" rows="4" style="width:100%;padding:8px" placeholder="example.com&#10;gmail.com">'+esc(s.allowedDomains||"")+'</textarea><div style="margin-top:.4rem"><button class="btn btn-promote" onclick="saveDomains()">'+t("save")+'</button></div></div>'}
var mspu=s.maxSitesPerUser||30;
h+='<div class="num-setting"><div class="num-info"><div class="num-label">'+t("maxUploadSize")+'</div><div class="num-desc">'+t("maxUploadSizeDesc")+'</div></div><div class="num-control"><input type="number" id="max-upload-size" value="'+mus+'" min="1" max="1024"><span class="unit">MB</span><button class="btn btn-promote" onclick="saveUploadSize()">'+t("save")+'</button></div></div>';
h+='<div class="num-setting"><div class="num-info"><div class="num-label">'+t("maxSitesPerUser")+'</div><div class="num-desc">'+t("maxSitesPerUserDesc")+'</div></div><div class="num-control"><input type="number" id="max-sites-per-user" value="'+mspu+'" min="0" max="10000"><button class="btn btn-promote" onclick="saveMaxSites()">'+t("save")+'</button></div></div>';
document.getElementById("settings").innerHTML=h}).catch(function(e){toast(e.message,"error")})
}
function toggleReg(){api("/admin/settings").then(function(d){var c=d.data;return api("/admin/settings",{method:"PUT",body:JSON.stringify({openRegistration:!c.openRegistration,allowPublicAccess:c.allowPublicAccess!==false,domainRestriction:c.domainRestriction,allowedDomains:c.allowedDomains||"",maxUploadSize:c.maxUploadSize||50,maxSitesPerUser:c.maxSitesPerUser||30})}).then(function(){toast(t("settingsSaved"));loadSettings()})}).catch(function(e){toast(e.message,"error")})}
function togglePub(){api("/admin/settings").then(function(d){var c=d.data;return api("/admin/settings",{method:"PUT",body:JSON.stringify({openRegistration:c.openRegistration,allowPublicAccess:!(c.allowPublicAccess!==false),domainRestriction:c.domainRestriction,allowedDomains:c.allowedDomains||"",maxUploadSize:c.maxUploadSize||50,maxSitesPerUser:c.maxSitesPerUser||30})}).then(function(){toast(t("settingsSaved"));loadSettings()})}).catch(function(e){toast(e.message,"error")})}
function toggleDomainRestriction(){api("/admin/settings").then(function(d){var c=d.data;return api("/admin/settings",{method:"PUT",body:JSON.stringify({openRegistration:c.openRegistration,allowPublicAccess:c.allowPublicAccess!==false,domainRestriction:!c.domainRestriction,allowedDomains:c.allowedDomains||"",maxUploadSize:c.maxUploadSize||50,maxSitesPerUser:c.maxSitesPerUser||30})}).then(function(){toast(t("settingsSaved"));loadSettings()})}).catch(function(e){toast(e.message,"error")})}
function saveDomains(){var v=document.getElementById("allowed-domains").value;api("/admin/settings").then(function(d){return api("/admin/settings",{method:"PUT",body:JSON.stringify({openRegistration:d.data.openRegistration,allowPublicAccess:d.data.allowPublicAccess!==false,domainRestriction:d.data.domainRestriction,allowedDomains:v,maxUploadSize:d.data.maxUploadSize||50,maxSitesPerUser:d.data.maxSitesPerUser||30})}).then(function(){toast(t("settingsSaved"));loadSettings()})}).catch(function(e){toast(e.message,"error")})}
function saveUploadSize(){var v=parseInt(document.getElementById("max-upload-size").value)||50;if(v<1)v=1;if(v>1024)v=1024;api("/admin/settings").then(function(d){return api("/admin/settings",{method:"PUT",body:JSON.stringify({openRegistration:d.data.openRegistration,allowPublicAccess:d.data.allowPublicAccess!==false,domainRestriction:d.data.domainRestriction,allowedDomains:d.data.allowedDomains||"",maxUploadSize:v})}).then(function(){toast(t("settingsSaved"));loadSettings()})}).catch(function(e){toast(e.message,"error")})}
function saveMaxSites(){var v=parseInt(document.getElementById("max-sites-per-user").value)||30;if(v<0)v=0;if(v>10000)v=10000;api("/admin/settings").then(function(d){return api("/admin/settings",{method:"PUT",body:JSON.stringify({openRegistration:d.data.openRegistration,allowPublicAccess:d.data.allowPublicAccess!==false,domainRestriction:d.data.domainRestriction,allowedDomains:d.data.allowedDomains||"",maxUploadSize:d.data.maxUploadSize||50,maxSitesPerUser:v})}).then(function(){toast(t("settingsSaved"));loadSettings()})}).catch(function(e){toast(e.message,"error")})}
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
h+='<button class="btn btn-promote" style="margin-top:.6rem" onclick="restartServer()">'+t("updateRestart")+'</button>';
el.innerHTML=h;
toast(t("updateSuccess"),"success")
}).catch(function(e){
var h='<div style="color:var(--danger);font-size:.8rem">✗ '+t("updateFailed")+': '+esc(e.message)+'</div><button class="btn btn-outline" style="margin-top:.4rem" onclick="checkUpdate()">'+t("updateCheck")+'</button>';
el.innerHTML=h;toast(e.message,"error")
})
})
}
function restartServer(){
var el=document.getElementById("update-result");
el.innerHTML+='<div style="margin-top:.6rem"><span style="color:var(--accent);font-size:.8rem;font-family:var(--mono)">⟳ '+t("updateRestarting")+'</span></div>';
api("/admin/update/restart",{method:"POST"}).then(function(){
toast(t("updateRestartSuccess"),"success");
setTimeout(function(){location.reload()},3000)
}).catch(function(e){
el.innerHTML+='<div style="color:var(--danger);font-size:.8rem">'+esc(e.message)+'</div>'
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
var h='<table><thead><tr><th>ID</th><th>'+t("name")+'</th><th>'+t("slug")+'</th><th>'+t("owner")+'</th><th>'+t("protected")+'</th><th>'+t("password")+'</th><th>'+t("visits")+'</th><th>'+t("url")+'</th><th>'+t("actions")+'</th></tr></thead><tbody>';
for(var i=0;i<sites.length;i++){var s=sites[i],b="";
if(s.publicAccessDisabled&&!s.protected){b='<span class="badge badge-disabled">'+t("accessDisabled")+'</span>'}
else if(s.protected){b='<span class="badge badge-protected">'+t("protected")+'</span>'}
else{b='<span class="badge badge-public">'+t("public")+'</span>'}
var pwd=s.protected?'<code style="font-family:var(--mono);font-size:.7rem">'+esc(s.password)+'</code>':'<span style="color:var(--dim)">'+t("none")+'</span>';
var vs=s.visits||{today:0,month:0,total:0};
var vis='<span style="font-family:var(--mono);font-size:.7rem">'+vs.today+'<span style="color:var(--dim)">/</span>'+vs.month+'<span style="color:var(--dim)">/</span><b style="color:var(--accent)">'+vs.total+'</b></span>';
h+='<tr class="site-row" id="row-'+s.id+'" onclick="adminToggleDetail('+s.id+')"><td style="font-family:var(--mono);color:var(--dim)">'+s.id+'</td><td>'+esc(s.name)+'</td><td style="font-family:var(--mono);font-size:.75rem">'+esc(s.slug)+'</td><td style="font-size:.75rem">'+esc(s.ownerEmail||"-")+'</td><td>'+b+'</td><td>'+pwd+'</td><td title="'+t("visitsToday")+"/"+t("visitsMonth")+"/"+t("visitsTotal")+'">'+vis+'</td><td><a href="'+siteUrl(s.url)+'" target="_blank" style="font-family:var(--mono);font-size:.75rem" onclick="event.stopPropagation()">'+s.url+'</a></td><td onclick="event.stopPropagation()"><button class="btn btn-danger" onclick="delSite('+s.id+')">'+t("delete")+'</button></td></tr>';
h+='<tr class="admin-detail-row" id="detail-'+s.id+'"><td colspan="9"></td></tr>'}
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

var orgInfo=null,orgMembers=[],orgMemberPage=1,orgMemberPerPage=10,orgMemberTotal=0,orgMemberSearch="";
function loadOrg(){
api("/org").then(function(d){
var r=d.data||{};orgInfo=r.hasOrg?r:null;
updateOrgNavBtn();
renderOrgSection();
}).catch(function(){})
}
function renderOrgSection(){
var el=document.getElementById("org-section");if(!el)return;
if(!orgInfo){
el.innerHTML='<p class="desc" style="margin-bottom:.8rem">'+t("noOrg")+'</p><div class="form-field"><label>'+t("orgName")+'</label><input id="org-create-name" placeholder="'+t("orgNamePh")+'"></div><button class="btn btn-primary" style="width:100%;margin-bottom:.5rem" onclick="createOrg()">'+t("orgCreate")+'</button><div style="border-top:1px solid var(--line);margin:.8rem 0;padding-top:.8rem"><div class="form-field"><label>'+t("orgInviteCode")+'</label><input id="org-join-code" placeholder="'+t("orgInviteCodePh")+'"></div><button class="btn btn-outline" style="width:100%" onclick="joinOrg()">'+t("orgJoin")+'</button></div>';
return
}
var inviteCode=orgInfo.inviteCode||"";
var name=orgInfo.name||"";
var isOwner=orgInfo.isOwner||false;
var h='<div class="org-info"><div class="org-name">'+esc(name)+'</div>';
if(isOwner){h+='<div class="org-invite"><label>'+t("orgInviteCode")+'</label><div class="invite-row"><code class="invite-code">'+esc(inviteCode)+'</code><button class="copy-btn" onclick="copyText(\''+esc(inviteCode)+'\')" title="'+t("copy")+'"><svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg></button></div></div>'}
h+='</div><div class="org-actions">';
if(isOwner){h+='<button class="btn btn-danger btn-sm" onclick="deleteOrg()">'+t("orgDelete")+'</button>'}
else{h+='<button class="btn btn-outline btn-sm" onclick="leaveOrg()">'+t("orgLeave")+'</button>'}
h+='</div><div class="org-members-section"><h3>'+t("orgMembers")+'</h3><div class="list-toolbar"><input type="text" id="org-member-search" placeholder="'+t("search")+'..." onkeydown="if(event.key===\'Enter\')searchOrgMembers()" value="'+esc(orgMemberSearch)+'"></div><div id="org-member-list"></div></div>';
el.innerHTML=h;
loadOrgMembers()
}
function orgMemberPageGo(p){orgMemberPage=p;loadOrgMembers()}
function searchOrgMembers(){orgMemberSearch=document.getElementById("org-member-search").value;orgMemberPage=1;loadOrgMembers()}
function loadOrgMembers(){
var q=orgMemberSearch?"&q="+encodeURIComponent(orgMemberSearch):"";
api("/org/members?page="+orgMemberPage+"&perPage="+orgMemberPerPage+q).then(function(d){
var r=d.data||{},members=r.items||[];orgMemberTotal=r.total||0;
var el=document.getElementById("org-member-list");if(!el)return;
var tp=Math.ceil(orgMemberTotal/orgMemberPerPage)||1,pg=paginationHtml(orgMemberPage,tp,"orgMemberPageGo");
if(!members.length){el.innerHTML='<div class="empty">'+t("orgNoMembers")+'</div>'+pg;return}
var h='<ul class="site-list">';
for(var i=0;i<members.length;i++){var m=members[i];
var role=m.isOwner?'<span class="badge badge-protected">'+t("orgOwner")+'</span>':'<span class="badge badge-public">'+t("orgMember")+'</span>';
var rmv=m.isOwner?'':' <button class="btn btn-sm btn-danger" onclick="removeOrgMember('+m.userId+')">'+t("orgRemove")+'</button>';
h+='<li class="site-item"><div class="site-head"><div class="info"><div class="name">'+esc(m.email)+' '+role+'</div></div><div class="actions">'+rmv+'</div></div></li>'}
h+='</ul>'+pg;
el.innerHTML=h
}).catch(function(){})
}
function createOrg(){
var name=document.getElementById("org-create-name").value.trim();
api("/org",{method:"POST",body:JSON.stringify({name:name})}).then(function(){toast(t("orgCreated"));loadOrg()}).catch(function(e){toast(e.message,"error")})
}
function joinOrg(){
var code=document.getElementById("org-join-code").value.trim();
if(!code){toast(t("orgInviteCodePh"),"error");return}
api("/org/join",{method:"POST",body:JSON.stringify({inviteCode:code})}).then(function(){toast(t("orgJoined"));loadOrg()}).catch(function(e){toast(e.message,"error")})
}
function leaveOrg(){
customConfirm(t("orgLeaveConfirm"),function(){api("/org/leave",{method:"POST"}).then(function(){toast(t("orgLeft"));loadOrg()}).catch(function(e){toast(e.message,"error")})})
}
function deleteOrg(){
customConfirm(t("orgDeleteConfirm"),function(){api("/org",{method:"DELETE"}).then(function(){toast(t("orgDeleted"));loadOrg()}).catch(function(e){toast(e.message,"error")})})
}
function removeOrgMember(uid){
api("/org/members/"+uid,{method:"DELETE"}).then(function(){toast(t("orgMemberRemoved"));loadOrgMembers()}).catch(function(e){toast(e.message,"error")})
}
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
%s
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
input[type=password],input[type=text]{width:100%%;padding:10px 12px;border:1px solid var(--line);border-radius:6px;font-size:.85rem;margin-bottom:.8rem;outline:none;background:var(--ink);color:var(--text)}
input[type=password]:focus,input[type=text]:focus{border-color:var(--accent)}
input[type=password]::placeholder,input[type=text]::placeholder{color:var(--placeholder)}
input[type=password]::-ms-reveal,input[type=text]::-ms-reveal,input[type=password]::-ms-clear,input[type=text]::-ms-clear{display:none}
.pwd-wrap{position:relative;display:flex;align-items:center}
.pwd-wrap input{width:100%%;padding-right:38px!important;box-sizing:border-box;margin-bottom:0}
.pwd-toggle{position:absolute;right:8px;top:50%%;transform:translateY(-50%%);background:none;border:none;cursor:pointer;display:flex;align-items:center;justify-content:center;width:28px;height:28px;border-radius:5px;color:var(--dim);padding:0;transition:color .15s,background .15s}
.pwd-toggle:hover{color:var(--text);background:var(--ink)}
.pwd-toggle svg{width:16px;height:16px;flex-shrink:0}
button[type=submit]{width:100%%;padding:10px;background:var(--accent);color:var(--ink);border:none;border-radius:6px;font-family:var(--mono);font-size:.85rem;font-weight:700;cursor:pointer;margin-top:.8rem}
button[type=submit]:hover{background:var(--accent-dim)}
.err{color:var(--danger);margin-bottom:.8rem;font-size:.8rem;font-family:var(--mono)}
</style>
<script>var _t="light";try{_t=localStorage.getItem("theme")||"light"}catch(e){}document.documentElement.setAttribute("data-theme",_t)
function togglePwd(btn){var inp=btn.parentElement.querySelector("input");if(!inp)return;var show=inp.type==="password";inp.type=show?"text":"password";if(show){btn.innerHTML='<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-10-8-10-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 10 8 10 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/><line x1="1" y1="1" x2="23" y2="23"/></svg>'}else{btn.innerHTML='<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7z"/><circle cx="12" cy="12" r="3"/></svg>'}}
</script>
</head>
<body>
<div class="card">
<h1>🔒 <span class="site-name">%s</span></h1>
<p>This site is password-protected. Enter the password to continue.</p>
<form method="POST" action="">
<div class="pwd-wrap"><input type="password" id="password" name="password" placeholder="Password" autofocus required><button type="button" class="pwd-toggle" onclick="togglePwd(this)"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7z"/><circle cx="12" cy="12" r="3"/></svg></button></div>
<button type="submit">Enter Site</button>
</form>
</div>
</body>
</html>`, logoFavicon, escHTML(siteName), escHTML(siteName))
}

// passwordPageHTMLWithErr returns the password gate page with an error message.
func passwordPageHTMLWithErr(slug, siteName, errMsg string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
%s
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
input[type=password],input[type=text]{width:100%%;padding:10px 12px;border:1px solid var(--line);border-radius:6px;font-size:.85rem;margin-bottom:.8rem;outline:none;background:var(--ink);color:var(--text)}
input[type=password]:focus,input[type=text]:focus{border-color:var(--accent)}
input[type=password]::placeholder,input[type=text]::placeholder{color:var(--placeholder)}
input[type=password]::-ms-reveal,input[type=text]::-ms-reveal,input[type=password]::-ms-clear,input[type=text]::-ms-clear{display:none}
.pwd-wrap{position:relative;display:flex;align-items:center}
.pwd-wrap input{width:100%%;padding-right:38px!important;box-sizing:border-box;margin-bottom:0}
.pwd-toggle{position:absolute;right:8px;top:50%%;transform:translateY(-50%%);background:none;border:none;cursor:pointer;display:flex;align-items:center;justify-content:center;width:28px;height:28px;border-radius:5px;color:var(--dim);padding:0;transition:color .15s,background .15s}
.pwd-toggle:hover{color:var(--text);background:var(--ink)}
.pwd-toggle svg{width:16px;height:16px;flex-shrink:0}
button[type=submit]{width:100%%;padding:10px;background:var(--accent);color:var(--ink);border:none;border-radius:6px;font-family:var(--mono);font-size:.85rem;font-weight:700;cursor:pointer;margin-top:.8rem}
button[type=submit]:hover{background:var(--accent-dim)}
.err{color:var(--danger);margin-bottom:.8rem;font-size:.8rem;font-family:var(--mono)}
</style>
<script>var _t="light";try{_t=localStorage.getItem("theme")||"light"}catch(e){}document.documentElement.setAttribute("data-theme",_t)
function togglePwd(btn){var inp=btn.parentElement.querySelector("input");if(!inp)return;var show=inp.type==="password";inp.type=show?"text":"password";if(show){btn.innerHTML='<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-10-8-10-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 10 8 10 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/><line x1="1" y1="1" x2="23" y2="23"/></svg>'}else{btn.innerHTML='<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7z"/><circle cx="12" cy="12" r="3"/></svg>'}}
</script>
</head>
<body>
<div class="card">
<h1>🔒 <span class="site-name">%s</span></h1>
<p>This site is password-protected. Enter the password to continue.</p>
<div class="err">%s</div>
<form method="POST" action="">
<div class="pwd-wrap"><input type="password" name="password" placeholder="Password" autofocus required><button type="button" class="pwd-toggle" onclick="togglePwd(this)"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7z"/><circle cx="12" cy="12" r="3"/></svg></button></div>
<button type="submit">Enter Site</button>
</form>
</div>
</body>
</html>`, logoFavicon, escHTML(siteName), escHTML(siteName), escHTML(errMsg))
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
