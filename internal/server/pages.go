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
<title>Vibecast &mdash; Build with vibe. Cast instantly.</title>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif;background:#0f0f0f;color:#e0e0e0;min-height:100vh;display:flex;flex-direction:column;align-items:center;justify-content:center}
.hero{text-align:center;padding:2rem}
.hero h1{font-size:3.5rem;font-weight:800;background:linear-gradient(135deg,#a78bfa,#60a5fa,#34d399);-webkit-background-clip:text;-webkit-text-fill-color:transparent;margin-bottom:1rem}
.tagline{font-size:1.25rem;color:#888;margin-bottom:2.5rem}
.cta{display:inline-block;padding:14px 40px;background:linear-gradient(135deg,#6366f1,#8b5cf6);color:#fff;border:none;border-radius:12px;font-size:1.1rem;font-weight:600;text-decoration:none;cursor:pointer;transition:transform .2s,box-shadow .2s}
.cta:hover{transform:translateY(-2px);box-shadow:0 8px 30px rgba(99,102,241,.4)}
.features{display:flex;gap:2rem;margin-top:4rem;flex-wrap:wrap;justify-content:center}
.feature{text-align:center;max-width:200px}
.feature .icon{font-size:2rem;margin-bottom:.5rem}
.feature h3{font-size:1rem;color:#a78bfa;margin-bottom:.3rem}
.feature p{font-size:.85rem;color:#666;line-height:1.4}
footer{margin-top:4rem;color:#444;font-size:.8rem}
</style>
</head>
<body>
<div class="hero">
<h1>Vibecast</h1>
<p class="tagline">Build with vibe. Cast instantly.</p>
<a href="/dashboard" class="cta">Get Started</a>
<div class="features">
<div class="feature"><div class="icon">Z</div><h3>Instant Deploy</h3><p>Upload a ZIP, get a live URL in seconds.</p></div>
<div class="feature"><div class="icon">*</div><h3>Password Protect</h3><p>Keep your site private with password gating.</p></div>
<div class="feature"><div class="icon">0</div><h3>Zero Nginx</h3><p>Pure Go application server. No dependencies.</p></div>
<div class="feature"><div class="icon">S</div><h3>Self-Hosted</h3><p>Your data, your rules. SQLite + filesystem.</p></div>
</div>
</div>
<footer>Vibecast &mdash; Self-hosted static site hosting</footer>
</body>
</html>`

// dashboardHTML is the admin dashboard SPA.
var dashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Vibecast Dashboard</title>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif;background:#0f0f0f;color:#e0e0e0;min-height:100vh}
.navbar{display:flex;justify-content:space-between;align-items:center;padding:1rem 2rem;background:#161616;border-bottom:1px solid #222}
.navbar .logo{font-size:1.3rem;font-weight:700;background:linear-gradient(135deg,#a78bfa,#60a5fa);-webkit-background-clip:text;-webkit-text-fill-color:transparent}
.navbar .user{font-size:.9rem;color:#888}
.navbar .user button{margin-left:1rem;padding:6px 14px;background:#333;color:#e0e0e0;border:none;border-radius:6px;cursor:pointer;font-size:.85rem}
.navbar .user button:hover{background:#444}
.container{max-width:900px;margin:2rem auto;padding:0 1rem}
.section{background:#161616;border-radius:12px;padding:1.5rem;margin-bottom:1.5rem;border:1px solid #222}
.section h2{font-size:1.1rem;margin-bottom:1rem;color:#a78bfa}
.form-row{display:flex;gap:.75rem;margin-bottom:.75rem;flex-wrap:wrap}
.form-row input{flex:1;padding:10px 14px;background:#222;border:1px solid #333;border-radius:8px;color:#e0e0e0;font-size:.9rem;outline:none}
.form-row input:focus{border-color:#6366f1}
.btn{padding:10px 20px;background:linear-gradient(135deg,#6366f1,#8b5cf6);color:#fff;border:none;border-radius:8px;font-size:.9rem;font-weight:600;cursor:pointer}
.btn:hover{opacity:.9}
.btn-sm{padding:6px 14px;font-size:.8rem;border-radius:6px}
.btn-danger{background:#dc2626}
.btn-success{background:#059669}
.site-list{list-style:none}
.site-item{display:flex;align-items:center;justify-content:space-between;padding:1rem;background:#1a1a1a;border-radius:8px;margin-bottom:.5rem;border:1px solid #222}
.site-item .info{flex:1}
.site-item .info .name{font-weight:600;font-size:.95rem}
.site-item .info .url{font-size:.8rem;color:#6366f1;margin-top:.2rem}
.site-item .info .url a{color:#6366f1;text-decoration:none}
.site-item .info .badge{display:inline-block;font-size:.7rem;padding:2px 8px;border-radius:4px;margin-left:.5rem;background:#374151;color:#9ca3af}
.site-item .actions{display:flex;gap:.5rem;align-items:center}
.upload-label{display:inline-block;padding:6px 14px;background:#059669;color:#fff;border-radius:6px;font-size:.8rem;cursor:pointer;font-weight:500}
.upload-label:hover{opacity:.9}
.upload-input{display:none}
.empty{text-align:center;color:#555;padding:2rem;font-size:.9rem}
.toast{position:fixed;bottom:2rem;right:2rem;padding:12px 24px;border-radius:8px;font-size:.9rem;z-index:999;opacity:0;transition:opacity .3s}
.toast.show{opacity:1}
.toast.success{background:#059669;color:#fff}
.toast.error{background:#dc2626;color:#fff}
.auth-screen{display:flex;align-items:center;justify-content:center;min-height:80vh}
.auth-card{background:#161616;padding:2.5rem;border-radius:16px;width:100%;max-width:400px;border:1px solid #222}
.auth-card h1{font-size:1.6rem;margin-bottom:.5rem;text-align:center}
.auth-card .subtitle{color:#666;text-align:center;margin-bottom:1.5rem;font-size:.9rem}
.auth-card input{width:100%;padding:12px 16px;background:#222;border:1px solid #333;border-radius:8px;color:#e0e0e0;font-size:1rem;margin-bottom:.75rem;outline:none}
.auth-card input:focus{border-color:#6366f1}
.auth-card .btn{width:100%;margin-top:.5rem}
.auth-card .switch{text-align:center;margin-top:1rem;font-size:.85rem;color:#666}
.auth-card .switch a{color:#6366f1;text-decoration:none;cursor:pointer}
</style>
</head>
<body>
<div id="app"></div>
<script>
var API="/api";
var currentUser=null;

function api(path,opts){
  opts=opts||{};
  return fetch(API+path,Object.assign({},opts,{headers:{"Content-Type":"application/json"},credentials:"same-origin"})).then(function(r){
    return r.json().catch(function(){return{error:"network error"}}).then(function(data){
      if(!r.ok)throw new Error(data.error||"request failed");
      return data;
    });
  });
}
function toast(msg,type){
  type=type||"success";
  var el=document.createElement("div");
  el.className="toast "+type+" show";
  el.textContent=msg;
  document.body.appendChild(el);
  setTimeout(function(){el.classList.remove("show");setTimeout(function(){el.remove()},300)},2500);
}
function esc(s){
  return String(s||"").replace(/[&<>"']/g,function(c){
    return{"&":"&amp;","<":"&lt;",">":"&gt;",'"':"&quot;","'":"&#39;"}[c];
  });
}
function checkAuth(){
  return api("/auth/me").then(function(d){currentUser=d.data;return true}).catch(function(){return false});
}
function renderAuth(){
  document.getElementById("app").innerHTML='<div class="auth-screen"><div class="auth-card"><h1>Vibecast</h1><p class="subtitle">Build with vibe. Cast instantly.</p><div id="auth-form"></div></div></div>';
  showLogin();
}
function showLogin(){
  document.getElementById("auth-form").innerHTML='<input id="email" type="email" placeholder="Email" autocomplete="email"><input id="password" type="password" placeholder="Password" autocomplete="current-password"><button class="btn" onclick="doLogin()">Login</button><p class="switch">No account? <a onclick="showRegister()">Register</a></p>';
}
function showRegister(){
  document.getElementById("auth-form").innerHTML='<input id="email" type="email" placeholder="Email" autocomplete="email"><input id="password" type="password" placeholder="Password (min 6 chars)" autocomplete="new-password"><button class="btn" onclick="doRegister()">Register</button><p class="switch">Have an account? <a onclick="showLogin()">Login</a></p>';
}
function doLogin(){
  api("/auth/login",{method:"POST",body:JSON.stringify({email:email.value,password:password.value})}).then(function(){location.reload()}).catch(function(e){toast(e.message,"error")});
}
function doRegister(){
  api("/auth/register",{method:"POST",body:JSON.stringify({email:email.value,password:password.value})}).then(function(){location.reload()}).catch(function(e){toast(e.message,"error")});
}
function doLogout(){
  api("/auth/logout",{method:"POST"}).then(function(){location.reload()});
}
function renderDashboard(){
  document.getElementById("app").innerHTML='<nav class="navbar"><div class="logo">Vibecast</div><div class="user">'+esc(currentUser.email)+' <button onclick="doLogout()">Logout</button></div></nav><div class="container"><div class="section"><h2>Create New Site</h2><div class="form-row"><input id="site-name" placeholder="Site name (e.g. My Portfolio)"><input id="site-slug" placeholder="custom-slug (optional)"><input id="site-pwd" type="password" placeholder="Access password (optional)"><button class="btn" onclick="createSite()">Create</button></div></div><div class="section"><h2>Your Sites</h2><div id="site-list"></div></div></div>';
  loadSites();
}
function loadSites(){
  api("/sites").then(function(d){
    var sites=d.data||[];
    var el=document.getElementById("site-list");
    if(!sites.length){el.innerHTML='<div class="empty">No sites yet. Create one above.</div>';return}
    var html='<ul class="site-list">';
    for(var i=0;i<sites.length;i++){
      var s=sites[i];
      var badge=s.protected?'<span class="badge">Protected</span>':'';
      html+='<li class="site-item"><div class="info"><span class="name">'+esc(s.name)+badge+'</span><div class="url"><a href="'+s.url+'" target="_blank">'+s.url+'</a></div></div><div class="actions"><label class="upload-label">Deploy ZIP<input type="file" accept=".zip" class="upload-input" onchange="deploy('+s.id+',this.files[0])"></label><button class="btn btn-sm btn-danger" onclick="delSite('+s.id+',\''+esc(s.name)+'\')">Delete</button></div></li>';
    }
    html+='</ul>';
    el.innerHTML=html;
  }).catch(function(e){toast(e.message,"error")});
}
function createSite(){
  api("/sites",{method:"POST",body:JSON.stringify({name:siteName.value,slug:siteSlug.value,password:sitePwd.value})}).then(function(){
    siteName.value="";siteSlug.value="";sitePwd.value="";
    toast("Site created");
    loadSites();
  }).catch(function(e){toast(e.message,"error")});
}
function deploy(id,file){
  if(!file)return;
  var fd=new FormData();
  fd.append("file",file);
  fetch("/api/sites/"+id+"/deploy",{method:"POST",body:fd,credentials:"same-origin"}).then(function(r){
    return r.json().catch(function(){return{error:"network error"}});
  }).then(function(d){
    if(d.error)throw new Error(d.error);
    toast("Deployed! Live at "+d.data.url);
    loadSites();
  }).catch(function(e){toast(e.message,"error")});
}
function delSite(id,name){
  if(!confirm("Delete \""+name+"\"? This removes all files."))return;
  api("/sites/"+id,{method:"DELETE"}).then(function(){toast("Deleted");loadSites()}).catch(function(e){toast(e.message,"error")});
}
checkAuth().then(function(ok){
  if(ok)renderDashboard();
  else renderAuth();
});
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
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif;background:#0f0f0f;color:#e0e0e0;min-height:100vh;display:flex;align-items:center;justify-content:center}
.card{background:#1a1a1a;padding:2.5rem;border-radius:16px;width:100%%;max-width:400px;text-align:center;border:1px solid #222}
.card h1{font-size:1.5rem;margin-bottom:.5rem}
.card .site-name{color:#a78bfa;font-weight:600}
.card p{color:#666;margin-bottom:1.5rem;font-size:.9rem}
input[type=password]{width:100%%;padding:12px 16px;background:#222;border:1px solid #333;border-radius:8px;color:#e0e0e0;font-size:1rem;margin-bottom:1rem;outline:none}
input[type=password]:focus{border-color:#6366f1}
button{width:100%%;padding:12px;background:linear-gradient(135deg,#6366f1,#8b5cf6);color:#fff;border:none;border-radius:8px;font-size:1rem;font-weight:600;cursor:pointer}
button:hover{opacity:.9}
.err{color:#f87171;margin-bottom:1rem;font-size:.85rem}
</style>
</head>
<body>
<div class="card">
<h1>Lock <span class="site-name">%s</span></h1>
<p>This site is password-protected. Enter the password to continue.</p>
<form method="POST" action="/p/%s">
<input type="password" name="password" placeholder="Password" autofocus required>
<button type="submit">Enter Site</button>
</form>
</div>
</body>
</html>`, escHTML(siteName), escHTML(siteName), slug)
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
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif;background:#0f0f0f;color:#e0e0e0;min-height:100vh;display:flex;align-items:center;justify-content:center}
.card{background:#1a1a1a;padding:2.5rem;border-radius:16px;width:100%%;max-width:400px;text-align:center;border:1px solid #222}
.card h1{font-size:1.5rem;margin-bottom:.5rem}
.card .site-name{color:#a78bfa;font-weight:600}
.card p{color:#666;margin-bottom:1.5rem;font-size:.9rem}
input[type=password]{width:100%%;padding:12px 16px;background:#222;border:1px solid #333;border-radius:8px;color:#e0e0e0;font-size:1rem;margin-bottom:1rem;outline:none}
input[type=password]:focus{border-color:#6366f1}
button{width:100%%;padding:12px;background:linear-gradient(135deg,#6366f1,#8b5cf6);color:#fff;border:none;border-radius:8px;font-size:1rem;font-weight:600;cursor:pointer}
button:hover{opacity:.9}
.err{color:#f87171;margin-bottom:1rem;font-size:.85rem}
</style>
</head>
<body>
<div class="card">
<h1>Lock <span class="site-name">%s</span></h1>
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
