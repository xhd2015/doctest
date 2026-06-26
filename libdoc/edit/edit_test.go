package edit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEditAddLabelAndExplanation(t *testing.T) {
	dir := t.TempDir()
	assertPath := filepath.Join(dir, "ASSERT.md")
	if err := os.WriteFile(assertPath, []byte("## Expected\n\n```go\nfunc Assert() {}\n```\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := Edit(dir, "ui-automation", "window test"); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(assertPath)
	if err != nil {
		t.Fatal(err)
	}
	got := string(data)
	wantPrefix := "---\nlabel: ui-automation\nexplanation: window test\n---\n\n## Expected\n"
	if got != wantPrefix+"\n```go\nfunc Assert() {}\n```\n" && !strings.HasPrefix(got, wantPrefix) {
		t.Fatalf("got:\n%s", got)
	}
}

func TestEditAppendExplanation(t *testing.T) {
	dir := t.TempDir()
	assertPath := filepath.Join(dir, "ASSERT.md")
	content := "---\nlabel: ui\nexplanation: first\n---\n\n## Expected\n"
	if err := os.WriteFile(assertPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	if err := Edit(dir, "", "second"); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(assertPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "explanation: first; second") {
		t.Fatalf("got:\n%s", data)
	}
}

func TestEditRejectsDotDotDot(t *testing.T) {
	if err := Edit("./tests/...", "ui", ""); err == nil {
		t.Fatal("expected error")
	}
}