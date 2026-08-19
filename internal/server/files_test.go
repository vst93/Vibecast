package server

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"vibecast/internal/db"
	"vibecast/internal/storage"
)

func TestSiteFileTreeListsNestedFilesAndTotals(t *testing.T) {
	s := newTestServer(t)
	_, site := createTestUserAndSite(t, s, false)
	writeServerTestFile(t, filepath.Join(s.config.StorageDir, site.Slug), "assets/app.css", "css")

	recorder := httptest.NewRecorder()
	s.siteFileTree(recorder, httptest.NewRequest(http.MethodGet, "/api/sites/1/files", nil), site)
	if recorder.Code != http.StatusOK {
		t.Fatalf("file tree failed: %d %s", recorder.Code, recorder.Body.String())
	}
	var response struct {
		Data struct {
			Items     []storage.FileEntry `json:"items"`
			FileCount int                 `json:"fileCount"`
			TotalSize int64               `json:"totalSize"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.Data.FileCount != 2 || response.Data.TotalSize != 37 {
		t.Fatalf("unexpected file stats: %#v", response.Data)
	}
	var foundNested bool
	for _, item := range response.Data.Items {
		if item.Name == "assets/app.css" && !item.Dir {
			foundNested = true
		}
	}
	if !foundNested {
		t.Fatalf("nested file missing: %#v", response.Data.Items)
	}
}

func TestClearSiteFilesRemovesAllContent(t *testing.T) {
	s := newTestServer(t)
	_, site := createTestUserAndSite(t, s, false)
	siteDir := filepath.Join(s.config.StorageDir, site.Slug)
	writeServerTestFile(t, siteDir, "assets/app.css", "css")

	recorder := httptest.NewRecorder()
	s.clearSiteFiles(recorder, httptest.NewRequest(http.MethodDelete, "/api/sites/1/files", nil), site)
	if recorder.Code != http.StatusOK {
		t.Fatalf("clear failed: %d %s", recorder.Code, recorder.Body.String())
	}
	entries, _, err := storage.ListFiles(siteDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("site still contains files: %#v", entries)
	}
}

func TestDownloadSiteFilesRejectsEmptyDirectory(t *testing.T) {
	s := newTestServer(t)
	_, site := createTestUserAndSite(t, s, false)
	if err := storage.ClearDirectory(filepath.Join(s.config.StorageDir, site.Slug)); err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	s.downloadSiteFiles(recorder, httptest.NewRequest(http.MethodGet, "/api/sites/1/download", nil), site)
	if recorder.Code != http.StatusConflict {
		t.Fatalf("expected empty download rejection, got %d %s", recorder.Code, recorder.Body.String())
	}
}

func TestDownloadSiteFilesReturnsZipWithContentLength(t *testing.T) {
	s := newTestServer(t)
	_, site := createTestUserAndSite(t, s, false)
	siteDir := filepath.Join(s.config.StorageDir, site.Slug)
	writeServerTestFile(t, siteDir, "assets/app.css", "css")

	recorder := httptest.NewRecorder()
	s.downloadSiteFiles(recorder, httptest.NewRequest(http.MethodGet, "/api/sites/1/download", nil), site)
	if recorder.Code != http.StatusOK {
		t.Fatalf("download failed: %d %s", recorder.Code, recorder.Body.String())
	}
	if recorder.Header().Get("Content-Type") != "application/zip" || recorder.Header().Get("Content-Length") == "" {
		t.Fatalf("missing download headers: %#v", recorder.Header())
	}
	zr, err := zip.NewReader(bytes.NewReader(recorder.Body.Bytes()), int64(recorder.Body.Len()))
	if err != nil {
		t.Fatal(err)
	}
	files := make(map[string]string)
	for _, file := range zr.File {
		if file.FileInfo().IsDir() {
			continue
		}
		rc, err := file.Open()
		if err != nil {
			t.Fatal(err)
		}
		contents, err := io.ReadAll(rc)
		_ = rc.Close()
		if err != nil {
			t.Fatal(err)
		}
		files[file.Name] = string(contents)
	}
	if files["assets/app.css"] != "css" || !strings.Contains(files["index.html"], "<title>site</title>") {
		t.Fatalf("unexpected downloaded archive: %#v", files)
	}
}

func TestDeploySingleFilePreservesOtherFiles(t *testing.T) {
	s := newTestServer(t)
	user, site := createTestUserAndSite(t, s, false)
	siteDir := filepath.Join(s.config.StorageDir, site.Slug)
	writeServerTestFile(t, siteDir, "keep.txt", "keep")

	req := multipartUploadRequest(t, "index.html", "new index")
	recorder := httptest.NewRecorder()
	s.deploySite(recorder, req, user, site)
	if recorder.Code != http.StatusOK {
		t.Fatalf("deploy failed: %d %s", recorder.Code, recorder.Body.String())
	}
	assertServerTestFile(t, siteDir, "index.html", "new index")
	assertServerTestFile(t, siteDir, "keep.txt", "keep")
}

func TestDeployRejectsAggregateSiteSize(t *testing.T) {
	s := newTestServer(t)
	user, site := createTestUserAndSite(t, s, false)
	if err := db.SetSetting(s.database, "max_upload_size", "1"); err != nil {
		t.Fatal(err)
	}
	siteDir := filepath.Join(s.config.StorageDir, site.Slug)
	writeServerTestFile(t, siteDir, "keep.bin", strings.Repeat("a", 700*1024))

	req := multipartUploadRequest(t, "new.txt", strings.Repeat("b", 400*1024))
	recorder := httptest.NewRecorder()
	s.deploySite(recorder, req, user, site)
	if recorder.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected aggregate size rejection, got %d %s", recorder.Code, recorder.Body.String())
	}
	if _, err := os.Stat(filepath.Join(siteDir, "new.txt")); !os.IsNotExist(err) {
		t.Fatalf("over-limit upload was committed: %v", err)
	}
}

func multipartUploadRequest(t *testing.T, filename, contents string) *http.Request {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.WriteString(part, contents); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/sites/1/deploy", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func writeServerTestFile(t *testing.T, root, name, contents string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0644); err != nil {
		t.Fatal(err)
	}
}

func assertServerTestFile(t *testing.T, root, name, want string) {
	t.Helper()
	got, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(name)))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Fatalf("%s = %q, want %q", name, got, want)
	}
}
