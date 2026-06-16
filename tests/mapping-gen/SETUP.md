# Scenario

**Feature**: the doctest module root is two levels above this test tree (`DOCTEST_ROOT/../..`)

```
# maps doc-style leaves to Go test packages
doctest build/test -> mapping-gen -> leaf <-> Go package

# gen-dir modes
auto gen-dir -> one package per leaf | explicit gen-dir -> user-specified layout
```

## Preconditions
- The doctest module root is two levels above this test tree (`DOCTEST_ROOT/../..`).
- The mapping-gen feature mirrors source directory structure under the cache root.
- Each leaf sets up a temporary project and verifies generated output.

## Steps
1. Build the doctest binary from the module root.
2. Execute the binary given by `req.Bin`.
3. Capture stdout, stderr, exit code, and the raw execution error.

## Context
- Each test creates a temp Go project with a doctest tree.
- The generated output is inspected for per-leaf directory structure.
- Shared go.mod is at project root level in the generated tree.

```go
import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	libdocbuild "github.com/xhd2015/doctest/libdoc/build"
)

type Request struct {
	Args    []string
	Env     []string
	WorkDir string
	Timeout time.Duration
	Bin     string
}

type Response struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Err      error
}

func Setup(t *testing.T, req *Request) error {
	req.Timeout = 120 * time.Second

	tmp := t.TempDir()
	doctestBin := filepath.Join(tmp, "doctest")
	buildDir := filepath.Join(DOCTEST_ROOT, "..", "..")
	buildArgs := []string{"build", "-o", doctestBin}
	if libdocbuild.NeedsBuildVCSFlag(buildDir) {
		buildArgs = append(buildArgs, "-buildvcs=false")
	}
	buildArgs = append(buildArgs, "./cmd/doctest")
	build := exec.Command("go", buildArgs...)
	build.Dir = buildDir
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build doctest: %v\n%s", err, string(out))
	}
	req.Bin = doctestBin
	return nil
}

func Run(t *testing.T, req *Request) (*Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), req.Timeout)
	defer cancel()

	bin := req.Bin
	if bin == "" {
		return nil, fmt.Errorf("req.Bin is not set")
	}
	cmd := exec.CommandContext(ctx, bin, req.Args...)
	cmd.Dir = req.WorkDir
	cmd.Env = append(os.Environ(), req.Env...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	resp := &Response{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
		Err:    err,
	}
	if err == nil {
		return resp, nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		resp.ExitCode = exitErr.ExitCode()
		return resp, nil
	}
	if ctx.Err() != nil {
		return resp, ctx.Err()
	}
	return resp, err
}

var bt = "`" + "`" + "`"

func doctestGoBlock(code string) string {
	return "## Test\n\n" + bt + "go\n" + code + bt + "\n"
}

func rootSetupContent(extraCode string) string {
	code := "import \"testing\"\n\ntype Request struct{ Name string }\ntype Response struct{ Message string }\n\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }\n" + extraCode
	return doctestGoBlock(code)
}

func leafSetupContent() string {
	return doctestGoBlock("import \"testing\"\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }")
}

func leafAssertContent() string {
	return doctestGoBlock("import \"testing\"\nfunc Assert(t *testing.T, req *Request, resp *Response, err error) {}")
}

func createDoctestRoot(dir, runCode string) error {
	if err := os.WriteFile(filepath.Join(dir, "DOCTEST.md"), []byte("# Test\n"), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "SETUP.md"), []byte(rootSetupContent(runCode)), 0644); err != nil {
		return err
	}
	return nil
}

func createDoctestLeaf(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "SETUP.md"), []byte(leafSetupContent()), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "ASSERT.md"), []byte(leafAssertContent()), 0644); err != nil {
		return err
	}
	return nil
}

func createTempProject(t *testing.T, dirName string) string {
	t.Helper()
	tmp := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module testproj\ngo 1.21\n"), 0644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	testDir := filepath.Join(tmp, dirName)
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("mkdir test dir: %v", err)
	}
	runCode := "func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }"
	if err := createDoctestRoot(testDir, runCode); err != nil {
		t.Fatalf("create doctest root: %v", err)
	}
	return testDir
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file %s to exist: %v", path, err)
	}
}

func assertFileNotExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err == nil {
		t.Fatalf("expected file %s to not exist", path)
	}
}

func parseGenDir(stderr string) string {
	idx := strings.Index(stderr, "→ ")
	if idx < 0 {
		return ""
	}
	rest := stderr[idx+len("→ "):]
	end := strings.IndexFunc(rest, func(r rune) bool { return r == '\n' || r == ' ' })
	if end < 0 {
		return strings.TrimSpace(rest)
	}
	return strings.TrimSpace(rest[:end])
}

func findGoModDir(startDir string) string {
	dir := startDir
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

```
