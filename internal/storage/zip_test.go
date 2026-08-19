package storage

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveSingleFileCleansTemporaryDirectoryOnFailure(t *testing.T) {
	parent := t.TempDir()
	dest := filepath.Join(parent, "site")
	if err := os.MkdirAll(dest, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dest, "keep.txt"), []byte("keep"), 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := SaveSingleFile(strings.NewReader("too large"), "index.html", dest, 4); err == nil {
		t.Fatal("expected size limit error")
	}

	matches, err := filepath.Glob(filepath.Join(dest, ".upload-*"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 0 {
		t.Fatalf("temporary directories were not cleaned: %v", matches)
	}
	if got, err := os.ReadFile(filepath.Join(dest, "keep.txt")); err != nil || string(got) != "keep" {
		t.Fatalf("existing content changed after failed upload: %q, %v", got, err)
	}
}

func TestSaveSingleFileOverwritesOnlyMatchingFile(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "site")
	if err := os.MkdirAll(dest, 0755); err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, dest, "index.html", "old")
	writeTestFile(t, dest, "keep.css", "keep")

	if _, err := SaveSingleFile(strings.NewReader("new index"), "index.html", dest, 100); err != nil {
		t.Fatal(err)
	}
	assertTestFile(t, dest, "index.html", "new index")
	assertTestFile(t, dest, "keep.css", "keep")
}

func TestSaveSingleFileRejectsAggregateSiteSize(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "site")
	if err := os.MkdirAll(dest, 0755); err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, dest, "keep.txt", "1234")

	_, err := SaveSingleFile(strings.NewReader("567"), "new.txt", dest, 6)
	if !errors.Is(err, ErrSiteSizeLimit) {
		t.Fatalf("expected site size error, got %v", err)
	}
	assertTestFile(t, dest, "keep.txt", "1234")
	if _, err := os.Stat(filepath.Join(dest, "new.txt")); !os.IsNotExist(err) {
		t.Fatalf("over-limit file was committed: %v", err)
	}
}

func TestExtractZipCleansTemporaryDirectoryOnInvalidArchive(t *testing.T) {
	parent := t.TempDir()
	dest := filepath.Join(parent, "site")
	data := strings.NewReader("not a zip")

	if _, err := ExtractZip(data, int64(data.Len()), dest, 1024); err == nil {
		t.Fatal("expected invalid zip error")
	}

	matches, err := filepath.Glob(dest + ".tmp-*")
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 0 {
		t.Fatalf("temporary directories were not cleaned: %v", matches)
	}
}

func TestExtractZipMergesAndOverwritesMatchingPaths(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "site")
	if err := os.MkdirAll(filepath.Join(dest, "assets"), 0755); err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, dest, "index.html", "old")
	writeTestFile(t, dest, "keep.txt", "keep")
	writeTestFile(t, dest, "assets/old.css", "old css")
	data := testZip(t, map[string]string{
		"site/index.html":     "new",
		"site/assets/new.css": "new css",
	})

	result, err := ExtractZip(bytes.NewReader(data), int64(len(data)), dest, 100)
	if err != nil {
		t.Fatal(err)
	}
	if result.TotalFiles != 2 {
		t.Fatalf("unexpected extracted count: %d", result.TotalFiles)
	}
	assertTestFile(t, dest, "index.html", "new")
	assertTestFile(t, dest, "keep.txt", "keep")
	assertTestFile(t, dest, "assets/old.css", "old css")
	assertTestFile(t, dest, "assets/new.css", "new css")
}

func TestExtractZipRejectsAggregateSizeWithoutChangingSite(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "site")
	if err := os.MkdirAll(dest, 0755); err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, dest, "keep.txt", "1234")
	data := testZip(t, map[string]string{"new.txt": "567"})

	_, err := ExtractZip(bytes.NewReader(data), int64(len(data)), dest, 6)
	if !errors.Is(err, ErrSiteSizeLimit) {
		t.Fatalf("expected site size error, got %v", err)
	}
	assertTestFile(t, dest, "keep.txt", "1234")
	if _, err := os.Stat(filepath.Join(dest, "new.txt")); !os.IsNotExist(err) {
		t.Fatalf("over-limit zip changed site: %v", err)
	}
}

func TestZipDirectoryRejectsEmptyAndPackagesAllFiles(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "site")
	if err := os.MkdirAll(dest, 0755); err != nil {
		t.Fatal(err)
	}
	if _, _, err := ZipDirectory(io.Discard, dest); !errors.Is(err, ErrSiteEmpty) {
		t.Fatalf("expected empty site error, got %v", err)
	}
	writeTestFile(t, dest, "index.html", "index")
	writeTestFile(t, dest, "assets/app.css", "css")

	var archive bytes.Buffer
	count, total, err := ZipDirectory(&archive, dest)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 || total != 8 {
		t.Fatalf("unexpected archive stats: count=%d total=%d", count, total)
	}
	zr, err := zip.NewReader(bytes.NewReader(archive.Bytes()), int64(archive.Len()))
	if err != nil {
		t.Fatal(err)
	}
	got := make(map[string]string)
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
		got[file.Name] = string(contents)
	}
	if got["index.html"] != "index" || got["assets/app.css"] != "css" {
		t.Fatalf("unexpected archive contents: %#v", got)
	}
}

func testZip(t *testing.T, files map[string]string) []byte {
	t.Helper()
	var data bytes.Buffer
	zw := zip.NewWriter(&data)
	for name, contents := range files {
		writer, err := zw.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := io.WriteString(writer, contents); err != nil {
			t.Fatal(err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return data.Bytes()
}

func writeTestFile(t *testing.T, root, name, contents string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0644); err != nil {
		t.Fatal(err)
	}
}

func assertTestFile(t *testing.T, root, name, want string) {
	t.Helper()
	got, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(name)))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Fatalf("%s = %q, want %q", name, got, want)
	}
}
