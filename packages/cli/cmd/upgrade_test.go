package cmd

import (
	"archive/tar"
	"compress/gzip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSafeExtractPath(t *testing.T) {
	dir := t.TempDir()

	cases := []struct {
		name  string
		entry string
		ok    bool
	}{
		{"flat file", "autodev", true},
		{"nested", "autodev_linux_amd64/autodev", true},
		{"deep nested", "a/b/c/autodev", true},
		{"parent traversal", "../autodev", false},
		{"double parent", "../../etc/passwd", false},
		{"wrapped traversal", "a/../../autodev", false},
		{"absolute path", "/etc/passwd", false},
		{"stays in dir after clean", "./autodev", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			target, ok := safeExtractPath(dir, tc.entry)
			if ok != tc.ok {
				t.Fatalf("ok = %v, want %v", ok, tc.ok)
			}
			if ok {
				rel, err := filepath.Rel(dir, target)
				if err != nil || strings.HasPrefix(rel, "..") {
					t.Fatalf("target %q escapes dir", target)
				}
				if !filepath.IsAbs(target) {
					t.Fatalf("target %q not absolute", target)
				}
			}
		})
	}
}

func TestExtractTarGzRejectsTraversal(t *testing.T) {
	dir := t.TempDir()
	archive := filepath.Join(dir, "evil.tar.gz")
	f, err := os.Create(archive)
	if err != nil {
		t.Fatal(err)
	}

	gz := gzip.NewWriter(f)
	tw := tar.NewWriter(gz)
	payload := []byte("root:x:0:0:root:/root:/bin/sh")
	hdr := &tar.Header{
		Typeflag: tar.TypeReg,
		Name:     "../../../../etc/passwd",
		Size:     int64(len(payload)),
		Mode:     0644,
	}
	if err := tw.WriteHeader(hdr); err != nil {
		t.Fatal(err)
	}
	if _, err := tw.Write(payload); err != nil {
		t.Fatal(err)
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	out := filepath.Join(dir, "out")
	if err := os.MkdirAll(out, 0755); err != nil {
		t.Fatal(err)
	}
	if err := extractTarGz(archive, out); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(out, "etc", "passwd")); err == nil {
		t.Fatalf("traversal entry should have been skipped, but was written under %s", out)
	}
}
