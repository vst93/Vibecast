package server

import (
	"fmt"
	"net/url"
	"strings"
)

// logoFavicon is the SVG logo as a data URI for use as favicon.
const logoFavicon = `<link rel="icon" href="data:image/svg+xml,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20viewBox%3D%220%200%20128%20128%22%3E%3Crect%20x%3D%228%22%20y%3D%228%22%20width%3D%22112%22%20height%3D%22112%22%20rx%3D%2224%22%20fill%3D%22%2300e5a0%22%2F%3E%3Cpolyline%20points%3D%2238%2C43%2064%2C97%2090%2C43%22%20fill%3D%22none%22%20stroke%3D%22%230a0a0c%22%20stroke-width%3D%2210%22%20stroke-linecap%3D%22round%22%20stroke-linejoin%3D%22round%22%2F%3E%3Ccircle%20cx%3D%2264%22%20cy%3D%2231%22%20r%3D%226%22%20fill%3D%22%230a0a0c%22%2F%3E%3C%2Fsvg%3E">`

// logoIcon is the inline SVG logo for navbar display.
const logoIcon = `<svg viewBox="0 0 128 128" width="20" height="20" style="vertical-align:middle;margin-right:6px;flex-shrink:0"><rect x="8" y="8" width="112" height="112" rx="24" fill="#00e5a0"/><polyline points="38,43 64,97 90,43" fill="none" stroke="#0a0a0c" stroke-width="10" stroke-linecap="round" stroke-linejoin="round"/><circle cx="64" cy="31" r="6" fill="#0a0a0c"/></svg>`

// logoIconLarge is the inline SVG logo for hero/auth display.
const logoIconLarge = `<svg viewBox="0 0 128 128" width="48" height="48" style="vertical-align:middle;margin-right:10px;flex-shrink:0"><rect x="8" y="8" width="112" height="112" rx="24" fill="#00e5a0"/><polyline points="38,43 64,97 90,43" fill="none" stroke="#0a0a0c" stroke-width="10" stroke-linecap="round" stroke-linejoin="round"/><circle cx="64" cy="31" r="6" fill="#0a0a0c"/></svg>`

// logoIconMedium is the inline SVG logo for auth card display.
const logoIconMedium = `<svg viewBox="0 0 128 128" width="32" height="32" style="vertical-align:middle;margin-right:8px;flex-shrink:0"><rect x="8" y="8" width="112" height="112" rx="24" fill="#00e5a0"/><polyline points="38,43 64,97 90,43" fill="none" stroke="#0a0a0c" stroke-width="10" stroke-linecap="round" stroke-linejoin="round"/><circle cx="64" cy="31" r="6" fill="#0a0a0c"/></svg>`

// projectRepoURL is the canonical source of the project. It is referenced from
// quiet "about" spots only (visitor footers, admin system tab), never as a
// promotional banner.
const projectRepoURL = "https://github.com/vst93/Vibecast"

// Icon constants - flat stroke style, 24x24 viewBox, stroke-width=2.
// Using display:inline-flex for proper vertical alignment inside flex containers
const iconMoon = `<svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="display:inline-flex;align-self:center;flex-shrink:0"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/></svg>`
const iconSun = `<svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="display:inline-flex;align-self:center;flex-shrink:0"><circle cx="12" cy="12" r="5"/><line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/><line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/></svg>`
const iconBuilding = `<svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="display:inline-flex;align-self:center;flex-shrink:0"><rect x="4" y="2" width="16" height="20" rx="2"/><line x1="9" y1="7" x2="9.01" y2="7"/><line x1="15" y1="7" x2="15.01" y2="7"/><line x1="9" y1="12" x2="9.01" y2="12"/><line x1="15" y1="12" x2="15.01" y2="12"/><line x1="9" y1="17" x2="9.01" y2="17"/><line x1="15" y1="17" x2="15.01" y2="17"/></svg>`
const iconLock = `<svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="display:inline-flex;align-self:center;flex-shrink:0"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>`
const iconFolder = `<svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="display:inline-flex;align-self:center;flex-shrink:0"><path d="M4 20h16a2 2 0 0 0 2-2V8a2 2 0 0 0-2-2h-7.93a2 2 0 0 1-1.66-.9l-.82-1.2a2 2 0 0 0-1.66-.9H4a2 2 0 0 0-2 2v13a2 2 0 0 0 2 2z"/></svg>`
const iconFile = `<svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="display:inline-flex;align-self:center;flex-shrink:0"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>`
const iconEye = `<svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="display:inline-flex;align-self:center;flex-shrink:0"><path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7z"/><circle cx="12" cy="12" r="3"/></svg>`
const iconCheck = `<svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="display:inline-flex;align-self:center;flex-shrink:0"><polyline points="20 6 9 17 4 12"/></svg>`
const iconLogout = `<svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="display:inline-flex;align-self:center;flex-shrink:0"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>`
const iconLayout = `<svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="display:inline-flex;align-self:center;flex-shrink:0"><rect x="3" y="3" width="18" height="18" rx="2"/><line x1="9" y1="3" x2="9" y2="21"/></svg>`

// Additional UI glyphs. Kept as constants so every page shares one icon set:
// 24x24 viewBox, 2px stroke, currentColor, rounded joins.
const iconTerminal = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="4 17 10 11 4 5"/><line x1="12" y1="19" x2="20" y2="19"/></svg>`
const iconServer = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="3" width="20" height="7" rx="2"/><rect x="2" y="14" width="20" height="7" rx="2"/><line x1="6" y1="6.5" x2="6.01" y2="6.5"/><line x1="6" y1="17.5" x2="6.01" y2="17.5"/></svg>`
const iconChevronRight = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 6 15 12 9 18"/></svg>`
const iconDots = `<svg viewBox="0 0 24 24" fill="currentColor" stroke="none"><circle cx="12" cy="5" r="1.7"/><circle cx="12" cy="12" r="1.7"/><circle cx="12" cy="19" r="1.7"/></svg>`
const iconExternal = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 4h6v6"/><path d="M20 4 10 14"/><path d="M18 14v5a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h5"/></svg>`
const iconCopy = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="12" height="12" rx="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg>`
const iconShare = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="18" cy="5" r="3"/><circle cx="6" cy="12" r="3"/><circle cx="18" cy="19" r="3"/><line x1="8.6" y1="13.5" x2="15.4" y2="17.5"/><line x1="15.4" y1="6.5" x2="8.6" y2="10.5"/></svg>`
const iconDownload = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 3v12"/><path d="m7 11 5 5 5-5"/><path d="M5 21h14"/></svg>`
const iconTrash = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 6h18"/><path d="M8 6V4h8v2"/><path d="m5 6 1 15h12l1-15"/><path d="M10 10v7"/><path d="M14 10v7"/></svg>`
const iconEdit = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"/><path d="M16.5 3.5a2.1 2.1 0 0 1 3 3L7 19l-4 1 1-4Z"/></svg>`
const iconUsers = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9.5" cy="7" r="4"/><path d="M22 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>`
const iconGlobe = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="9"/><path d="M3.6 9h16.8M3.6 15h16.8"/><path d="M12 3a15 15 0 0 1 0 18a15 15 0 0 1 0-18Z"/></svg>`
const iconShield = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 3 5 6v6c0 4.2 3 7.6 7 9 4-1.4 7-4.8 7-9V6Z"/><path d="m9 12 2 2 4-4"/></svg>`
const iconSearch = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="7"/><path d="m20 20-3.5-3.5"/></svg>`
const iconX = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"><path d="M18 6 6 18M6 6l12 12"/></svg>`
const iconGithub = `<svg viewBox="0 0 24 24" fill="currentColor"><path d="M12 2a10 10 0 0 0-3.16 19.49c.5.09.68-.22.68-.48l-.01-1.7c-2.78.6-3.37-1.34-3.37-1.34-.45-1.16-1.11-1.47-1.11-1.47-.9-.62.07-.6.07-.6 1 .07 1.53 1.03 1.53 1.03.89 1.52 2.34 1.08 2.91.83.09-.65.35-1.09.63-1.34-2.22-.25-4.56-1.11-4.56-4.95 0-1.09.39-1.98 1.03-2.68-.1-.25-.45-1.27.1-2.65 0 0 .84-.27 2.75 1.02a9.5 9.5 0 0 1 5 0c1.91-1.29 2.75-1.02 2.75-1.02.55 1.38.2 2.4.1 2.65.64.7 1.03 1.59 1.03 2.68 0 3.85-2.35 4.69-4.58 4.94.36.31.68.92.68 1.85l-.01 2.75c0 .27.18.58.69.48A10 10 0 0 0 12 2Z"/></svg>`
const iconInbox = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M3 12h5l2 3h4l2-3h5"/><path d="M5.5 5h13l2.5 7v5a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-5Z"/></svg>`

// iconLockLarge is the oversized lock used as the password gate badge.
const iconLockLarge = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="11" width="18" height="11" rx="2.5"/><path d="M7.5 11V7a4.5 4.5 0 0 1 9 0v4"/></svg>`

// iconAlert is the inline warning glyph used by the password gate error block.
const iconAlert = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="9"/><line x1="12" y1="8" x2="12" y2="13"/><line x1="12" y1="16.5" x2="12.01" y2="16.5"/></svg>`

// pwdIconHidden / pwdIconShown are the password visibility toggle glyphs:
// slashed eye while the value is hidden, open eye while it is visible.
const pwdIconHidden = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 10 8 10 8a18.5 18.5 0 0 1-2.16 3.19M6.06 6.06A18.45 18.45 0 0 0 2 12s3 8 10 8a10.07 10.07 0 0 0 5.94-1.94"/><path d="M14.12 14.12a3 3 0 1 1-4.24-4.24"/><line x1="2" y1="2" x2="22" y2="22"/></svg>`
const pwdIconShown = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2 12s3-8 10-8 10 8 10 8-3 8-10 8-10-8-10-8z"/><circle cx="12" cy="12" r="3"/></svg>`

// landingPageHTML is the public landing page.
var landingPageHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
` + logoFavicon + `
<meta name="description" content="Self-hosted static site hosting: upload files or a ZIP, get a stable URL with optional password protection.">
<meta name="theme-color" content="#f2f5f3" media="(prefers-color-scheme: light)">
<meta name="theme-color" content="#0b1210" media="(prefers-color-scheme: dark)">
<meta property="og:type" content="website">
<meta property="og:title" content="Vibecast - self-hosted static site hosting">
<meta property="og:description" content="Upload files or a ZIP, get a stable URL with optional password protection. One Go binary, SQLite, no Nginx.">
<meta property="og:image" content="assets/dashboard.png">
<meta name="twitter:card" content="summary_large_image">
<title>Vibecast - Self-hosted static site hosting</title>
<style>
` + uiTokensCSS + uiComponentsCSS + uiLandingCSS + `
</style>
<script>var _t="light";try{_t=localStorage.getItem("theme")||"light"}catch(e){}document.documentElement.setAttribute("data-theme",_t)</script>
</head>
<body>
<nav class="site-nav"><a href="/" class="brand">` + logoIconMedium + `<span class="brand-text">Vibe<span class="accent">cast</span></span></a><div class="nav-actions"><div class="lang-toggle"><a id="theme-toggle" onclick="toggleTheme()" aria-label="Toggle theme">` + iconMoon + `</a><a id="langEn" class="active" onclick="setLang('en')">EN</a><a id="langZh" onclick="setLang('zh')">中文</a></div></div></nav>
<main>
<section class="hero">
<div class="hero-copy"><div class="eyebrow" data-i="eyebrow">SELF-HOSTED STATIC PUBLISHING</div><h1 data-i="headline">Your files.<br><span>Live instantly.</span></h1><p class="tagline" data-i="tagline">Deploy static sites and shareable files from one compact Go service.</p><div class="cta-row"><a href="dashboard" class="cta cta-primary" data-i="getStarted">Open Dashboard</a><a href="#features" class="cta cta-secondary" data-i="learnMore">See how it works</a></div></div>
<div class="hero-visual">
<div class="app-window" role="img" aria-label="Vibecast dashboard preview: three deployed sites with public and protected status plus visit counts, and an expanded file list">
<div class="app-bar"><span class="app-dots"><i></i><i></i><i></i></span><span class="app-addr">sites.vibecast.app/dashboard</span></div>
<div class="app-body">
<div class="app-tabs"><span class="on" data-i="mTabMy">My Sites</span><span data-i="mTabOrg">Org Pinned</span></div>
<div class="app-row"><span class="dot pub"></span><b>Launch Notes</b><span class="chip pub" data-i="mPublic">Public</span><span class="app-url">/s/dopx19518lhm/</span><span class="app-visits" data-i="mVisitsA">132 visits</span><span class="app-cta" data-i="mUpload">Upload</span></div>
<div class="app-row"><span class="dot pro"></span><b>Studio Portfolio</b><span class="chip pro" data-i="mProtected">Protected</span><span class="app-url">/s/5s9ffg99o60a/</span><span class="app-visits" data-i="mVisitsB">86 visits</span><span class="app-cta" data-i="mUpload">Upload</span></div>
<div class="app-fold">
<div class="app-row is-open"><span class="dot pub"></span><b>Northstar Docs</b><span class="chip pub" data-i="mPublic">Public</span><span class="app-url">/s/ijzz8ddzj6xh/</span><span class="app-visits" data-i="mVisitsC">326 visits</span><span class="app-cta" data-i="mUpload">Upload</span></div>
<div class="app-files">
<div class="app-files-head"><span data-i="mFiles">Files</span><span class="app-mono" data-i="mFilesMeta">4 files · 212 KB</span></div>
<div class="app-file app-mono"><span>index.html</span><span>63 B</span></div>
<div class="app-file app-mono"><span>assets/app.css</span><span>12 KB</span></div>
<div class="app-file app-mono"><span>assets/guide.pdf</span><span>200 KB</span></div>
</div>
</div>
</div>
</div>
<div class="frame-cap"><span data-i="frameCap">Dashboard · sites, files and visit counts</span><span>Go + SQLite</span></div>
</div>
</section>
<section class="features-wrap" id="features"><h2 class="section-title" data-i="sectionTitle">A small server with the essentials built in.</h2><div class="features">
<div class="feature feature-hero"><div class="feature-icon">` + iconTerminal + `</div><div class="term" aria-hidden="true"><div class="term-bar"><span></span><span></span><span></span></div><div class="term-body"><div><b>$</b> vibecast up ./site</div><div class="ok">zipped  12 files &middot; 340 KB</div><div class="ok">pushed  site/&lt;slug&gt;/</div><div class="ok">live    https://&lt;slug&gt;.sites.example</div></div></div><div class="hero-copy"><h3 data-i="feat1Title">Instant Deploy</h3><p data-i="feat1Desc">Upload a ZIP, get a live URL in seconds.</p></div></div>
<div class="feature"><div class="feature-icon">` + iconLock + `</div><div class="feature-copy"><h3 data-i="feat2Title">Password Protect</h3><p data-i="feat2Desc">Keep your site private with password gating.</p></div></div>
<div class="feature"><div class="feature-icon">` + iconServer + `</div><div class="feature-copy"><h3 data-i="feat3Title">Zero Nginx</h3><p data-i="feat3Desc">Pure Go application server. No dependencies.</p></div></div>
<div class="feature"><div class="feature-icon">` + iconFolder + `</div><div class="feature-copy"><h3 data-i="feat4Title">Self-Hosted</h3><p data-i="feat4Desc">Your data, your rules. SQLite + filesystem.</p></div></div>
</div></section>
</main>
<footer><div class="footer-inner"><span data-i="footer">Vibecast · self-hosted static site hosting</span><span class="foot-links"><a href="dashboard" data-i="footDashboard">Dashboard</a><a href="https://github.com/earendil-works/vibecast" target="_blank" rel="noopener">GitHub</a></span></div></footer>
<script>
var i18n={en:{eyebrow:"SELF-HOSTED STATIC PUBLISHING",headline:"Your files. Live instantly.",tagline:"Deploy static sites and shareable files from one compact Go service.",getStarted:"Open Dashboard",learnMore:"See how it works",sectionTitle:"A small server with the essentials built in.",feat1Title:"Instant deploy",feat1Desc:"Upload a ZIP or a single file and get a stable URL in seconds.",feat2Title:"Private when needed",feat2Desc:"Add password protection or share access with your organization.",feat3Title:"One Go service",feat3Desc:"Authentication, deployment, updates and file serving without Nginx.",feat4Title:"Your infrastructure",feat4Desc:"Keep sites in your filesystem and metadata in SQLite.",footer:"Vibecast - self-hosted static site hosting",frameCap:"Dashboard · sites, files and visit counts",mTabMy:"My Sites",mTabOrg:"Org Pinned",mPublic:"Public",mProtected:"Protected",mVisitsA:"132 visits",mVisitsB:"86 visits",mVisitsC:"326 visits",mUpload:"Upload",mFiles:"Files",mFilesMeta:"4 files · 212 KB",mZip:"Download ZIP",footDashboard:"Dashboard"},zh:{eyebrow:"自托管静态内容发布",headline:"文件上传，即刻上线。",tagline:"用一个轻量 Go 服务部署静态站点和可分享文件。",getStarted:"打开控制台",learnMore:"了解工作方式",sectionTitle:"小巧的服务，完整的发布能力。",feat1Title:"即时部署",feat1Desc:"上传 ZIP 或单个文件，数秒内获得稳定访问地址。",feat2Title:"按需保护",feat2Desc:"设置访问密码，或向同一组织的成员开放。",feat3Title:"单一 Go 服务",feat3Desc:"认证、部署、更新和文件服务，无需 Nginx。",feat4Title:"掌控基础设施",feat4Desc:"站点保存在文件系统，元数据保存在 SQLite。",footer:"Vibecast - 自托管静态站点托管平台",frameCap:"控制台 · 站点、文件与访问统计",mTabMy:"我的站点",mTabOrg:"组织钉选",mPublic:"公开",mProtected:"已保护",mVisitsA:"132 次访问",mVisitsB:"86 次访问",mVisitsC:"326 次访问",mUpload:"上传",mFiles:"文件",mFilesMeta:"4 个文件 · 212 KB",mZip:"下载 ZIP",footDashboard:"控制台"}};
function getTheme(){return document.documentElement.getAttribute("data-theme")||"light"}
function setTheme(t){document.documentElement.setAttribute("data-theme",t);try{localStorage.setItem("theme",t)}catch(e){}var b=document.getElementById("theme-toggle");if(b)b.innerHTML=t==="dark"?'` + iconSun + `':'` + iconMoon + `'}
function toggleTheme(){setTheme(getTheme()==="dark"?"light":"dark")}
var lang="en";
function setLang(l){lang=l;document.documentElement.lang=l;document.querySelectorAll("[data-i]").forEach(function(e){var k=e.getAttribute("data-i");if(i18n[l][k]){if(k==="headline"){var p=i18n[l][k].split(l==="zh"?"，":". ");e.innerHTML=escText(p[0])+(l==="zh"?"，<br><span>":".<br><span>")+escText(p.slice(1).join(l==="zh"?"，":". "))+"</span>"}else{e.textContent=i18n[l][k]}}});document.getElementById("langEn").className=l==="en"?"active":"";document.getElementById("langZh").className=l==="zh"?"active":"";try{localStorage.setItem("lang",l)}catch(e){}}
function escText(s){return String(s||"").replace(/[&<>]/g,function(c){return{"&":"&amp;","<":"&lt;",">":"&gt;"}[c]})}
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
<style>
` + uiTokensCSS + uiComponentsCSS + uiDashboardCSS + `
</style>
<script>var _t="light";try{_t=localStorage.getItem("theme")||"light"}catch(e){}document.documentElement.setAttribute("data-theme",_t)</script>
</head>
<body>
<div id="app"></div>
<script>
var BASE=location.pathname.replace(/\/(?:dashboard|admin)\/?$/,"/")||"/";var API=BASE+"api";
var currentUser=null;
function getTheme(){return document.documentElement.getAttribute("data-theme")||"light"}
function setTheme(t){document.documentElement.setAttribute("data-theme",t);try{localStorage.setItem("theme",t)}catch(e){}var b=document.getElementById("theme-toggle");if(b)b.innerHTML=t==="dark"?'` + iconSun + `':'` + iconMoon + `'}
function toggleTheme(){setTheme(getTheme()==="dark"?"light":"dark")}
var lang="en";
var sitePage=1,sitePerPage=10,siteTotal=0,siteSearch="";
var maxUploadMB=50;
var maxSites=30;
var siteTab="my"; // "my" or "org"
var orgPinnedPage=1,orgPinnedTotal=0,orgPinnedSearch="";
var i18n={en:{siteName:"Site Name",siteNamePh:"e.g. My Portfolio",siteDesc:"Description",siteDescPh:"Brief description (optional)",slug:"URL Slug",slugPh:"my-portfolio",sitePwd:"Access Password",sitePwdPh:"Leave empty for public",create:"Create Site",yourSites:"Your Sites",noSites:"No sites yet. Create one above.",deployBtn:"Upload",delete:"Delete",deleteConfirm:"Delete this site? This removes all files.",protected:"Protected",public:"Public",login:"Login",register:"Register",email:"Email",emailPh:"you@example.com",password:"Password",pwdHint:"At least 6 characters",noAccount:"No account?",haveAccount:"Have an account?",logout:"Logout",adminPanel:"Admin",deployed:"Deployed!",siteCreated:"Site created",deleted:"Deleted",loginFailed:"Login failed",registerFailed:"Registration failed",slugDesc:"Auto-generated from name if blank. a-z, 0-9, hyphens only.",pwdDesc:"If set, visitors need this password.",sitesHint:"Click to expand",deployHint:"Upload ZIP or single file (PDF, Word, images, etc.)",storagePath:"Storage",accessPassword:"Password",none:"None",accessDisabled:"Disabled",pwdRequired:"Public access disabled - password required",pwdOptional:"Optional password protection",search:"Search",searchSitesPh:"Search sites...",prev:"Prev",next:"Next",page:"Page",of:"of",captcha:"Captcha",captchaPh:"Answer",captchaLabel:"Verification",confirmPassword:"Confirm Password",pwdMismatch:"Passwords do not match",changePassword:"Change Password",currentPassword:"Current Password",newPassword:"New Password",newPasswordPh:"New password (min 6)",passwordChanged:"Password changed",cancel:"Cancel",save:"Save",confirm:"Confirm",emailRequired:"Please enter your email",pwdRequired:"Please enter your password",confirmRequired:"Please confirm your password",captchaRequired:"Please solve the captcha",copyUrl:"Copy URL",copied:"Copied!",copy:"Copy",copyFailed:"Copy failed",visit:"Visit",files:"Files",noFiles:"No files",visitsToday:"Today",visitsMonth:"Month",visitsTotal:"Total",visits:"Visits",uploading:"Uploading...",dragDrop:"Drag file here or click to upload",dragDropHover:"Drop to upload",fileTooLarge:"File too large",fileTypeBlocked:"This file type is not allowed",uploadFailed:"Upload failed",networkError:"Network error",allowedTypes:"ZIP, PDF, Word, Excel, PPT, images, audio, video, text",share:"Share",shareCopied:"Share text copied!",shareTemplatePwd:"Site: {name}\nURL: {url}\nPassword: {password}",shareTemplatePublic:"Site: {name}\nURL: {url}",organization:"Organization",orgManage:"Organization",orgCreate:"Create Organization",orgJoin:"Join Organization",orgName:"Organization Name",orgNamePh:"e.g. My Team",orgInviteCode:"Invite Code",orgInviteCodePh:"Enter 12-char code",orgCreateDesc:"Create a new organization. You will get an invite code to share with your team.",orgJoinDesc:"Enter an invite code to join an existing organization.",orgLeave:"Leave Organization",orgDelete:"Delete Organization",orgMembers:"Members",orgNoMembers:"No members",orgOwner:"Owner",orgMember:"Member",orgRemove:"Remove",orgLeaveConfirm:"Leave this organization?",orgDeleteConfirm:"Delete this organization? This cannot be undone.",orgCreated:"Organization created",orgJoined:"Joined organization",orgLeft:"Left organization",orgDeleted:"Organization deleted",orgMemberRemoved:"Member removed",orgInviteCodeCopied:"Invite code copied!",orgOpenLabel:"Open to org members",orgOpenDesc:"If enabled, logged-in users in the same organization can access this site without a password.",orgBadge:"Org",orgFeatureDesc:"Organizations let you group users together. When a site is set to open to org members, anyone in the same org can access it without a password. Each user can only belong to one organization.",noOrg:"Not in an organization",orgPinnedLabel:"Pin to org list",orgPinnedDesc:"If enabled, this site appears in the org pinned sites tab. Org members can see its name and visit it, but cannot view its config or password.",tabMySites:"My Sites",tabOrgPinned:"Org Pinned",noOrgPinned:"No pinned sites in your organization.",orgPinnedBadge:"Pinned",noOrgForPinned:"Join an organization to see pinned sites.",editSite:"Edit Site",editName:"Site Name",editDesc:"Description",editSiteDesc:"Edit site name and description"},zh:{siteName:"站点名称",siteNamePh:"例如：我的作品集",siteDesc:"描述",siteDescPh:"简短描述（可选）",slug:"URL Slug",slugPh:"my-portfolio",sitePwd:"访问密码",sitePwdPh:"留空则公开访问",create:"创建站点",yourSites:"我的站点",noSites:"还没有站点，在左侧创建一个。",deployBtn:"上传",delete:"删除",deleteConfirm:"确定删除此站点？所有文件将被移除。",protected:"已保护",public:"公开",login:"登录",register:"注册",email:"邮箱",emailPh:"you@example.com",password:"密码",pwdHint:"至少 6 个字符",noAccount:"没有账号？",haveAccount:"已有账号？",logout:"退出",adminPanel:"管理",deployed:"部署成功！",siteCreated:"站点已创建",deleted:"已删除",loginFailed:"登录失败",registerFailed:"注册失败",slugDesc:"留空则自动生成。仅限 a-z、0-9、连字符。",pwdDesc:"设置后，访问者需要输入此密码。",sitesHint:"点击展开详情",deployHint:"上传 ZIP 压缩包或单个文件（PDF、Word、图片等）",storagePath:"存储路径",accessPassword:"访问密码",none:"无",accessDisabled:"已禁用",pwdRequired:"公开访问已关闭 - 必须设置密码",pwdOptional:"可选的密码保护",search:"搜索",searchSitesPh:"搜索站点...",prev:"上一页",next:"下一页",page:"第",of:"/ 共",captcha:"验证码",captchaPh:"输入答案",captchaLabel:"验证",confirmPassword:"确认密码",pwdMismatch:"两次密码不一致",changePassword:"修改密码",currentPassword:"当前密码",newPassword:"新密码",newPasswordPh:"新密码（至少 6 位）",passwordChanged:"密码已修改",cancel:"取消",save:"保存",confirm:"确认",emailRequired:"请输入邮箱",pwdRequired:"请输入密码",confirmRequired:"请确认密码",captchaRequired:"请输入验证码",copyUrl:"复制链接",copied:"已复制！",copy:"复制",copyFailed:"复制失败",visit:"访问",files:"文件",noFiles:"暂无文件",visitsToday:"今日",visitsMonth:"本月",visitsTotal:"总计",visits:"访问",uploading:"上传中...",dragDrop:"拖拽文件到此处或点击上传",dragDropHover:"松开以上传",fileTooLarge:"文件过大",fileTypeBlocked:"不允许此文件类型",uploadFailed:"上传失败",networkError:"网络错误",allowedTypes:"ZIP、PDF、Word、Excel、PPT、图片、音频、视频、文本",share:"分享",shareCopied:"分享文本已复制！",shareTemplatePwd:"站点：{name}\n地址：{url}\n密码：{password}",shareTemplatePublic:"站点：{name}\n地址：{url}",organization:"组织",orgManage:"组织管理",orgCreate:"创建组织",orgJoin:"加入组织",orgName:"组织名称",orgNamePh:"例如：我的团队",orgInviteCode:"邀请码",orgInviteCodePh:"输入 12 位邀请码",orgCreateDesc:"创建一个新组织，你将获得一个邀请码来分享给团队成员。",orgJoinDesc:"输入邀请码加入已有组织。",orgLeave:"退出组织",orgDelete:"删除组织",orgMembers:"成员",orgNoMembers:"暂无成员",orgOwner:"创建者",orgMember:"成员",orgRemove:"移除",orgLeaveConfirm:"确定退出此组织？",orgDeleteConfirm:"确定删除此组织？此操作不可撤销。",orgCreated:"组织已创建",orgJoined:"已加入组织",orgLeft:"已退出组织",orgDeleted:"组织已删除",orgMemberRemoved:"成员已移除",orgInviteCodeCopied:"邀请码已复制！",orgOpenLabel:"对组织成员开放",orgOpenDesc:"开启后，同一组织内已登录的用户可以直接访问此站点，无需密码。",orgBadge:"组织",orgFeatureDesc:"组织用于将用户分组。当站点设置为对组织成员开放时，同一组织内的已登录用户无需密码即可访问。每个用户只能加入一个组织。",noOrg:"未加入任何组织",orgPinnedLabel:"钉到组织列表",orgPinnedDesc:"开启后，此站点会出现在组织钉选列表中。组织成员可以看到名称并访问，但无法查看配置和密码。",tabMySites:"我的站点",tabOrgPinned:"组织钉选",noOrgPinned:"组织中暂无钉选站点。",orgPinnedBadge:"已钉选",noOrgForPinned:"加入组织后可查看钉选站点。",editSite:"编辑站点",editName:"站点名称",editDesc:"描述",editSiteDesc:"编辑站点名称和描述"}};
Object.assign(i18n.en,{close:"Close",downloadZip:"Download ZIP",clearFiles:"Clear files",clearFilesConfirm:"Clear all {count} files ({size}) from this site? This cannot be undone.",filesCleared:"Site files cleared",noFilesToClear:"No files to clear",preparingDownload:"Preparing ZIP...",downloadingFile:"Downloading {size}",downloadReady:"ZIP download ready",downloadFailed:"Download failed",uploadPreparing:"Preparing upload...",uploadProcessing:"Processing files and checking site size...",uploadComplete:"Upload complete",fileSummary:"{count} files · {size}",deployHintDetailed:"Single files and ZIP archives must be no larger than {size} MB. Uploads overwrite matching file paths and keep all other files. ZIP contents are merged after extraction. The final total site size must also stay within {size} MB."});
Object.assign(i18n.zh,{close:"关闭",downloadZip:"下载 ZIP",clearFiles:"清空文件",clearFilesConfirm:"确定清空该站点的全部 {count} 个文件（{size}）吗？此操作不可撤销。",filesCleared:"站点文件已清空",noFilesToClear:"没有可清空的文件",preparingDownload:"正在打包 ZIP...",downloadingFile:"正在下载 {size}",downloadReady:"ZIP 下载已就绪",downloadFailed:"下载失败",uploadPreparing:"正在准备上传...",uploadProcessing:"正在处理文件并检查站点总大小...",uploadComplete:"上传完成",fileSummary:"{count} 个文件 · {size}",deployHintDetailed:"单个文件和 ZIP 压缩包均不得超过 {size} MB。上传只覆盖路径相同的文件，其他文件保持不变；ZIP 解压后按路径合并。合并后的站点文件总大小也不得超过 {size} MB。"});
Object.assign(i18n.en,{menu:"Menu",expandCreate:"Expand",collapseCreate:"Collapse"});
Object.assign(i18n.zh,{menu:"菜单",expandCreate:"展开",collapseCreate:"收起"});
Object.assign(i18n.en,{toggleTheme:"Toggle theme",copyPassword:"Copy password",openSite:"Open site in a new tab",siteActions:"Site actions",visitsCount:"{count} visits",visitsOne:"1 visit",noMatch:"No sites match your search",noMatchHint:"Try a shorter name, or clear the search field.",noSitesTitle:"No sites yet",noSitesHint:"Create your first site in the panel, then upload a ZIP or a single file.",clearSearch:"Clear search",deleteFileConfirm:"Delete {name}? Only this file is removed from the site.",fileDeleted:"File deleted",orgPinnedJoinHint:"Create or join an organization to share sites and see pinned ones.",orgCta:"Create or join",orgPinnedEmptyTitle:"Nothing pinned yet",showPassword:"Show password",hidePassword:"Hide password",updateReady:"Update ready",statSitesNote:"of {max} allowed",statUsersNote:"registered accounts",statAdminsNote:"with admin access"});
Object.assign(i18n.zh,{toggleTheme:"切换主题",copyPassword:"复制密码",openSite:"在新标签页打开站点",siteActions:"站点操作",visitsCount:"{count} 次访问",visitsOne:"1 次访问",noMatch:"没有匹配的站点",noMatchHint:"换一个更短的关键词，或清空搜索框。",noSitesTitle:"还没有站点",noSitesHint:"先在左侧面板创建站点，然后上传 ZIP 包或单个文件。",clearSearch:"清空搜索",deleteFileConfirm:"确定删除 {name} 吗？仅移除该文件，站点其他内容不变。",fileDeleted:"文件已删除",orgPinnedJoinHint:"创建或加入组织后，即可共享站点并查看钉选内容。",orgCta:"创建或加入",orgPinnedEmptyTitle:"还没有钉选站点",showPassword:"显示密码",hidePassword:"隐藏密码",updateReady:"更新已就绪",statSitesNote:"上限 {max} 个",statUsersNote:"已注册账号",statAdminsNote:"拥有管理权限"});
function t(k){return(i18n[lang]||i18n.en)[k]||(i18n.en[k]||k)}
function setLang(l){lang=l;try{localStorage.setItem("lang",l)}catch(e){}if(currentUser)renderDashboard();else renderAuth();}
function esc(s){return String(s||"").replace(/[&<>"']/g,function(c){return{"&":"&amp;","<":"&lt;",">":"&gt;",'"':"&quot;","'":"&#39;"}[c]})}
function siteUrl(u){return BASE+(u||"").replace(/^\//,"")}
function toast(msg,type){
type=type||"success";
var host=document.getElementById("toast-host");
if(!host){host=document.createElement("div");host.id="toast-host";host.setAttribute("aria-live","polite");host.setAttribute("aria-atomic","false");document.body.appendChild(host)}
var el=document.createElement("div");el.className="toast "+type+" show";el.textContent=msg;host.appendChild(el);
var ms=type==="error"?5000:2600;
setTimeout(function(){el.classList.remove("show");setTimeout(function(){el.remove()},200)},ms)
}
function togglePwd(btn){var inp=btn.parentElement.querySelector("input");if(!inp)return;var show=inp.type=="password";inp.type=show?"text":"password";btn.innerHTML=show?PWD_SHOW_ICON:PWD_HIDE_ICON;btn.setAttribute("aria-label",t(show?"hidePassword":"showPassword"));btn.setAttribute("aria-pressed",show?"true":"false")}
var PWD_SHOW_ICON='<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7z"/><circle cx="12" cy="12" r="3"/></svg>';
var PWD_HIDE_ICON='<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-10-8-10-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 10 8 10 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/><line x1="1" y1="1" x2="23" y2="23"/></svg>';
// Authentication is carried by the HttpOnly vibecast_session cookie.
// Keep these helpers only to remove tokens left by older frontend versions.
function getToken(){return"cookie"}
function setToken(tk){try{localStorage.removeItem("vibecast_token")}catch(e){}}
function clearToken(){try{localStorage.removeItem("vibecast_token")}catch(e){}}
function customConfirm(msg,onConfirm,opts){
opts=opts||{};
var o=document.createElement("div");o.className="modal-overlay";o.style.zIndex="320";
var danger=opts.danger!==false;
var title=opts.title||t("confirm");
o.innerHTML='<div class="modal" role="alertdialog" aria-modal="true" aria-labelledby="cf-title" aria-describedby="cf-desc"><h3 id="cf-title">'+esc(title)+'</h3><p id="cf-desc">'+esc(msg)+'</p><div class="modal-actions"><button type="button" class="btn btn-outline" id="cf-no">'+t("cancel")+'</button><button type="button" class="btn '+(danger?"btn-danger":"btn-primary")+'" id="cf-yes">'+(opts.confirmLabel||t("confirm"))+'</button></div></div>';
document.body.appendChild(o);
var prevFocus=document.activeElement;
setTimeout(function(){o.classList.add("show")},10);
var closed=false;
function close(){if(closed)return;closed=true;document.removeEventListener("keydown",onKey);o.classList.remove("show");setTimeout(function(){o.remove()},180);if(prevFocus&&prevFocus.focus)try{prevFocus.focus()}catch(e){}}
function onKey(e){if(e.key==="Escape"){e.preventDefault();close()}}
document.addEventListener("keydown",onKey);
o.addEventListener("click",function(e){if(e.target===o)close()});
var no=o.querySelector("#cf-no"),yes=o.querySelector("#cf-yes");
no.addEventListener("click",close);
yes.addEventListener("click",function(){close();onConfirm()});
setTimeout(function(){no.focus()},30)
}
function api(path,opts){
opts=opts||{};
	var headers={"Content-Type":"application/json","Accept-Language":lang||"en"};
	if(opts.headers)Object.assign(headers,opts.headers);
	return fetch(API+path,Object.assign({},opts,{headers:headers,credentials:"same-origin"})).then(function(r){
	if(r.status===401){clearToken()}
return r.json().catch(function(){return{error:"network error"}}).then(function(data){if(!r.ok)throw new Error(data.error||"request failed");return data})
})
}
function checkAuth(){return api("/auth/me").then(function(d){currentUser=d.data;return true}).catch(function(){clearToken();return false})}
var loginCaptchaId="",registerCaptchaId="";var regOpen=true;
function loadCaptcha(v){
return fetch(BASE+"api/auth/captcha").then(function(r){return r.json()}).then(function(d){
if(v==="login"){loginCaptchaId=d.data.id;var el=document.getElementById("login-captcha-img");if(el)el.innerHTML=d.data.image}
else{registerCaptchaId=d.data.id;var el=document.getElementById("register-captcha-img");if(el)el.innerHTML=d.data.image}
}).catch(function(){})
}
function renderAuth(){
var lh='<div class="auth-utility"><button type="button" id="theme-toggle" class="btn-icon" onclick="toggleTheme()" aria-label="'+t("toggleTheme")+'">'+(getTheme()==="dark"?'` + iconSun + `':'` + iconMoon + `')+'</button><div class="lang-toggle" role="group" aria-label="Language"><button type="button" class="'+(lang==="en"?"active":"")+'" onclick="setLang(\'en\')">EN</button><button type="button" class="'+(lang==="zh"?"active":"")+'" onclick="setLang(\'zh\')">中文</button></div></div>';
document.getElementById("app").innerHTML='<div class="auth-screen">'+lh+'<div class="auth-card"><h1>` + logoIconMedium + `<span class="brand-text">Vibe<span class="accent">cast</span></span></h1><p class="subtitle">Build with vibe. Cast instantly.</p><div id="auth-form"></div></div></div>';
fetch(BASE+"api/settings").then(function(r){return r.json()}).then(function(d){regOpen=d.data&&d.data.openRegistration!==false;showLogin()}).catch(function(){regOpen=true;showLogin()});
}
function showLogin(){
document.getElementById("auth-form").innerHTML='<form onsubmit="doLogin();return false"><div class="auth-field"><label>'+t("email")+'</label><input id="email" type="email" placeholder="'+t("emailPh")+'" autocomplete="email"></div><div class="auth-field"><label>'+t("password")+'</label><div class="pwd-wrap"><input id="password" type="password" placeholder="'+t("password")+'" autocomplete="current-password"><button type="button" class="pwd-toggle" onclick="togglePwd(this)" aria-label="'+t("showPassword")+'" aria-pressed="false"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-10-8-10-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 10 8 10 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/><line x1="1" y1="1" x2="23" y2="23"/></svg></button></div></div><div class="auth-field"><label class="captcha-label">'+t("captchaLabel")+'</label><div class="captcha-row"><div class="captcha-img" id="login-captcha-img" onclick="loadCaptcha(\'login\')" title="'+t("captcha")+'"></div><input id="login-captcha" type="text" placeholder="'+t("captchaPh")+'"></div></div><button class="btn btn-primary" type="submit">'+t("login")+'</button></form>'+(regOpen?'<p class="switch">'+t("noAccount")+' <a onclick="showRegister()">'+t("register")+'</a></p>':'');
loadCaptcha("login");
}
function showRegister(){
document.getElementById("auth-form").innerHTML='<form onsubmit="doRegister();return false"><div class="auth-field"><label>'+t("email")+'</label><input id="email" type="email" placeholder="'+t("emailPh")+'" autocomplete="email"></div><div class="auth-field"><label>'+t("password")+'</label><div class="pwd-wrap"><input id="password" type="password" placeholder="'+t("password")+'" autocomplete="new-password"><button type="button" class="pwd-toggle" onclick="togglePwd(this)" aria-label="'+t("showPassword")+'" aria-pressed="false"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-10-8-10-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 10 8 10 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/><line x1="1" y1="1" x2="23" y2="23"/></svg></button></div><div class="field-hint">'+t("pwdHint")+'</div></div><div class="auth-field"><label>'+t("confirmPassword")+'</label><div class="pwd-wrap"><input id="confirm-pwd" type="password" placeholder="'+t("confirmPassword")+'" autocomplete="new-password"><button type="button" class="pwd-toggle" onclick="togglePwd(this)" aria-label="'+t("showPassword")+'" aria-pressed="false"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-10-8-10-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 10 8 10 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/><line x1="1" y1="1" x2="23" y2="23"/></svg></button></div></div><div class="auth-field"><label class="captcha-label">'+t("captchaLabel")+'</label><div class="captcha-row"><div class="captcha-img" id="register-captcha-img" onclick="loadCaptcha(\'register\')" title="'+t("captcha")+'"></div><input id="register-captcha" type="text" placeholder="'+t("captchaPh")+'"></div></div><button class="btn btn-primary" type="submit">'+t("register")+'</button></form><p class="switch">'+t("haveAccount")+' <a onclick="showLogin()">'+t("login")+'</a></p>';
loadCaptcha("register");
}
function doLogin(){
var em=document.getElementById("email").value.trim(),pw=document.getElementById("password").value,ca=document.getElementById("login-captcha").value;
if(!em){toast(t("emailRequired"),"error");return}if(!pw){toast(t("pwdRequired"),"error");return}if(!ca){toast(t("captchaRequired"),"error");return}
api("/auth/login",{method:"POST",body:JSON.stringify({email:em,password:pw,captchaId:loginCaptchaId,captchaCode:ca})}).then(function(d){if(d.data){setToken("");location.reload()}else{toast(t("loginFailed"),"error")}}).catch(function(e){toast(e.message,"error");loadCaptcha("login")})
}
function doRegister(){
var em=document.getElementById("email").value.trim(),pw=document.getElementById("password").value,cf=document.getElementById("confirm-pwd").value,ca=document.getElementById("register-captcha").value;
if(!em){toast(t("emailRequired"),"error");return}if(!pw){toast(t("pwdRequired"),"error");return}if(!cf){toast(t("confirmRequired"),"error");return}if(!ca){toast(t("captchaRequired"),"error");return}if(pw!==cf){toast(t("pwdMismatch"),"error");return}
api("/auth/register",{method:"POST",body:JSON.stringify({email:em,password:pw,confirm:cf,captchaId:registerCaptchaId,captchaCode:ca})}).then(function(d){if(d.data){setToken("");location.reload()}else{toast(t("registerFailed"),"error")}}).catch(function(e){toast(e.message,"error");loadCaptcha("register")})
}
function doLogout(){api("/auth/logout",{method:"POST"}).then(function(){clearToken();location.href=BASE}).catch(function(){clearToken();location.href=BASE})}
function toggleMobileMenu(btn){var nav=document.querySelector(".navbar .nav-right");if(!nav)return;var open=nav.classList.toggle("open");btn.setAttribute("aria-expanded",open?"true":"false")}
function closeMobileMenu(){var nav=document.querySelector(".navbar .nav-right"),btn=document.querySelector(".mobile-menu-toggle");if(nav)nav.classList.remove("open");if(btn)btn.setAttribute("aria-expanded","false")}
function toggleCreatePanel(btn){var card=document.getElementById("create-card");if(!card)return;var open=card.classList.toggle("mobile-open");btn.setAttribute("aria-expanded",open?"true":"false");var label=btn.querySelector("span");if(label)label.textContent=t(open?"collapseCreate":"expandCreate")}
function enhanceDashboardShell(){
var nav=document.querySelector(".navbar"),navRight=nav&&nav.querySelector(".nav-right");
if(nav&&navRight&&!nav.querySelector(".mobile-menu-toggle")){
var menu=document.createElement("button");menu.type="button";menu.className="mobile-menu-toggle";menu.setAttribute("aria-label",t("menu"));menu.setAttribute("aria-expanded","false");menu.setAttribute("aria-controls","nav-menu");
menu.onclick=function(){var open=navRight.classList.toggle("open");menu.setAttribute("aria-expanded",open?"true":"false")};
menu.innerHTML='<svg class="menu-open-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M4 7h16M4 12h16M4 17h16"/></svg><svg class="menu-close-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M18 6 6 18M6 6l12 12"/></svg>';
navRight.id="nav-menu";
nav.insertBefore(menu,navRight)
}
var card=document.querySelector(".sidebar .card"),header=card&&card.querySelector(".card-header"),body=card&&card.querySelector(".card-body");
if(card&&header&&body){card.id="create-card";card.classList.add("create-card");header.classList.add("create-card-header");body.classList.add("create-form-body");var toggle=document.createElement("button");toggle.type="button";toggle.className="create-panel-toggle";toggle.setAttribute("aria-expanded","false");toggle.setAttribute("aria-controls","create-card-body");toggle.onclick=function(){toggleCreatePanel(toggle)};toggle.innerHTML='<span>'+t("expandCreate")+'</span><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m6 9 6 6 6-6"/></svg>';header.appendChild(toggle)}
}
function renderDashboard(){
var al=currentUser.isAdmin?'<a class="admin-link" href="'+BASE+'admin">` + iconLayout + `<span>'+t("adminPanel")+'</span></a>':'';
var th='<button type="button" id="theme-toggle" class="btn-icon" onclick="toggleTheme()" aria-label="'+t("toggleTheme")+'" title="'+t("toggleTheme")+'">'+(getTheme()==="dark"?'` + iconSun + `':'` + iconMoon + `')+'</button>';
var lh='<div class="lang-toggle" role="group" aria-label="Language"><button type="button" class="'+(lang==="en"?"active":"")+'" onclick="setLang(\'en\')" aria-pressed="'+(lang==="en")+'">EN</button><button type="button" class="'+(lang==="zh"?"active":"")+'" onclick="setLang(\'zh\')" aria-pressed="'+(lang==="zh")+'">中文</button></div>';
document.getElementById("app").innerHTML='<nav class="navbar"><div class="logo">` + logoIcon + `<span class="logo-text">Vibe<span class="accent">cast</span></span></div><div class="nav-right"><div class="nav-section nav-workspace">'+al+'<button type="button" class="btn-icon" id="nav-org-btn" onclick="openOrgModal()">` + iconBuilding + `<span id="nav-org-label">'+t("noOrg")+'</span></button><button type="button" class="btn-icon" onclick="openChangePwdModal()">` + iconLock + `<span>'+t("changePassword")+'</span></button></div><div class="nav-section nav-preferences">'+th+lh+'</div><div class="nav-section nav-account"><span class="email">'+esc(currentUser.email)+'</span><button type="button" class="btn-link" onclick="doLogout()">` + iconLogout + `<span>'+t("logout")+'</span></button></div></div></nav><div class="container"><div class="dashboard-grid"><div class="sidebar"><div class="card"><div class="card-header"><h2>'+t("create")+'</h2></div><div class="card-body"><div class="form-field"><label>'+t("siteName")+'</label><input id="site-name" placeholder="'+t("siteNamePh")+'"></div><div class="form-field"><label>'+t("siteDesc")+'</label><input id="site-desc" placeholder="'+t("siteDescPh")+'"></div><div class="form-field"><label>'+t("sitePwd")+'</label><div class="pwd-wrap"><input id="site-pwd" type="password" placeholder="'+t("sitePwdPh")+'"><button type="button" class="pwd-toggle" onclick="togglePwd(this)" aria-label="'+t("showPassword")+'" aria-pressed="false"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-10-8-10-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 10 8 10 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/><line x1="1" y1="1" x2="23" y2="23"/></svg></button></div><div class="desc" id="pwd-desc">'+t("pwdDesc")+'</div></div><label class="switch-field"><span class="switch-copy"><span class="switch-title">'+t("orgOpenLabel")+'</span><span class="desc">'+t("orgOpenDesc")+'</span></span><span class="switch-toggle"><input type="checkbox" id="site-org-open"><span class="switch-track"></span></span></label><label class="switch-field"><span class="switch-copy"><span class="switch-title">'+t("orgPinnedLabel")+'</span><span class="desc">'+t("orgPinnedDesc")+'</span></span><span class="switch-toggle"><input type="checkbox" id="site-org-pinned"><span class="switch-track"></span></span></label><div class="form-actions"><button type="button" class="btn btn-primary btn-block" id="create-site-btn" onclick="createSite()">'+t("create")+'</button></div></div></div></div><div class="main-content"><div class="card"><div class="site-tabs"><div class="site-tab '+(siteTab==="my"?"active":"")+'" onclick="switchSiteTab(\'my\')">'+t("tabMySites")+'</div><div class="site-tab '+(siteTab==="org"?"active":"")+'" onclick="switchSiteTab(\'org\')">'+t("tabOrgPinned")+'</div></div><div class="card-body"><div id="site-tab-content"></div></div></div></div></div></div><div class="modal-overlay" id="pwd-modal" onclick="if(event.target===this)closeChangePwdModal()"><div class="modal" role="dialog" aria-modal="true" aria-labelledby="pwd-modal-title"><h3 id="pwd-modal-title">'+t("changePassword")+'</h3><div class="modal-field"><label>'+t("currentPassword")+'</label><div class="pwd-wrap"><input id="old-pwd" type="password" placeholder="'+t("currentPassword")+'"><button type="button" class="pwd-toggle" onclick="togglePwd(this)" aria-label="'+t("showPassword")+'" aria-pressed="false"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-10-8-10-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 10 8 10 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/><line x1="1" y1="1" x2="23" y2="23"/></svg></button></div></div><div class="modal-field"><label>'+t("newPassword")+'</label><div class="pwd-wrap"><input id="new-pwd" type="password" placeholder="'+t("newPasswordPh")+'"><button type="button" class="pwd-toggle" onclick="togglePwd(this)" aria-label="'+t("showPassword")+'" aria-pressed="false"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-10-8-10-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 10 8 10 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/><line x1="1" y1="1" x2="23" y2="23"/></svg></button></div></div><div class="modal-field"><label>'+t("confirmPassword")+'</label><div class="pwd-wrap"><input id="confirm-new-pwd" type="password" placeholder="'+t("confirmPassword")+'"><button type="button" class="pwd-toggle" onclick="togglePwd(this)" aria-label="'+t("showPassword")+'" aria-pressed="false"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-10-8-10-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 10 8 10 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/><line x1="1" y1="1" x2="23" y2="23"/></svg></button></div></div><div class="modal-actions"><button class="btn btn-outline" onclick="closeChangePwdModal()">'+t("cancel")+'</button><button class="btn btn-primary" onclick="changePassword()">'+t("save")+'</button></div></div></div><div class="modal-overlay" id="org-modal" onclick="if(event.target===this)closeOrgModal()"><div class="modal modal-lg" role="dialog" aria-modal="true" aria-labelledby="org-modal-title"><div class="org-modal-header"><div class="org-modal-heading"><h3 id="org-modal-title">'+t("orgManage")+'</h3><p class="org-feature-desc">'+t("orgFeatureDesc")+'</p></div><button type="button" class="modal-close" onclick="closeOrgModal()" title="'+t("cancel")+'" aria-label="'+t("cancel")+'"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M18 6 6 18M6 6l12 12"/></svg></button></div><div id="org-section"></div></div></div><div class="modal-overlay" id="edit-site-modal" onclick="if(event.target===this)closeEditSiteModal()"><div class="modal" role="dialog" aria-modal="true" aria-labelledby="edit-modal-title"><h3 id="edit-modal-title">'+t("editSite")+'</h3><div class="modal-field"><label>'+t("editName")+'</label><input id="edit-site-name" type="text"></div><div class="modal-field"><label>'+t("editDesc")+'</label><input id="edit-site-desc" type="text"></div><div class="modal-actions"><button class="btn btn-outline" onclick="closeEditSiteModal()">'+t("cancel")+'</button><button class="btn btn-primary" onclick="saveEditSite()">'+t("save")+'</button></div></div></div>';
enhanceDashboardShell();
var orgClose=document.querySelector("#org-modal .modal-close");orgClose.title=t("close");orgClose.setAttribute("aria-label",t("close"));
loadSites();
loadOrg();
fetch(BASE+"api/settings").then(function(r){return r.json()}).then(function(d){if(d.data&&d.data.maxUploadSize)maxUploadMB=d.data.maxUploadSize;if(d.data&&d.data.maxSitesPerUser)maxSites=d.data.maxSitesPerUser}).catch(function(){});
}
var modalReturnFocus={};
function openModal(id){
var m=document.getElementById(id);if(!m)return;
modalReturnFocus[id]=document.activeElement;
m.classList.add("show");
var first=m.querySelector("[autofocus]")||m.querySelector("input")||m.querySelector("button");
if(first)setTimeout(function(){try{first.focus()}catch(e){}},40)
}
function closeModal(id){
var m=document.getElementById(id);if(!m)return;
m.classList.remove("show");
var back=modalReturnFocus[id];delete modalReturnFocus[id];
if(back&&back.focus)try{back.focus()}catch(e){}
}
document.addEventListener("keydown",function(e){
if(e.key!=="Escape")return;
var menus=document.querySelectorAll("details.row-menu[open]");
if(menus.length){for(var mi=0;mi<menus.length;mi++)menus[mi].removeAttribute("open");e.preventDefault();return}
var open=document.querySelectorAll(".modal-overlay.show");
if(!open.length)return;
e.preventDefault();
var top=open[open.length-1];
if(top.id){closeModal(top.id);return}
top.classList.remove("show")
});
function openChangePwdModal(){openModal("pwd-modal")}
function closeChangePwdModal(){var m=document.getElementById("pwd-modal");if(m)m.querySelectorAll("input").forEach(function(i){i.value=""});closeModal("pwd-modal")}
function openOrgModal(){openModal("org-modal");loadOrg()}
function closeOrgModal(){closeModal("org-modal")}
function updateOrgNavBtn(){var btn=document.getElementById("nav-org-btn");if(!btn)return;var label=orgInfo&&orgInfo.name?esc(orgInfo.name):t("noOrg");btn.innerHTML='` + iconBuilding + `<span> '+label+'</span>'}
function sitePageGo(p){sitePage=p;loadSites()}
function switchSiteTab(tab){siteTab=tab;var tabs=document.querySelectorAll(".site-tab");for(var i=0;i<tabs.length;i++){tabs[i].classList.remove("active");if(tabs[i].getAttribute("onclick").indexOf("'"+tab+"'")>=0)tabs[i].classList.add("active")}if(siteTab==="my"){loadSites()}else{loadOrgPinnedSites()}}
function orgPinnedPageGo(p){orgPinnedPage=p;loadOrgPinnedSites()}
function searchOrgPinned(){orgPinnedSearch=document.getElementById("site-search").value;orgPinnedPage=1;loadOrgPinnedSites()}
function loadOrgPinnedSites(){
var q=orgPinnedSearch?"&q="+encodeURIComponent(orgPinnedSearch):"";
api("/org/sites?page="+orgPinnedPage+"&perPage="+sitePerPage+q).then(function(d){
var r=d.data||{},sites=r.items||[];orgPinnedTotal=r.total||0;
var el=document.getElementById("site-tab-content");if(!el)return;
if(sites.length===0&&r.total===0&&!orgInfo){el.innerHTML='<div class="empty"><span class="empty-icon">` + iconUsers + `</span><strong>'+t("noOrg")+'</strong><span>'+t("orgPinnedJoinHint")+'</span><button type="button" class="btn btn-outline btn-sm" onclick="openOrgModal()">'+t("orgCta")+'</button></div>';return}
var tp=Math.ceil(orgPinnedTotal/sitePerPage)||1,pg=paginationHtml(orgPinnedPage,tp,"orgPinnedPageGo");
var toolbar='<div class="list-toolbar"><input type="text" id="site-search" placeholder="'+t("searchSitesPh")+'" onkeydown="if(event.key===\'Enter\')searchOrgPinned()" value="'+esc(orgPinnedSearch)+'"></div>';
if(!sites.length){el.innerHTML=toolbar+'<div class="empty"><span class="empty-icon">` + iconInbox + `</span><strong>'+t("orgPinnedEmptyTitle")+'</strong><span>'+t("noOrgPinned")+'</span></div>'+pg;return}
var h=toolbar;
for(var i=0;i<sites.length;i++){var s=sites[i];
var desc=s.description?'<span class="pinned-desc">'+esc(s.description)+'</span>':'';
h+='<div class="org-pinned-item"><span class="badge badge-pinned">'+t("orgPinnedBadge")+'</span><div class="pinned-text"><div class="pinned-title"><span class="pinned-name" title="'+esc(s.name)+'">'+esc(s.name)+'</span>'+desc+'</div></div><span class="pinned-url" title="'+esc(s.url)+'">'+esc(s.url)+'</span><div class="pinned-actions"><a class="btn btn-outline btn-sm" href="'+esc(s.url)+'" target="_blank" rel="noopener">'+t("visit")+'</a></div></div>'}
h+=pg;el.innerHTML=h}).catch(function(e){var el=document.getElementById("site-tab-content");if(!el)return;
if(!orgInfo){el.innerHTML='<div class="empty"><span class="empty-icon">` + iconUsers + `</span><strong>'+t("noOrg")+'</strong><span>'+t("orgPinnedJoinHint")+'</span><button type="button" class="btn btn-outline btn-sm" onclick="openOrgModal()">'+t("orgCta")+'</button></div>';return}
el.innerHTML='<div class="empty"><span class="empty-icon">` + iconInbox + `</span><strong>'+t("orgPinnedEmptyTitle")+'</strong><span>'+esc(e.message)+'</span><button type="button" class="btn btn-outline btn-sm" onclick="loadOrgPinnedSites()">'+t("confirm")+'</button></div>'})
}
function openSiteQuery(){var el=document.getElementById("site-search");return el?el.value:siteSearch}
function searchSites(){siteSearch=openSiteQuery();sitePage=1;loadSites()}
function clearSiteSearch(){siteSearch="";sitePage=1;loadSites()}
function siteHeadKey(e,id){
var t=e.target;
if(e.key==="Enter"||e.key===" "){if(t!==e.currentTarget)return;e.preventDefault();toggleDetail(id)}
}
function closeRowMenus(){var ds=document.querySelectorAll("details.row-menu[open]");for(var i=0;i<ds.length;i++)ds[i].removeAttribute("open")}
function copySiteUrl(url){var full=url&&url.indexOf("http")===0?url:location.origin+siteUrl(url);copyText(full);closeRowMenus()}
function emptySitesHtml(){
var filtered=!!siteSearch;
var title=filtered?t("noMatch"):t("noSitesTitle");
var hint=filtered?t("noMatchHint"):t("noSitesHint");
return '<div class="empty"><span class="empty-icon">` + iconInbox + `</span><strong>'+esc(title)+'</strong><span>'+esc(hint)+'</span>'+(filtered?'<button class="btn btn-outline btn-sm" onclick="clearSiteSearch()">'+t("clearSearch")+'</button>':'')+'</div>'
}
function sitePath(u){
try{var p=new URL(u,location.origin);return p.pathname+(p.search||"")}catch(e){return u||""}
}
function siteRowHtml(s){
var dot="",badge="";
if(s.publicAccessDisabled&&!s.protected){dot="disabled";badge='<span class="badge badge-disabled">'+t("accessDisabled")+'</span>'}
else if(s.protected){dot="protected";badge='<span class="badge badge-protected">'+t("protected")+'</span>'}
else{dot="public";badge='<span class="badge badge-public">'+t("public")+'</span>'}
if(s.orgOpen){badge+='<span class="badge badge-org">'+t("orgBadge")+'</span>'}
if(s.orgPinned){badge+='<span class="badge badge-pinned">'+t("orgPinnedBadge")+'</span>'}
var vs=s.visits||{today:0,month:0,total:0};
var visText=t(vs.total===1?"visitsOne":"visitsCount").replace("{count}",vs.total||0);
var visTitle=t("visitsToday")+" "+vs.today+" · "+t("visitsMonth")+" "+vs.month+" · "+t("visitsTotal")+" "+vs.total;
var pwd=s.protected
  ?'<span class="pwd-value" id="pwd-'+s.id+'">••••</span><span class="pwd-actions"><button type="button" class="copy-btn" onclick="event.stopPropagation();copySitePwd('+s.id+')" title="'+t("copyPassword")+'" aria-label="'+t("copyPassword")+'">` + iconCopy + `</button><button type="button" class="copy-btn" onclick="event.stopPropagation();shareSite('+s.id+',\''+esc(s.name)+'\',\''+esc(s.url)+'\')" title="'+t("share")+'" aria-label="'+t("share")+'">` + iconShare + `</button></span>'
  :'<span class="pwd-none">'+t("none")+'</span>';
var menu='<details class="row-menu" onclick="event.stopPropagation()"><summary aria-label="'+t("siteActions")+'" title="'+t("siteActions")+'" aria-haspopup="true">` + iconDots + `</summary><div class="menu-panel" role="menu">'
  +'<button type="button" role="menuitem" onclick="copySiteUrl(\''+esc(s.url)+'\')">` + iconExternal + `<span>'+t("copyUrl")+'</span></button>'
  +(s.protected?'<button type="button" role="menuitem" onclick="copySitePwd('+s.id+');closeRowMenus()">` + iconCopy + `<span>'+t("copyPassword")+'</span></button>':'')
  +'<button type="button" role="menuitem" onclick="shareSite('+s.id+',\''+esc(s.name)+'\',\''+esc(s.url)+'\');closeRowMenus()">` + iconShare + `<span>'+t("share")+'</span></button>'
  +'<div class="menu-sep" role="separator"></div>'
  +'<button type="button" role="menuitem" onclick="openEditSiteModal('+s.id+');closeRowMenus()">` + iconEdit + `<span>'+t("editSite")+'</span></button>'
  +'<button type="button" role="menuitem" onclick="downloadSite('+s.id+');closeRowMenus()">` + iconDownload + `<span>'+t("downloadZip")+'</span></button>'
  +'<button type="button" role="menuitem" onclick="clearSiteFiles('+s.id+');closeRowMenus()">` + iconTrash + `<span>'+t("clearFiles")+'</span></button>'
  +'<div class="menu-sep" role="separator"></div>'
  +'<button type="button" role="menuitem" class="is-danger" onclick="delSite('+s.id+');closeRowMenus()">` + iconTrash + `<span>'+t("delete")+'</span></button>'
  +'</div></details>';
var detail='<div class="site-detail" id="detail-'+s.id+'"><div class="detail-clip"><div class="detail-inner"><div class="detail-grid">'
  +'<div class="detail-section"><div class="detail-label">'+t("siteDesc")+'</div><div class="detail-value'+(s.description?"":" is-muted")+'">'+(s.description?esc(s.description):t("none"))+'</div></div>'
  +'<div class="detail-section"><div class="detail-label">'+t("accessPassword")+'</div><div class="detail-value">'+pwd+'</div></div>'
  +'<div class="detail-section"><div class="detail-label">'+t("visits")+'</div><div class="visits-grid">'
  +'<span class="visit-cell"><span class="visit-key">'+t("visitsToday")+'</span><span class="visit-num">'+vs.today+'</span></span>'
  +'<span class="visit-cell"><span class="visit-key">'+t("visitsMonth")+'</span><span class="visit-num">'+vs.month+'</span></span>'
  +'<span class="visit-cell"><span class="visit-key">'+t("visitsTotal")+'</span><span class="visit-num">'+vs.total+'</span></span>'
  +'</div></div></div>'
  +'<div class="detail-toggles">'
  +'<label class="switch-field"><span class="switch-copy"><span class="switch-title">'+t("orgOpenLabel")+'</span><span class="desc">'+t("orgOpenDesc")+'</span></span><span class="switch-toggle"><input type="checkbox" id="org-open-'+s.id+'" '+(s.orgOpen?"checked":"")+' onchange="toggleOrgOpen('+s.id+')"><span class="switch-track"></span></span></label>'
  +'<label class="switch-field"><span class="switch-copy"><span class="switch-title">'+t("orgPinnedLabel")+'</span><span class="desc">'+t("orgPinnedDesc")+'</span></span><span class="switch-toggle"><input type="checkbox" id="org-pinned-'+s.id+'" '+(s.orgPinned?"checked":"")+' onchange="toggleOrgPinned('+s.id+')"><span class="switch-track"></span></span></label>'
  +'</div><div class="site-files" id="files-'+s.id+'"></div></div></div></div>';
return '<li class="site-item" id="site-'+s.id+'">'
  +'<div class="site-head" role="button" tabindex="0" aria-expanded="false" aria-controls="detail-'+s.id+'" onclick="toggleDetail('+s.id+')" onkeydown="siteHeadKey(event,'+s.id+')">'
  +'<span class="site-chevron" aria-hidden="true">` + iconChevronRight + `</span>'
  +'<div class="site-main"><div class="site-title"><span class="status-dot '+dot+'" aria-hidden="true"></span><span class="site-name" title="'+esc(s.name)+'">'+esc(s.name)+'</span>'+badge+'</div>'
  +'<div class="site-sub"><a class="site-url" href="'+esc(s.url)+'" target="_blank" rel="noopener" onclick="event.stopPropagation()" title="'+esc(s.url)+'"><span class="url-text"><span class="url-full">'+esc(s.url)+'</span><span class="url-path">'+esc(sitePath(s.url))+'</span></span>` + iconExternal + `</a>'
  +'<span class="sep" aria-hidden="true">·</span><span class="site-visits" title="'+esc(visTitle)+'">'+esc(visText)+'</span></div></div>'
  +'<div class="site-actions"><label class="btn btn-primary btn-sm upload-btn" id="upload-btn-'+s.id+'" onclick="event.stopPropagation()">'+t("deployBtn")+'<input type="file" accept=".zip,.pdf,.doc,.docx,.xls,.xlsx,.ppt,.pptx,.odt,.ods,.odp,.rtf,.txt,.csv,.json,.xml,.html,.htm,.css,.js,.mjs,.svg,.png,.jpg,.jpeg,.gif,.webp,.avif,.ico,.mp4,.webm,.mp3,.ogg,.wav,.flac,.wasm,.woff,.woff2,.ttf,.otf,.eot,.webmanifest" onchange="deploy('+s.id+',this.files[0])"></label>'+menu+'</div>'
  +'</div>'+detail+'</li>'
}
function loadSites(){
var q=siteSearch?"&q="+encodeURIComponent(siteSearch):"";
api("/sites?page="+sitePage+"&perPage="+sitePerPage+q).then(function(d){
var r=d.data||{},sites=r.items||[];siteTotal=r.total||0;
var el=document.getElementById("site-tab-content");if(!el)return;
var pd=document.getElementById("pwd-desc");
if(sites.length>0&&sites[0].publicAccessDisabled&&pd){pd.textContent=t("pwdRequired");pd.classList.add("is-danger")}
var tp=Math.ceil(siteTotal/sitePerPage)||1,pg=paginationHtml(sitePage,tp,"sitePageGo");
var toolbar='<div class="list-toolbar"><input type="text" id="site-search" placeholder="'+t("searchSitesPh")+'" aria-label="'+t("searchSitesPh")+'" value="'+esc(siteSearch)+'" onkeydown="if(event.key===\'Enter\')searchSites();if(event.key===\'Escape\'){this.value=\'\';clearSiteSearch()}"></div>';
if(!sites.length){el.innerHTML=toolbar+emptySitesHtml()+pg;return}
var h=toolbar+'<ul class="site-list">';
for(var i=0;i<sites.length;i++){h+=siteRowHtml(sites[i])}
h+='</ul><div class="deploy-hint">'+t("deployHintDetailed").replaceAll("{size}",maxUploadMB)+'</div>'+pg;
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
function toggleDetail(id){
var e=document.getElementById("detail-"+id);if(!e)return;
var open=e.classList.toggle("show");
var item=document.getElementById("site-"+id);
if(item){item.classList.toggle("is-open",open);var head=item.querySelector(".site-head");if(head)head.setAttribute("aria-expanded",open?"true":"false")}
if(open){loadFileTree(id);loadSitePwd(id)}
}
function loadFileTree(id,preload){
var req=preload?Promise.resolve({data:preload}):api("/sites/"+id+"/files");
req.then(function(d){
var data=d.data||{},files=Array.isArray(data)?data:(data.items||[]);
var fileCount=typeof data.fileCount==="number"?data.fileCount:files.filter(function(f){return!f.dir}).length;
var totalSize=typeof data.totalSize==="number"?data.totalSize:files.reduce(function(n,f){return n+(f.dir?0:(f.size||0))},0);
var host=document.getElementById("files-"+id);if(!host)return;
var summary=t("fileSummary").replace("{count}",fileCount).replace("{size}",formatSize(totalSize));
var h='<div class="file-tree-section"><div class="file-tree-head"><div class="file-tree-title"><span class="detail-label">'+t("files")+'</span>';
if(fileCount>0){h+='<span class="file-tree-summary">'+esc(summary)+'</span>'}
h+='</div>';
if(fileCount>0){h+='<div class="file-tree-actions"><button type="button" class="btn btn-outline btn-sm" onclick="event.stopPropagation();downloadSite('+id+')">` + iconDownload + `'+t("downloadZip")+'</button><button type="button" class="btn btn-outline btn-sm file-clear-btn" onclick="event.stopPropagation();clearSiteFiles('+id+')">` + iconTrash + `'+t("clearFiles")+'</button></div>'}
h+='</div>';
if(!files.length){h+='<div class="file-empty">` + iconFile + `'+t("noFiles")+'</div>'}
else{h+='<div class="file-list">';
for(var i=0;i<files.length;i++){var f=files[i],icon=f.dir?'` + iconFolder + `':'` + iconFile + `',sz=f.dir?"—":formatSize(f.size);
var del=f.dir?'<span></span>':'<button type="button" class="file-del" title="'+t("delete")+'" aria-label="'+t("delete")+' '+esc(f.name)+'" onclick="event.stopPropagation();deleteSiteFile('+id+',\''+esc(f.name)+'\')">` + iconX + `</button>';
h+='<div class="file-row'+(f.dir?" file-dir":"")+'"><span class="file-icon">'+icon+'</span><span class="file-name" title="'+esc(f.name)+'">'+esc(f.name)+(f.dir?"/":"")+'</span><span class="file-size">'+sz+'</span>'+del+'</div>'}
h+='</div>'}
h+='</div>';
host.innerHTML=h}).catch(function(e){var host=document.getElementById("files-"+id);if(host)host.innerHTML='<div class="file-empty">'+esc(e.message)+'</div>'})
}
function formatSize(n){if(n<1024)return n+" B";if(n<1048576)return(n/1024).toFixed(1)+" KB";if(n<1073741824)return(n/1048576).toFixed(1)+" MB";return(n/1073741824).toFixed(1)+" GB"}
function deleteSiteFile(id,path){
var msg=t("deleteFileConfirm").replace("{name}",path);
customConfirm(msg,function(){
api("/sites/"+id+"/files?path="+encodeURIComponent(path),{method:"DELETE"}).then(function(d){
toast(t("fileDeleted"));
var data=(d&&d.data)||{};
loadFileTree(id,data);
}).catch(function(e){toast(e.message,"error")})
},{confirmLabel:t("delete")})
}
function clearSiteFiles(id){
api("/sites/"+id+"/files").then(function(d){var data=d.data||{},count=data.fileCount||0,total=data.totalSize||0;if(!count){toast(t("noFilesToClear"),"error");return}var msg=t("clearFilesConfirm").replace("{count}",count).replace("{size}",formatSize(total));customConfirm(msg,function(){api("/sites/"+id+"/files",{method:"DELETE"}).then(function(){toast(t("filesCleared"));loadFileTree(id)}).catch(function(e){toast(e.message,"error")})})}).catch(function(e){toast(e.message,"error")})
}
function ensureTransferStatus(id){var el=document.getElementById("transfer-"+id);if(el)return el;var detail=document.getElementById("detail-"+id);if(!detail)return null;el=document.createElement("div");el.id="transfer-"+id;el.className="transfer-status";el.innerHTML='<div class="transfer-head"><span class="transfer-label"></span><span class="transfer-percent"></span></div><div class="transfer-track"><span class="transfer-fill"></span></div>';detail.parentElement.insertBefore(el,detail);return el}
function setTransferStatus(id,label,pct,state){var el=ensureTransferStatus(id);if(!el)return;delete el.dataset.hideToken;state=state||"active";el.className="transfer-status show "+state;el.querySelector(".transfer-label").textContent=label;var percent=el.querySelector(".transfer-percent"),fill=el.querySelector(".transfer-fill");if(typeof pct==="number"){pct=Math.max(0,Math.min(100,pct));percent.textContent=pct+"%";fill.style.width=pct+"%"}else{percent.textContent="";fill.style.width=""}}
function hideTransferStatus(id,delay){var el=document.getElementById("transfer-"+id);if(!el)return;var token=String(Date.now());el.dataset.hideToken=token;setTimeout(function(){if(el.dataset.hideToken!==token)return;el.classList.remove("show");setTimeout(function(){if(el.parentElement)el.remove()},220)},delay||1500)}
function downloadSite(id){
var xhr=new XMLHttpRequest();xhr.open("GET",API+"/sites/"+id+"/download");xhr.responseType="blob";xhr.setRequestHeader("Accept-Language",lang||"en");setTransferStatus(id,t("preparingDownload"),null,"indeterminate");
xhr.onprogress=function(e){if(e.lengthComputable){var pct=Math.round(e.loaded/e.total*100);setTransferStatus(id,t("downloadingFile").replace("{size}",formatSize(e.loaded)+" / "+formatSize(e.total)),pct,"active")}else{setTransferStatus(id,t("downloadingFile").replace("{size}",formatSize(e.loaded)),null,"indeterminate")}};
xhr.onload=function(){if(xhr.status>=200&&xhr.status<300){var disposition=xhr.getResponseHeader("Content-Disposition")||"",match=/filename="([^"]+)"/.exec(disposition),name=match?match[1]:"site.zip";var url=URL.createObjectURL(xhr.response),a=document.createElement("a");a.href=url;a.download=name;document.body.appendChild(a);a.click();a.remove();setTimeout(function(){URL.revokeObjectURL(url)},1000);setTransferStatus(id,t("downloadReady"),100,"success");toast(t("downloadReady"));hideTransferStatus(id,1800);return}xhr.response.text().then(function(text){var message=t("downloadFailed");try{message=JSON.parse(text).error||message}catch(e){}setTransferStatus(id,message,100,"error");toast(message,"error");hideTransferStatus(id,3500)}).catch(function(){setTransferStatus(id,t("downloadFailed"),100,"error");toast(t("downloadFailed"),"error");hideTransferStatus(id,3500)})};
xhr.onerror=function(){setTransferStatus(id,t("downloadFailed"),100,"error");toast(t("downloadFailed"),"error");hideTransferStatus(id,3500)};xhr.send()
}
function copyText(text){if(navigator.clipboard){navigator.clipboard.writeText(text).then(function(){toast(t("copied"))}).catch(function(){fallbackCopy(text)})}else{fallbackCopy(text)}}
function fallbackCopy(text){var ta=document.createElement("textarea");ta.value=text;ta.style.position="fixed";ta.style.opacity="0";document.body.appendChild(ta);ta.select();try{document.execCommand("copy");toast(t("copied"))}catch(e){toast(t("copyFailed"),"error")}document.body.removeChild(ta)}
function loadSitePwd(id){api("/sites/"+id+"/password").then(function(d){var p=(d.data||{}).password||"";var el=document.getElementById("pwd-"+id);if(el)el.textContent=p;el=document.querySelector('[data-pwd-btn="'+id+'"]');if(el)el.dataset.pwd=p}).catch(function(){})}
function copySitePwd(id){api("/sites/"+id+"/password").then(function(d){var p=(d.data||{}).password||"";copyText(p)}).catch(function(e){toast(e.message,"error")})}
function shareSite(siteId,name,url){
var fullUrl=url.startsWith("http")?url:location.origin+"/"+url.replace(/^\//,"");
if(!siteId){var key="shareTemplatePublic";var txt=t(key).replace("{name}",name).replace("{url}",fullUrl);copyText(txt);toast(t("shareCopied"));return}
api("/sites/"+siteId+"/password").then(function(d){var pwd=(d.data||{}).password||"";var key=pwd?"shareTemplatePwd":"shareTemplatePublic";var txt=t(key).replace("{name}",name).replace("{url}",fullUrl);if(pwd)txt=txt.replace("{password}",pwd);copyText(txt);toast(t("shareCopied"))}).catch(function(e){toast(e.message,"error")})
}
function createSite(){
api("/sites",{method:"POST",body:JSON.stringify({name:document.getElementById("site-name").value,description:document.getElementById("site-desc").value,password:document.getElementById("site-pwd").value,orgOpen:document.getElementById("site-org-open").checked,orgPinned:document.getElementById("site-org-pinned").checked})}).then(function(){document.getElementById("site-name").value="";document.getElementById("site-desc").value="";document.getElementById("site-pwd").value="";toast(t("siteCreated"));loadSites()}).catch(function(e){toast(e.message,"error")})
}
function deploy(id,file){
if(!file)return;
var btn=document.getElementById("upload-btn-"+id);
var input=btn?btn.querySelector("input[type=file]"):null;
var origText=btn?btn.firstChild.textContent:"";
var allowedExts=[".zip",".pdf",".doc",".docx",".xls",".xlsx",".ppt",".pptx",".odt",".ods",".odp",".rtf",".txt",".csv",".json",".xml",".html",".htm",".css",".js",".mjs",".svg",".png",".jpg",".jpeg",".gif",".webp",".avif",".ico",".mp4",".webm",".mp3",".ogg",".wav",".flac",".wasm",".woff",".woff2",".ttf",".otf",".eot",".webmanifest"];
var fName=file.name.toLowerCase();
var ext="."+fName.split(".").pop();
if(allowedExts.indexOf(ext)<0){toast(t("fileTypeBlocked")+" ("+ext+")","error");if(input)input.value="";return}
if(file.size>maxUploadMB*1024*1024){toast(t("fileTooLarge")+" ("+maxUploadMB+"MB)","error");if(input)input.value="";return}
if(btn){btn.firstChild.textContent=t("uploading");btn.style.opacity="0.6";btn.style.pointerEvents="none";if(input)input.disabled=true}
setTransferStatus(id,t("uploadPreparing"),0,"active");
var fd=new FormData();fd.append("file",file);
	var xhr=new XMLHttpRequest();
	xhr.open("POST",API+"/sites/"+id+"/deploy");
	xhr.setRequestHeader("Accept-Language",lang||"en");
xhr.upload.onprogress=function(e){if(e.lengthComputable){var pct=Math.round(e.loaded/e.total*100);if(btn)btn.firstChild.textContent=t("uploading")+" "+pct+"%";setTransferStatus(id,t("uploading")+" "+formatSize(e.loaded)+" / "+formatSize(e.total),pct,"active")}};
xhr.upload.onload=function(){if(btn)btn.firstChild.textContent=t("uploading");setTransferStatus(id,t("uploadProcessing"),null,"indeterminate")};
xhr.onload=function(){
if(btn){btn.firstChild.textContent=origText;btn.style.opacity="";btn.style.pointerEvents="";if(input){input.disabled=false;input.value=""}}
try{var data=JSON.parse(xhr.responseText);if(xhr.status>=200&&xhr.status<300){setTransferStatus(id,t("uploadComplete"),100,"success");toast(t("deployed"));hideTransferStatus(id,1400);setTimeout(loadSites,700)}else{var message=data.error||t("uploadFailed");setTransferStatus(id,message,100,"error");toast(message,"error");hideTransferStatus(id,3500)}}catch(e){setTransferStatus(id,t("networkError"),100,"error");toast(t("networkError"),"error");hideTransferStatus(id,3500)}
};
xhr.onerror=function(){if(btn){btn.firstChild.textContent=origText;btn.style.opacity="";btn.style.pointerEvents="";if(input){input.disabled=false;input.value=""}}setTransferStatus(id,t("networkError"),100,"error");toast(t("networkError"),"error");hideTransferStatus(id,3500)};
xhr.send(fd);
}
function delSite(id){
customConfirm(t("deleteConfirm"),function(){api("/sites/"+id,{method:"DELETE"}).then(function(){toast(t("deleted"));loadSites()}).catch(function(e){toast(e.message,"error")})})
}
function toggleOrgOpen(id){
var cb=document.getElementById("org-open-"+id);if(!cb)return;
api("/sites/"+id,{method:"PUT",body:JSON.stringify({orgOpen:cb.checked})}).then(function(){toast(t("save"));loadSites()}).catch(function(e){toast(e.message,"error");cb.checked=!cb.checked})
}
function toggleOrgPinned(id){
var cb=document.getElementById("org-pinned-"+id);if(!cb)return;
api("/sites/"+id,{method:"PUT",body:JSON.stringify({orgPinned:cb.checked})}).then(function(){toast(t("save"));loadSites()}).catch(function(e){toast(e.message,"error");cb.checked=!cb.checked})
}
var editSiteId=0;
function openEditSiteModal(id){
editSiteId=id;
api("/sites?page=1&perPage=100").then(function(d){var sites=(d.data||{}).items||[];for(var i=0;i<sites.length;i++){if(sites[i].id===id){document.getElementById("edit-site-name").value=sites[i].name||"";document.getElementById("edit-site-desc").value=sites[i].description||"";break}}openModal("edit-site-modal")}).catch(function(e){toast(e.message,"error")})
}
function closeEditSiteModal(){closeModal("edit-site-modal")}
function saveEditSite(){
var n=document.getElementById("edit-site-name").value,d=document.getElementById("edit-site-desc").value;
api("/sites/"+editSiteId,{method:"PUT",body:JSON.stringify({name:n,description:d})}).then(function(){toast(t("save"));closeEditSiteModal();loadSites()}).catch(function(e){toast(e.message,"error")})
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
el.innerHTML='<div class="org-empty"><section class="org-empty-panel"><div class="org-panel-title"><span class="org-panel-icon create"><svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg></span>'+t("orgCreate")+'</div><p class="org-panel-desc">'+t("orgCreateDesc")+'</p><div class="org-panel-form"><input class="org-panel-input" id="org-create-name" placeholder="'+t("orgNamePh")+'" onkeydown="if(event.key===\'Enter\')createOrg()"><button class="btn btn-primary" onclick="createOrg()">'+t("orgCreate")+'</button></div></section><section class="org-empty-panel"><div class="org-panel-title"><span class="org-panel-icon join"><svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 3h5v5M10 14 21 3M21 14v5a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5"/></svg></span>'+t("orgJoin")+'</div><p class="org-panel-desc">'+t("orgJoinDesc")+'</p><div class="org-panel-form"><input class="org-panel-input" id="org-join-code" placeholder="'+t("orgInviteCodePh")+'" onkeydown="if(event.key===\'Enter\')joinOrg()"><button class="btn btn-outline" onclick="joinOrg()">'+t("orgJoin")+'</button></div></section></div>';
return
}
var inviteCode=orgInfo.inviteCode||"";
var name=orgInfo.name||"";
var isOwner=orgInfo.isOwner||false;
var initials=name.substring(0,2).toUpperCase();
var roleBadge=isOwner?'<span class="org-role-badge owner">'+t("orgOwner")+'</span>':'<span class="org-role-badge member">'+t("orgMember")+'</span>';
var h='<div class="org-bound"><div class="org-summary"><div class="org-header"><div class="org-avatar">'+esc(initials)+'</div><div><div class="org-title">'+esc(name)+'</div>'+roleBadge+'</div></div>';
if(isOwner){h+='<div class="org-invite-section"><button type="button" class="org-invite-toggle" aria-expanded="true" onclick="toggleOrgInvite(this)"><span>'+t("orgInviteCode")+'</span><svg class="org-invite-chevron" viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="6 9 12 15 18 9"/></svg></button><div class="org-invite-detail"><div class="org-invite-row"><code class="org-invite-code">'+esc(inviteCode)+'</code><button type="button" class="org-copy-btn" onclick="copyText(\''+esc(inviteCode)+'\')" title="'+t("copy")+'" aria-label="'+t("copy")+'"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg></button></div></div></div>'}
h+='<div class="org-summary-action">';
if(isOwner){h+='<button class="btn btn-outline btn-sm org-danger-btn" onclick="deleteOrg()">'+t("orgDelete")+'</button>'}else{h+='<button class="btn btn-outline btn-sm" onclick="leaveOrg()">'+t("orgLeave")+'</button>'}
h+='</div></div><section class="org-bound-right"><div class="org-members-header"><h3>'+t("orgMembers")+'</h3><span class="org-members-count" id="org-members-count">-</span></div><div class="org-member-search"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><circle cx="11" cy="11" r="7"/><path d="m20 20-3.6-3.6"/></svg><input type="text" id="org-member-search" placeholder="'+t("search")+'..." onkeydown="if(event.key===\'Enter\')searchOrgMembers()" value="'+esc(orgMemberSearch)+'"></div><div id="org-member-list"></div></section></div>';
el.innerHTML=h;
loadOrgMembers()
}
function toggleOrgInvite(btn){var section=btn.parentElement;section.classList.toggle("org-invite-collapsed");btn.setAttribute("aria-expanded",section.classList.contains("org-invite-collapsed")?"false":"true")}
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
<style>
` + uiTokensCSS + uiComponentsCSS + uiAdminCSS + `
</style>
<script>var _t="light";try{_t=localStorage.getItem("theme")||"light"}catch(e){}document.documentElement.setAttribute("data-theme",_t)</script>
</head>
<body>
<div id="app"></div>
<script>
var BASE=location.pathname.replace(/\/(?:dashboard|admin)\/?$/,"/")||"/";var API=BASE+"api";
function getTheme(){return document.documentElement.getAttribute("data-theme")||"light"}
function setTheme(t){document.documentElement.setAttribute("data-theme",t);try{localStorage.setItem("theme",t)}catch(e){}var b=document.getElementById("theme-toggle");if(b)b.innerHTML=t==="dark"?'` + iconSun + `':'` + iconMoon + `'}
function toggleTheme(){setTheme(getTheme()==="dark"?"light":"dark")}
var lang="en";
var userPage=1,userPerPage=10,userTotal=0,userSearch="";
var adminSitePage=1,adminSitePerPage=10,adminSiteTotal=0,adminSiteSearch="";
var i18n={en:{overview:"Overview",settings:"Settings",users:"Users",allSites:"All Sites",userCount:"Users",siteCount:"Sites",adminCount:"Admins",openReg:"Open Registration",openRegDesc:"When enabled, anyone can register a new account.",publicAccess:"Public Site Access",publicAccessDesc:"When disabled, all deployed sites return 403 to visitors.",dashboard:"Dashboard",logout:"Logout",role:"Role",created:"Created",actions:"Actions",promote:"Promote",demote:"Demote",delete:"Delete",deleteUserConfirm:"Delete this user and all their sites?",deleteSiteConfirm:"Delete this site?",name:"Name",slug:"Slug",owner:"Owner",protected:"Protected",public:"Public",password:"Password",storagePath:"Storage",url:"URL",noUsers:"No users",noSites:"No sites",updated:"Updated",deleted:"Deleted",roleUpdated:"Role updated",userDeleted:"User deleted",siteDeleted:"Site deleted",settingsSaved:"Settings saved",none:"None",accessDisabled:"Disabled",search:"Search",searchUsersPh:"Search users...",searchSitesPh:"Search sites...",prev:"Prev",next:"Next",page:"Page",of:"of",ownerEmail:"Owner",domainRestriction:"Email Domain Restriction",domainRestrictionDesc:"Only allow registration from specified domains (one per line).",save:"Save",cancel:"Cancel",confirm:"Confirm",cleanup:"Storage Cleanup",cleanupDesc:"Remove orphaned directories that exist on disk but have no corresponding site in the database.",checkOrphans:"Check for Orphans",noOrphans:"No orphaned directories found",orphanFound:"Orphaned directories found",confirmDelete:"Delete All",cleanupConfirmWarning:"Delete all orphaned directories? This cannot be undone.",cleanupDone:"Cleanup complete",scanning:"Scanning...",files:"Files",noFiles:"No files",visits:"Visits",visitsToday:"Today",visitsMonth:"Month",visitsTotal:"Total",systemUpdate:"System Update",updateCheck:"Check for Updates",updateAvailable:"Update Available",updateToDate:"Up to Date",currentVersion:"Current Version",latestVersion:"Latest Version",updateNow:"Update Now",updating:"Updating...",updateSuccess:"Update successful! Please restart vibecast.",updateFailed:"Update failed",noUpdate:"You are running the latest version.",checkingUpdate:"Checking...",releaseNotes:"Release Notes",viewRelease:"View on GitHub",updateRestart:"Restart Now",updateRestarting:"Restarting...",updateRestartSuccess:"Server is restarting. Page will reload automatically.",updateRestartFailed:"Server did not come back online. Please refresh manually.",updateInProgress:"Update already in progress",updateVerifyFailed:"Checksum verification failed",updateNoChecksum:"No checksum available (skipped verification)",updateDownloadProgress:"Downloading",updateVerifying:"Verifying checksum...",updateInstalling:"Installing...",systemInfo:"System Info",systemInfoDesc:"Runtime environment and resource paths",tabOverview:"Overview",tabUsers:"Users",tabSites:"Sites",tabSettings:"Settings",tabMaintenance:"Maintenance",tabSystem:"System",siVersion:"Version",siStoragePath:"Storage Path",siDbPath:"Database Path",siListenAddr:"Listen Address",siGoVersion:"Go Runtime",siPlatform:"Platform",siStartTime:"Started At",siUptime:"Uptime",siLanURLs:"LAN Access URLs",maxUploadSize:"Max Upload Size (MB)",maxUploadSizeDesc:"Maximum file size for uploads. Default: 50 MB.",maxSitesPerUser:"Max Sites per User",maxSitesPerUserDesc:"Maximum number of sites each user can create. Set to 0 for unlimited.",siteBaseUrl:"Site Base URL",siteBaseUrlDesc:"Configure a dedicated domain for serving uploaded site content (e.g. https://sites.example.com). Isolates user-uploaded HTML from the admin domain to prevent cookie/session injection. Include http:// or https:// (defaults to https if omitted).",siteBaseUrlWarn:"Not configured - uploaded sites share the same domain as the admin panel. This is a security risk: user-uploaded HTML could access admin cookies and sessions. Configure a separate domain for site content (e.g. https://sites.example.com).",siteBaseUrlPh:"https://sites.example.com"},zh:{overview:"概览",settings:"设置",users:"用户",allSites:"所有站点",userCount:"用户数",siteCount:"站点数",adminCount:"管理员数",openReg:"开放注册",openRegDesc:"启用后，任何人都可以注册新账号。",publicAccess:"公开访问",publicAccessDesc:"禁用后，所有已部署站点对访客返回 403。",dashboard:"仪表盘",logout:"退出",role:"角色",created:"创建时间",actions:"操作",promote:"提升",demote:"降级",delete:"删除",deleteUserConfirm:"删除该用户及其所有站点？",deleteSiteConfirm:"删除此站点？",name:"名称",slug:"Slug",owner:"所有者",protected:"已加密",public:"公开",password:"密码",storagePath:"存储",url:"URL",noUsers:"暂无用户",noSites:"暂无站点",updated:"已更新",deleted:"已删除",roleUpdated:"角色已更新",userDeleted:"用户已删除",siteDeleted:"站点已删除",settingsSaved:"设置已保存",none:"无",accessDisabled:"已禁用",search:"搜索",searchUsersPh:"搜索用户...",searchSitesPh:"搜索站点...",prev:"上一页",next:"下一页",page:"第",of:"页 共",ownerEmail:"所有者",domainRestriction:"邮箱域名限制",domainRestrictionDesc:"仅允许指定域名的邮箱注册（每行一个）。",save:"保存",cancel:"取消",confirm:"确认",cleanup:"存储清理",cleanupDesc:"删除磁盘上存在但数据库中无对应站点的孤立目录。",checkOrphans:"检查孤立目录",noOrphans:"未发现孤立目录",orphanFound:"发现孤立目录",confirmDelete:"全部删除",cleanupConfirmWarning:"确认删除所有孤立目录？此操作不可撤销。",cleanupDone:"清理完成",scanning:"扫描中...",files:"文件",noFiles:"暂无文件",visits:"访问",visitsToday:"今日",visitsMonth:"本月",visitsTotal:"总计",systemUpdate:"系统更新",updateCheck:"检查更新",updateAvailable:"有可用更新",updateToDate:"已是最新版本",currentVersion:"当前版本",latestVersion:"最新版本",updateNow:"立即更新",updating:"更新中...",updateSuccess:"更新成功！请重启 vibecast。",updateFailed:"更新失败",noUpdate:"当前已是最新版本。",checkingUpdate:"检查中...",releaseNotes:"更新日志",viewRelease:"在 GitHub 查看",updateRestart:"立即重启",updateRestarting:"正在重启...",updateRestartSuccess:"服务器正在重启，页面将自动刷新。",updateRestartFailed:"服务器未在30秒内恢复，请手动刷新页面。",updateInProgress:"更新正在进行中",updateVerifyFailed:"校验和验证失败",updateNoChecksum:"无校验和（已跳过验证）",updateDownloadProgress:"下载中",updateVerifying:"正在校验...",updateInstalling:"正在安装...",systemInfo:"系统信息",systemInfoDesc:"运行环境与资源路径",tabOverview:"概览",tabUsers:"用户",tabSites:"站点",tabSettings:"设置",tabMaintenance:"维护",tabSystem:"系统",siVersion:"版本",siStoragePath:"存储路径",siDbPath:"数据库路径",siListenAddr:"监听地址",siGoVersion:"Go 运行时",siPlatform:"平台",siStartTime:"启动时间",siUptime:"运行时长",siLanURLs:"内网访问地址",maxUploadSize:"最大上传大小 (MB)",maxUploadSizeDesc:"上传文件的最大大小。默认：50 MB。",maxSitesPerUser:"每用户站点上限",maxSitesPerUserDesc:"每个用户可创建的最大站点数。设为 0 则不限制。",siteBaseUrl:"站点访问地址",siteBaseUrlDesc:"配置专用域名来提供上传站点内容（如 https://sites.example.com）。将用户上传的 HTML 与管理面板域名隔离，防止 Cookie/Session 注入。需带 http:// 或 https:// 协议（省略默认 https）。",siteBaseUrlWarn:"未配置 - 上传的站点与管理面板共用同一域名。存在安全风险：用户上传的 HTML 可能访问管理员的 Cookie 和会话。建议为站点内容配置独立域名（如 https://sites.example.com）。",siteBaseUrlPh:"https://sites.example.com"}};
Object.assign(i18n.en,{adminPanel:"Admin",recentSites:"Latest sites",roleAdmin:"Admin",roleUser:"User",access:"Access",sourceRepo:"Source code on GitHub",statUsersNote:"registered accounts",statSitesNote:"of {max} per user allowed",statAdminsNote:"full access",toggleTheme:"Toggle theme",noMatch:"Nothing matches your search",email:"Email",files:"Files"});
Object.assign(i18n.zh,{adminPanel:"管理",recentSites:"最新站点",roleAdmin:"管理员",roleUser:"普通用户",access:"访问控制",sourceRepo:"源码（GitHub）",statUsersNote:"已注册账号",statSitesNote:"每用户上限 {max} 个",statAdminsNote:"拥有完整权限",toggleTheme:"切换主题",noMatch:"没有匹配结果",email:"邮箱",files:"文件"});
function t(k){return(i18n[lang]||i18n.en)[k]||(i18n.en[k]||k)}
function setLang(l){lang=l;try{localStorage.setItem("lang",l)}catch(e){}renderAdmin()}
function esc(s){return String(s||"").replace(/[&<>"']/g,function(c){return{"&":"&amp;","<":"&lt;",">":"&gt;",'"':"&quot;","'":"&#39;"}[c]})}
function siteUrl(u){return BASE+(u||"").replace(/^\//,"")}
function fmtDate(s){if(!s)return"-";var d=new Date(s);return d.toLocaleDateString()+" "+d.toLocaleTimeString([],{hour:"2-digit",minute:"2-digit"})}
function toast(msg,type){
type=type||"success";
var host=document.getElementById("toast-host");
if(!host){host=document.createElement("div");host.id="toast-host";host.setAttribute("aria-live","polite");host.setAttribute("aria-atomic","false");document.body.appendChild(host)}
var el=document.createElement("div");el.className="toast "+type+" show";el.textContent=msg;host.appendChild(el);
var ms=type==="error"?5000:2600;
setTimeout(function(){el.classList.remove("show");setTimeout(function(){el.remove()},200)},ms)
}
function togglePwd(btn){var inp=btn.parentElement.querySelector("input");if(!inp)return;var show=inp.type=="password";inp.type=show?"text":"password";btn.innerHTML=show?PWD_SHOW_ICON:PWD_HIDE_ICON;btn.setAttribute("aria-label",t(show?"hidePassword":"showPassword"));btn.setAttribute("aria-pressed",show?"true":"false")}
var PWD_SHOW_ICON='<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7z"/><circle cx="12" cy="12" r="3"/></svg>';
var PWD_HIDE_ICON='<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-10-8-10-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 10 8 10 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/><line x1="1" y1="1" x2="23" y2="23"/></svg>';
function getToken(){return"cookie"}
function clearToken(){try{localStorage.removeItem("vibecast_token")}catch(e){}}
function customConfirm(msg,onConfirm,opts){
opts=opts||{};
var o=document.createElement("div");o.className="modal-overlay";o.style.zIndex="320";
var danger=opts.danger!==false;
o.innerHTML='<div class="modal" role="alertdialog" aria-modal="true" aria-labelledby="cf-title" aria-describedby="cf-desc"><h3 id="cf-title">'+esc(opts.title||t("confirm"))+'</h3><p id="cf-desc">'+esc(msg)+'</p><div class="modal-actions"><button type="button" class="btn btn-outline" id="cf-no">'+t("cancel")+'</button><button type="button" class="btn '+(danger?"btn-danger":"btn-primary")+'" id="cf-yes">'+(opts.confirmLabel||t("confirm"))+'</button></div></div>';
document.body.appendChild(o);
var prevFocus=document.activeElement,closed=false;
function close(){if(closed)return;closed=true;document.removeEventListener("keydown",onKey);o.classList.remove("show");setTimeout(function(){o.remove()},180);try{if(prevFocus&&prevFocus.focus)prevFocus.focus()}catch(e){}}
function onKey(e){if(e.key==="Escape"){e.preventDefault();close()}}
document.addEventListener("keydown",onKey);
setTimeout(function(){o.classList.add("show")},10);
o.addEventListener("click",function(e){if(e.target===o)close()});
var no=o.querySelector("#cf-no"),yes=o.querySelector("#cf-yes");
no.addEventListener("click",close);
yes.addEventListener("click",function(){close();onConfirm()});
setTimeout(function(){no.focus()},30)
}
function api(path,opts){
opts=opts||{};
	var headers={"Content-Type":"application/json","Accept-Language":lang||"en"};
	if(opts.headers)Object.assign(headers,opts.headers);
	return fetch(API+path,Object.assign({},opts,{headers:headers,credentials:"same-origin"})).then(function(r){
	if(r.status===401){clearToken();location.href=BASE+"dashboard"}
return r.json().catch(function(){return{error:"network error"}}).then(function(data){if(!r.ok)throw new Error(data.error||"request failed");return data})
})
}
function checkAuth(){return api("/auth/me").then(function(d){if(!d.data||!d.data.isAdmin)throw new Error("not admin");return d.data}).catch(function(){clearToken();location.href=BASE+"dashboard"})}
function doLogout(){api("/auth/logout",{method:"POST"}).then(function(){clearToken();location.href=BASE}).catch(function(){clearToken();location.href=BASE})}
var adminTab="overview";
function renderAdmin(){
var th='<button type="button" id="theme-toggle" class="btn-icon" onclick="toggleTheme()" aria-label="'+t("toggleTheme")+'" title="'+t("toggleTheme")+'">'+(getTheme()==="dark"?'` + iconSun + `':'` + iconMoon + `')+'</button>';
var lh='<div class="lang-toggle" role="group" aria-label="Language"><button type="button" class="'+(lang==="en"?"active":"")+'" onclick="setLang(\'en\')" aria-pressed="'+(lang==="en")+'">EN</button><button type="button" class="'+(lang==="zh"?"active":"")+'" onclick="setLang(\'zh\')" aria-pressed="'+(lang==="zh")+'">中文</button></div>';
var icons={overview:'<svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="7" height="7" rx="1"/><rect x="14" y="3" width="7" height="7" rx="1"/><rect x="3" y="14" width="7" height="7" rx="1"/><rect x="14" y="14" width="7" height="7" rx="1"/></svg>',users:'<svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>',sites:'<svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="3" width="20" height="14" rx="2"/><line x1="8" y1="21" x2="16" y2="21"/><line x1="12" y1="17" x2="12" y2="21"/></svg>',settings:'<svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>',maintenance:'<svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/></svg>',system:'<svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>'};
var navItems=[["overview","tabOverview"],["users","tabUsers"],["sites","tabSites"],["settings","tabSettings"],["maintenance","tabMaintenance"],["system","tabSystem"]];
var navHtml='';
for(var i=0;i<navItems.length;i++){var k=navItems[i][0],lk=navItems[i][1];navHtml+='<button type="button" class="admin-nav-item'+(adminTab===k?" active":"")+'" onclick="switchTab(\''+k+'\')">'+(icons[k]||'')+'<span>'+t(lk)+'</span></button>'}
document.getElementById("app").innerHTML='<nav class="navbar"><div class="logo">` + logoIcon + `<span class="logo-text">Vibe<span class="accent">cast</span></span><span class="admin-tag">'+t("adminPanel")+'</span></div><div class="nav-right"><div class="nav-section nav-preferences">'+th+lh+'</div><div class="nav-section nav-account"><a href="'+BASE+'dashboard" class="btn-link">` + iconLayout + `<span>'+t("dashboard")+'</span></a><button type="button" class="btn-link" onclick="doLogout()">` + iconLogout + `<span>'+t("logout")+'</span></button></div></div></nav>'
+'<div class="admin-layout"><div class="admin-sidebar">'+navHtml+'<div class="admin-sidebar-foot"><span class="ver">Vibecast <span id="ver"></span></span></div></div>'
+'<div class="admin-main">'
+'<div class="admin-tab-content" id="tab-overview"><div class="card"><div class="card-header"><h2>'+t("overview")+'</h2></div><div class="card-body"><div id="stats" class="stats-grid"></div></div></div><div class="card"><div class="card-header"><h2>'+t("recentSites")+'</h2><a class="btn btn-ghost btn-sm" href="#" onclick="switchTab(\'sites\');return false">'+t("allSites")+'</a></div><div class="card-body" id="overview-recent"></div></div></div>'
+'<div class="admin-tab-content" id="tab-users"><div class="card"><div class="card-header"><h2>'+t("users")+'</h2></div><div class="card-body"><div class="list-toolbar"><input type="text" id="user-search" placeholder="'+t("searchUsersPh")+'" aria-label="'+t("searchUsersPh")+'" onkeydown="if(event.key===\'Enter\')searchUsers()" value="'+esc(userSearch)+'"></div><div id="users"></div></div></div></div>'
+'<div class="admin-tab-content" id="tab-sites"><div class="card"><div class="card-header"><h2>'+t("allSites")+'</h2></div><div class="card-body"><div class="list-toolbar"><input type="text" id="admin-site-search" placeholder="'+t("searchSitesPh")+'" aria-label="'+t("searchSitesPh")+'" onkeydown="if(event.key===\'Enter\')searchAdminSites()" value="'+esc(adminSiteSearch)+'"></div><div id="sites"></div></div></div></div>'
+'<div class="admin-tab-content" id="tab-settings"><div class="card"><div class="card-header"><h2>'+t("settings")+'</h2></div><div class="card-body"><div id="settings"></div></div></div></div>'
+'<div class="admin-tab-content" id="tab-maintenance"><div class="card"><div class="card-header"><h2>'+t("cleanup")+'</h2></div><div class="card-body"><div id="cleanup-section"><p class="tab-desc">'+t("cleanupDesc")+'</p><button type="button" class="btn btn-outline" onclick="checkCleanup()">'+t("checkOrphans")+'</button><div id="cleanup-result" class="result-box"></div></div></div></div></div>'
+'<div class="admin-tab-content" id="tab-system"><div class="card"><div class="card-header"><h2>'+t("systemInfo")+'</h2></div><div class="card-body"><div id="sysinfo" class="info-grid"></div><div class="info-about"><a class="about-link" href="` + projectRepoURL + `" target="_blank" rel="noopener">` + iconGithub + `<span>'+t("sourceRepo")+'</span></a></div><div class="update-inline"><div id="update-section"><p>'+t("currentVersion")+': <span id="update-current-ver" class="ver-inline">-</span></p><button type="button" class="btn btn-outline" id="update-check-btn" onclick="checkUpdate()">'+t("updateCheck")+'</button><div id="update-result" class="result-box"></div></div></div></div></div>'
+'</div></div>';
switchTab(adminTab);
loadStats();loadSettings();loadUsers();loadSites();loadSystemInfo();
fetch(BASE+"api/version").then(function(r){return r.json()}).then(function(d){var ver=(d.data?d.data.version:"dev");var v=document.getElementById("ver");if(v)v.textContent=ver;var uv=document.getElementById("update-current-ver");if(uv)uv.textContent=ver}).catch(function(){})
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
var items=[[t("siVersion"),r.version||"dev","accent"],[t("siPlatform"),(r.os||"?")+"/"+(r.arch||"?"),""],[t("siGoVersion"),r.goVersion||"-",""],[t("siListenAddr"),r.listenAddr||"-",""],[t("siStoragePath"),r.storagePath||"-",""],[t("siDbPath"),r.dbPath||"-",""],[t("siStartTime"),r.startTime?fmtDateTime(r.startTime):"-",""],[t("siUptime"),r.startTime?fmtUptime(r.startTime):"-",""]];
var h='';
for(var i=0;i<items.length;i++){h+='<div class="info-item"><div class="info-label">'+items[i][0]+'</div><div class="info-value'+(items[i][2]?" "+items[i][2]:"")+'">'+esc(items[i][1])+'</div></div>'}
var ips=r.localIPs||[];var port=r.listenAddr?r.listenAddr.split(":").pop():"8080";
if(ips.length>0){var urls="";for(var i=0;i<ips.length;i++){urls+='<a href="http://'+ips[i]+':'+port+'" target="_blank" style="display:block;font-family:var(--mono);font-size:.8rem;color:var(--accent);margin-bottom:.2rem">http://'+ips[i]+':'+port+'</a>'}h+='<div class="info-item"><div class="info-label">'+t("siLanURLs")+'</div><div class="info-value">'+urls+'</div></div>'}
el.innerHTML=h}).catch(function(e){var el=document.getElementById("sysinfo");if(el)el.innerHTML='<div style="color:var(--danger);font-size:.8rem">'+esc(e.message)+'</div>'})
}
function fmtDateTime(v){try{var d=new Date(v);if(isNaN(d.getTime()))return String(v||"");return d.toLocaleString(lang==="zh"?"zh-CN":"en-GB",{year:"numeric",month:"2-digit",day:"2-digit",hour:"2-digit",minute:"2-digit"})}catch(e){return String(v||"")}}
function fmtUptime(startStr){var start=new Date(startStr);var now=new Date();var diff=Math.floor((now-start)/1000);if(diff<60)return diff+"s";var m=Math.floor(diff/60),s=diff%60;if(m<60)return m+"m "+s+"s";var h=Math.floor(m/60),mm=m%60;if(h<24)return h+"h "+mm+"m";var d=Math.floor(h/24),hh=h%24;return d+"d "+hh+"h"}
function loadStats(){
api("/admin/stats").then(function(d){var s=d.data||{},el=document.getElementById("stats");if(!el)return;
el.innerHTML='<div class="stat-card"><div class="stat-top"><span class="label">'+t("userCount")+'</span><span class="stat-icon">` + iconUsers + `</span></div><div class="num">'+s.users+'</div><div class="stat-note">'+t("statUsersNote")+'</div></div>'
+'<div class="stat-card"><div class="stat-top"><span class="label">'+t("siteCount")+'</span><span class="stat-icon">` + iconGlobe + `</span></div><div class="num">'+s.sites+'</div><div class="stat-note">'+t("statSitesNote").replace("{max}",maxSitesLimit)+'</div></div>'
+'<div class="stat-card"><div class="stat-top"><span class="label">'+t("adminCount")+'</span><span class="stat-icon">` + iconShield + `</span></div><div class="num">'+s.admins+'</div><div class="stat-note">'+t("statAdminsNote")+'</div></div>'}).catch(function(e){toast(e.message,"error")});
loadRecentSites()
}
function loadRecentSites(){
var el=document.getElementById("overview-recent");if(!el)return;
api("/admin/sites?page=1&perPage=5").then(function(d){
var sites=(d.data||{}).items||[];
if(!sites.length){el.innerHTML='<div class="empty"><span class="empty-icon">` + iconInbox + `</span><strong>'+t("noSites")+'</strong></div>';return}
var h='<div class="table-wrap"><table><thead><tr><th>'+t("name")+'</th><th>'+t("owner")+'</th><th>'+t("visits")+'</th><th>'+t("created")+'</th></tr></thead><tbody>';
for(var i=0;i<sites.length;i++){var s=sites[i],vs=s.visits||{total:0};
h+='<tr><td><div class="cell-strong">'+esc(s.name)+'</div><a class="cell-link" href="'+esc(s.url)+'" target="_blank" rel="noopener">'+esc(s.slug)+'</a></td>'
+'<td class="cell-email" title="'+esc(s.ownerEmail||"")+'">'+esc(s.ownerEmail||"—")+'</td>'
+'<td class="num">'+vs.total+'</td><td class="num">'+esc(fmtDate(s.createdAt||""))+'</td></tr>'}
h+='</tbody></table></div>';el.innerHTML=h}).catch(function(e){el.innerHTML='<div class="empty">'+esc(e.message)+'</div>'})
}
var maxSitesLimit=30;
function loadSettings(){
api("/admin/settings").then(function(d){var s=d.data;
maxSitesLimit=s.maxSitesPerUser||30;
var regOn=s.openRegistration,pubOn=s.allowPublicAccess!==false,drOn=s.domainRestriction;
var mus=s.maxUploadSize||50;
var h='<div class="toggle-row"><div class="toggle-info"><div class="toggle-label">'+t("openReg")+'</div><div class="toggle-desc">'+t("openRegDesc")+'</div></div><button type="button" class="toggle-switch '+(regOn?"on":"")+'" role="switch" aria-checked="'+(regOn?"true":"false")+'" aria-label="'+t("openReg")+'" onclick="toggleReg()"></button></div><div class="toggle-row"><div class="toggle-info"><div class="toggle-label">'+t("publicAccess")+'</div><div class="toggle-desc">'+t("publicAccessDesc")+'</div></div><button type="button" class="toggle-switch '+(pubOn?"on":"")+'" role="switch" aria-checked="'+(pubOn?"true":"false")+'" aria-label="'+t("publicAccess")+'" onclick="togglePub()"></button></div><div class="toggle-row"><div class="toggle-info"><div class="toggle-label">'+t("domainRestriction")+'</div><div class="toggle-desc">'+t("domainRestrictionDesc")+'</div></div><button type="button" class="toggle-switch '+(drOn?"on":"")+'" role="switch" aria-checked="'+(drOn?"true":"false")+'" aria-label="'+t("domainRestriction")+'" onclick="toggleDomainRestriction()"></button></div>';
if(drOn){h+='<div class="settings-block"><label class="num-label" for="allowed-domains">'+t("domainRestriction")+'</label><textarea id="allowed-domains" rows="4" placeholder="example.com&#10;gmail.com">'+esc(s.allowedDomains||"")+'</textarea><div style="margin-top:8px"><button type="button" class="btn btn-outline btn-sm" onclick="saveDomains()">'+t("save")+'</button></div></div>'}
var mspu=s.maxSitesPerUser||30;
h+='<div class="num-setting"><div class="num-info"><div class="num-label">'+t("maxUploadSize")+'</div><div class="num-desc">'+t("maxUploadSizeDesc")+'</div></div><div class="num-control"><input type="number" id="max-upload-size" value="'+mus+'" min="1" max="1024"><span class="unit">MB</span><button class="btn btn-outline btn-sm" onclick="saveUploadSize()">'+t("save")+'</button></div></div>';
h+='<div class="num-setting"><div class="num-info"><div class="num-label">'+t("maxSitesPerUser")+'</div><div class="num-desc">'+t("maxSitesPerUserDesc")+'</div></div><div class="num-control"><input type="number" id="max-sites-per-user" value="'+mspu+'" min="0" max="10000"><button class="btn btn-outline btn-sm" onclick="saveMaxSites()">'+t("save")+'</button></div></div>';
var sad=s.siteBaseUrl||"";
h+='<div class="url-setting"><div class="num-info"><div class="num-label">'+t("siteBaseUrl")+'</div><div class="num-desc">'+t("siteBaseUrlDesc")+'</div></div>';
if(!sad){h+='<div class="url-warn">'+t("siteBaseUrlWarn")+'</div>'}
h+='<div class="num-control url-control"><input type="text" id="site-base-url" placeholder="'+t("siteBaseUrlPh")+'" value="'+esc(sad)+'"><button class="btn btn-outline btn-sm" onclick="saveSiteBaseUrl()">'+t("save")+'</button></div></div>';
document.getElementById("settings").innerHTML=h}).catch(function(e){toast(e.message,"error")})
}
function toggleReg(){api("/admin/settings").then(function(d){var c=d.data;return api("/admin/settings",{method:"PUT",body:JSON.stringify({openRegistration:!c.openRegistration,allowPublicAccess:c.allowPublicAccess!==false,domainRestriction:c.domainRestriction,allowedDomains:c.allowedDomains||"",maxUploadSize:c.maxUploadSize||50,maxSitesPerUser:c.maxSitesPerUser||30,siteBaseUrl:c.siteBaseUrl||""})}).then(function(){toast(t("settingsSaved"));loadSettings()})}).catch(function(e){toast(e.message,"error")})}
function togglePub(){api("/admin/settings").then(function(d){var c=d.data;return api("/admin/settings",{method:"PUT",body:JSON.stringify({openRegistration:c.openRegistration,allowPublicAccess:!(c.allowPublicAccess!==false),domainRestriction:c.domainRestriction,allowedDomains:c.allowedDomains||"",maxUploadSize:c.maxUploadSize||50,maxSitesPerUser:c.maxSitesPerUser||30,siteBaseUrl:c.siteBaseUrl||""})}).then(function(){toast(t("settingsSaved"));loadSettings()})}).catch(function(e){toast(e.message,"error")})}
function toggleDomainRestriction(){api("/admin/settings").then(function(d){var c=d.data;return api("/admin/settings",{method:"PUT",body:JSON.stringify({openRegistration:c.openRegistration,allowPublicAccess:c.allowPublicAccess!==false,domainRestriction:!c.domainRestriction,allowedDomains:c.allowedDomains||"",maxUploadSize:c.maxUploadSize||50,maxSitesPerUser:c.maxSitesPerUser||30,siteBaseUrl:c.siteBaseUrl||""})}).then(function(){toast(t("settingsSaved"));loadSettings()})}).catch(function(e){toast(e.message,"error")})}
function saveDomains(){var v=document.getElementById("allowed-domains").value;api("/admin/settings").then(function(d){return api("/admin/settings",{method:"PUT",body:JSON.stringify({openRegistration:d.data.openRegistration,allowPublicAccess:d.data.allowPublicAccess!==false,domainRestriction:d.data.domainRestriction,allowedDomains:v,maxUploadSize:d.data.maxUploadSize||50,maxSitesPerUser:d.data.maxSitesPerUser||30,siteBaseUrl:d.data.siteBaseUrl||""})}).then(function(){toast(t("settingsSaved"));loadSettings()})}).catch(function(e){toast(e.message,"error")})}
function saveUploadSize(){var v=parseInt(document.getElementById("max-upload-size").value)||50;if(v<1)v=1;if(v>1024)v=1024;api("/admin/settings").then(function(d){return api("/admin/settings",{method:"PUT",body:JSON.stringify({openRegistration:d.data.openRegistration,allowPublicAccess:d.data.allowPublicAccess!==false,domainRestriction:d.data.domainRestriction,allowedDomains:d.data.allowedDomains||"",maxUploadSize:v,siteBaseUrl:d.data.siteBaseUrl||""})}).then(function(){toast(t("settingsSaved"));loadSettings()})}).catch(function(e){toast(e.message,"error")})}
function saveMaxSites(){var v=parseInt(document.getElementById("max-sites-per-user").value)||30;if(v<0)v=0;if(v>10000)v=10000;api("/admin/settings").then(function(d){return api("/admin/settings",{method:"PUT",body:JSON.stringify({openRegistration:d.data.openRegistration,allowPublicAccess:d.data.allowPublicAccess!==false,domainRestriction:d.data.domainRestriction,allowedDomains:d.data.allowedDomains||"",maxUploadSize:d.data.maxUploadSize||50,maxSitesPerUser:v,siteBaseUrl:d.data.siteBaseUrl||""})}).then(function(){toast(t("settingsSaved"));loadSettings()})}).catch(function(e){toast(e.message,"error")})}
function saveSiteBaseUrl(){var v=document.getElementById("site-base-url").value.trim();if(v&&v.indexOf("://")<0){v="https://"+v;document.getElementById("site-base-url").value=v}api("/admin/settings").then(function(d){return api("/admin/settings",{method:"PUT",body:JSON.stringify({openRegistration:d.data.openRegistration,allowPublicAccess:d.data.allowPublicAccess!==false,domainRestriction:d.data.domainRestriction,allowedDomains:d.data.allowedDomains||"",maxUploadSize:d.data.maxUploadSize||50,maxSitesPerUser:d.data.maxSitesPerUser||30,siteBaseUrl:v})}).then(function(){toast(t("settingsSaved"));loadSettings()})}).catch(function(e){toast(e.message,"error")})}
function checkUpdate(){
var el=document.getElementById("update-result");var btn=document.getElementById("update-check-btn");
if(btn){btn.disabled=true;btn.textContent=t("checkingUpdate")}
el.innerHTML='<span style="color:var(--dim);font-size:.8rem">'+t("checkingUpdate")+'</span>';
api("/admin/update/check").then(function(d){
var r=d.data||{};
if(btn){btn.disabled=false;btn.textContent=t("updateCheck")}
if(r.error){el.innerHTML='<span style="color:var(--danger);font-size:.8rem">` + iconCheck + ` '+esc(r.error)+'</span>';return}
if(r.updateAvailable){
var h='<div style="margin-bottom:.6rem"><span style="color:var(--accent);font-size:.8rem;font-family:var(--mono)">` + iconCheck + ` '+t("updateAvailable")+'</span></div>';
h+='<div style="font-size:.8rem;color:var(--dim);margin-bottom:.4rem">'+t("currentVersion")+': <span style="font-family:var(--mono);color:var(--text)">v'+esc(r.currentVersion)+'</span> → '+t("latestVersion")+': <span style="font-family:var(--mono);color:var(--accent)">v'+esc(r.latestVersion)+'</span></div>';
if(r.releaseNotes){h+='<details style="margin-bottom:.6rem"><summary style="font-size:.75rem;cursor:pointer;color:var(--dim)">'+t("releaseNotes")+'</summary><div style="font-size:.75rem;color:var(--dim);margin-top:.4rem;white-space:pre-wrap;max-height:200px;overflow-y:auto">'+esc(r.releaseNotes)+'</div></details>'}
h+='<button class="btn btn-primary btn-sm" id="update-apply-btn" onclick="applyUpdate()">'+t("updateNow")+'</button>';
el.innerHTML=h}else{
var h2='<span style="color:var(--dim);font-size:.8rem">` + iconCheck + ` '+t("noUpdate")+'</span>';
if(r.latestVersion){h2+=' <span style="color:var(--dim);font-size:.75rem;font-family:var(--mono)">v'+esc(r.latestVersion)+'</span>'}
el.innerHTML=h2}
}).catch(function(e){if(btn){btn.disabled=false;btn.textContent=t("updateCheck")}el.innerHTML='<span style="color:var(--danger);font-size:.8rem">'+esc(e.message)+'</span>'})
}
function applyUpdate(){
var btn=document.getElementById("update-apply-btn");var el=document.getElementById("update-result");
if(btn){btn.disabled=true;btn.textContent=t("updating")}
el.innerHTML='<div style="margin-top:.6rem"><div style="font-size:.8rem;color:var(--dim);margin-bottom:.5rem">'+t("updateDownloadProgress")+'</div><div style="width:100%;height:6px;background:var(--surface-2);border-radius:3px;overflow:hidden"><div id="update-progress-bar" style="width:0%;height:100%;background:var(--accent);border-radius:3px;transition:width .3s"></div></div><div id="update-progress-text" style="font-size:.7rem;color:var(--dim);font-family:var(--mono);margin-top:.3rem">0%</div></div>';
customConfirm(t("updateNow")+"?",function(){
api("/admin/update/apply",{method:"POST"}).then(function(){
// Start polling
var pollTimer=setInterval(function(){
api("/admin/update/status").then(function(d){
var st=d.data||{};
var bar=document.getElementById("update-progress-bar");
var txt=document.getElementById("update-progress-text");
if(bar&&st.status==="downloading"){bar.style.width=st.progress+"%";txt.textContent=t("updateDownloadProgress")+" "+st.progress+"%"}
else if(bar&&st.status==="verifying"){bar.style.width="100%";txt.textContent=t("updateVerifying")}
else if(bar&&st.status==="installing"){bar.style.width="100%";txt.textContent=t("updateInstalling")}
else if(st.status==="done"){clearInterval(pollTimer);
var h='<div style="color:var(--accent);font-size:.8rem;font-family:var(--mono)">` + iconCheck + ` '+t("updateSuccess")+'</div>';
if(st.prevVersion&&st.newVersion){h+='<div style="font-size:.75rem;color:var(--dim);margin-top:.3rem">v'+esc(st.prevVersion)+' → v'+esc(st.newVersion)+'</div>'}
h+='<button class="btn btn-primary btn-sm" style="margin-top:.6rem" onclick="restartServer()">'+t("updateRestart")+'</button>';
el.innerHTML=h;toast(t("updateSuccess"),"success")}
else if(st.status==="error"){clearInterval(pollTimer);
var h='<div style="color:var(--danger);font-size:.8rem">` + iconCheck + ` '+t("updateFailed")+': '+esc(st.error)+'</div><button class="btn btn-outline" style="margin-top:.4rem" onclick="checkUpdate()">'+t("updateCheck")+'</button>';
el.innerHTML=h;toast(st.error,"error")}
}).catch(function(){})
},1000)
}).catch(function(e){
clearInterval(pollTimer);
var h='<div style="color:var(--danger);font-size:.8rem">` + iconCheck + ` '+t("updateFailed")+': '+esc(e.message)+'</div><button class="btn btn-outline" style="margin-top:.4rem" onclick="checkUpdate()">'+t("updateCheck")+'</button>';
el.innerHTML=h;toast(e.message,"error")
})
})
}
function restartServer(){
var el=document.getElementById("update-result");
el.innerHTML+='<div style="margin-top:.6rem"><span style="color:var(--accent);font-size:.8rem;font-family:var(--mono)">` + iconCheck + ` '+t("updateRestarting")+'</span></div>';
api("/admin/update/restart",{method:"POST"}).then(function(){
toast(t("updateRestartSuccess"),"success");
function tryReload(attempt){
var max=15;
if(attempt>max){
el.innerHTML+='<div style="color:var(--danger);font-size:.8rem;margin-top:.4rem">` + iconCheck + ` '+t("updateRestartFailed")+'</div>';
return
}
fetch(BASE+"api/version",{cache:"no-store"}).then(function(r){
if(!r.ok)throw new Error("HTTP "+r.status);
return r.json()
}).then(function(){
location.reload()
}).catch(function(){
setTimeout(function(){tryReload(attempt+1)},2000)
})
}
setTimeout(function(){tryReload(1)},2500)
}).catch(function(e){
el.innerHTML+='<div style="color:var(--danger);font-size:.8rem">'+esc(e.message)+'</div>'
})
}
function checkCleanup(){
var el=document.getElementById("cleanup-result");
el.innerHTML='<span class="result-loading">'+t("scanning")+'</span>';
api("/admin/cleanup").then(function(d){
var r=d.data||{},orphans=r.orphans||[];
if(!orphans.length){el.innerHTML='<span class="result-ok">` + iconCheck + `'+t("noOrphans")+'</span>';return}
var h='<div class="result-warn">'+t("orphanFound")+' · '+orphans.length+'</div>';
h+='<div class="orphan-list">';
for(var i=0;i<orphans.length;i++){h+='<div>` + iconFolder + `'+esc(orphans[i])+'</div>'}
h+='</div><button type="button" class="btn btn-danger btn-sm" onclick="confirmCleanup()">'+t("confirmDelete")+'</button>';
el.innerHTML=h}).catch(function(e){el.innerHTML='<span class="result-err">` + iconAlert + `'+esc(e.message)+'</span>'})
}
function confirmCleanup(){customConfirm(t("cleanupConfirmWarning"),function(){api("/admin/cleanup",{method:"POST",body:JSON.stringify({confirm:true})}).then(function(d){var r=d.data||{};toast(t("cleanupDone")+" ("+r.deleted+"/"+r.total+")");document.getElementById("cleanup-result").innerHTML='<span class="result-ok">` + iconCheck + `'+t("cleanupDone")+' · '+r.deleted+'/'+r.total+'</span>';loadStats()}).catch(function(e){toast(e.message,"error")})})}
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
if(!users.length){el.innerHTML='<div class="empty"><span class="empty-icon">` + iconInbox + `</span><strong>'+t("noUsers")+'</strong></div>'+pg;return}
var h='<div class="table-wrap"><table><thead><tr><th class="col-hide-sm">ID</th><th>'+t("email")+'</th><th class="col-hide-sm">'+t("role")+'</th><th class="col-hide-sm">'+t("created")+'</th><th>'+t("actions")+'</th></tr></thead><tbody>';
for(var i=0;i<users.length;i++){var u=users[i],b=u.isAdmin?'<span class="badge badge-admin">'+t("roleAdmin")+'</span>':'<span class="badge badge-user">'+t("roleUser")+'</span>',btn=u.isAdmin?'<button type="button" class="btn btn-ghost btn-sm" onclick="toggleAdmin('+u.id+')">'+t("demote")+'</button>':'<button type="button" class="btn btn-ghost btn-sm" onclick="toggleAdmin('+u.id+')">'+t("promote")+'</button>';
h+='<tr><td class="num col-hide-sm">'+u.id+'</td><td class="cell-email" title="'+esc(u.email)+'"><span class="cell-strong">'+esc(u.email)+'</span><span class="cell-sub">'+b+'</span></td><td class="col-hide-sm">'+b+'</td><td class="num col-hide-sm">'+fmtDate(u.createdAt)+'</td><td class="cell-actions">'+btn+'<button type="button" class="btn btn-danger-ghost btn-sm" onclick="delUser('+u.id+')">'+t("delete")+'</button></td></tr>'}
h+='</tbody></table></div>'+pg;el.innerHTML=h}).catch(function(e){toast(e.message,"error")})
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
if(!sites.length){el.innerHTML='<div class="empty"><span class="empty-icon">` + iconInbox + `</span><strong>'+t("noSites")+'</strong></div>'+pg;return}
var h='<div class="table-wrap"><table class="sites-table"><thead><tr><th>'+t("name")+'</th><th class="col-hide-sm">'+t("owner")+'</th><th>'+t("access")+'</th><th class="col-hide-sm">'+t("visits")+'</th><th>'+t("actions")+'</th></tr></thead><tbody>';
for(var i=0;i<sites.length;i++){var s=sites[i],b="";
if(s.publicAccessDisabled&&!s.protected){b='<span class="badge badge-disabled">'+t("accessDisabled")+'</span>'}
else if(s.protected){b='<span class="badge badge-protected">'+t("protected")+'</span>'}
else{b='<span class="badge badge-public">'+t("public")+'</span>'}
if(s.orgOpen){b+=' <span class="badge badge-org">'+t("orgBadge")+'</span>'}
var vs=s.visits||{today:0,month:0,total:0};
h+='<tr class="site-row" id="row-'+s.id+'" onclick="adminToggleDetail('+s.id+')">'
+'<td><div class="cell-lead"><span class="row-caret" aria-hidden="true">` + iconChevronRight + `</span><div class="cell-lead-text"><div class="cell-strong">'+esc(s.name)+'</div><a class="cell-link" href="'+esc(s.url)+'" target="_blank" rel="noopener" onclick="event.stopPropagation()">'+esc(s.slug)+'</a></div></div></td>'
+'<td class="cell-email col-hide-sm" title="'+esc(s.ownerEmail||"")+'">'+esc(s.ownerEmail||"—")+'</td>'
+'<td>'+b+'</td>'
+'<td class="num col-hide-sm" title="'+t("visitsToday")+' '+vs.today+' · '+t("visitsMonth")+' '+vs.month+'">'+vs.total+'</td>'
+'<td class="cell-actions" onclick="event.stopPropagation()"><button type="button" class="btn btn-ghost btn-sm" onclick="delSite('+s.id+')">'+t("delete")+'</button></td></tr>'
+'<tr class="admin-detail-row" id="detail-'+s.id+'"><td colspan="5"></td></tr>'}
h+='</tbody></table></div>'+pg;el.innerHTML=h}).catch(function(e){toast(e.message,"error")})
}
function delSite(id){customConfirm(t("deleteSiteConfirm"),function(){api("/admin/sites/"+id,{method:"DELETE"}).then(function(){toast(t("siteDeleted"));loadSites();loadStats()}).catch(function(e){toast(e.message,"error")})})}
function formatSize(n){if(n<1024)return n+" B";if(n<1048576)return(n/1024).toFixed(1)+" KB";if(n<1073741824)return(n/1048576).toFixed(1)+" MB";return(n/1073741824).toFixed(1)+" GB"}
function adminToggleDetail(id){var r=document.getElementById("row-"+id);var d=document.getElementById("detail-"+id);if(!r||!d)return;var wasShown=d.classList.contains("show");d.classList.toggle("show");r.classList.toggle("expanded");if(!wasShown){adminLoadFileTree(id);adminLoadSitePwd(id)}}
function adminLoadFileTree(id){api("/admin/sites/"+id+"/files").then(function(d){var data=d.data||{},files=Array.isArray(data)?data:(data.items||[]);var el=document.getElementById("detail-"+id);if(!el)return;var td=el.querySelector("td");if(!td)return;
var pwd=el.querySelector(".admin-detail-body");
var h='<div class="admin-detail-body"><div class="detail-label">'+t("files")+'</div>';
if(!files.length){h+='<div class="file-empty">` + iconFile + `'+t("noFiles")+'</div></div>';td.innerHTML=h;return}
h+='<div class="file-list admin-file-list">';
for(var i=0;i<files.length;i++){var f=files[i],icon=f.dir?'` + iconFolder + `':'` + iconFile + `',sz=f.dir?"—":formatSize(f.size);
h+='<div class="file-row"><span class="file-icon">'+icon+'</span><span class="file-name">'+esc(f.name)+(f.dir?"/":"")+'</span><span class="file-size">'+sz+'</span></div>'}
h+='</div></div>';td.innerHTML=h}).catch(function(){})}
function adminLoadSitePwd(id){api("/admin/sites/"+id+"/password").then(function(d){var p=(d.data||{}).password||"";var el=document.getElementById("admin-pwd-"+id);if(el)el.textContent=p}).catch(function(){})}

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
if(isOwner){h+='<button class="btn btn-outline btn-sm org-danger-btn" onclick="deleteOrg()">'+t("orgDelete")+'</button>'}
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

// passwordPageHTML renders the password gate for a protected site. It carries
// the requested sub-path (next) so deep links resume where the visitor meant to
// go, and it reports an error message inline when a submission is rejected.
func passwordPageHTML(slug, siteName, next, errMsg string) string {
	errBlock := ""
	if errMsg != "" {
		errBlock = `<div class="gate-err" role="alert">` + iconAlert + `<span>` + escHTML(errMsg) + `</span></div>`
	}
	action := "?next=" + url.QueryEscape(next)
	if next == "" {
		action = ""
	}
	css := uiTokensCSS + uiGateCSS
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<meta name="robots" content="noindex, nofollow">
%s
<title>%s - Password Required</title>
<style>%s</style>
<script>var _t="light";try{_t=localStorage.getItem("theme")||"light"}catch(e){}document.documentElement.setAttribute("data-theme",_t)
var PWD_SHOWN=%q,PWD_HIDDEN=%q;
function togglePwd(btn){var f=btn.parentElement,inp=f&&f.querySelector("input");if(!inp)return;var show=inp.type==="password";inp.type=show?"text":"password";btn.setAttribute("aria-label",show?"Hide password":"Show password");btn.setAttribute("aria-pressed",show?"true":"false");btn.innerHTML=show?PWD_SHOWN:PWD_HIDDEN}
document.addEventListener("DOMContentLoaded",function(){var f=document.querySelector("form");if(!f)return;f.addEventListener("submit",function(){var b=f.querySelector("button[type=submit]");if(!b)return;b.disabled=true;b.textContent="Verifying..."})})
</script>
</head>
<body>
<main class="gate-card">
<div class="gate-badge">`+iconLockLarge+`</div>
<div class="eyebrow">Protected site</div>
<h1><span class="site">%s</span></h1>
<p class="hint">This site is password-protected. Enter the access password to continue.</p>
%s
<form method="POST" action="%s" autocomplete="off">
<div class="field"><input type="password" name="password" placeholder="Access password" aria-label="Access password" autofocus required autocomplete="current-password"><button type="button" class="toggle" onclick="togglePwd(this)" aria-label="Show password" aria-pressed="false">`+pwdIconHidden+`</button></div>
<button type="submit" class="submit">Enter site</button>
</form>
<div class="foot">Secured by <a href="`+projectRepoURL+`" target="_blank" rel="noopener">Vibecast</a></div>
</main>
</body>
</html>`, logoFavicon, escHTML(siteName), css, pwdIconShown, pwdIconHidden, escHTML(siteName), errBlock, escHTML(action))
}

// visitorErrorPageHTML renders a branded error page for people visiting a
// deployed site. It is intentionally dependency-free: tokens only, no page JS,
// so it renders even when the dashboard shell is unavailable.
func visitorErrorPageHTML(status int, title, detail string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<meta name="robots" content="noindex, nofollow">
<meta name="theme-color" content="#f2f5f3" media="(prefers-color-scheme: light)">
<meta name="theme-color" content="#0b1210" media="(prefers-color-scheme: dark)">
%s
<title>%s - Vibecast</title>
<style>
%s
body{display:grid;grid-template-rows:1fr auto;min-height:100dvh;padding:24px;background:var(--bg)}
.err-wrap{display:grid;place-items:center}
.err-card{width:100%%;max-width:448px;padding:30px;background:var(--surface);border:1px solid var(--line);border-radius:var(--r-lg);box-shadow:var(--shadow-md)}
.err-badge{display:inline-flex;align-items:center;justify-content:center;width:38px;height:38px;margin-bottom:16px;border-radius:var(--r-md);background:var(--surface-2);color:var(--warn)}
.err-badge svg{width:19px;height:19px}
.err-code{font-family:var(--ff-mono);font-size:.6875rem;letter-spacing:.09em;text-transform:uppercase;color:var(--dim)}
.err-card h1{margin:6px 0 8px;font-size:1.375rem;line-height:1.25;font-weight:600;letter-spacing:-.01em}
.err-card p{font-size:var(--fs-sm);line-height:1.65;color:var(--dim);margin-bottom:20px}
.err-cta{display:inline-flex;align-items:center;justify-content:center;min-height:38px;padding:0 16px;border-radius:var(--r-md);background:var(--accent);color:var(--on-accent);font-size:var(--fs-sm);font-weight:550;text-decoration:none}
.err-foot{max-width:448px;margin:0 auto;padding-top:14px;text-align:center;font-family:var(--ff-mono);font-size:.6875rem;color:var(--dim)}
</style>
</head>
<body>
<div class="err-wrap"><main class="err-card">
<span class="err-badge">%s</span>
<div class="err-code">Error %d</div>
<h1>%s</h1>
<p>%s</p>
<a class="err-cta" href="/">Open Vibecast</a>
</main></div>
<div class="err-foot">Served by <a href="`+projectRepoURL+`" target="_blank" rel="noopener">Vibecast</a></div>
</body>
</html>`, logoFavicon, escHTML(title), uiTokensCSS, iconAlert, status, escHTML(title), escHTML(detail))
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
