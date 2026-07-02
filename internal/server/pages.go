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
<style>
:root{--bg:#f8f9fa;--card:#fff;--border:#e2e8f0;--text:#1a202c;--muted:#718096;--primary:#6366f1;--primary-hover:#5563d1;--primary-light:#eef2ff;--green:#059669;--red:#dc2626;--amber:#d97706;--radius:10px;--shadow:0 1px 3px rgba(0,0,0,.08),0 1px 2px rgba(0,0,0,.04)}
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,"Helvetica Neue",sans-serif;background:var(--bg);color:var(--text);min-height:100vh}
a{color:var(--primary);text-decoration:none}
.wrap{max-width:720px;margin:0 auto;padding:2rem 1.5rem}
.hero{text-align:center;padding:3rem 0}
.hero h1{font-size:2.8rem;font-weight:800;letter-spacing:-.02em;margin-bottom:.5rem}
.hero h1 span{background:linear-gradient(135deg,#6366f1,#8b5cf6);-webkit-background-clip:text;-webkit-text-fill-color:transparent}
.hero .tagline{font-size:1.15rem;color:var(--muted);margin-bottom:2rem}
.hero .cta{display:inline-block;padding:12px 32px;background:var(--primary);color:#fff;border:none;border-radius:var(--radius);font-size:1rem;font-weight:600;cursor:pointer;transition:background .2s;text-decoration:none}
.hero .cta:hover{background:var(--primary-hover)}
.features{display:grid;grid-template-columns:repeat(auto-fit,minmax(160px,1fr));gap:1rem;margin-top:3rem}
.feature{text-align:center;padding:1.5rem 1rem;background:var(--card);border:1px solid var(--border);border-radius:var(--radius);box-shadow:var(--shadow)}
.feature .icon{font-size:1.8rem;margin-bottom:.5rem}
.feature h3{font-size:.95rem;margin-bottom:.3rem}
.feature p{font-size:.8rem;color:var(--muted);line-height:1.5}
footer{text-align:center;padding:2rem 0;color:var(--muted);font-size:.8rem;border-top:1px solid var(--border);margin-top:3rem}
.lang-toggle{position:absolute;top:1.5rem;right:1.5rem}
.lang-toggle a{font-size:.85rem;cursor:pointer;padding:4px 10px;border-radius:6px}
.lang-toggle a.active{background:var(--primary-light);color:var(--primary)}
</style>
</head>
<body>
<div class="lang-toggle"><a id="langEn" class="active" onclick="setLang('en')">EN</a> <a id="langZh" onclick="setLang('zh')">中文</a></div>
<div class="wrap">
<div class="hero">
<h1><span data-i="title">Vibecast</span></h1>
<p class="tagline" data-i="tagline">Build with vibe. Cast instantly.</p>
<a href="/dashboard" class="cta" data-i="getStarted">Get Started</a>
<div class="features">
<div class="feature"><div class="icon">&#9889;</div><h3 data-i="feat1Title">Instant Deploy</h3><p data-i="feat1Desc">Upload a ZIP, get a live URL in seconds.</p></div>
<div class="feature"><div class="icon">&#128274;</div><h3 data-i="feat2Title">Password Protect</h3><p data-i="feat2Desc">Keep your site private with password gating.</p></div>
<div class="feature"><div class="icon">&#128230;</div><h3 data-i="feat3Title">Zero Nginx</h3><p data-i="feat3Desc">Pure Go application server. No dependencies.</p></div>
<div class="feature"><div class="icon">&#127379;</div><h3 data-i="feat4Title">Self-Hosted</h3><p data-i="feat4Desc">Your data, your rules. SQLite + filesystem.</p></div>
</div>
</div>
<footer data-i="footer">Vibecast — Self-hosted static site hosting</footer>
</div>
<script>
var i18n={en:{title:"Vibecast",tagline:"Build with vibe. Cast instantly.",getStarted:"Get Started",feat1Title:"Instant Deploy",feat1Desc:"Upload a ZIP, get a live URL in seconds.",feat2Title:"Password Protect",feat2Desc:"Keep your site private with password gating.",feat3Title:"Zero Nginx",feat3Desc:"Pure Go application server. No dependencies.",feat4Title:"Self-Hosted",feat4Desc:"Your data, your rules. SQLite + filesystem.",footer:"Vibecast — Self-hosted static site hosting"},zh:{title:"Vibecast",tagline:"Build with vibe. Cast instantly.",getStarted:"开始使用",feat1Title:"即时部署",feat1Desc:"上传 ZIP，秒级获取线上 URL。",feat2Title:"密码保护",feat2Desc:"通过密码门禁保护你的站点隐私。",feat3Title:"零依赖",feat3Desc:"纯 Go 应用服务器，无需 Nginx。",feat4Title:"自托管",feat4Desc:"数据归你所有，SQLite + 文件系统。",footer:"Vibecast — 自托管静态站点托管平台"}};
var lang="en";
function setLang(l){lang=l;document.querySelectorAll("[data-i]").forEach(function(e){var k=e.getAttribute("data-i");if(i18n[l][k])e.textContent=i18n[l][k]});document.getElementById("langEn").className=l==="en"?"active":"";document.getElementById("langZh").className=l==="zh"?"active":"";try{localStorage.setItem("lang",l)}catch(e){}}
var saved="en";try{saved=localStorage.getItem("lang")||"en"}catch(e){}setLang(saved);
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
<style>
:root{--bg:#f1f5f9;--card:#fff;--border:#e2e8f0;--text:#1e293b;--muted:#64748b;--primary:#6366f1;--primary-hover:#5563d1;--primary-light:#eef2ff;--green:#059669;--green-light:#d1fae5;--red:#dc2626;--red-light:#fee2e2;--amber:#d97706;--amber-light:#fef3c7;--radius:10px;--shadow:0 1px 3px rgba(0,0,0,.06),0 1px 2px rgba(0,0,0,.04)}
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,"Helvetica Neue",sans-serif;background:var(--bg);color:var(--text);min-height:100vh}
a{color:var(--primary);text-decoration:none}
input,button,select{font-family:inherit}
.navbar{display:flex;justify-content:space-between;align-items:center;padding:0 1.5rem;height:56px;background:#fff;border-bottom:1px solid var(--border);position:sticky;top:0;z-index:100}
.navbar .brand .logo{font-size:1.2rem;font-weight:800;background:linear-gradient(135deg,#6366f1,#8b5cf6);-webkit-background-clip:text;-webkit-text-fill-color:transparent}
.navbar .nav-right{display:flex;align-items:center;gap:1rem}
.navbar .nav-right .email{font-size:.85rem;color:var(--muted)}
.navbar .nav-right .btn-link{font-size:.85rem;cursor:pointer;color:var(--muted);background:none;border:none}
.navbar .nav-right .btn-link:hover{color:var(--text)}
.navbar .admin-link{font-size:.85rem;color:var(--amber);font-weight:600;cursor:pointer}
.lang-toggle a{font-size:.8rem;cursor:pointer;padding:3px 8px;border-radius:5px}
.lang-toggle a.active{background:var(--primary-light);color:var(--primary)}
.container{max-width:860px;margin:1.5rem auto;padding:0 1.5rem}
.card{background:var(--card);border:1px solid var(--border);border-radius:var(--radius);box-shadow:var(--shadow);margin-bottom:1.5rem}
.card-header{padding:1rem 1.25rem;border-bottom:1px solid var(--border);display:flex;justify-content:space-between;align-items:center}
.card-header h2{font-size:1rem;font-weight:600}
.card-header .hint{font-size:.75rem;color:var(--muted)}
.card-body{padding:1.25rem}
.form-grid{display:grid;grid-template-columns:1fr 1fr;gap:1rem;margin-bottom:1rem}
.form-field{display:flex;flex-direction:column;gap:.3rem}
.form-field.full{grid-column:1/-1}
.form-field label{font-size:.8rem;font-weight:600;color:var(--text)}
.form-field input{padding:9px 12px;border:1px solid var(--border);border-radius:8px;font-size:.9rem;outline:none;transition:border-color .15s;background:#fff;color:var(--text)}
.form-field input:focus{border-color:var(--primary);box-shadow:0 0 0 3px var(--primary-light)}
.form-field input::placeholder{color:#94a3b8}
.form-field .desc{font-size:.72rem;color:var(--muted)}
.form-actions{display:flex;justify-content:flex-end;margin-top:.5rem}
.btn{padding:9px 18px;border:none;border-radius:8px;font-size:.9rem;font-weight:600;cursor:pointer;transition:opacity .15s}
.btn-primary{background:var(--primary);color:#fff}
.btn-primary:hover{opacity:.9}
.btn-sm{padding:6px 12px;font-size:.8rem;border-radius:6px}
.btn-danger{background:var(--red);color:#fff}
.btn-danger:hover{opacity:.9}
.btn-outline{background:#fff;border:1px solid var(--border);color:var(--text)}
.btn-outline:hover{background:#f8fafc}
.btn-green{background:var(--green);color:#fff}
.site-list{list-style:none}
.site-item{border:1px solid var(--border);border-radius:8px;margin-bottom:.75rem;overflow:hidden;transition:border-color .15s}
.site-item:hover{border-color:#cbd5e1}
.site-item .site-head{display:flex;justify-content:space-between;align-items:center;padding:.85rem 1rem;cursor:pointer}
.site-item .site-head .name{font-weight:600;font-size:.9rem;display:flex;align-items:center;gap:.4rem}
.site-item .site-head .url{font-size:.8rem;color:var(--primary);margin-top:.15rem;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.site-item .site-head .actions{display:flex;gap:.4rem;align-items:center;flex-shrink:0}
.site-item .site-detail{padding:.75rem 1rem;border-top:1px solid var(--border);background:#f8fafc;font-size:.8rem;color:var(--muted);display:none}
.site-item .site-detail.show{display:block}
.site-item .site-detail .detail-row{display:flex;gap:.5rem;padding:.2rem 0}
.site-item .site-detail .detail-row .label{font-weight:600;min-width:90px;color:var(--text)}
.site-item .site-detail .detail-row .value{word-break:break-all}
.badge{display:inline-block;font-size:.7rem;padding:2px 7px;border-radius:4px;font-weight:500}
.badge-protected{background:var(--amber-light);color:var(--amber)}
.badge-public{background:var(--green-light);color:var(--green)}
.badge-disabled{background:var(--red-light);color:var(--red)}
.upload-btn{display:inline-block;padding:6px 12px;background:var(--green);color:#fff;border-radius:6px;font-size:.8rem;cursor:pointer;font-weight:600;position:relative;overflow:hidden}
.upload-btn:hover{opacity:.9}
.upload-btn input[type=file]{position:absolute;top:0;left:0;width:100%;height:100%;opacity:0;cursor:pointer}
.empty{text-align:center;color:var(--muted);padding:2rem;font-size:.9rem}
.toast{position:fixed;bottom:1.5rem;right:1.5rem;padding:10px 20px;border-radius:8px;font-size:.9rem;z-index:999;opacity:0;transform:translateY(10px);transition:all .25s}
.toast.show{opacity:1;transform:translateY(0)}
.toast.success{background:var(--green);color:#fff}
.toast.error{background:var(--red);color:#fff}
.auth-screen{display:flex;align-items:center;justify-content:center;min-height:100vh;padding:1rem}
.auth-card{background:var(--card);padding:2.5rem;border:1px solid var(--border);border-radius:14px;width:100%;max-width:380px;box-shadow:0 4px 6px -1px rgba(0,0,0,.05),0 2px 4px -2px rgba(0,0,0,.03)}
.auth-card h1{font-size:1.5rem;font-weight:800;text-align:center;margin-bottom:.3rem}
.auth-card h1 span{background:linear-gradient(135deg,#6366f1,#8b5cf6);-webkit-background-clip:text;-webkit-text-fill-color:transparent}
.auth-card .subtitle{color:var(--muted);text-align:center;margin-bottom:1.5rem;font-size:.85rem}
.auth-card input{width:100%;padding:11px 14px;border:1px solid var(--border);border-radius:8px;font-size:.9rem;margin-bottom:.6rem;outline:none;transition:border-color .15s}
.auth-card input:focus{border-color:var(--primary);box-shadow:0 0 0 3px var(--primary-light)}
.auth-card .btn{width:100%;margin-top:.3rem}
.auth-card .switch{text-align:center;margin-top:1rem;font-size:.85rem;color:var(--muted)}
.auth-card .switch a{color:var(--primary);cursor:pointer;font-weight:500}
.field-hint{font-size:.75rem;color:var(--muted);margin-top:.3rem;margin-bottom:.6rem}
.list-toolbar{display:flex;align-items:center;gap:.6rem;margin-bottom:.75rem;flex-wrap:wrap}
.list-toolbar input[type=text]{flex:1;min-width:120px;padding:7px 11px;border:1px solid var(--border);border-radius:7px;font-size:.85rem;outline:none;transition:border-color .15s;background:#fff;color:var(--text)}
.list-toolbar input[type=text]:focus{border-color:var(--primary);box-shadow:0 0 0 3px var(--primary-light)}
.pagination{display:flex;align-items:center;gap:.4rem;justify-content:center;padding:.75rem 0 0;flex-wrap:wrap}
.pagination .page-info{font-size:.78rem;color:var(--muted);margin:0 .4rem}
.pagination button{padding:4px 10px;border:1px solid var(--border);background:#fff;color:var(--text);border-radius:5px;font-size:.78rem;cursor:pointer;font-weight:500}
.pagination button:hover:not(:disabled){background:#f8fafc;border-color:#cbd5e1}
.pagination button.active{background:var(--primary);color:#fff;border-color:var(--primary)}
.pagination button:disabled{opacity:.4;cursor:default}
</style>
</head>
<body>
<div id="app"></div>
<script>
var API="/api";
var currentUser=null;
var lang="en";
var sitePage=1,sitePerPage=10,siteTotal=0,siteSearch="";
var i18n={en:{siteName:"Site Name",siteNamePh:"e.g. My Portfolio",slug:"URL Slug",slugPh:"my-portfolio",sitePwd:"Access Password",sitePwdPh:"Leave empty for public",create:"Create Site",yourSites:"Your Sites",noSites:"No sites yet. Create one above.",deployZip:"Deploy ZIP",delete:"Delete",deleteConfirm:"Delete this site? This removes all files.",protected:"Protected",public:"Public",login:"Login",register:"Register",email:"Email",password:"Password",pwdHint:"At least 6 characters",noAccount:"No account?",haveAccount:"Have an account?",logout:"Logout",adminPanel:"Admin Panel",deployed:"Deployed!",siteCreated:"Site created",deleted:"Deleted",incorrectPwd:"Incorrect password",loginFailed:"Login failed",registerFailed:"Registration failed",slugDesc:"The URL path of your site. Auto-generated from name if left blank. Only lowercase letters, numbers, and hyphens.",pwdDesc:"If set, visitors must enter this password to access your site.",sitesHint:"Click a site to expand details.",deployHint:"Upload a .zip file to deploy or update your site.",storagePath:"Storage Path",accessPassword:"Access Password",none:"None",details:"Details",expand:"Expand",filesSkipped:"Dangerous files skipped",accessDisabled:"Access Disabled",pwdRequired:"Public access is disabled — password is required",pwdOptional:"If set, visitors must enter this password to access your site.",search:"Search",searchSitesPh:"Search sites...",prev:"Prev",next:"Next",page:"Page",of:"of"},zh:{siteName:"站点名称",siteNamePh:"例如：我的作品集",slug:"URL Slug",slugPh:"my-site",sitePwd:"访问密码",sitePwdPh:"留空则公开访问",create:"创建站点",yourSites:"我的站点",noSites:"还没有站点，在上方创建一个。",deployZip:"部署 ZIP",delete:"删除",deleteConfirm:"确定删除此站点？所有文件将被移除。",protected:"已保护",public:"公开",login:"登录",register:"注册",email:"邮箱",password:"密码",pwdHint:"至少 6 个字符",noAccount:"没有账号？",haveAccount:"已有账号？",logout:"退出",adminPanel:"管理后台",deployed:"部署成功！",siteCreated:"站点已创建",deleted:"已删除",incorrectPwd:"密码错误",loginFailed:"登录失败",registerFailed:"注册失败",slugDesc:"站点的 URL 路径。留空则根据名称自动生成。只允许小写字母、数字和连字符。",pwdDesc:"设置后，访问者需要输入此密码才能查看站点。",sitesHint:"点击站点展开详情。",deployHint:"上传 .zip 文件来部署或更新站点。",storagePath:"存储路径",accessPassword:"访问密码",none:"无",details:"详情",expand:"展开",filesSkipped:"已跳过危险文件",accessDisabled:"已禁止访问",pwdRequired:"公开访问已关闭 — 必须设置密码",pwdOptional:"设置后，访问者需要输入此密码才能查看站点。",search:"搜索",searchSitesPh:"搜索站点...",prev:"上一页",next:"下一页",page:"第",of:"/ 共"}};
function t(k){return(i18n[lang]||i18n.en)[k]||(i18n.en[k]||k)}
function setLang(l){lang=l;try{localStorage.setItem("lang",l)}catch(e){}if(currentUser)renderDashboard();}
function esc(s){return String(s||"").replace(/[&<>"']/g,function(c){return{"&":"&amp;","<":"&lt;",">":"&gt;",'"':"&quot;","'":"&#39;"}[c]})}
function toast(msg,type){type=type||"success";var el=document.createElement("div");el.className="toast "+type+" show";el.textContent=msg;document.body.appendChild(el);setTimeout(function(){el.classList.remove("show");setTimeout(function(){el.remove()},300)},2500)}
function getToken(){try{return localStorage.getItem("vibecast_token")}catch(e){return""}}
function setToken(tk){try{localStorage.setItem("vibecast_token",tk)}catch(e){}}
function clearToken(){try{localStorage.removeItem("vibecast_token")}catch(e){}}
function api(path,opts){
  opts=opts||{};
  var token=localStorage.getItem("vibecast_token")||"";
  var headers={"Content-Type":"application/json"};
  if(token)headers["Authorization"]="Bearer "+token;
  if(opts.headers)Object.assign(headers,opts.headers);
  return fetch(API+path,Object.assign({},opts,{headers:headers,credentials:"same-origin"})).then(function(r){
    if(r.status===401){try{localStorage.removeItem("vibecast_token")}catch(e){}location.reload();}
    return r.json().catch(function(){return{error:"network error"}}).then(function(data){if(!r.ok)throw new Error(data.error||"request failed");return data})
  })
}
function checkAuth(){if(!getToken())return Promise.resolve(false);return api("/auth/me").then(function(d){currentUser=d.data;return true}).catch(function(){clearToken();return false})}

function renderAuth(){
var langHtml='<div class="lang-toggle" style="position:absolute;top:1rem;right:1rem"><a class="'+(lang==="en"?"active":"")+'" onclick="setLang(\'en\')">EN</a> <a class="'+(lang==="zh"?"active":"")+'" onclick="setLang(\'zh\')">中文</a></div>';
document.getElementById("app").innerHTML=langHtml+'<div class="auth-screen"><div class="auth-card"><h1><span>Vibecast</span></h1><p class="subtitle">Build with vibe. Cast instantly.</p><div id="auth-form"></div></div></div>';
showLogin();}
function showLogin(){
document.getElementById("auth-form").innerHTML='<input id="email" type="email" placeholder="'+t("email")+'" autocomplete="email"><input id="password" type="password" placeholder="'+t("password")+'" autocomplete="current-password"><button class="btn btn-primary" onclick="doLogin()">'+t("login")+'</button><p class="switch">'+t("noAccount")+' <a onclick="showRegister()">'+t("register")+'</a></p>';}
function showRegister(){
document.getElementById("auth-form").innerHTML='<input id="email" type="email" placeholder="'+t("email")+'" autocomplete="email"><input id="password" type="password" placeholder="'+t("password")+'" autocomplete="new-password"><div class="field-hint">'+t("pwdHint")+'</div><button class="btn btn-primary" onclick="doRegister()">'+t("register")+'</button><p class="switch">'+t("haveAccount")+' <a onclick="showLogin()">'+t("login")+'</a></p>';}
function doLogin(){
api("/auth/login",{method:"POST",body:JSON.stringify({email:document.getElementById("email").value,password:document.getElementById("password").value})}).then(function(d){if(d.data&&d.data.token){setToken(d.data.token);location.reload()}else{toast("Login failed","error")}}).catch(function(e){toast(e.message,"error")});}
function doRegister(){
api("/auth/register",{method:"POST",body:JSON.stringify({email:document.getElementById("email").value,password:document.getElementById("password").value})}).then(function(d){if(d.data&&d.data.token){setToken(d.data.token);location.reload()}else{toast("Registration failed","error")}}).catch(function(e){toast(e.message,"error")});}
function doLogout(){api("/auth/logout",{method:"POST"}).then(function(){clearToken();location.href="/"}).catch(function(){clearToken();location.href="/"})}

function renderDashboard(){
var adminLink=currentUser.isAdmin?'<a class="admin-link" href="/admin">'+t("adminPanel")+'</a>':'';
var langHtml='<div class="lang-toggle"><a class="'+(lang==="en"?"active":"")+'" onclick="setLang(\'en\')">EN</a> <a class="'+(lang==="zh"?"active":"")+'" onclick="setLang(\'zh\')">中文</a></div>';
document.getElementById("app").innerHTML='<nav class="navbar"><div class="brand"><div class="logo">Vibecast</div></div><div class="nav-right">'+adminLink+langHtml+'<span class="email">'+esc(currentUser.email)+'</span><button class="btn-link" onclick="doLogout()">'+t("logout")+'</button></div></nav><div class="container"><div class="card"><div class="card-header"><h2>'+t("create")+'</h2></div><div class="card-body"><div class="form-grid"><div class="form-field full"><label>'+t("siteName")+'</label><input id="site-name" placeholder="'+t("siteNamePh")+'"></div><div class="form-field"><label>'+t("slug")+'</label><input id="site-slug" placeholder="'+t("slugPh")+'"><div class="desc">'+t("slugDesc")+'</div></div><div class="form-field"><label>'+t("sitePwd")+'</label><input id="site-pwd" type="password" placeholder="'+t("sitePwdPh")+'"><div class="desc" id="pwd-desc">'+t("pwdDesc")+'</div></div></div><div class="form-actions"><button class="btn btn-primary" onclick="createSite()">'+t("create")+'</button></div></div></div><div class="card"><div class="card-header"><h2>'+t("yourSites")+'</h2><span class="hint">'+t("sitesHint")+'</span></div><div class="card-body"><div class="list-toolbar"><input type="text" id="site-search" placeholder="'+t("searchSitesPh")+'" onkeydown="if(event.key===\'Enter\')searchSites()" value="'+esc(siteSearch)+'"></div><div id="site-list"></div></div></div></div>';
loadSites();}

function searchSites(){siteSearch=document.getElementById("site-search").value;sitePage=1;loadSites()}
function sitePageGo(p){sitePage=p;loadSites()}
function loadSites(){
var q=siteSearch?"&q="+encodeURIComponent(siteSearch):"";
api("/sites?page="+sitePage+"&perPage="+sitePerPage+q).then(function(d){
var r=d.data||{};var sites=r.items||[];siteTotal=r.total||0;
var el=document.getElementById("site-list");
// Update password field description based on public access setting
var pwdDesc=document.getElementById("pwd-desc");
if(sites.length>0){
var pubDisabled=sites[0].publicAccessDisabled;
var pwdInput=document.getElementById("site-pwd");
if(pwdDesc){
pwdDesc.textContent=pubDisabled?t("pwdRequired"):t("pwdOptional");
pwdDesc.style.color=pubDisabled?"var(--red)":"var(--muted)";
}
if(pwdInput&&pubDisabled){
pwdInput.placeholder=t("pwdRequired");
}
}
var totalPages=Math.ceil(siteTotal/sitePerPage)||1;
var pgHtml=paginationHtml(sitePage,totalPages,"sitePageGo");
if(!sites.length){
el.innerHTML='<div class="empty">'+t("noSites")+'</div>'+pgHtml;
return;
}
var html='<ul class="site-list">';
for(var i=0;i<sites.length;i++){
var s=sites[i];
var badge='';
if(s.publicAccessDisabled&&!s.protected){
badge='<span class="badge badge-disabled">'+t("accessDisabled")+'</span>';
}else if(s.protected){
badge='<span class="badge badge-protected">'+t("protected")+'</span>';
}else{
badge='<span class="badge badge-public">'+t("public")+'</span>';
}
var pwdDisplay=s.protected?'<code style="background:#e2e8f0;padding:1px 5px;border-radius:3px;font-size:.8rem">'+esc(s.password)+'</code>':'<span style="color:var(--muted)">'+t("none")+'</span>';
html+='<li class="site-item"><div class="site-head" onclick="toggleDetail('+s.id+')"><div class="info" style="flex:1;min-width:0"><div class="name">'+esc(s.name)+' '+badge+'</div><div class="url"><a href="'+s.url+'" target="_blank" onclick="event.stopPropagation()">'+s.url+'</a></div></div><div class="actions"><label class="upload-btn" onclick="event.stopPropagation()">'+t("deployZip")+'<input type="file" accept=".zip" onchange="deploy('+s.id+',this.files[0])"></label><button class="btn btn-sm btn-danger" onclick="event.stopPropagation();delSite('+s.id+')">'+t("delete")+'</button></div></div><div class="site-detail" id="detail-'+s.id+'"><div class="detail-row"><span class="label">'+t("storagePath")+'</span><span class="value"><code>./data/sites/'+esc(s.slug)+'/</code></span></div><div class="detail-row"><span class="label">'+t("accessPassword")+'</span><span class="value">'+pwdDisplay+'</span></div></div></li>';
}
html+='</ul><div class="field-hint" style="margin-top:.5rem">'+t("deployHint")+'</div>';
html+=pgHtml;
el.innerHTML=html;
}).catch(function(e){toast(e.message,"error")})
}
function paginationHtml(page,totalPages,goFn){
if(totalPages<=1)return '';
var html='<div class="pagination">';
html+='<button '+(page<=1?'disabled':'')+' onclick="'+goFn+'('+(page-1)+')">'+t("prev")+'</button>';
var start=Math.max(1,page-2),end=Math.min(totalPages,page+2);
if(start>1){html+='<button onclick="'+goFn+'(1)">1</button>';if(start>2)html+='<span class="page-info">...</span>'}
for(var i=start;i<=end;i++){
html+='<button class="'+(i===page?'active':'')+'" onclick="'+goFn+'('+i+')">'+i+'</button>';
}
if(end<totalPages){if(end<totalPages-1)html+='<span class="page-info">...</span>';html+='<button onclick="'+goFn+'('+totalPages+')">'+totalPages+'</button>'}
html+='<button '+(page>=totalPages?'disabled':'')+' onclick="'+goFn+'('+(page+1)+')">'+t("next")+'</button>';
html+='<span class="page-info">'+t("page")+' '+page+' '+t("of")+' '+totalPages+'</span>';
html+='</div>';
return html;
}

function toggleDetail(id){
var el=document.getElementById("detail-"+id);
if(el)el.classList.toggle("show");}

function createSite(){
api("/sites",{method:"POST",body:JSON.stringify({name:document.getElementById("site-name").value,slug:document.getElementById("site-slug").value,password:document.getElementById("site-pwd").value})}).then(function(){
document.getElementById("site-name").value="";document.getElementById("site-slug").value="";document.getElementById("site-pwd").value="";
toast(t("siteCreated"));loadSites();
}).catch(function(e){toast(e.message,"error")});}

function deploy(id,file){
if(!file)return;
var fd=new FormData();
fd.append("file",file);
api("/sites/"+id+"/deploy",{method:"POST",body:fd,headers:{}}).then(function(){
toast(t("deployed"));
loadSites();
}).catch(function(e){toast(e.message,"error")});}

function delSite(id){
if(!confirm(t("deleteConfirm")))return;
api("/sites/"+id,{method:"DELETE"}).then(function(){toast(t("deleted"));loadSites()}).catch(function(e){toast(e.message,"error")});}

var savedLang="en";try{savedLang=localStorage.getItem("lang")||"en"}catch(e){}
lang=savedLang;
checkAuth().then(function(ok){if(ok)renderDashboard();else renderAuth()});
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
<style>
:root{--bg:#f1f5f9;--card:#fff;--border:#e2e8f0;--text:#1e293b;--muted:#64748b;--primary:#6366f1;--primary-hover:#5563d1;--primary-light:#eef2ff;--green:#059669;--green-light:#d1fae5;--red:#dc2626;--red-light:#fee2e2;--amber:#d97706;--amber-light:#fef3c7;--radius:10px;--shadow:0 1px 3px rgba(0,0,0,.06),0 1px 2px rgba(0,0,0,.04)}
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,"Helvetica Neue",sans-serif;background:var(--bg);color:var(--text);min-height:100vh}
a{color:var(--primary);text-decoration:none}
button{font-family:inherit;cursor:pointer}
.navbar{display:flex;justify-content:space-between;align-items:center;padding:0 1.5rem;height:56px;background:#fff;border-bottom:1px solid var(--border);position:sticky;top:0;z-index:100}
.navbar .brand .logo{font-size:1.2rem;font-weight:800;color:var(--amber)}
.navbar .nav-right{display:flex;align-items:center;gap:1rem}
.navbar .nav-right .btn-link{font-size:.85rem;cursor:pointer;color:var(--muted);background:none;border:none}
.navbar .nav-right .btn-link:hover{color:var(--text)}
.lang-toggle a{font-size:.8rem;cursor:pointer;padding:3px 8px;border-radius:5px}
.lang-toggle a.active{background:var(--primary-light);color:var(--primary)}
.container{max-width:960px;margin:1.5rem auto;padding:0 1.5rem}
.card{background:var(--card);border:1px solid var(--border);border-radius:var(--radius);box-shadow:var(--shadow);margin-bottom:1.5rem}
.card-header{padding:1rem 1.25rem;border-bottom:1px solid var(--border)}
.card-header h2{font-size:1rem;font-weight:600}
.card-body{padding:1.25rem}
.stats-grid{display:grid;grid-template-columns:repeat(3,1fr);gap:1rem}
.stat-card{background:#fff;border:1px solid var(--border);border-radius:8px;padding:1.25rem;text-align:center}
.stat-card .num{font-size:1.8rem;font-weight:700;color:var(--text)}
.stat-card .label{font-size:.75rem;color:var(--muted);margin-top:.2rem}
table{width:100%;border-collapse:collapse}
th,td{padding:.55rem .8rem;text-align:left;border-bottom:1px solid var(--border)}
th{color:var(--muted);font-size:.75rem;font-weight:600;text-transform:uppercase;letter-spacing:.03em}
td{font-size:.85rem}
.badge{display:inline-block;font-size:.7rem;padding:2px 7px;border-radius:4px;font-weight:500}
.badge-admin{background:var(--amber-light);color:var(--amber)}
.badge-user{background:#f1f5f9;color:var(--muted)}
.badge-protected{background:var(--amber-light);color:var(--amber)}
.badge-public{background:var(--green-light);color:var(--green)}
.badge-disabled{background:var(--red-light);color:var(--red)}
.btn{padding:7px 14px;border:none;border-radius:6px;font-size:.8rem;font-weight:600}
.btn-promote{background:var(--green);color:#fff}
.btn-demote{background:var(--amber);color:#fff}
.btn-danger{background:var(--red);color:#fff}
.btn:hover{opacity:.85}
.toggle-row{display:flex;align-items:center;justify-content:space-between;padding:.6rem 0;border-bottom:1px solid var(--border)}
.toggle-row:last-child{border-bottom:none}
.toggle-row .toggle-info{flex:1}
.toggle-row .toggle-label{font-size:.9rem;font-weight:500}
.toggle-row .toggle-desc{font-size:.75rem;color:var(--muted);margin-top:.15rem}
.toggle-switch{position:relative;width:40px;height:22px;background:#cbd5e1;border-radius:11px;cursor:pointer;transition:background .2s;flex-shrink:0}
.toggle-switch.on{background:var(--green)}
.toggle-switch::after{content:"";position:absolute;top:2px;left:2px;width:18px;height:18px;background:#fff;border-radius:50%;transition:transform .2s}
.toggle-switch.on::after{transform:translateX(18px)}
.empty{text-align:center;color:var(--muted);padding:1.5rem;font-size:.85rem}
.toast{position:fixed;bottom:1.5rem;right:1.5rem;padding:10px 20px;border-radius:8px;font-size:.9rem;z-index:999;opacity:0;transform:translateY(10px);transition:all .25s}
.toast.show{opacity:1;transform:translateY(0)}
.toast.success{background:var(--green);color:#fff}
.toast.error{background:var(--red);color:#fff}
.field-hint{font-size:.75rem;color:var(--muted);margin-top:.3rem}
.list-toolbar{display:flex;align-items:center;gap:.6rem;margin-bottom:.75rem;flex-wrap:wrap}
.list-toolbar input[type=text]{flex:1;min-width:120px;padding:7px 11px;border:1px solid var(--border);border-radius:7px;font-size:.85rem;outline:none;transition:border-color .15s;background:#fff;color:var(--text)}
.list-toolbar input[type=text]:focus{border-color:var(--primary);box-shadow:0 0 0 3px var(--primary-light)}
.pagination{display:flex;align-items:center;gap:.4rem;justify-content:center;padding:.75rem 0 0;flex-wrap:wrap}
.pagination .page-info{font-size:.78rem;color:var(--muted);margin:0 .4rem}
.pagination button{padding:4px 10px;border:1px solid var(--border);background:#fff;color:var(--text);border-radius:5px;font-size:.78rem;cursor:pointer;font-weight:500}
.pagination button:hover:not(:disabled){background:#f8fafc;border-color:#cbd5e1}
.pagination button.active{background:var(--primary);color:#fff;border-color:var(--primary)}
.pagination button:disabled{opacity:.4;cursor:default}
</style>
</head>
<body>
<div id="app"></div>
<script>
var API="/api";
var lang="en";
var userPage=1,userPerPage=10,userTotal=0,userSearch="";
var adminSitePage=1,adminSitePerPage=10,adminSiteTotal=0,adminSiteSearch="";
var i18n={en:{overview:"Overview",settings:"Settings",users:"Users",allSites:"All Sites",userCount:"Users",siteCount:"Sites",adminCount:"Admins",openReg:"Open Registration",openRegDesc:"When enabled, anyone can register a new account.",publicAccess:"Public Site Access",publicAccessDesc:"When disabled, all deployed sites return 403 to visitors.",dashboard:"Dashboard",logout:"Logout",role:"Role",created:"Created",actions:"Actions",promote:"Promote",demote:"Demote",delete:"Delete",deleteUserConfirm:"Delete this user and all their sites?",deleteSiteConfirm:"Delete this site?",name:"Name",slug:"Slug",owner:"Owner",protected:"Protected",public:"Public",password:"Password",storagePath:"Storage",url:"URL",noUsers:"No users",noSites:"No sites",updated:"Updated",deleted:"Deleted",roleUpdated:"Role updated",userDeleted:"User deleted",siteDeleted:"Site deleted",settingsSaved:"Settings saved",none:"None",accessDisabled:"Access Disabled",search:"Search",searchUsersPh:"Search users...",searchSitesPh:"Search sites...",prev:"Prev",next:"Next",page:"Page",of:"of",ownerEmail:"Owner Email"},zh:{overview:"概览",settings:"设置",users:"用户",allSites:"全部站点",userCount:"用户",siteCount:"站点",adminCount:"管理员",openReg:"开放注册",openRegDesc:"开启后，任何人都可以注册新账号。",publicAccess:"公开站点访问",publicAccessDesc:"关闭后，所有已部署站点对访问者返回 403。",dashboard:"Dashboard",logout:"退出",role:"角色",created:"创建时间",actions:"操作",promote:"提升",demote:"降级",delete:"删除",deleteUserConfirm:"删除此用户及其所有站点？",deleteSiteConfirm:"删除此站点？",name:"名称",slug:"Slug",owner:"所有者",protected:"已保护",public:"公开",password:"密码",storagePath:"存储路径",url:"URL",noUsers:"暂无用户",noSites:"暂无站点",updated:"已更新",deleted:"已删除",roleUpdated:"角色已更新",userDeleted:"用户已删除",siteDeleted:"站点已删除",settingsSaved:"设置已保存",none:"无",accessDisabled:"已禁止访问",search:"搜索",searchUsersPh:"搜索用户...",searchSitesPh:"搜索站点...",prev:"上一页",next:"下一页",page:"第",of:"/ 共",ownerEmail:"所有者邮箱"}};
function t(k){return(i18n[lang]||i18n.en)[k]||(i18n.en[k]||k)}
function setLang(l){lang=l;try{localStorage.setItem("lang",l)}catch(e){}renderAdmin();}
function esc(s){return String(s||"").replace(/[&<>"']/g,function(c){return{"&":"&amp;","<":"&lt;",">":"&gt;",'"':"&quot;","'":"&#39;"}[c]})}
function fmtDate(s){if(!s)return"-";var d=new Date(s);return d.toLocaleDateString()+" "+d.toLocaleTimeString([],{hour:"2-digit",minute:"2-digit"})}
function toast(msg,type){type=type||"success";var el=document.createElement("div");el.className="toast "+type+" show";el.textContent=msg;document.body.appendChild(el);setTimeout(function(){el.classList.remove("show");setTimeout(function(){el.remove()},300)},2500)}
function getToken(){try{return localStorage.getItem("vibecast_token")}catch(e){return""}}
function clearToken(){try{localStorage.removeItem("vibecast_token")}catch(e){}}
function api(path,opts){
  opts=opts||{};
  var token=localStorage.getItem("vibecast_token")||"";
  var headers={"Content-Type":"application/json"};
  if(token)headers["Authorization"]="Bearer "+token;
  if(opts.headers)Object.assign(headers,opts.headers);
  return fetch(API+path,Object.assign({},opts,{headers:headers,credentials:"same-origin"})).then(function(r){
    if(r.status===401){try{localStorage.removeItem("vibecast_token")}catch(e){}location.href="/dashboard";}
    return r.json().catch(function(){return{error:"network error"}}).then(function(data){if(!r.ok)throw new Error(data.error||"request failed");return data})
  })
}
function checkAuth(){if(!getToken())return Promise.reject(new Error("no token"));return api("/auth/me").then(function(d){if(!d.data||!d.data.isAdmin)throw new Error("not admin");return d.data}).catch(function(){clearToken();location.href="/dashboard"})}
function doLogout(){api("/auth/logout",{method:"POST"}).then(function(){clearToken();location.href="/"}).catch(function(){clearToken();location.href="/"})}

function renderAdmin(){
var langHtml='<div class="lang-toggle"><a class="'+(lang==="en"?"active":"")+'" onclick="setLang(\'en\')">EN</a> <a class="'+(lang==="zh"?"active":"")+'" onclick="setLang(\'zh\')">中文</a></div>';
document.getElementById("app").innerHTML='<nav class="navbar"><div class="brand"><div class="logo">Vibecast Admin</div></div><div class="nav-right">'+langHtml+'<a href="/dashboard" style="font-size:.85rem;color:var(--muted)">'+t("dashboard")+'</a><button class="btn-link" onclick="doLogout()">'+t("logout")+'</button></div></nav><div class="container"><div class="card"><div class="card-header"><h2>'+t("overview")+'</h2></div><div class="card-body"><div id="stats" class="stats-grid"></div></div></div><div class="card"><div class="card-header"><h2>'+t("settings")+'</h2></div><div class="card-body"><div id="settings"></div></div></div><div class="card"><div class="card-header"><h2>'+t("users")+'</h2></div><div class="card-body"><div class="list-toolbar"><input type="text" id="user-search" placeholder="'+t("searchUsersPh")+'" onkeydown="if(event.key===\'Enter\')searchUsers()" value="'+esc(userSearch)+'"></div><div id="users"></div></div></div><div class="card"><div class="card-header"><h2>'+t("allSites")+'</h2></div><div class="card-body"><div class="list-toolbar"><input type="text" id="admin-site-search" placeholder="'+t("searchSitesPh")+'" onkeydown="if(event.key===\'Enter\')searchAdminSites()" value="'+esc(adminSiteSearch)+'"></div><div id="sites"></div></div></div></div>';
loadStats();loadSettings();loadUsers();loadSites();}

function loadStats(){
api("/admin/stats").then(function(d){var s=d.data;document.getElementById("stats").innerHTML='<div class="stat-card"><div class="num">'+s.users+'</div><div class="label">'+t("userCount")+'</div></div><div class="stat-card"><div class="num">'+s.sites+'</div><div class="label">'+t("siteCount")+'</div></div><div class="stat-card"><div class="num">'+s.admins+'</div><div class="label">'+t("adminCount")+'</div></div>'}).catch(function(e){toast(e.message,"error")});}

function loadSettings(){
api("/admin/settings").then(function(d){var s=d.data;
var regOn=s.openRegistration;
var pubOn=s.allowPublicAccess!==false;
document.getElementById("settings").innerHTML='<div class="toggle-row"><div class="toggle-info"><div class="toggle-label">'+t("openReg")+'</div><div class="toggle-desc">'+t("openRegDesc")+'</div></div><div class="toggle-switch '+(regOn?"on":"")+'" onclick="toggleReg()"></div></div><div class="toggle-row"><div class="toggle-info"><div class="toggle-label">'+t("publicAccess")+'</div><div class="toggle-desc">'+t("publicAccessDesc")+'</div></div><div class="toggle-switch '+(pubOn?"on":"")+'" onclick="togglePub()"></div></div>';
}).catch(function(e){toast(e.message,"error")});}

function toggleReg(){
api("/admin/settings").then(function(d){var newVal=!d.data.openRegistration;return api("/admin/settings",{method:"PUT",body:JSON.stringify({openRegistration:newVal,allowPublicAccess:d.data.allowPublicAccess!==false})}).then(function(){toast(t("settingsSaved"));loadSettings()})}).catch(function(e){toast(e.message,"error")});}

function togglePub(){
api("/admin/settings").then(function(d){var cur=d.data;var pubOn=cur.allowPublicAccess!==false;var newVal=!pubOn;return api("/admin/settings",{method:"PUT",body:JSON.stringify({openRegistration:cur.openRegistration,allowPublicAccess:newVal})}).then(function(){toast(t("settingsSaved"));loadSettings()})}).catch(function(e){toast(e.message,"error")});}

function searchUsers(){userSearch=document.getElementById("user-search").value;userPage=1;loadUsers()}
function userPageGo(p){userPage=p;loadUsers()}
function paginationHtml(page,totalPages,goFn){
if(totalPages<=1)return '';
var html='<div class="pagination">';
html+='<button '+(page<=1?'disabled':'')+' onclick="'+goFn+'('+(page-1)+')">'+t("prev")+'</button>';
var start=Math.max(1,page-2),end=Math.min(totalPages,page+2);
if(start>1){html+='<button onclick="'+goFn+'(1)">1</button>';if(start>2)html+='<span class="page-info">...</span>'}
for(var i=start;i<=end;i++){
html+='<button class="'+(i===page?'active':'')+'" onclick="'+goFn+'('+i+')">'+i+'</button>';
}
if(end<totalPages){if(end<totalPages-1)html+='<span class="page-info">...</span>';html+='<button onclick="'+goFn+'('+totalPages+')">'+totalPages+'</button>'}
html+='<button '+(page>=totalPages?'disabled':'')+' onclick="'+goFn+'('+(page+1)+')">'+t("next")+'</button>';
html+='<span class="page-info">'+t("page")+' '+page+' '+t("of")+' '+totalPages+'</span>';
html+='</div>';
return html;
}
function loadUsers(){
var q=userSearch?"&q="+encodeURIComponent(userSearch):"";
api("/admin/users?page="+userPage+"&perPage="+userPerPage+q).then(function(d){
var r=d.data||{};var users=r.items||[];userTotal=r.total||0;
var el=document.getElementById("users");
var totalPages=Math.ceil(userTotal/userPerPage)||1;
var pgHtml=paginationHtml(userPage,totalPages,"userPageGo");
if(!users.length){el.innerHTML='<div class="empty">'+t("noUsers")+'</div>'+pgHtml;return}
var html='<table><thead><tr><th>ID</th><th>Email</th><th>'+t("role")+'</th><th>'+t("created")+'</th><th>'+t("actions")+'</th></tr></thead><tbody>';
for(var i=0;i<users.length;i++){
var u=users[i];
var badge=u.isAdmin?'<span class="badge badge-admin">Admin</span>':'<span class="badge badge-user">User</span>';
var btn=u.isAdmin?'<button class="btn btn-demote" onclick="toggleAdmin('+u.id+')">'+t("demote")+'</button>':'<button class="btn btn-promote" onclick="toggleAdmin('+u.id+')">'+t("promote")+'</button>';
html+='<tr><td>'+u.id+'</td><td>'+esc(u.email)+'</td><td>'+badge+'</td><td>'+fmtDate(u.createdAt)+'</td><td>'+btn+' <button class="btn btn-danger" onclick="delUser('+u.id+')">'+t("delete")+'</button></td></tr>';
}
html+='</tbody></table>';
html+=pgHtml;
el.innerHTML=html;
}).catch(function(e){toast(e.message,"error")})
}

function toggleAdmin(id){api("/admin/users/"+id,{method:"PUT"}).then(function(){toast(t("roleUpdated"));loadUsers();loadStats()}).catch(function(e){toast(e.message,"error")})}
function delUser(id){if(!confirm(t("deleteUserConfirm")))return;api("/admin/users/"+id,{method:"DELETE"}).then(function(){toast(t("userDeleted"));loadUsers();loadStats();loadSites()}).catch(function(e){toast(e.message,"error")})}

function searchAdminSites(){adminSiteSearch=document.getElementById("admin-site-search").value;adminSitePage=1;loadSites()}
function adminSitePageGo(p){adminSitePage=p;loadSites()}
function loadSites(){
var q=adminSiteSearch?"&q="+encodeURIComponent(adminSiteSearch):"";
api("/admin/sites?page="+adminSitePage+"&perPage="+adminSitePerPage+q).then(function(d){
var r=d.data||{};var sites=r.items||[];adminSiteTotal=r.total||0;
var el=document.getElementById("sites");
var totalPages=Math.ceil(adminSiteTotal/adminSitePerPage)||1;
var pgHtml=paginationHtml(adminSitePage,totalPages,"adminSitePageGo");
if(!sites.length){el.innerHTML='<div class="empty">'+t("noSites")+'</div>'+pgHtml;return}
var html='<table><thead><tr><th>ID</th><th>'+t("name")+'</th><th>'+t("slug")+'</th><th>'+t("owner")+'</th><th>'+t("protected")+'</th><th>'+t("password")+'</th><th>'+t("url")+'</th><th>'+t("actions")+'</th></tr></thead><tbody>';
for(var i=0;i<sites.length;i++){
var s=sites[i];
var badge='';
if(s.publicAccessDisabled&&!s.protected){badge='<span class="badge badge-disabled">'+t("accessDisabled")+'</span>'}
else if(s.protected){badge='<span class="badge badge-protected">'+t("protected")+'</span>'}
else{badge='<span class="badge badge-public">'+t("public")+'</span>'}
var pwd=s.protected?'<code style="font-size:.8rem">'+esc(s.password)+'</code>':'<span style="color:var(--muted)">'+t("none")+'</span>';
html+='<tr><td>'+s.id+'</td><td>'+esc(s.name)+'</td><td>'+esc(s.slug)+'</td><td>'+esc(s.ownerEmail||"-")+'</td><td>'+badge+'</td><td>'+pwd+'</td><td><a href="'+s.url+'" target="_blank">'+s.url+'</a></td><td><button class="btn btn-danger" onclick="delSite('+s.id+')">'+t("delete")+'</button></td></tr>';
}
html+='</tbody></table>';
html+=pgHtml;
el.innerHTML=html;
}).catch(function(e){toast(e.message,"error")})
}

function delSite(id){if(!confirm(t("deleteSiteConfirm")))return;api("/admin/sites/"+id,{method:"DELETE"}).then(function(){toast(t("siteDeleted"));loadSites();loadStats()}).catch(function(e){toast(e.message,"error")})}

var savedLang="en";try{savedLang=localStorage.getItem("lang")||"en"}catch(e){}
lang=savedLang;
checkAuth().then(function(){renderAdmin()});
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
<style>
:root{--bg:#f1f5f9;--card:#fff;--border:#e2e8f0;--text:#1e293b;--muted:#64748b;--primary:#6366f1;--primary-light:#eef2ff;--radius:10px;--shadow:0 1px 3px rgba(0,0,0,.06)}
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif;background:var(--bg);color:var(--text);min-height:100vh;display:flex;align-items:center;justify-content:center}
.card{background:var(--card);padding:2.5rem;border:1px solid var(--border);border-radius:14px;width:100%%;max-width:380px;box-shadow:var(--shadow)}
.card h1{font-size:1.3rem;margin-bottom:.3rem}
.card .site-name{color:var(--primary);font-weight:600}
.card p{color:var(--muted);margin-bottom:1.2rem;font-size:.85rem}
input[type=password]{width:100%%;padding:11px 14px;border:1px solid var(--border);border-radius:8px;font-size:.9rem;margin-bottom:.8rem;outline:none;transition:border-color .15s}
input[type=password]:focus{border-color:var(--primary);box-shadow:0 0 0 3px var(--primary-light)}
button{width:100%%;padding:11px;background:var(--primary);color:#fff;border:none;border-radius:8px;font-size:.9rem;font-weight:600;cursor:pointer}
button:hover{opacity:.9}
.err{color:#dc2626;margin-bottom:.8rem;font-size:.85rem}
</style>
</head>
<body>
<div class="card">
<h1>&#128274; <span class="site-name">%s</span></h1>
<p>This site is password-protected. Enter the password to continue.</p>
<div id="err" class="err" style="display:none"></div>
<form onsubmit="submitPassword(event)">
<input type="password" id="password" name="password" placeholder="Password" autofocus required>
<button type="submit">Enter Site</button>
</form>
</div>
<script>
function submitPassword(e){
e.preventDefault();
var pwd=document.getElementById("password").value;
var errEl=document.getElementById("err");
errEl.style.display="none";
fetch("/p/%s",{method:"POST",headers:{"Content-Type":"application/json"},body:JSON.stringify({password:pwd})}).then(function(r){return r.json()}).then(function(d){
if(d.error){errEl.textContent=d.error;errEl.style.display="block";return}
if(d.data&&d.data.token){window.location.href="/s/%s/?token="+d.data.token}
}).catch(function(){errEl.textContent="Network error";errEl.style.display="block"});
}
</script>
</body>
</html>`, escHTML(siteName), escHTML(siteName), slug, slug)
}

// passwordPageHTMLWithErr returns the password gate page with an error message.
func passwordPageHTMLWithErr(slug, siteName, errMsg string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>%s — Password Required</title>
<style>
:root{--bg:#f1f5f9;--card:#fff;--border:#e2e8f0;--text:#1e293b;--muted:#64748b;--primary:#6366f1;--primary-light:#eef2ff;--radius:10px;--shadow:0 1px 3px rgba(0,0,0,.06)}
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif;background:var(--bg);color:var(--text);min-height:100vh;display:flex;align-items:center;justify-content:center}
.card{background:var(--card);padding:2.5rem;border:1px solid var(--border);border-radius:14px;width:100%%;max-width:380px;box-shadow:var(--shadow)}
.card h1{font-size:1.3rem;margin-bottom:.3rem}
.card .site-name{color:var(--primary);font-weight:600}
.card p{color:var(--muted);margin-bottom:1.2rem;font-size:.85rem}
input[type=password]{width:100%%;padding:11px 14px;border:1px solid var(--border);border-radius:8px;font-size:.9rem;margin-bottom:.8rem;outline:none;transition:border-color .15s}
input[type=password]:focus{border-color:var(--primary);box-shadow:0 0 0 3px var(--primary-light)}
button{width:100%%;padding:11px;background:var(--primary);color:#fff;border:none;border-radius:8px;font-size:.9rem;font-weight:600;cursor:pointer}
button:hover{opacity:.9}
.err{color:#dc2626;margin-bottom:.8rem;font-size:.85rem}
</style>
</head>
<body>
<div class="card">
<h1>&#128274; <span class="site-name">%s</span></h1>
<p>This site is password-protected. Enter the password to continue.</p>
<div class="err">%s</div>
<form method="POST" action="/p/%s">
<input type="password" name="password" placeholder="Password" autofocus required>
<button type="submit">Enter Site</button>
</form>
</div>
</body>
</html>`, escHTML(siteName), escHTML(siteName), escHTML(errMsg), slug)
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
