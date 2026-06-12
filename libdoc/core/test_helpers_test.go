package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTreeFile(t *testing.T, root string, rel string, content string) {
	t.Helper()
	path := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func setupDoc(code string) string {
	code = strings.TrimSpace(code)
	if !strings.Contains(code, "func Setup") {
		setupLine := "func Setup(t *testing.T, req *Request) error { _ = req; return nil }"
		idx := strings.Index(code, "\")\n")
		if idx >= 0 && strings.Contains(code[:idx], "import") {
			code = code[:idx+3] + "\n" + setupLine + "\n" + code[idx+3:]
		} else {
			idx = strings.Index(code, "\"\n")
			if idx >= 0 && strings.Contains(code[:idx], "import") {
				code = code[:idx+2] + "\n" + setupLine + "\n" + code[idx+2:]
			} else {
				code = setupLine + "\n" + code
			}
		}
	}
	return "# Setup\n\nAny section names are allowed.\n\n```go\n" + strings.TrimSpace(code) + "\n```\n"
}

func assertDoc(code string) string {
	return "# Assert\n\nAny section names are allowed.\n\n```go\n" + strings.TrimSpace(code) + "\n```\n"
}
