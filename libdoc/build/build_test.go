package build

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
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

func doctestDoc(code string) string {
	return "# Tests\n\n## Version\n0.0.2\n\n## DSN (Domain Specific Notion)\n\n### Participants\n- **system** — under test.\n\n### Behaviors\n- **run** — executes.\n\n```go\n" + strings.TrimSpace(code) + "\n```\n"
}

func writeRootHarness(t *testing.T, root, doctestCode, rootSetupCode string) {
	t.Helper()
	writeTreeFile(t, root, "DOCTEST.md", doctestDoc(doctestCode))
	if strings.TrimSpace(rootSetupCode) != "" {
		writeTreeFile(t, root, "SETUP.md", setupDoc(rootSetupCode))
	}
}

func TestBuildBasicTree(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }
`, "")
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
`))

	genDir := filepath.Join(t.TempDir(), "gen")
	if err := Build(root, core.Options{GenDir: genDir}); err != nil {
		t.Fatalf("build: %v", err)
	}
}

func TestBuildWithGenDirCreatesGoFile(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }
`, "")
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
`))

	genDir := filepath.Join(t.TempDir(), "gen")
	if err := Build(root, core.Options{GenDir: genDir}); err != nil {
		t.Fatalf("build: %v", err)
	}

	entries, err := os.ReadDir(genDir)
	if err != nil {
		t.Fatalf("read gen dir: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("expected generated files")
	}
}

func TestBuildVerboseOutput(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }
`, "")
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
`))

	var stderr bytes.Buffer
	genDir := filepath.Join(t.TempDir(), "gen")
	if err := Build(root, core.Options{GenDir: genDir, Verbose: true, Stderr: &stderr}); err != nil {
		t.Fatalf("build verbose: %v", err)
	}
	out := stderr.String()
	if !strings.Contains(out, "build") {
		t.Fatalf("expected build output, got %q", out)
	}
}

func TestBuildFailInvalidTree(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeTreeFile(t, root, "SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { return nil }
`))
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
`))

	genDir := filepath.Join(t.TempDir(), "gen")
	err := Build(root, core.Options{GenDir: genDir})
	if err == nil {
		t.Fatal("expected build failure for tree missing Run")
	}
}
