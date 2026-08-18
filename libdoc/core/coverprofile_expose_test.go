package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCoverFileKeyFromExposeAbs(t *testing.T) {
	product := t.TempDir()
	if err := os.WriteFile(filepath.Join(product, "go.mod"), []byte("module example.com/app\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}
	abs := filepath.Join(product, DoctestInternalExposeDir, "greet", "expose.go")
	if err := os.MkdirAll(filepath.Dir(abs), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(abs, []byte("package greet\n"), 0644); err != nil {
		t.Fatal(err)
	}
	key, err := coverFileKeyFromExposeAbs(abs)
	if err != nil {
		t.Fatal(err)
	}
	want := "example.com/app/" + DoctestInternalExposeDir + "/greet/expose.go"
	if key != want {
		t.Fatalf("key = %q, want %q", key, want)
	}
}

func TestStripExposeFromCoverProfile_onlyListedKeys(t *testing.T) {
	product := t.TempDir()
	if err := os.WriteFile(filepath.Join(product, "go.mod"), []byte("module example.com/app\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}
	genRoot := t.TempDir()
	abs := filepath.Join(product, DoctestInternalExposeDir, "greet", "expose.go")
	if err := os.MkdirAll(filepath.Dir(abs), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(abs, []byte("package greet\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := recordExposeMaterialized(genRoot, abs); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = CleanupExposeMaterialized(genRoot) })

	profile := filepath.Join(t.TempDir(), "cover.out")
	raw := strings.Join([]string{
		"mode: set",
		"example.com/app/" + DoctestInternalExposeDir + "/greet/expose.go:6.21,8.2 1 1",
		"example.com/app/internal/greet/greet.go:3.21,3.36 1 1",
		// Unlisted fake path — must stay (no blind reserved-name strip).
		"other.com/mod/" + DoctestInternalExposeDir + "/x/expose.go:1.1,2.2 1 0",
		"",
	}, "\n")
	if err := os.WriteFile(profile, []byte(raw), 0644); err != nil {
		t.Fatal(err)
	}

	removed, err := StripExposeFromCoverProfile(profile, []string{genRoot})
	if err != nil {
		t.Fatal(err)
	}
	if removed != 1 {
		t.Fatalf("removed = %d, want 1", removed)
	}
	got, err := os.ReadFile(profile)
	if err != nil {
		t.Fatal(err)
	}
	s := string(got)
	if strings.Contains(s, "example.com/app/"+DoctestInternalExposeDir+"/greet/expose.go") {
		t.Fatalf("listed expose line still present:\n%s", s)
	}
	if !strings.Contains(s, "example.com/app/internal/greet/greet.go") {
		t.Fatalf("product internal line missing:\n%s", s)
	}
	if !strings.Contains(s, "other.com/mod/"+DoctestInternalExposeDir+"/x/expose.go") {
		t.Fatalf("unlisted path must not be stripped blindly:\n%s", s)
	}
	if !strings.HasPrefix(strings.TrimSpace(s), "mode: set") {
		t.Fatalf("mode not preserved:\n%s", s)
	}
}

func TestStripExposeFromCoverProfile_noopEmptyList(t *testing.T) {
	profile := filepath.Join(t.TempDir(), "cover.out")
	raw := "mode: set\nexample.com/app/internal/greet/greet.go:3.21,3.36 1 1\n"
	if err := os.WriteFile(profile, []byte(raw), 0644); err != nil {
		t.Fatal(err)
	}
	removed, err := StripExposeFromCoverProfile(profile, []string{t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	if removed != 0 {
		t.Fatalf("removed = %d, want 0", removed)
	}
	got, _ := os.ReadFile(profile)
	if string(got) != raw {
		t.Fatalf("profile mutated without list:\n%s", got)
	}
}
