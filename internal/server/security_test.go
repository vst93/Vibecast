package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"vibecast/internal/auth"
	"vibecast/internal/db"
)

func newTestServer(t *testing.T) *Server {
	t.Helper()
	root := t.TempDir()
	database, err := db.Open(filepath.Join(root, "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = database.Close() })
	storageDir := filepath.Join(root, "sites")
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		t.Fatal(err)
	}
	return &Server{config: &Config{StorageDir: storageDir}, database: database}
}

func createTestUserAndSite(t *testing.T, s *Server, protected bool) (*db.User, *db.Site) {
	t.Helper()
	passwordHash, err := auth.HashPassword("old-password")
	if err != nil {
		t.Fatal(err)
	}
	user, err := db.CreateUser(s.database, "owner@example.com", passwordHash, false)
	if err != nil {
		t.Fatal(err)
	}
	siteHash := ""
	if protected {
		siteHash = "protected"
	}
	site, err := db.CreateSite(s.database, user.ID, "owner-site", "Owner Site", "", siteHash, "", false, false)
	if err != nil {
		t.Fatal(err)
	}
	siteDir := filepath.Join(s.config.StorageDir, site.Slug)
	if err := os.MkdirAll(siteDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(siteDir, "index.html"), []byte("<!doctype html><title>site</title>"), 0644); err != nil {
		t.Fatal(err)
	}
	return user, site
}

func TestStaticHandlerAllowsLoggedInOwnerWithoutSitePassword(t *testing.T) {
	s := newTestServer(t)
	user, site := createTestUserAndSite(t, s, true)
	if err := db.CreateSession(s.database, user.ID, "owner-session", time.Now().Add(time.Hour)); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/s/"+site.Slug+"/", nil)
	req.AddCookie(&http.Cookie{Name: "vibecast_session", Value: "owner-session"})
	recorder := httptest.NewRecorder()
	s.staticHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected owner access, got %d with Location %q", recorder.Code, recorder.Header().Get("Location"))
	}
}

func TestUploadedHTMLIsSandboxedOnManagementOrigin(t *testing.T) {
	s := newTestServer(t)
	_, site := createTestUserAndSite(t, s, false)
	req := httptest.NewRequest(http.MethodGet, "/s/"+site.Slug+"/", nil)
	recorder := httptest.NewRecorder()
	s.staticHandler(recorder, req)

	if got := recorder.Header().Get("Content-Security-Policy"); !strings.HasPrefix(got, "sandbox") {
		t.Fatalf("expected sandbox CSP, got %q", got)
	}
}

func TestSameOriginMiddlewareRejectsOpaqueOrigin(t *testing.T) {
	s := &Server{}
	called := false
	handler := s.sameOriginMiddleware(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { called = true }))
	req := httptest.NewRequest(http.MethodPost, "http://example.com/api/sites", nil)
	req.Header.Set("Origin", "null")
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusForbidden || called {
		t.Fatalf("expected request to be blocked, status=%d called=%v", recorder.Code, called)
	}
}

func TestChangePasswordRevokesOldSessionsAndRotatesCurrentSession(t *testing.T) {
	s := newTestServer(t)
	user, _ := createTestUserAndSite(t, s, false)
	if err := db.CreateSession(s.database, user.ID, "current-session", time.Now().Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	if err := db.CreateSession(s.database, user.ID, "stolen-session", time.Now().Add(time.Hour)); err != nil {
		t.Fatal(err)
	}

	body, _ := json.Marshal(map[string]string{"oldPassword": "old-password", "newPassword": "new-password"})
	req := httptest.NewRequest(http.MethodPut, "/api/auth/change-password", strings.NewReader(string(body)))
	recorder := httptest.NewRecorder()
	s.handleChangePassword(recorder, req, user)

	if recorder.Code != http.StatusOK {
		t.Fatalf("password change failed: %d %s", recorder.Code, recorder.Body.String())
	}
	for _, token := range []string{"current-session", "stolen-session"} {
		if session, err := db.GetSession(s.database, token); err != nil || session != nil {
			t.Fatalf("old session %q remains valid: session=%v err=%v", token, session, err)
		}
	}
	response := recorder.Result()
	var rotated *http.Cookie
	for _, cookie := range response.Cookies() {
		if cookie.Name == "vibecast_session" {
			rotated = cookie
		}
	}
	if rotated == nil || rotated.Value == "" {
		t.Fatal("expected a rotated session cookie")
	}
	if session, err := db.GetSession(s.database, rotated.Value); err != nil || session == nil || session.ID != user.ID {
		t.Fatalf("rotated session is invalid: session=%v err=%v", session, err)
	}
}

func TestSitePasswordCanBeRemovedAndOldAccessSessionsAreRevoked(t *testing.T) {
	s := newTestServer(t)
	user, site := createTestUserAndSite(t, s, true)
	if err := db.CreateSiteSession(s.database, site.ID, "old-site-session", time.Now().Add(time.Hour)); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPut, "/api/sites/1", strings.NewReader(`{"password":""}`))
	recorder := httptest.NewRecorder()
	s.updateSite(recorder, req, user, site)

	if recorder.Code != http.StatusOK {
		t.Fatalf("site update failed: %d %s", recorder.Code, recorder.Body.String())
	}
	updated, err := db.GetSiteByID(s.database, site.ID)
	if err != nil || updated == nil || updated.Password != "" {
		t.Fatalf("site password was not removed: site=%v err=%v", updated, err)
	}
	if session, err := db.GetSiteSession(s.database, "old-site-session"); err != nil || session != nil {
		t.Fatalf("old site session remains valid: session=%v err=%v", session, err)
	}
}

func TestSitePasswordCookieIsScopedPerSite(t *testing.T) {
	s := newTestServer(t)
	_, site := createTestUserAndSite(t, s, false)
	hash, err := auth.HashPassword("site-password")
	if err != nil {
		t.Fatal(err)
	}
	if err := db.UpdateSite(s.database, site.ID, site.Name, site.Description, hash, "site-password", false, false); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/p/"+site.Slug, strings.NewReader("password=site-password"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	recorder := httptest.NewRecorder()
	s.passwordPageHandler(recorder, req)

	if recorder.Code != http.StatusSeeOther {
		t.Fatalf("password login failed: %d %s", recorder.Code, recorder.Body.String())
	}
	var scoped bool
	for _, cookie := range recorder.Result().Cookies() {
		if cookie.Name == "site_token" && cookie.Value != "" && cookie.Path == "/s/"+site.Slug+"/" {
			scoped = true
		}
	}
	if !scoped {
		t.Fatal("site access cookie was not scoped to its site path")
	}
}

func TestAdminSettingsPartialUpdatePreservesOmittedValues(t *testing.T) {
	s := newTestServer(t)
	if err := db.SetSetting(s.database, "max_sites_per_user", "7"); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPut, "/api/admin/settings", strings.NewReader(`{"maxUploadSize":100}`))
	recorder := httptest.NewRecorder()
	s.adminHandleSettings(recorder, req, &db.User{IsAdmin: true})

	if recorder.Code != http.StatusOK {
		t.Fatalf("settings update failed: %d %s", recorder.Code, recorder.Body.String())
	}
	if got := db.GetSettingInt(s.database, "max_sites_per_user", -1); got != 7 {
		t.Fatalf("omitted max_sites_per_user was overwritten: got %d", got)
	}
	if got := db.GetSettingInt(s.database, "max_upload_size", -1); got != 100 {
		t.Fatalf("max_upload_size was not updated: got %d", got)
	}
}

// TestPasswordGateDeepLinkRedirectsWithSubPath guards the reported bug where
// opening a nested page of a protected site landed on a 404 instead of the
// password form: the old handler replied with the relative "../../p/{slug}",
// which the browser resolved against the requested file.
func TestPasswordGateDeepLinkRedirectsWithSubPath(t *testing.T) {
	s := newTestServer(t)
	_, site := createTestUserAndSite(t, s, true)
	hash, err := auth.HashPassword("site-password")
	if err != nil {
		t.Fatal(err)
	}
	if err := db.UpdateSite(s.database, site.ID, site.Name, site.Description, hash, "site-password", false, false); err != nil {
		t.Fatal(err)
	}
	sub := filepath.Join(s.config.StorageDir, site.Slug, "docs")
	if err := os.MkdirAll(sub, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sub, "page.html"), []byte("<!doctype html><title>doc</title>"), 0644); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/s/"+site.Slug+"/docs/page.html", nil)
	recorder := httptest.NewRecorder()
	s.staticHandler(recorder, req)

	if recorder.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect to the gate, got %d", recorder.Code)
	}
	location := recorder.Header().Get("Location")
	want := "/p/" + site.Slug + "?next=docs%2Fpage.html"
	if location != want {
		t.Fatalf("deep link redirect must keep the requested page: got %q want %q", location, want)
	}
}

// TestPasswordGateUnlockResumesRequestedPage verifies the return target carried
// through ?next= is applied after a successful unlock.
func TestPasswordGateUnlockResumesRequestedPage(t *testing.T) {
	s := newTestServer(t)
	_, site := createTestUserAndSite(t, s, true)
	hash, err := auth.HashPassword("site-password")
	if err != nil {
		t.Fatal(err)
	}
	if err := db.UpdateSite(s.database, site.ID, site.Name, site.Description, hash, "site-password", false, false); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/p/"+site.Slug+"?next=docs%2Fpage.html", strings.NewReader("password=site-password"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	recorder := httptest.NewRecorder()
	s.passwordPageHandler(recorder, req)

	if recorder.Code != http.StatusSeeOther {
		t.Fatalf("password login failed: %d %s", recorder.Code, recorder.Body.String())
	}
	if got, want := recorder.Header().Get("Location"), "/s/"+site.Slug+"/docs/page.html"; got != want {
		t.Fatalf("unlock did not resume the requested page: got %q want %q", got, want)
	}
}

// TestPasswordGateRejectsForeignReturnTarget makes sure ?next= cannot be used as
// an open redirect.
func TestPasswordGateRejectsForeignReturnTarget(t *testing.T) {
	s := newTestServer(t)
	_, site := createTestUserAndSite(t, s, true)
	hash, err := auth.HashPassword("site-password")
	if err != nil {
		t.Fatal(err)
	}
	if err := db.UpdateSite(s.database, site.ID, site.Name, site.Description, hash, "site-password", false, false); err != nil {
		t.Fatal(err)
	}

	for _, next := range []string{"//evil.example/x", "https://evil.example/", "../../../etc/passwd"} {
		body := "password=site-password"
		req := httptest.NewRequest(http.MethodPost, "/p/"+site.Slug+"?next="+url.QueryEscape(next), strings.NewReader(body))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		recorder := httptest.NewRecorder()
		s.passwordPageHandler(recorder, req)

		if recorder.Code != http.StatusSeeOther {
			t.Fatalf("next=%q: unexpected status %d", next, recorder.Code)
		}
		location := recorder.Header().Get("Location")
		if want := "/s/" + site.Slug + "/"; location != want {
			t.Fatalf("next=%q must be dropped: got %q want %q", next, location, want)
		}
	}
}

// TestPasswordGateScopesCookieToProxyPrefix covers sub-path deployments: the
// site cookie has to include the mount prefix or the browser never sends it
// back and the password submission looks like it did nothing.
func TestPasswordGateScopesCookieToProxyPrefix(t *testing.T) {
	s := newTestServer(t)
	_, site := createTestUserAndSite(t, s, true)
	hash, err := auth.HashPassword("site-password")
	if err != nil {
		t.Fatal(err)
	}
	if err := db.UpdateSite(s.database, site.ID, site.Name, site.Description, hash, "site-password", false, false); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/p/"+site.Slug, strings.NewReader("password=site-password"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Forwarded-Prefix", "/vibecast")
	recorder := httptest.NewRecorder()
	s.passwordPageHandler(recorder, req)

	if recorder.Code != http.StatusSeeOther {
		t.Fatalf("password login failed: %d %s", recorder.Code, recorder.Body.String())
	}
	if got, want := recorder.Header().Get("Location"), "/vibecast/s/"+site.Slug+"/"; got != want {
		t.Fatalf("redirect ignored the mount prefix: got %q want %q", got, want)
	}
	var cookiePath string
	for _, cookie := range recorder.Result().Cookies() {
		if cookie.Name == "site_token" && cookie.Value != "" {
			cookiePath = cookie.Path
		}
	}
	if want := "/vibecast/s/" + site.Slug + "/"; cookiePath != want {
		t.Fatalf("cookie scope ignored the mount prefix: got %q want %q", cookiePath, want)
	}
}

// TestPasswordGateReportsFormErrors verifies a rejected password re-renders the
// form with the error inline instead of failing silently.
func TestPasswordGateReportsFormErrors(t *testing.T) {
	s := newTestServer(t)
	_, site := createTestUserAndSite(t, s, true)
	hash, err := auth.HashPassword("site-password")
	if err != nil {
		t.Fatal(err)
	}
	if err := db.UpdateSite(s.database, site.ID, site.Name, site.Description, hash, "site-password", false, false); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/p/"+site.Slug, strings.NewReader("password=nope"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	recorder := httptest.NewRecorder()
	s.passwordPageHandler(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for a wrong password, got %d", recorder.Code)
	}
	body := recorder.Body.String()
	if !strings.Contains(body, "Incorrect password") {
		t.Fatal("wrong password did not render an inline error")
	}
	if !strings.Contains(body, `name="password"`) {
		t.Fatal("wrong password did not re-render the form")
	}
}
