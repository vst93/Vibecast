package server

// This file is the single source of truth for Vibecast's UI. Every page
// (landing, auth/dashboard, admin, password gate) composes stylesheets from
// the layers below instead of re-declaring tokens and components:
//
//	uiTokensCSS   – palette, type scale, radii, motion, reset, a11y baseline
//	uiComponentsCSS – chrome shared by dashboard and admin
//	uiDashboardCSS / uiAdminCSS / uiLandingCSS – page-specific layout
//
// Rules of the road:
//   - Typography: system UI stack only (Chinese-capable), mono for machine data.
//   - Never italic, never `transition: all`, never animate layout properties.
//   - Numbers get tabular figures; interactive rows expose focus-visible rings.

const uiTokensCSS = `
:root{
  color-scheme:light;
  --bg:#f2f5f3;
  --surface:#ffffff;
  --surface-2:#eef2f0;
  --surface-3:#e3eae6;
  --line:#d5ded9;
  --line-strong:#bccbc3;
  --text:#16201b;
  --text-soft:#31413a;
  --dim:#596b62;
  --dim-2:#61746c;
  --accent:#0a7d59;
  --accent-hover:#086549;
  --accent-soft:#dbefe6;
  --accent-ink:#076347;
  --accent-ring:rgba(10,125,89,.22);
  --on-accent:#ffffff;
  --warn:#8a5a00;
  --warn-soft:#f7ecd6;
  --danger:#b3242c;
  --danger-hover:#8e1c22;
  --danger-soft:#fbe7e8;
  --shadow-sm:0 1px 2px rgba(20,42,32,.06);
  --shadow-md:0 12px 32px rgba(20,42,32,.09),0 1px 2px rgba(20,42,32,.06);
  --shadow-lg:0 28px 72px rgba(16,38,28,.16);
  --field:#ffffff;
  /* Concrete families first, generic last. Do not lead with system-ui/ui-sans-serif:
     fontconfig can resolve those to a serif face at intermediate weights on Linux,
     which silently turns the whole UI serif. Order also covers CJK on macOS/Windows. */
  --ff-sans:-apple-system,BlinkMacSystemFont,"Segoe UI Variable Text","Segoe UI",Roboto,"Helvetica Neue","Noto Sans",Arial,"PingFang SC","Hiragino Sans GB","Microsoft YaHei UI","Microsoft YaHei","Noto Sans SC","Liberation Sans","DejaVu Sans",sans-serif;
  --ff-mono:"SF Mono",SFMono-Regular,"JetBrains Mono","Cascadia Mono",Menlo,Consolas,"Liberation Mono","DejaVu Sans Mono",monospace;
  --fs-xs:.75rem;--fs-sm:.8125rem;--fs-base:.875rem;--fs-md:.9375rem;--fs-lg:1.0625rem;--fs-xl:1.375rem;
  --r-sm:6px;--r-md:8px;--r-lg:12px;--r-pill:999px;
  --ease:cubic-bezier(.22,.61,.36,1);
  --dur:150ms;
  --nav-h:52px;
}
[data-theme="dark"]{
  color-scheme:dark;
  --bg:#0b1210;
  --surface:#121b17;
  --surface-2:#182521;
  --surface-3:#20302a;
  --line:#2b3d35;
  --line-strong:#3b5148;
  --text:#e8f0eb;
  --text-soft:#d2ded7;
  --dim:#9db1a7;
  --dim-2:#849790;
  --accent:#4bd6a1;
  --accent-hover:#6ee6b8;
  --accent-soft:#16342a;
  --accent-ink:#4bd6a1;
  --accent-ring:rgba(75,214,161,.26);
  --on-accent:#06211a;
  --warn:#e6bd6a;
  --warn-soft:#2f2617;
  --danger:#ff8a91;
  --danger-hover:#ffaab0;
  --danger-soft:#331b1e;
  --shadow-sm:0 1px 2px rgba(0,0,0,.4);
  --shadow-md:0 14px 36px rgba(0,0,0,.42),0 1px 2px rgba(0,0,0,.4);
  --shadow-lg:0 30px 80px rgba(0,0,0,.55);
  --field:#0b1512;
}
*,*::before,*::after{box-sizing:border-box;margin:0;padding:0}
html{-webkit-text-size-adjust:100%}
body{
  min-height:100dvh;
  font-family:var(--ff-sans);
  font-size:var(--fs-base);
  line-height:1.55;
  color:var(--text);
  background:var(--bg);
  -webkit-font-smoothing:antialiased;
  text-rendering:optimizeLegibility;
}
a{color:inherit;text-decoration:none}
button,input,select,textarea{font:inherit;color:inherit}
button{cursor:pointer;background:none;border:0}
h1,h2,h3,h4{font-weight:600;line-height:1.25;letter-spacing:0;text-wrap:balance}
p{text-wrap:pretty}
code,kbd,samp,.tnum{font-variant-numeric:tabular-nums}
svg{display:block}
[hidden]{display:none!important}
::selection{background:color-mix(in srgb,var(--accent) 26%,transparent)}
:where(a,button,input,select,textarea,summary,[tabindex]):focus-visible{
  outline:2px solid var(--accent);
  outline-offset:2px;
  border-radius:var(--r-sm);
}
:where(button,input,select,textarea,.site-head,summary){-webkit-tap-highlight-color:transparent}
::-webkit-scrollbar{width:11px;height:11px}
::-webkit-scrollbar-track{background:transparent}
::-webkit-scrollbar-thumb{background:var(--line-strong);border:3px solid var(--bg);border-radius:var(--r-pill)}
::-webkit-scrollbar-thumb:hover{background:var(--dim-2)}
@media(prefers-reduced-motion:reduce){
  *,*::before,*::after{animation-duration:.001ms!important;animation-iteration-count:1!important;transition-duration:.001ms!important;scroll-behavior:auto!important}
}
`

const uiComponentsCSS = `
/* ── navbar ─────────────────────────────────────────────────────────── */
.navbar{
  position:sticky;top:0;z-index:60;
  display:flex;align-items:center;justify-content:space-between;gap:14px;
  min-height:var(--nav-h);padding:0 clamp(12px,2.2vw,22px);
  background:color-mix(in srgb,var(--bg) 88%,transparent);
  border-bottom:1px solid var(--line);
  backdrop-filter:blur(14px) saturate(1.25);
}
.navbar .logo{display:inline-flex;align-items:center;gap:8px;font-size:.95rem;font-weight:600;color:var(--text);letter-spacing:-.01em;min-width:0}
.navbar .logo .accent{color:var(--accent)}
.navbar .logo svg{flex-shrink:0}
.navbar .nav-right{display:flex;align-items:center;gap:6px;min-width:0}
.nav-section{display:flex;align-items:center;gap:4px;min-width:0}
.nav-preferences{padding:0 8px;border-left:1px solid var(--line);border-right:1px solid var(--line)}
.navbar .email{
  display:inline-flex;align-items:center;max-width:200px;height:30px;padding:0 10px;
  font-size:var(--fs-xs);color:var(--dim);
  background:var(--surface-2);border:1px solid var(--line);border-radius:var(--r-pill);
  overflow:hidden;text-overflow:ellipsis;white-space:nowrap;
}
.btn-link,.btn-icon{
  display:inline-flex;align-items:center;gap:6px;
  min-height:32px;padding:0 9px;border-radius:var(--r-md);
  font-size:var(--fs-sm);font-weight:500;color:var(--dim);
  border:1px solid transparent;
  transition:background var(--dur) var(--ease),color var(--dur) var(--ease),border-color var(--dur) var(--ease);
}
.btn-link svg,.btn-icon svg{flex-shrink:0}
@media(hover:hover) and (pointer:fine){
  .btn-link:hover,.btn-icon:hover{background:var(--surface-2);color:var(--text)}
}
.btn-link:active,.btn-icon:active{transform:scale(.98)}
.admin-link{
  display:inline-flex;align-items:center;gap:6px;min-height:32px;padding:0 10px;
  border-radius:var(--r-md);font-size:var(--fs-sm);font-weight:550;
  color:var(--warn);background:color-mix(in srgb,var(--warn) 10%,var(--surface));
  border:1px solid color-mix(in srgb,var(--warn) 26%,var(--line));
  transition:background var(--dur) var(--ease),color var(--dur) var(--ease);
}
@media(hover:hover) and (pointer:fine){
  .admin-link:hover{background:color-mix(in srgb,var(--warn) 18%,var(--surface));color:var(--text)}
}
.lang-toggle{display:inline-flex;align-items:center;gap:2px;padding:2px;background:var(--surface-2);border:1px solid var(--line);border-radius:var(--r-md)}
.lang-toggle a,.lang-toggle button{display:inline-flex;align-items:center;justify-content:center;min-width:30px;height:26px;padding:0 7px;border-radius:5px;font-family:var(--ff-mono);font-size:.6875rem;color:var(--dim);cursor:pointer;transition:background var(--dur) var(--ease),color var(--dur) var(--ease)}
@media(hover:hover) and (pointer:fine){.lang-toggle a:hover,.lang-toggle button:hover{color:var(--text);background:color-mix(in srgb,var(--surface) 70%,transparent)}}
.lang-toggle a.active,.lang-toggle button.active{background:var(--surface);color:var(--accent);box-shadow:var(--shadow-sm)}
.mobile-menu-toggle{display:none}

/* ── layout ─────────────────────────────────────────────────────────── */
.container{max-width:1320px;margin:0 auto;padding:clamp(16px,2.4vw,28px) clamp(12px,2.4vw,24px) 64px}

/* ── card ───────────────────────────────────────────────────────────── */
.card{background:var(--surface);border:1px solid var(--line);border-radius:var(--r-lg);box-shadow:var(--shadow-sm);margin-bottom:16px}
.card-header{display:flex;align-items:center;justify-content:space-between;gap:12px;min-height:52px;padding:12px 16px;border-bottom:1px solid var(--line)}
.card-header h2{font-size:var(--fs-base);font-weight:600;color:var(--text)}
.card-header .hint{font-size:var(--fs-xs);color:var(--dim)}
.card-body{padding:16px}

/* ── buttons ────────────────────────────────────────────────────────── */
.btn{
  display:inline-flex;align-items:center;justify-content:center;gap:6px;
  min-height:36px;padding:0 13px;border-radius:var(--r-md);
  font-size:var(--fs-sm);font-weight:550;line-height:1;
  border:1px solid transparent;white-space:nowrap;
  transition:background var(--dur) var(--ease),color var(--dur) var(--ease),border-color var(--dur) var(--ease),transform var(--dur) var(--ease);
}
.btn:active{transform:scale(.985)}
.btn svg{flex-shrink:0;width:14px;height:14px}
.btn:disabled,.btn[disabled]{opacity:.5;cursor:not-allowed;transform:none}
.btn-primary{background:var(--accent);color:var(--on-accent)}
@media(hover:hover) and (pointer:fine){.btn-primary:hover:not(:disabled){background:var(--accent-hover)}}
.btn-outline{background:var(--surface);border-color:var(--line);color:var(--text-soft)}
@media(hover:hover) and (pointer:fine){.btn-outline:hover:not(:disabled){border-color:var(--line-strong);background:var(--surface-2);color:var(--text)}}
.btn-ghost{color:var(--dim);border-color:transparent}
@media(hover:hover) and (pointer:fine){.btn-ghost:hover:not(:disabled){background:var(--surface-2);color:var(--text)}}
.btn-danger{background:var(--danger);color:#fff}
@media(hover:hover) and (pointer:fine){.btn-danger:hover:not(:disabled){background:var(--danger-hover)}}
.btn-danger-ghost{color:var(--danger);border-color:transparent}
@media(hover:hover) and (pointer:fine){.btn-danger-ghost:hover:not(:disabled){background:var(--danger-soft);color:var(--danger)}}
.btn-sm{min-height:30px;padding:0 10px;font-size:var(--fs-xs);border-radius:var(--r-sm)}
.btn-sm svg{width:13px;height:13px}
.btn-block{width:100%}
.btn-promote{background:var(--accent);color:var(--on-accent)}
.btn-demote{background:var(--surface);border-color:var(--line);color:var(--text-soft)}
.icon-btn{
  display:inline-flex;align-items:center;justify-content:center;width:30px;height:30px;
  border-radius:var(--r-sm);color:var(--dim);border:1px solid transparent;
  transition:background var(--dur) var(--ease),color var(--dur) var(--ease),border-color var(--dur) var(--ease);
}
.icon-btn svg{width:15px;height:15px}
@media(hover:hover) and (pointer:fine){.icon-btn:hover{background:var(--surface-2);color:var(--text);border-color:var(--line)}}

/* ── forms ──────────────────────────────────────────────────────────── */
.form-field,.auth-field{display:flex;flex-direction:column;gap:5px;margin-bottom:14px}
.form-field label,.auth-field label,.modal-field label,.captcha-label{font-size:var(--fs-sm);font-weight:550;color:var(--text-soft)}
.form-field input,.auth-field input,.modal-field input,.captcha-row input,.list-toolbar input[type=text],.org-panel-input,.num-control input,textarea{
  width:100%;min-height:38px;padding:0 11px;
  background:var(--field);color:var(--text);
  border:1px solid var(--line);border-radius:var(--r-md);
  font-size:var(--fs-base);
  transition:border-color var(--dur) var(--ease),box-shadow var(--dur) var(--ease);
}
textarea{padding:9px 11px;line-height:1.55;resize:vertical}
.form-field input::placeholder,.auth-field input::placeholder,.modal-field input::placeholder,.captcha-row input::placeholder,.list-toolbar input[type=text]::placeholder,.org-panel-input::placeholder{color:var(--dim-2)}
.form-field input:focus,.auth-field input:focus,.modal-field input:focus,.captcha-row input:focus,.list-toolbar input[type=text]:focus,.org-panel-input:focus,.num-control input:focus,textarea:focus{
  outline:none;border-color:var(--accent);box-shadow:0 0 0 3px var(--accent-ring);
}
.form-field .desc,.auth-field .field-hint{font-size:var(--fs-xs);color:var(--dim);line-height:1.5}
.form-actions{display:flex;justify-content:flex-end;gap:8px;margin-top:6px}
.pwd-wrap{position:relative;display:flex;align-items:center}
.pwd-wrap input{padding-right:40px}
.pwd-toggle{
  position:absolute;right:5px;top:50%;transform:translateY(-50%);
  display:flex;align-items:center;justify-content:center;width:28px;height:28px;
  border-radius:var(--r-sm);color:var(--dim);border:1px solid transparent;
  transition:background var(--dur) var(--ease),color var(--dur) var(--ease);
}
.pwd-toggle svg{width:15px;height:15px}
@media(hover:hover) and (pointer:fine){.pwd-toggle:hover{background:var(--surface-2);color:var(--text)}}
input[type=password]::-ms-reveal,input[type=password]::-ms-clear{display:none}
/* switch: label wraps a checkbox, knob drawn by CSS */
.switch-field{display:flex;align-items:flex-start;justify-content:space-between;gap:14px;padding:12px 0;border-top:1px solid var(--line)}
.switch-field .switch-copy{display:block;min-width:0}
.switch-field .switch-title{display:block;font-size:var(--fs-sm);font-weight:550;color:var(--text);line-height:1.4}
.switch-field .desc{display:block;margin-top:2px;font-size:var(--fs-xs);color:var(--dim);line-height:1.5}
.switch-toggle{position:relative;flex:0 0 36px;width:36px;height:20px;margin-top:2px}
.switch-toggle input{position:absolute;inset:0;width:100%;height:100%;margin:0;opacity:0;cursor:pointer;z-index:2}
.switch-track{
  display:block;width:36px;height:20px;border-radius:var(--r-pill);
  background:var(--surface-3);border:1px solid var(--line-strong);
  transition:background var(--dur) var(--ease),border-color var(--dur) var(--ease);
  pointer-events:none;
}
.switch-track::after{
  content:"";position:absolute;top:3px;left:3px;width:14px;height:14px;border-radius:50%;
  background:var(--surface);box-shadow:var(--shadow-sm);
  transition:transform var(--dur) var(--ease);
}
.switch-toggle input:checked + .switch-track{background:var(--accent);border-color:var(--accent)}
.switch-toggle input:checked + .switch-track::after{transform:translateX(16px)}
.switch-toggle input:focus-visible + .switch-track{outline:2px solid var(--accent);outline-offset:2px}
.switch-toggle input:disabled{cursor:not-allowed}
.switch-toggle input:disabled + .switch-track{opacity:.45}

/* ── badges ─────────────────────────────────────────────────────────── */
.badge{
  display:inline-flex;align-items:center;gap:4px;height:20px;padding:0 7px;
  border-radius:var(--r-pill);font-size:.6875rem;font-weight:550;letter-spacing:.01em;white-space:nowrap;
  background:var(--surface-2);color:var(--dim);border:1px solid var(--line);
}
.badge-public{background:color-mix(in srgb,var(--accent) 10%,var(--surface));color:var(--accent-ink);border-color:color-mix(in srgb,var(--accent) 26%,var(--line))}
.badge-protected{background:var(--warn-soft);color:var(--warn);border-color:color-mix(in srgb,var(--warn) 30%,var(--line))}
.badge-disabled{background:var(--danger-soft);color:var(--danger);border-color:color-mix(in srgb,var(--danger) 28%,var(--line))}
.badge-org{background:color-mix(in srgb,#3b82f6 12%,var(--surface));color:#2563c9;border-color:color-mix(in srgb,#3b82f6 26%,var(--line))}
[data-theme="dark"] .badge-org{color:#8ab4ff}
.badge-pinned{background:color-mix(in srgb,#6366f1 12%,var(--surface));color:#4f46e5;border-color:color-mix(in srgb,#6366f1 26%,var(--line))}
[data-theme="dark"] .badge-pinned{color:#a5b4fc}
.badge-admin{background:var(--warn-soft);color:var(--warn);border-color:color-mix(in srgb,var(--warn) 30%,var(--line))}
.badge-user{background:var(--surface-2);color:var(--dim);border-color:var(--line)}
.status-dot{display:inline-block;width:7px;height:7px;border-radius:50%;flex-shrink:0}
.status-dot.public{background:var(--accent);box-shadow:0 0 0 3px color-mix(in srgb,var(--accent) 16%,transparent)}
.status-dot.protected{background:var(--warn);box-shadow:0 0 0 3px color-mix(in srgb,var(--warn) 16%,transparent)}
.status-dot.disabled{background:var(--danger);box-shadow:0 0 0 3px color-mix(in srgb,var(--danger) 16%,transparent)}

/* ── list toolbar / pagination ──────────────────────────────────────── */
.list-toolbar{display:flex;align-items:center;gap:8px;margin-bottom:12px}
.list-toolbar input[type=text]{flex:1;min-width:0;min-height:36px;font-size:var(--fs-sm)}
.pagination{display:flex;align-items:center;justify-content:center;gap:4px;flex-wrap:wrap;margin-top:4px;padding-top:12px}
.pagination button{
  min-width:30px;height:30px;padding:0 8px;border-radius:var(--r-sm);
  border:1px solid var(--line);background:var(--surface);color:var(--text-soft);
  font-size:var(--fs-xs);font-family:var(--ff-mono);font-variant-numeric:tabular-nums;
  transition:background var(--dur) var(--ease),border-color var(--dur) var(--ease),color var(--dur) var(--ease);
}
@media(hover:hover) and (pointer:fine){.pagination button:hover:not(:disabled):not(.active){background:var(--surface-2);border-color:var(--line-strong)}}
.pagination button.active{background:var(--accent);border-color:var(--accent);color:var(--on-accent);font-weight:600}
.pagination button:disabled{opacity:.4;cursor:default}
.pagination .page-info{font-size:var(--fs-xs);color:var(--dim);font-family:var(--ff-mono);padding:0 6px;font-variant-numeric:tabular-nums}

/* ── file list (shared by the site detail and the admin site detail) ── */
.file-list{border:1px solid var(--line);border-radius:var(--r-md);overflow:hidden;background:var(--surface)}
.file-row{display:grid;grid-template-columns:16px minmax(0,1fr) auto 22px;align-items:center;gap:8px;min-height:32px;padding:4px 6px 4px 10px;font-family:var(--ff-mono);font-size:.75rem;border-bottom:1px solid var(--line)}
.file-row:last-child{border-bottom:0}
@media(hover:hover) and (pointer:fine){.file-row:hover{background:var(--surface-2)}}
.file-icon{display:flex;align-items:center;color:var(--dim)}
.file-icon svg{width:13px;height:13px}
.file-name{overflow:hidden;text-overflow:ellipsis;white-space:nowrap;color:var(--text-soft)}
.file-dir .file-name{color:var(--text);font-weight:600}
.file-size{color:var(--dim);font-size:.6875rem;font-variant-numeric:tabular-nums;white-space:nowrap}
.file-del{display:inline-flex;align-items:center;justify-content:center;width:22px;height:22px;border-radius:var(--r-sm);border:1px solid transparent;background:none;color:var(--dim-2);opacity:0;cursor:pointer;transition:opacity var(--dur) var(--ease),color var(--dur) var(--ease),background var(--dur) var(--ease),border-color var(--dur) var(--ease)}
.file-del svg{width:12px;height:12px}
.file-row:hover .file-del,.file-row:focus-within .file-del{opacity:1}
.file-del:focus-visible{opacity:1;outline:2px solid var(--accent);outline-offset:1px}
@media(hover:hover) and (pointer:fine){.file-del:hover{color:var(--danger);background:var(--danger-soft);border-color:color-mix(in srgb,var(--danger) 26%,transparent)}}
@media(hover:none){.file-del{opacity:.55;width:26px;height:26px}}
.file-empty{display:flex;align-items:center;gap:8px;padding:14px;border:1px dashed var(--line-strong);border-radius:var(--r-md);font-size:var(--fs-sm);color:var(--dim);background:var(--surface)}

/* ── table ──────────────────────────────────────────────────────────── */
.table-wrap{overflow-x:auto;border:1px solid var(--line);border-radius:var(--r-md)}
table{width:100%;border-collapse:collapse;font-size:var(--fs-sm)}
thead th{
  padding:9px 12px;text-align:left;white-space:nowrap;
  font-size:.6875rem;font-weight:600;letter-spacing:.05em;text-transform:uppercase;
  color:var(--dim);background:var(--surface-2);border-bottom:1px solid var(--line);
}
tbody td{padding:10px 12px;border-bottom:1px solid var(--line);vertical-align:middle}
tbody tr:last-child td{border-bottom:0}
tbody tr{transition:background var(--dur) var(--ease)}
@media(hover:hover) and (pointer:fine){tbody tr:hover{background:color-mix(in srgb,var(--surface-2) 60%,var(--surface))}}
td.num,th.num{font-family:var(--ff-mono);font-variant-numeric:tabular-nums;color:var(--dim);font-size:var(--fs-xs);white-space:nowrap}
td.cell-actions{text-align:right;white-space:nowrap}
td.cell-actions .btn+.btn{margin-left:4px}
td.cell-email{max-width:280px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}

/* ── empty / loading ────────────────────────────────────────────────── */
.empty{
  display:flex;flex-direction:column;align-items:center;justify-content:center;gap:6px;
  min-height:190px;padding:28px 20px;text-align:center;
  color:var(--dim);font-size:var(--fs-sm);line-height:1.6;
}
.empty strong{display:block;color:var(--text-soft);font-size:var(--fs-md);font-weight:600}
.empty .empty-icon{display:flex;align-items:center;justify-content:center;width:38px;height:38px;margin-bottom:2px;border-radius:var(--r-md);background:var(--surface-2);color:var(--dim-2)}
.empty .empty-icon svg{width:18px;height:18px}

/* ── toast ──────────────────────────────────────────────────────────── */
#toast-host{position:fixed;right:20px;bottom:20px;z-index:400;display:flex;flex-direction:column;align-items:flex-end;gap:8px;pointer-events:none}
.toast{
  display:flex;align-items:center;gap:8px;max-width:min(360px,calc(100vw - 32px));
  padding:11px 14px;border-radius:var(--r-md);pointer-events:auto;
  font-size:var(--fs-sm);font-weight:500;
  background:var(--surface);color:var(--text);
  border:1px solid var(--line);
  box-shadow:var(--shadow-md);
  opacity:0;transform:translateY(6px);
  transition:opacity var(--dur) var(--ease),transform var(--dur) var(--ease);
}
.toast.show{opacity:1;transform:none}
.toast::before{content:"";flex:0 0 auto;width:7px;height:7px;border-radius:50%;background:var(--accent);box-shadow:0 0 0 3px color-mix(in srgb,var(--accent) 18%,transparent)}
.toast.error{border-color:color-mix(in srgb,var(--danger) 30%,var(--line));color:var(--danger)}
.toast.error::before{background:var(--danger);box-shadow:0 0 0 3px color-mix(in srgb,var(--danger) 16%,transparent)}

/* ── quiet brand references ─────────────────────────────────────────── */
.foot a,.err-foot a{color:inherit;text-decoration:none;border-bottom:1px solid transparent;transition:color var(--dur) var(--ease),border-color var(--dur) var(--ease)}
@media(hover:hover) and (pointer:fine){.foot a:hover,.err-foot a:hover{color:var(--accent);border-bottom-color:color-mix(in srgb,var(--accent) 40%,transparent)}}
.about-link{display:inline-flex;align-items:center;gap:5px;font-size:var(--fs-xs);color:var(--dim);text-decoration:none;transition:color var(--dur) var(--ease)}
@media(hover:hover) and (pointer:fine){.about-link:hover{color:var(--accent)}}
.about-link svg{width:13px;height:13px;flex-shrink:0}

/* ── modal ──────────────────────────────────────────────────────────── */
.modal-overlay{
  position:fixed;inset:0;z-index:300;display:flex;align-items:center;justify-content:center;padding:20px;
  background:rgba(8,18,13,.55);backdrop-filter:blur(4px);
  opacity:0;visibility:hidden;
  transition:opacity var(--dur) var(--ease),visibility var(--dur) var(--ease);
  overscroll-behavior:contain;
}
.modal-overlay.show{opacity:1;visibility:visible}
.modal{
  width:100%;max-width:400px;padding:20px;
  background:var(--surface);border:1px solid var(--line);border-radius:var(--r-lg);
  box-shadow:var(--shadow-lg);
  transform:translateY(6px) scale(.99);
  transition:transform var(--dur) var(--ease);
}
.modal-overlay.show .modal{transform:none}
.modal h3{font-size:var(--fs-md);font-weight:600;margin-bottom:4px}
.modal p{font-size:var(--fs-sm);color:var(--dim);line-height:1.6;margin-bottom:14px}
.modal .modal-field{margin-bottom:12px}
.modal .modal-field label{display:block;margin-bottom:5px}
.modal .modal-actions{display:flex;justify-content:flex-end;gap:8px;margin-top:16px}
.modal-close{
  display:inline-flex;align-items:center;justify-content:center;width:30px;height:30px;flex-shrink:0;
  border-radius:var(--r-sm);color:var(--dim);border:1px solid transparent;
  transition:background var(--dur) var(--ease),color var(--dur) var(--ease),border-color var(--dur) var(--ease);
}
.modal-close svg{width:16px;height:16px}
@media(hover:hover) and (pointer:fine){.modal-close:hover{background:var(--surface-2);color:var(--text);border-color:var(--line)}}

/* ── overflow menu (details/summary) ────────────────────────────────── */
.row-menu{position:relative;flex-shrink:0}
.row-menu>summary{
  display:inline-flex;align-items:center;justify-content:center;width:30px;height:30px;
  border-radius:var(--r-sm);color:var(--dim);list-style:none;cursor:pointer;
  border:1px solid transparent;
  transition:background var(--dur) var(--ease),color var(--dur) var(--ease),border-color var(--dur) var(--ease);
}
.row-menu>summary::-webkit-details-marker{display:none}
@media(hover:hover) and (pointer:fine){.row-menu>summary:hover{background:var(--surface-2);color:var(--text);border-color:var(--line)}}
.row-menu[open]>summary{background:var(--surface-2);color:var(--text);border-color:var(--line)}
.menu-panel{
  position:absolute;top:calc(100% + 6px);right:0;z-index:80;min-width:196px;padding:5px;
  background:var(--surface);border:1px solid var(--line);border-radius:var(--r-md);
  box-shadow:var(--shadow-md);
}
.menu-panel button{
  display:flex;align-items:center;gap:9px;width:100%;min-height:32px;padding:0 9px;
  border-radius:var(--r-sm);font-size:var(--fs-sm);color:var(--text-soft);text-align:left;
  transition:background var(--dur) var(--ease),color var(--dur) var(--ease);
}
.menu-panel button svg{width:14px;height:14px;flex-shrink:0;color:var(--dim)}
@media(hover:hover) and (pointer:fine){.menu-panel button:hover{background:var(--surface-2);color:var(--text)}}
.menu-panel button.is-danger{color:var(--danger)}
.menu-panel button.is-danger svg{color:var(--danger)}
@media(hover:hover) and (pointer:fine){.menu-panel button.is-danger:hover{background:var(--danger-soft)}}
.menu-panel .menu-sep{height:1px;margin:5px 4px;background:var(--line)}
.menu-panel .menu-label{padding:7px 9px 3px;font-size:.6875rem;font-weight:600;letter-spacing:.04em;text-transform:uppercase;color:var(--dim-2)}

/* ── auth ───────────────────────────────────────────────────────────── */
.auth-screen{
  position:relative;
  display:grid;place-items:center;min-height:100dvh;padding:clamp(24px,6vh,72px) 16px;
  background:
    radial-gradient(60% 50% at 50% -10%,color-mix(in srgb,var(--accent) 12%,transparent),transparent 70%),
    var(--bg);
}
.auth-card{
  width:100%;max-width:392px;padding:28px;
  background:var(--surface);border:1px solid var(--line);border-radius:var(--r-lg);
  box-shadow:var(--shadow-md);
}
.auth-card h1{display:flex;align-items:center;justify-content:center;gap:2px;font-size:var(--fs-xl);font-weight:600;letter-spacing:-.01em}
.auth-card h1 svg{margin-right:8px!important}
.auth-card h1 .accent{color:var(--accent)}
.auth-card .subtitle{margin:6px 0 22px;text-align:center;font-size:var(--fs-sm);color:var(--dim)}
.auth-card .btn[type=submit]{width:100%;min-height:42px;margin-top:4px}
.auth-card .switch{margin-top:16px;text-align:center;font-size:var(--fs-sm);color:var(--dim)}
.auth-card .switch a{color:var(--accent);font-weight:550;cursor:pointer}
.auth-utility{position:absolute;top:16px;right:16px;z-index:2;display:flex;align-items:center;gap:6px}
.captcha-row{display:flex;align-items:stretch;gap:8px}
.captcha-row .captcha-img{
  flex:0 0 auto;display:grid;place-items:center;width:124px;height:38px;overflow:hidden;
  border:1px solid var(--line);border-radius:var(--r-md);background:var(--surface-2);cursor:pointer;
  transition:border-color var(--dur) var(--ease);
}
.captcha-row .captcha-img svg{width:100%;height:100%}
@media(hover:hover) and (pointer:fine){.captcha-row .captcha-img:hover{border-color:var(--line-strong)}}
.captcha-row input{flex:1;min-width:0}
`

const uiDashboardCSS = `
.dashboard-grid{
  display:grid;grid-template-columns:minmax(272px,320px) minmax(0,1fr);
  gap:20px;align-items:start;
}
.sidebar{position:sticky;top:calc(var(--nav-h) + 20px);max-height:calc(100dvh - var(--nav-h) - 40px);overflow:auto;overscroll-behavior:contain;scrollbar-gutter:stable}
.create-card{margin-bottom:0}
.create-form-body .switch-field{padding:11px 0}
.create-form-body .switch-field:last-of-type{border-bottom:1px solid var(--line)}
.create-form-body .form-actions{margin-top:14px}
.desc.is-danger{color:var(--danger)}
.create-card-header{display:flex;align-items:center;justify-content:space-between;gap:10px}
.create-panel-toggle{display:none}
.main-content{min-width:0}

/* ── site list ──────────────────────────────────────────────────────── */
.site-tabs{display:flex;gap:2px;padding:8px 10px 0;border-bottom:1px solid var(--line)}
.site-tab{
  position:relative;display:inline-flex;align-items:center;gap:6px;
  min-height:36px;padding:0 12px;border-radius:var(--r-sm) var(--r-sm) 0 0;
  font-size:var(--fs-sm);font-weight:550;color:var(--dim);
  transition:color var(--dur) var(--ease),background var(--dur) var(--ease);
}
@media(hover:hover) and (pointer:fine){.site-tab:hover{color:var(--text);background:var(--surface-2)}}
.site-tab.active{color:var(--text)}
.site-tab.active::after{content:"";position:absolute;left:10px;right:10px;bottom:-1px;height:2px;border-radius:2px 2px 0 0;background:var(--accent)}
.site-tabs .tab-count{font-family:var(--ff-mono);font-size:.6875rem;color:var(--dim);font-variant-numeric:tabular-nums}

.site-list{list-style:none;display:flex;flex-direction:column;gap:8px}
.site-item{
  background:var(--surface);border:1px solid var(--line);border-radius:var(--r-md);
  transition:border-color var(--dur) var(--ease),box-shadow var(--dur) var(--ease);
}
@media(hover:hover) and (pointer:fine){.site-item:hover{border-color:var(--line-strong);box-shadow:var(--shadow-sm)}}
.site-item.is-open{border-color:color-mix(in srgb,var(--accent) 42%,var(--line));box-shadow:var(--shadow-sm)}
.site-head{
  display:grid;grid-template-columns:18px minmax(0,1fr) auto;align-items:center;gap:10px;
  padding:12px 12px 12px 10px;cursor:pointer;
}
.site-chevron{display:flex;align-items:center;justify-content:center;color:var(--dim-2);transition:transform var(--dur) var(--ease),color var(--dur) var(--ease)}
.site-chevron svg{width:14px;height:14px}
.site-item.is-open .site-chevron{transform:rotate(90deg);color:var(--accent)}
.site-main{min-width:0}
.site-title{display:flex;align-items:center;gap:7px;min-width:0}
.site-name{font-size:var(--fs-base);font-weight:600;color:var(--text);overflow:hidden;text-overflow:ellipsis;white-space:nowrap;max-width:100%}
.site-sub{display:flex;align-items:center;gap:8px;margin-top:3px;min-width:0;font-size:var(--fs-xs);color:var(--dim)}
.site-url{display:inline-flex;align-items:center;gap:4px;min-width:0;font-family:var(--ff-mono);font-size:.75rem;color:var(--dim);transition:color var(--dur) var(--ease)}
.site-url .url-text{overflow:hidden;text-overflow:ellipsis;white-space:nowrap;min-width:0}
.site-url .url-path{display:none}
.site-url svg{width:11px;height:11px;flex-shrink:0;opacity:0;transition:opacity var(--dur) var(--ease)}
@media(hover:hover) and (pointer:fine){.site-url:hover{color:var(--accent)}.site-url:hover svg{opacity:1}}
.site-sub .sep{color:var(--line-strong)}
.site-visits{flex-shrink:0;font-variant-numeric:tabular-nums;white-space:nowrap}
.site-actions{display:flex;align-items:center;gap:6px;flex-shrink:0}
.upload-btn{position:relative;overflow:hidden}
.upload-btn input[type=file]{position:absolute;inset:0;width:100%;height:100%;opacity:0;cursor:pointer}

/* ── site detail ────────────────────────────────────────────────────── */
.site-detail{
  display:grid;grid-template-rows:0fr;overflow:hidden;
  transition:grid-template-rows 220ms var(--ease);
  border-top:1px solid transparent;
}
.site-detail.show{grid-template-rows:1fr;border-top-color:var(--line)}
.detail-clip{overflow:hidden;min-height:0}
.detail-inner{padding:16px;background:color-mix(in srgb,var(--surface-2) 72%,var(--surface));border-radius:0 0 calc(var(--r-md) - 1px) calc(var(--r-md) - 1px)}
.detail-grid{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:14px 24px}
.detail-section{min-width:0}
.detail-label{font-size:.6875rem;font-weight:600;letter-spacing:.05em;text-transform:uppercase;color:var(--dim);margin-bottom:4px}
.detail-value{font-size:var(--fs-sm);color:var(--text-soft);line-height:1.55;overflow-wrap:anywhere}
.detail-value.is-muted{color:var(--dim)}
.detail-value .pwd-value{font-family:var(--ff-mono);font-size:var(--fs-sm);letter-spacing:.06em}
.pwd-actions{display:inline-flex;align-items:center;gap:3px;margin-left:6px;vertical-align:middle}
.visits-grid{display:flex;gap:18px;flex-wrap:wrap}
.visits-grid .visit-cell{display:flex;flex-direction:column;gap:1px}
.visits-grid .visit-key{font-size:.6875rem;color:var(--dim)}
.visits-grid .visit-num{font-family:var(--ff-mono);font-size:var(--fs-md);font-weight:600;color:var(--text);font-variant-numeric:tabular-nums}
.detail-toggles{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:0 24px;margin-top:12px;border-top:1px solid var(--line)}
.detail-actions{display:flex;flex-wrap:wrap;gap:8px;padding-top:14px}

/* ── file tree ──────────────────────────────────────────────────────── */
.site-files{display:block}
.file-tree-section{margin-top:16px}
.file-clear-btn{color:var(--danger);border-color:color-mix(in srgb,var(--danger) 32%,var(--line))}
@media(hover:hover) and (pointer:fine){.file-clear-btn:hover:not(:disabled){background:var(--danger-soft);border-color:var(--danger);color:var(--danger)}}
.file-tree-head{display:flex;align-items:center;justify-content:space-between;gap:12px;flex-wrap:wrap;margin-bottom:8px}
.file-tree-title{display:flex;align-items:center;gap:8px;min-width:0}
.file-tree-summary{font-family:var(--ff-mono);font-size:.6875rem;color:var(--dim);font-variant-numeric:tabular-nums;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.file-tree-actions{display:flex;align-items:center;gap:6px;flex-shrink:0}

/* ── transfer progress ──────────────────────────────────────────────── */
.transfer-status{
  display:none;padding:10px 14px;margin-top:8px;
  background:var(--surface-2);border:1px solid var(--line);border-radius:var(--r-md);
}
.transfer-status.show{display:block}
.transfer-head{display:flex;align-items:center;justify-content:space-between;gap:12px;margin-bottom:6px;font-size:var(--fs-xs);color:var(--dim)}
.transfer-label{overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.transfer-percent{flex-shrink:0;font-family:var(--ff-mono);color:var(--text-soft);font-variant-numeric:tabular-nums}
.transfer-track{height:4px;background:var(--line);border-radius:var(--r-pill);overflow:hidden}
.transfer-fill{display:block;width:0;height:100%;background:var(--accent);border-radius:var(--r-pill);transition:width 160ms linear}
.transfer-status.indeterminate .transfer-fill{width:32%;animation:transfer-slide 1.2s var(--ease) infinite}
.transfer-status.error .transfer-fill{background:var(--danger)}
.transfer-status.error .transfer-label,.transfer-status.error .transfer-percent{color:var(--danger)}
@keyframes transfer-slide{0%{transform:translateX(-110%)}100%{transform:translateX(320%)}}
@media(prefers-reduced-motion:reduce){.transfer-status.indeterminate .transfer-fill{animation:none;width:100%}}

.deploy-hint{margin-top:14px;font-size:var(--fs-xs);line-height:1.65;color:var(--dim);max-width:78ch}
.pwd-none{color:var(--dim)}

/* ── org pinned list ────────────────────────────────────────────────── */
.org-pinned-item{display:flex;align-items:center;gap:14px;padding:11px 12px;border:1px solid var(--line);border-radius:var(--r-md);margin-bottom:8px;background:var(--surface)}
.pinned-text{flex:1;min-width:0}
.pinned-title{display:flex;align-items:baseline;gap:8px;min-width:0}
.pinned-name{font-size:var(--fs-base);font-weight:600;color:var(--text);overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.pinned-desc{font-size:var(--fs-xs);color:var(--dim);overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.pinned-url{flex:0 1 auto;max-width:38%;font-family:var(--ff-mono);font-size:.6875rem;color:var(--dim);overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.pinned-actions{display:flex;align-items:center;gap:6px;flex-shrink:0}

/* ── copy button ────────────────────────────────────────────────────── */
.copy-btn{
  display:inline-flex;align-items:center;justify-content:center;width:24px;height:24px;
  border-radius:var(--r-sm);border:1px solid transparent;color:var(--dim);
  transition:background var(--dur) var(--ease),color var(--dur) var(--ease),border-color var(--dur) var(--ease);
}
.copy-btn svg{width:12px;height:12px}
@media(hover:hover) and (pointer:fine){.copy-btn:hover{background:var(--surface);color:var(--accent);border-color:var(--line)}}

/* ── organization modal ─────────────────────────────────────────────── */
#org-modal .modal{width:100%;max-width:840px;padding:0;max-height:min(680px,calc(100dvh - 40px));overflow:hidden;display:flex;flex-direction:column}
.org-modal-header{display:flex;align-items:flex-start;justify-content:space-between;gap:16px;padding:18px 20px;border-bottom:1px solid var(--line);flex-shrink:0}
.org-modal-heading{min-width:0}
.org-modal-header h3{margin:0;font-size:var(--fs-md)}
.org-feature-desc{margin-top:4px;font-size:var(--fs-sm);color:var(--dim);line-height:1.55;max-width:66ch}
#org-section{display:flex;flex-direction:column;min-height:0;overflow:hidden}
.org-empty{display:grid;grid-template-columns:1fr 1fr;min-height:0;overflow:auto;overscroll-behavior:contain}
.org-empty-panel{display:flex;flex-direction:column;gap:10px;padding:20px 22px;min-height:0}
.org-empty-panel+.org-empty-panel{border-left:1px solid var(--line)}
.org-panel-title{display:flex;align-items:center;gap:10px;font-size:var(--fs-md);font-weight:600}
.org-panel-icon{display:flex;align-items:center;justify-content:center;width:32px;height:32px;border-radius:var(--r-md);flex-shrink:0}
.org-panel-icon svg{width:16px;height:16px}
.org-panel-icon.create{background:var(--accent-soft);color:var(--accent-ink)}
.org-panel-icon.join{background:color-mix(in srgb,#3b82f6 12%,var(--surface));color:#2563c9}
[data-theme="dark"] .org-panel-icon.join{color:#8ab4ff}
.org-panel-desc{font-size:var(--fs-sm);color:var(--dim);line-height:1.6}
.org-panel-form{margin-top:auto;display:flex;flex-direction:column;gap:8px}
.org-bound{display:flex;flex-direction:column;min-height:0;flex:1}
.org-summary{display:grid;grid-template-columns:minmax(0,1fr) minmax(210px,300px) auto;align-items:center;gap:16px;padding:14px 18px;border-bottom:1px solid var(--line);background:color-mix(in srgb,var(--surface-2) 45%,var(--surface));flex-shrink:0}
.org-header{display:flex;align-items:center;gap:10px;min-width:0}
.org-avatar{display:flex;align-items:center;justify-content:center;width:36px;height:36px;border-radius:var(--r-md);background:var(--accent);color:var(--on-accent);font-family:var(--ff-mono);font-size:var(--fs-base);font-weight:700;flex-shrink:0}
.org-title{display:flex;align-items:center;gap:8px;min-width:0;font-size:var(--fs-md);font-weight:600;overflow-wrap:anywhere}
.org-role-badge{display:inline-flex;align-items:center;height:19px;padding:0 7px;border-radius:var(--r-pill);font-size:.6875rem;font-weight:600;letter-spacing:.03em;text-transform:uppercase;flex-shrink:0}
.org-role-badge.owner{background:var(--warn-soft);color:var(--warn)}
.org-role-badge.member{background:var(--accent-soft);color:var(--accent-ink)}
.org-invite-section{display:flex;flex-direction:column;gap:5px;min-width:0}
.org-invite-toggle{display:flex;align-items:center;justify-content:space-between;gap:8px;padding:0;font-size:.6875rem;font-weight:600;letter-spacing:.05em;text-transform:uppercase;color:var(--dim);cursor:pointer}
@media(hover:hover) and (pointer:fine){.org-invite-toggle:hover{color:var(--text)}}
.org-invite-chevron{transition:transform var(--dur) var(--ease);color:var(--dim)}
.org-invite-collapsed .org-invite-chevron{transform:rotate(-90deg)}
.org-invite-collapsed .org-invite-detail{display:none}
.org-invite-row{display:flex;align-items:center;gap:4px;padding:2px 2px 2px 9px;border:1px solid var(--line);border-radius:var(--r-sm);background:var(--surface)}
.org-invite-code{flex:1;min-width:0;font-family:var(--ff-mono);font-size:var(--fs-sm);color:var(--accent);letter-spacing:.06em;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.org-copy-btn{display:inline-flex;align-items:center;justify-content:center;width:26px;height:26px;border-radius:var(--r-sm);border:1px solid transparent;color:var(--dim);flex-shrink:0;transition:background var(--dur) var(--ease),color var(--dur) var(--ease),border-color var(--dur) var(--ease)}
.org-copy-btn svg{width:13px;height:13px}
@media(hover:hover) and (pointer:fine){.org-copy-btn:hover{color:var(--accent);border-color:var(--line);background:var(--surface-2)}}
.org-summary-action{min-width:0}
.org-summary-action .btn{white-space:nowrap}
.org-danger-btn{color:var(--danger);border-color:color-mix(in srgb,var(--danger) 34%,var(--line))}
@media(hover:hover) and (pointer:fine){.org-danger-btn:hover:not(:disabled){background:var(--danger-soft);border-color:var(--danger);color:var(--danger)}}
.org-bound-right{display:flex;flex-direction:column;flex:1;min-height:0;overflow:hidden;padding:16px 18px}
.org-members-header{display:flex;align-items:center;justify-content:space-between;gap:8px;margin-bottom:10px}
.org-members-header h3{font-size:var(--fs-sm);font-weight:600}
.org-members-count{padding:1px 8px;border-radius:var(--r-pill);background:var(--surface-2);color:var(--dim);font-family:var(--ff-mono);font-size:.6875rem;font-variant-numeric:tabular-nums}
.org-member-search{position:relative;margin-bottom:10px}
.org-member-search svg{position:absolute;left:9px;top:50%;transform:translateY(-50%);width:13px;height:13px;color:var(--dim-2);pointer-events:none}
.org-member-search input{width:100%;min-height:34px;padding:0 11px 0 28px;border:1px solid var(--line);border-radius:var(--r-md);background:var(--surface);font-size:var(--fs-sm);transition:border-color var(--dur) var(--ease),box-shadow var(--dur) var(--ease)}
.org-member-search input:focus{outline:none;border-color:var(--accent);box-shadow:0 0 0 3px var(--accent-ring)}
#org-member-list{flex:1;min-height:0;overflow-y:auto;overscroll-behavior:contain;padding-right:4px;padding-bottom:2px;scrollbar-gutter:stable}
.org-member-item{display:flex;align-items:center;gap:10px;min-height:46px;padding:6px 8px;border-bottom:1px solid var(--line)}
.org-member-item:first-child{border-top:1px solid var(--line)}
@media(hover:hover) and (pointer:fine){.org-member-item:hover{background:var(--surface-2)}}
.org-member-avatar{display:flex;align-items:center;justify-content:center;width:28px;height:28px;border-radius:var(--r-sm);background:var(--surface-2);color:var(--dim);font-family:var(--ff-mono);font-size:.6875rem;font-weight:600;flex-shrink:0}
.org-member-avatar.owner{background:var(--warn-soft);color:var(--warn)}
.org-member-info{flex:1;min-width:0}
.org-member-email{font-size:var(--fs-sm);color:var(--text);overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.org-member-role{font-size:.6875rem;color:var(--dim);font-family:var(--ff-mono)}
.org-member-remove{padding:4px 8px;border-radius:var(--r-sm);font-size:var(--fs-xs);font-weight:550;color:var(--danger);border:1px solid transparent;flex-shrink:0;white-space:nowrap;transition:background var(--dur) var(--ease),color var(--dur) var(--ease)}
@media(hover:hover) and (pointer:fine){.org-member-remove:hover{background:var(--danger-soft)}}

@media(max-width:1024px){
  .dashboard-grid{grid-template-columns:minmax(0,1fr)}
  .sidebar{position:static;max-height:none;overflow:visible}
  .create-card .create-form-body{display:none}
  .create-card.mobile-open .create-form-body{display:block}
  .create-panel-toggle{display:inline-flex;align-items:center;gap:5px;min-height:30px;padding:0 9px;border:1px solid var(--line);border-radius:var(--r-sm);background:var(--surface);color:var(--dim);font-size:var(--fs-xs);font-weight:550}
  .create-panel-toggle svg{width:14px;height:14px;transition:transform var(--dur) var(--ease)}
  .create-card.mobile-open .create-panel-toggle svg{transform:rotate(180deg)}
  .detail-grid{grid-template-columns:minmax(0,1fr)}
}
@media(max-width:720px){
  .navbar{min-height:var(--nav-h);padding:0 12px;flex-wrap:nowrap}
  .navbar .logo .logo-text{display:none}
  .navbar .nav-right{display:none;position:absolute;top:calc(100% + 6px);right:10px;width:min(320px,calc(100vw - 20px));max-height:calc(100dvh - var(--nav-h) - 20px);overflow:auto;overscroll-behavior:contain;padding:8px;flex-direction:column;align-items:stretch;gap:6px;background:var(--surface);border:1px solid var(--line);border-radius:var(--r-md);box-shadow:var(--shadow-md)}
  .navbar .nav-right.open{display:flex}
  .navbar .nav-section{width:100%;flex-direction:column;align-items:stretch;gap:2px}
  .navbar .nav-preferences{flex-direction:row;align-items:center;justify-content:space-between;padding:6px 0;border-left:0;border-right:0;border-top:1px solid var(--line);border-bottom:1px solid var(--line)}
  .navbar .nav-account{flex-direction:row;align-items:center;gap:6px}
  .navbar .nav-account .btn-link{width:auto;flex:0 0 auto}
  .navbar .email{max-width:100%;flex:1 1 auto;min-width:0}
  .navbar .btn-link,.navbar .btn-icon,.navbar .admin-link{justify-content:flex-start;width:100%}
  .navbar .nav-preferences .btn-icon{width:auto}
  .mobile-menu-toggle{display:inline-flex;align-items:center;justify-content:center;width:38px;height:38px;flex:0 0 38px;border:1px solid var(--line);border-radius:var(--r-md);background:var(--surface);color:var(--text)}
  .mobile-menu-toggle svg{width:18px;height:18px}
  .mobile-menu-toggle .menu-close-icon{display:none}
  .mobile-menu-toggle[aria-expanded="true"] .menu-open-icon{display:none}
  .mobile-menu-toggle[aria-expanded="true"] .menu-close-icon{display:block}
  .container{padding:12px 12px 40px}
  .card-header{padding:12px 14px;min-height:48px}
  .card-body{padding:14px}
  .site-head{grid-template-columns:14px minmax(0,1fr);row-gap:8px;padding:11px 10px}
  .site-chevron svg{width:13px;height:13px}
  .site-name{font-size:var(--fs-md)}
  .site-actions{grid-column:1/-1;grid-row:2;justify-content:flex-end;gap:8px;margin-top:2px}
  .site-actions .upload-btn{min-height:34px;padding:0 14px}
  .site-sub{flex-wrap:nowrap;gap:8px}
  .site-sub .sep{display:none}
  .site-url{flex:1 1 auto;min-width:0}
  .site-url .url-full{display:none}
  .site-url .url-path{display:inline}
  .site-visits{flex:0 0 auto}
  .row-menu>summary{width:34px;height:34px}
  .detail-toggles{grid-template-columns:minmax(0,1fr)}
  .detail-inner{padding:12px}
  .detail-actions .btn{flex:1}
  .file-tree-actions{width:100%}
  .file-tree-actions .btn{flex:1}
  .org-empty{grid-template-columns:minmax(0,1fr)}
  .org-empty-panel+.org-empty-panel{border-left:0;border-top:1px solid var(--line)}
  .org-summary{grid-template-columns:minmax(0,1fr);gap:12px;padding:12px 14px}
  .org-summary-action .btn{width:100%}
  .org-bound-right{padding:12px 14px}
  #org-modal{padding:0}
  #org-modal .modal{max-width:none;max-height:100dvh;height:100dvh;border-radius:0;border-left:0;border-right:0}
  .modal{max-width:none;border-radius:var(--r-lg) var(--r-lg) 0 0;align-self:flex-end}
  .modal-overlay{padding:0;align-items:flex-end}
  #toast-host{right:10px;left:10px;bottom:calc(10px + env(safe-area-inset-bottom));align-items:stretch}
  .toast{max-width:none}
  .auth-card{padding:22px 18px}
  .auth-utility{top:10px;right:10px}
  .captcha-row{display:grid;grid-template-columns:112px minmax(0,1fr)}
  .captcha-row .captcha-img{width:112px}
  .pagination{padding-top:10px}
}
@media(max-width:420px){
  .visits-grid{gap:14px}
  .detail-inner{padding:12px 11px}
}
`

const uiAdminCSS = `
.admin-layout{display:grid;grid-template-columns:212px minmax(0,1fr);align-items:start}
.admin-sidebar{
  position:sticky;top:var(--nav-h);align-self:start;
  display:flex;flex-direction:column;gap:2px;
  height:calc(100dvh - var(--nav-h));padding:14px 10px;
  overflow-y:auto;overscroll-behavior:contain;
  border-right:1px solid var(--line);background:color-mix(in srgb,var(--surface) 60%,var(--bg));
}
.admin-nav-group{display:flex;flex-direction:column;gap:2px}
.admin-nav-label{padding:8px 10px 4px;font-size:.6875rem;font-weight:600;letter-spacing:.06em;text-transform:uppercase;color:var(--dim-2)}
.admin-nav-item{
  display:flex;align-items:center;gap:9px;min-height:36px;padding:0 10px;
  border-radius:var(--r-md);font-size:var(--fs-sm);font-weight:500;color:var(--text-soft);
  cursor:pointer;user-select:none;
  transition:background var(--dur) var(--ease),color var(--dur) var(--ease);
}
.admin-nav-item svg{width:15px;height:15px;flex-shrink:0;color:var(--dim)}
@media(hover:hover) and (pointer:fine){.admin-nav-item:hover{background:var(--surface-2);color:var(--text)}}
.admin-nav-item.active{background:var(--accent);color:var(--on-accent);font-weight:550}
.admin-nav-item.active svg{color:currentColor}
@media(hover:hover) and (pointer:fine){.admin-nav-item.active:hover{background:var(--accent-hover);color:var(--on-accent)}}
.admin-sidebar-foot{margin-top:auto;padding:10px;border-top:1px solid var(--line);display:flex;flex-direction:column;gap:2px}
.admin-sidebar-foot .ver{font-family:var(--ff-mono);font-size:.6875rem;color:var(--dim-2);padding:0 2px}
.admin-main{min-width:0;max-width:1180px;padding:clamp(16px,2.4vw,28px) clamp(14px,2.6vw,32px) 64px}

.admin-tab-content{display:none}
.admin-tab-content.active{display:block}

/* stats */
.stats-grid{display:grid;grid-template-columns:repeat(3,minmax(0,1fr));gap:12px}
.stat-card{
  display:flex;flex-direction:column;gap:10px;padding:14px 16px;
  background:var(--surface);border:1px solid var(--line);border-radius:var(--r-md);
}
.stat-card .stat-top{display:flex;align-items:center;justify-content:space-between;gap:8px}
.stat-card .label{font-size:var(--fs-xs);font-weight:550;color:var(--dim)}
.stat-card .stat-icon{display:flex;align-items:center;justify-content:center;width:26px;height:26px;border-radius:var(--r-sm);background:var(--surface-2);color:var(--dim)}
.stat-card .stat-icon svg{width:14px;height:14px}
.stat-card .num{font-family:var(--ff-mono);font-size:1.625rem;font-weight:600;line-height:1;color:var(--text);font-variant-numeric:tabular-nums;letter-spacing:-.02em}
.stat-card .stat-note{font-size:var(--fs-xs);color:var(--dim-2)}

/* settings rows */
.toggle-row{
  display:flex;align-items:center;justify-content:space-between;gap:16px;
  padding:14px 0;border-bottom:1px solid var(--line);
}
.toggle-row:first-child{padding-top:0}
.toggle-info{min-width:0}
.toggle-label{font-size:var(--fs-sm);font-weight:550;color:var(--text)}
.toggle-desc{margin-top:2px;font-size:var(--fs-xs);color:var(--dim);line-height:1.55;max-width:62ch}
.toggle-switch{
  position:relative;flex:0 0 38px;width:38px;height:22px;border-radius:var(--r-pill);
  background:var(--surface-3);border:1px solid var(--line-strong);cursor:pointer;
  transition:background var(--dur) var(--ease),border-color var(--dur) var(--ease);
}
.toggle-switch::after{
  content:"";position:absolute;top:2px;left:2px;width:16px;height:16px;border-radius:50%;
  background:var(--surface);box-shadow:var(--shadow-sm);
  transition:transform var(--dur) var(--ease);
}
.toggle-switch.on{background:var(--accent);border-color:var(--accent)}
.toggle-switch.on::after{transform:translateX(16px)}
.num-setting,.url-setting{
  display:grid;grid-template-columns:minmax(0,1fr) auto;align-items:center;gap:12px 20px;
  padding:14px 0;border-bottom:1px solid var(--line);
}
.num-label,.url-setting .num-label{font-size:var(--fs-sm);font-weight:550;color:var(--text)}
.num-desc,.url-setting .num-desc{margin-top:2px;font-size:var(--fs-xs);color:var(--dim);line-height:1.55;max-width:62ch}
.num-control{display:flex;align-items:center;gap:8px;flex-shrink:0}
.num-control input[type=number]{width:88px;font-variant-numeric:tabular-nums;text-align:right}
.num-control input[type=text]{min-width:min(320px,44vw)}
.num-control .unit{font-size:var(--fs-xs);color:var(--dim)}
.url-setting .num-control{grid-column:1/-1;justify-content:flex-start;width:100%}
.url-setting .num-control input{flex:1;min-width:0}
.url-warn{
  grid-column:1/-1;margin:2px 0 0;padding:9px 11px;border-radius:var(--r-sm);
  font-size:var(--fs-xs);line-height:1.55;
  color:var(--warn);background:var(--warn-soft);border:1px solid color-mix(in srgb,var(--warn) 28%,var(--line));
}
.settings-block{margin-top:4px}
.settings-block>h3{font-size:var(--fs-sm);font-weight:600;margin-bottom:2px}
.settings-block>.sub{font-size:var(--fs-xs);color:var(--dim);margin-bottom:8px}

/* admin tables */
.site-row{cursor:pointer}
.cell-lead{display:flex;align-items:center;gap:8px;min-width:0}
.cell-lead-text{min-width:0}
.row-caret{display:inline-flex;align-items:center;color:var(--dim-2);flex-shrink:0;transition:transform var(--dur) var(--ease),color var(--dur) var(--ease)}
.row-caret svg{width:13px;height:13px}
.site-row.expanded .row-caret{transform:rotate(90deg);color:var(--accent)}
.cell-strong{display:block;font-weight:550;color:var(--text);overflow:hidden;text-overflow:ellipsis;white-space:nowrap;max-width:34ch}
.cell-sub{display:none;margin-top:3px}
.cell-link{display:inline-block;max-width:26ch;margin-top:1px;font-family:var(--ff-mono);font-size:var(--fs-xs);color:var(--dim);overflow:hidden;text-overflow:ellipsis;white-space:nowrap;transition:color var(--dur) var(--ease)}
@media(hover:hover) and (pointer:fine){.cell-link:hover{color:var(--accent)}}
.admin-detail-row{display:none}
.admin-detail-row.show{display:table-row}
.admin-detail-row>td{padding:0;background:var(--surface-2)}
.admin-detail-body{padding:12px 14px}
.admin-detail-body .detail-label{margin-bottom:6px}
.admin-detail-body .file-list{background:var(--surface)}
.admin-file-list{max-height:260px;overflow:auto;overscroll-behavior:contain}
.admin-tag{margin-left:6px;padding:1px 7px;border-radius:var(--r-pill);background:var(--surface-2);border:1px solid var(--line);font-size:.6875rem;font-weight:600;letter-spacing:.02em;color:var(--dim)}
.tab-desc{font-size:var(--fs-sm);color:var(--dim);line-height:1.6;margin-bottom:12px;max-width:72ch}
.ver-inline{font-family:var(--ff-mono);color:var(--text)}
.result-loading{font-size:var(--fs-sm);color:var(--dim)}
.result-warn{display:inline-flex;align-items:center;gap:6px;margin-bottom:8px;font-size:var(--fs-sm);color:var(--warn);font-weight:550}

/* system info */
.info-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(232px,1fr));border:1px solid var(--line);border-radius:var(--r-md);overflow:hidden;background:var(--surface)}
.info-item{display:flex;flex-direction:column;gap:3px;min-height:64px;padding:11px 13px;background:var(--surface);border-right:1px solid var(--line);border-bottom:1px solid var(--line)}
.info-label{font-size:.6875rem;font-weight:600;letter-spacing:.05em;text-transform:uppercase;color:var(--dim)}
.info-value{font-size:var(--fs-sm);color:var(--text-soft);overflow-wrap:anywhere}
.info-value.accent{font-family:var(--ff-mono);color:var(--accent)}
.info-value a{display:block;font-family:var(--ff-mono);font-size:var(--fs-xs);color:var(--accent);margin-top:2px}
.info-about{margin-top:12px;display:flex;justify-content:flex-end}
.update-inline{margin-top:16px;padding-top:16px;border-top:1px solid var(--line)}
.update-inline p{font-size:var(--fs-sm);color:var(--dim)}
.result-box{margin-top:10px;padding:12px;border:1px solid var(--line);border-radius:var(--r-md);background:var(--surface-2);font-size:var(--fs-sm)}
.result-box:empty{display:none}
.result-ok{display:inline-flex;align-items:center;gap:6px;font-size:var(--fs-sm);color:var(--accent)}
.result-ok svg{width:14px;height:14px}
.result-err{display:inline-flex;align-items:center;gap:6px;font-size:var(--fs-sm);color:var(--danger)}
.result-err svg{width:14px;height:14px}
.orphan-list{max-height:220px;overflow:auto;overscroll-behavior:contain;margin:10px 0;border:1px solid var(--line);border-radius:var(--r-sm);background:var(--surface)}
.orphan-list div{display:flex;align-items:center;gap:7px;padding:5px 9px;font-family:var(--ff-mono);font-size:var(--fs-xs);color:var(--dim);border-bottom:1px solid var(--line)}
.orphan-list div:last-child{border-bottom:0}
.orphan-list svg{width:12px;height:12px;flex-shrink:0}
.progress-inline{margin-top:12px}
.progress-inline .p-label{display:flex;justify-content:space-between;gap:10px;font-size:var(--fs-xs);color:var(--dim);margin-bottom:5px}
.progress-inline .p-track{height:5px;border-radius:var(--r-pill);background:var(--surface-3);overflow:hidden}
.progress-inline .p-fill{height:100%;background:var(--accent);border-radius:var(--r-pill);transition:width 200ms var(--ease)}

@media(max-width:900px){
  .admin-layout{grid-template-columns:minmax(0,1fr)}
  .admin-sidebar{
    position:sticky;top:var(--nav-h);z-index:40;
    flex-direction:row;align-items:center;gap:4px;
    height:auto;padding:8px 10px;overflow-x:auto;overflow-y:hidden;
    border-right:0;border-bottom:1px solid var(--line);
    background:color-mix(in srgb,var(--bg) 92%,transparent);backdrop-filter:blur(10px);
    scrollbar-width:none;
    -webkit-mask-image:linear-gradient(to right,#000 calc(100% - 16px),transparent);
    mask-image:linear-gradient(to right,#000 calc(100% - 16px),transparent);
  }
  .admin-sidebar::-webkit-scrollbar{display:none}
  .admin-nav-group{flex-direction:row;gap:4px}
  .admin-nav-label{display:none}
  .admin-nav-item{min-height:32px;padding:0 10px;white-space:nowrap;border-radius:var(--r-sm)}
  .admin-sidebar-foot{display:none}
  .admin-main{padding:16px 14px 48px}
  .stats-grid{grid-template-columns:repeat(2,minmax(0,1fr))}
  .stat-card:first-child{grid-column:1/-1}
}
@media(max-width:720px){
  .navbar{min-height:var(--nav-h);padding:0 12px;gap:8px}
  .navbar .logo{font-size:.9375rem;gap:6px}
  .navbar .logo .logo-text{display:none}
  .navbar .admin-tag{display:none}
  .navbar .nav-right,.navbar .nav-section{display:flex;align-items:center;gap:4px;min-width:0}
  .navbar .nav-preferences{padding:0;border:0}
  .navbar .nav-preferences .btn-icon{width:34px;min-width:34px;padding:0}
  .navbar .nav-account .btn-link{min-width:34px;padding:0 8px;justify-content:center}
  .navbar .nav-account .btn-link span{display:none}
  .admin-sidebar{top:var(--nav-h)}
  .col-hide-sm{display:none}
  .cell-sub{display:block}
  .cell-strong{max-width:52vw}
  td.cell-email{max-width:none;white-space:nowrap}
  td.cell-actions .btn+.btn{margin-left:0}
  .sites-table .cell-strong{max-width:36vw}
  .stats-grid{grid-template-columns:minmax(0,1fr)}
  .stat-card:first-child{grid-column:auto}
  .stat-card{padding:12px 14px}
  .stat-card .num{font-size:1.5rem}
  .num-setting,.url-setting{grid-template-columns:minmax(0,1fr)}
  .num-control{justify-content:flex-start}
  .num-control input[type=number]{width:80px}
  .toggle-row{align-items:flex-start;gap:12px}
  .table-wrap{margin:0}
  th,td{padding-left:10px;padding-right:10px}
}
`

const uiLandingCSS = `
body{overflow-x:hidden}
body::before{
  content:"";position:fixed;inset:0;z-index:-1;pointer-events:none;
  background:radial-gradient(56% 42% at 18% 0%,color-mix(in srgb,var(--accent) 9%,transparent),transparent 72%);
}
.site-nav{display:flex;align-items:center;justify-content:space-between;gap:16px;max-width:1200px;margin:0 auto;padding:0 clamp(14px,3vw,28px);min-height:60px}
.brand{display:inline-flex;align-items:center;gap:2px;font-size:1rem;font-weight:600;letter-spacing:-.01em}
.brand .accent{color:var(--accent)}
.brand svg{margin-right:6px!important}
.nav-actions{display:flex;align-items:center;gap:8px}
.nav-actions .lang-toggle{background:transparent;border-color:transparent}
.nav-actions .lang-toggle a.active{background:var(--surface-2);border:1px solid var(--line)}

.hero{
  display:grid;grid-template-columns:minmax(0,.86fr) minmax(0,1.14fr);
  align-items:center;gap:clamp(32px,5vw,72px);
  max-width:1200px;margin:0 auto;
  padding:clamp(28px,4vh,48px) clamp(14px,3vw,28px) clamp(36px,5vh,60px);
  min-height:min(620px,calc(100dvh - 60px));
}
.hero-copy{max-width:34rem}
.eyebrow{font-family:var(--ff-mono);font-size:.6875rem;letter-spacing:.09em;text-transform:uppercase;color:var(--accent);margin-bottom:16px}
.hero h1{font-size:clamp(2.5rem,5vw,4.2rem);line-height:1.04;font-weight:600;letter-spacing:-.035em;max-width:17ch}
.hero h1 span{color:var(--accent)}
.tagline{margin:18px 0 26px;font-size:clamp(.95rem,1.15vw,1.0625rem);line-height:1.65;color:var(--dim);max-width:34ch}
.cta-row{display:flex;align-items:center;gap:10px;flex-wrap:wrap}
.cta{display:inline-flex;align-items:center;justify-content:center;gap:7px;min-height:44px;padding:0 20px;border-radius:var(--r-md);font-size:var(--fs-md);font-weight:550;border:1px solid transparent;transition:background var(--dur) var(--ease),border-color var(--dur) var(--ease),color var(--dur) var(--ease),transform var(--dur) var(--ease)}
.cta:active{transform:scale(.99)}
.cta-primary{background:var(--accent);color:var(--on-accent)}
@media(hover:hover) and (pointer:fine){.cta-primary:hover{background:var(--accent-hover)}}
.cta-secondary{background:var(--surface);border-color:var(--line);color:var(--text)}
@media(hover:hover) and (pointer:fine){.cta-secondary:hover{border-color:var(--line-strong);background:var(--surface-2)}}

.hero-visual{position:relative;min-width:0}
.app-window{
  position:relative;overflow:hidden;background:var(--surface);
  border:1px solid var(--line);border-radius:var(--r-lg);box-shadow:var(--shadow-lg);
}
.app-bar{display:flex;align-items:center;gap:10px;padding:9px 12px;border-bottom:1px solid var(--line);background:color-mix(in srgb,var(--surface-2) 60%,var(--surface))}
.app-dots{display:inline-flex;gap:4px}
.app-dots i{width:7px;height:7px;border-radius:50%;background:var(--line-strong)}
.app-addr{font-family:var(--ff-mono);font-size:.6875rem;color:var(--dim);overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.app-body{padding:12px}
.app-tabs{display:flex;gap:14px;padding:0 2px 9px;border-bottom:1px solid var(--line);font-size:.8125rem;color:var(--dim)}
.app-tabs .on{position:relative;color:var(--text);font-weight:550}
.app-tabs .on::after{content:"";position:absolute;left:0;right:0;bottom:-10px;height:2px;background:var(--accent);border-radius:2px}
.app-row{display:grid;grid-template-columns:9px minmax(0,1fr) auto auto;grid-template-rows:auto auto;align-items:center;column-gap:8px;row-gap:1px;padding:9px 2px;border-bottom:1px solid var(--line);font-size:.8125rem}
.app-row b{grid-area:1/2/2/3;font-weight:550;color:var(--text);overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.app-row .dot{grid-area:1/1/2/2;width:7px;height:7px;border-radius:50%}
.app-row .dot.pub{background:var(--accent)}
.app-row .dot.pro{background:var(--warn)}
.app-row .chip{grid-area:1/3/2/4;justify-self:end;font-size:.6875rem;font-weight:550;padding:1px 7px;border-radius:var(--r-pill);border:1px solid transparent;white-space:nowrap}
.app-row .chip.pub{color:var(--accent-ink);background:var(--accent-soft);border-color:color-mix(in srgb,var(--accent) 24%,transparent)}
.app-row .chip.pro{color:var(--warn);background:var(--warn-soft);border-color:color-mix(in srgb,var(--warn) 24%,transparent)}
.app-url{grid-area:2/2/3/3;font-family:var(--ff-mono);font-size:.6875rem;color:var(--dim);overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.app-visits{grid-area:2/3/3/4;justify-self:end;font-size:.6875rem;color:var(--dim);white-space:nowrap;font-variant-numeric:tabular-nums}
.app-cta{grid-area:1/4/3/5;padding:5px 11px;border-radius:var(--r-sm);background:var(--accent);color:var(--on-accent);font-size:.6875rem;font-weight:550;white-space:nowrap}
.app-fold{background:color-mix(in srgb,var(--surface-2) 70%,var(--surface))}
.app-row.is-open{border-bottom:1px solid var(--line)}
.app-files{padding:10px 12px 12px}
.app-files-head{display:flex;align-items:center;gap:8px;margin-bottom:7px;font-size:.6875rem;font-weight:600;letter-spacing:.05em;text-transform:uppercase;color:var(--dim)}
.app-files-head .app-mono{text-transform:none;letter-spacing:0;font-weight:400}
.app-files-head .app-btn{margin-left:auto;padding:3px 8px;border:1px solid var(--line);border-radius:var(--r-sm);background:var(--surface);font-size:.6875rem;font-weight:500;letter-spacing:0;text-transform:none;color:var(--text-soft)}
.app-file{display:flex;align-items:center;justify-content:space-between;gap:10px;padding:5px 8px;border:1px solid var(--line);border-radius:var(--r-sm);background:var(--surface);font-size:.6875rem;color:var(--text-soft)}
.app-file+.app-file{margin-top:3px}
.app-mono{font-family:var(--ff-mono);font-variant-numeric:tabular-nums}
.frame-cap{display:flex;align-items:center;justify-content:space-between;gap:12px;padding:9px 2px 0;font-family:var(--ff-mono);font-size:.6875rem;color:var(--dim)}
.frame-cap{display:flex;align-items:center;justify-content:space-between;gap:12px;padding:9px 4px 2px;font-family:var(--ff-mono);font-size:.6875rem;color:var(--dim)}

.features-wrap{max-width:1200px;margin:0 auto;padding:0 clamp(14px,3vw,28px) clamp(48px,8vh,88px)}
.section-title{font-size:clamp(1.5rem,2.6vw,2.25rem);line-height:1.18;font-weight:600;letter-spacing:-.02em;max-width:22ch;margin-bottom:26px}
.features{display:grid;grid-template-columns:1.1fr .9fr;gap:12px;align-items:stretch}
.feature{
  display:flex;flex-direction:column;gap:12px;padding:22px;
  background:var(--surface);border:1px solid var(--line);border-radius:var(--r-lg);
}
.feature .feature-icon{display:flex;align-items:center;justify-content:center;width:30px;height:30px;border-radius:var(--r-sm);background:var(--surface-2);color:var(--accent)}
.feature .feature-icon svg{width:16px;height:16px}
.feature h3{font-size:var(--fs-md);font-weight:600}
.feature p{font-size:var(--fs-sm);line-height:1.6;color:var(--dim);max-width:38ch}
.feature:last-child{grid-column:1/-1;flex-direction:row;align-items:center;gap:16px}
.feature:last-child .feature-copy{min-width:0}
.feature-hero{grid-row:span 2;justify-content:flex-start;gap:16px;background:color-mix(in srgb,var(--accent) 7%,var(--surface))}
.feature-hero h3{font-size:var(--fs-lg)}
.term{
  margin-top:auto;border-radius:var(--r-md);overflow:hidden;
  background:color-mix(in srgb,var(--text) 92%,var(--accent));
  border:1px solid color-mix(in srgb,var(--text) 70%,transparent);
}
.term-bar{display:flex;gap:5px;padding:9px 11px;background:rgba(255,255,255,.07)}
.term-bar span{width:7px;height:7px;border-radius:50%;background:rgba(255,255,255,.32)}
.term-body{padding:12px 14px 16px;font-family:var(--ff-mono);font-size:.75rem;line-height:1.9;color:#e8f3ec}
.term-body b{color:#8de9c4;margin-right:5px}
.term-body .dim{color:rgba(232,243,236,.6)}
.term-body .ok{padding-left:12px;color:rgba(232,243,236,.82)}
.term-body .ok::before{content:"→ ";opacity:.55}

footer{border-top:1px solid var(--line);margin-top:8px}
.footer-inner{display:flex;align-items:center;justify-content:space-between;gap:14px;flex-wrap:wrap;max-width:1200px;margin:0 auto;padding:20px clamp(14px,3vw,28px) calc(24px + env(safe-area-inset-bottom));font-size:var(--fs-xs);color:var(--dim)}
.footer-inner .foot-links{display:flex;align-items:center;gap:14px}
.footer-inner a{color:var(--dim);transition:color var(--dur) var(--ease)}
@media(hover:hover) and (pointer:fine){.footer-inner a:hover{color:var(--accent)}}

@media(max-width:980px){
  .hero{grid-template-columns:minmax(0,1fr);gap:32px;min-height:0;padding-top:20px}
  .hero h1{max-width:14ch}
  .tagline{max-width:none}
  .features{grid-template-columns:minmax(0,1fr)}
  .feature-hero{grid-row:auto}
}
@media(max-width:640px){
  .brand{font-size:.9375rem}
  .hero h1{font-size:clamp(2.2rem,11vw,3rem)}
  .cta-row{display:grid;grid-template-columns:minmax(0,1fr)}
  .cta{width:100%}
  .app-window{border-radius:var(--r-md)}
  .app-cta{display:none}
  .frame-cap span:last-child{display:none}
  .frame-cap span:first-child{overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
  .feature{padding:18px}
  .footer-inner{padding-bottom:calc(20px + env(safe-area-inset-bottom))}
}
`

// uiGateCSS styles the password gate for a protected site. It uses the shared
// tokens only and stays dependency-free.
const uiGateCSS = `
body{
  display:grid;place-items:center;min-height:100dvh;padding:24px;
  background:radial-gradient(58% 46% at 50% 0%,color-mix(in srgb,var(--accent) 11%,transparent),transparent 68%),var(--bg);
}
.gate-card{
  width:100%;max-width:404px;padding:30px;background:var(--surface);
  border:1px solid var(--line);border-radius:var(--r-lg);box-shadow:var(--shadow-md);
}
.gate-badge{display:inline-flex;align-items:center;justify-content:center;width:40px;height:40px;margin-bottom:16px;border-radius:var(--r-md);background:var(--accent-soft);color:var(--accent)}
.gate-badge svg{width:19px;height:19px}
.eyebrow{font-family:var(--ff-mono);font-size:.6875rem;letter-spacing:.09em;text-transform:uppercase;color:var(--dim);margin-bottom:4px}
.gate-card h1{font-size:1.3125rem;line-height:1.3;font-weight:600;letter-spacing:-.01em;overflow-wrap:anywhere}
.gate-card h1 .site{color:var(--accent)}
.hint{margin:8px 0 20px;font-size:var(--fs-sm);line-height:1.65;color:var(--dim)}
.field{position:relative;display:flex;align-items:center}
.field input{
  width:100%;min-height:44px;padding:0 44px 0 13px;
  background:var(--field);color:var(--text);
  border:1px solid var(--line);border-radius:var(--r-md);font-size:var(--fs-md);
  transition:border-color var(--dur) var(--ease),box-shadow var(--dur) var(--ease);
}
.field input::placeholder{color:var(--dim-2)}
.field input:focus{outline:none;border-color:var(--accent);box-shadow:0 0 0 3px var(--accent-ring)}
.field input::-ms-reveal,.field input::-ms-clear{display:none}
.toggle{position:absolute;right:5px;top:50%;transform:translateY(-50%);display:flex;align-items:center;justify-content:center;width:32px;height:32px;padding:0;border:1px solid transparent;border-radius:var(--r-sm);color:var(--dim)}
.toggle svg{width:16px;height:16px}
@media(hover:hover) and (pointer:fine){.toggle:hover{color:var(--text);background:var(--surface-2)}}
.submit{
  width:100%;min-height:44px;margin-top:12px;padding:0 14px;
  border-radius:var(--r-md);background:var(--accent);color:var(--on-accent);
  font-size:var(--fs-md);font-weight:550;border:1px solid transparent;
  transition:background var(--dur) var(--ease),transform var(--dur) var(--ease),opacity var(--dur) var(--ease);
}
@media(hover:hover) and (pointer:fine){.submit:hover{background:var(--accent-hover)}}
.submit:active{transform:scale(.99)}
.submit[disabled]{opacity:.6;cursor:progress}
.gate-err{
  display:flex;align-items:flex-start;gap:9px;margin-bottom:14px;padding:10px 12px;
  font-size:var(--fs-sm);line-height:1.55;color:var(--danger);
  background:var(--danger-soft);border:1px solid color-mix(in srgb,var(--danger) 32%,var(--line));
  border-radius:var(--r-md);
}
.gate-err svg{width:15px;height:15px;flex-shrink:0;margin-top:1px;fill:none;stroke:currentColor}
.foot{margin-top:18px;padding-top:14px;border-top:1px solid var(--line);text-align:center;font-family:var(--ff-mono);font-size:.6875rem;color:var(--dim)}
@media(max-width:480px){body{padding:14px;align-items:start;padding-top:max(14px,7vh)}.gate-card{padding:24px 18px}}
`
