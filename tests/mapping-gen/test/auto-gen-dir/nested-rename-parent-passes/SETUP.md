# Scenario

**Feature**: a parent doctest tree contains a nested doctest boundary that is renamed after an initial failing run

```
# maps doc-style leaves to Go test packages
doctest build/test -> mapping-gen -> leaf <-> Go package

# gen-dir modes
auto gen-dir -> one package per leaf | explicit gen-dir -> user-specified layout
```

## Preconditions
- A parent doctest tree has one passing leaf and a nested doctest with a stale failing leaf.
- The nested directory is renamed and replaced with a fixed tree before re-running the parent tree.

## Steps
1. Create `parent/parent_ok` (passing) and `parent/nested/verbose_leaf` (always-failing stale leaf).
2. Run `doctest test parent/nested` to populate mapping-gen cache; expect failure.
3. Rename `parent/nested` to `parent/nested-renamed`, recreate the nested tree with a passing leaf.
4. Run `doctest test parent` (parent scope only) via `req.Args`.

```go
import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var nestedRenameState struct {
	NestedFailResp *Response
	StaleCachePath string
}

func Setup(t *testing.T, req *Request) error {
	nestedRenameState.NestedFailResp = nil
	nestedRenameState.StaleCachePath = ""

	if req.Bin == "" {
		t.Fatalf("req.Bin is not set")
	}

	proj := t.TempDir()
	if err := os.WriteFile(filepath.Join(proj, "go.mod"), []byte("module testproj\ngo 1.21\n"), 0644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	parentDir := filepath.Join(proj, "parent")
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		t.Fatalf("mkdir parent dir: %v", err)
	}
	runCode := "func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }"
	if err := createDoctestRoot(parentDir, runCode); err != nil {
		t.Fatalf("create parent doctest root: %v", err)
	}
	if err := createDoctestLeaf(filepath.Join(parentDir, "parent_ok")); err != nil {
		t.Fatalf("create parent_ok leaf: %v", err)
	}

	nestedDir := filepath.Join(parentDir, "nested")
	createNestedTree(t, nestedDir, "verbose_leaf", staleLeafAssertGo())

	nestedFailResp := doRun(t, req.Bin, []string{"test", nestedDir})
	nestedRenameState.NestedFailResp = nestedFailResp
	staleLeafDir := parseGoTestRunDir(nestedFailResp.Stderr)
	if staleLeafDir == "" {
		t.Fatalf("could not parse stale leaf gen dir from nested run stderr:\n%s", nestedFailResp.Stderr)
	}
	nestedRenameState.StaleCachePath = expandDisplayPath(staleLeafDir)

	renamedDir := filepath.Join(parentDir, "nested-renamed")
	if err := os.Rename(nestedDir, renamedDir); err != nil {
		t.Fatalf("rename nested -> nested-renamed: %v", err)
	}
	createNestedTree(t, renamedDir, "verbose_ok", passingLeafAssertGo())

	req.Args = append(req.Args, parentDir)
	return nil
}

func createNestedTree(t *testing.T, nestedDir, leafName, assertGo string) {
	t.Helper()
	if err := os.MkdirAll(nestedDir, 0755); err != nil {
		t.Fatalf("mkdir nested dir: %v", err)
	}
	runCode := "func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }"
	if err := createDoctestRoot(nestedDir, runCode); err != nil {
		t.Fatalf("create nested doctest root: %v", err)
	}
	leafDir := filepath.Join(nestedDir, leafName)
	if err := os.MkdirAll(leafDir, 0755); err != nil {
		t.Fatalf("mkdir nested leaf: %v", err)
	}
	if err := os.WriteFile(filepath.Join(leafDir, "SETUP.md"), []byte(leafSetupContent()), 0644); err != nil {
		t.Fatalf("write nested leaf SETUP.md: %v", err)
	}
	fence := bt + "go\n"
	assertDoc := "## Expected\n\n" + fence + "import (\n\t\"testing\"\n)\n\n" + assertGo + "\n" + bt + "\n"
	if err := os.WriteFile(filepath.Join(leafDir, "ASSERT.md"), []byte(assertDoc), 0644); err != nil {
		t.Fatalf("write nested leaf ASSERT.md: %v", err)
	}
}

func staleLeafAssertGo() string {
	return `func Assert(t *testing.T, req *Request, resp *Response, err error) {
	t.Fatal("stale nested verbose_leaf package")
}`
}

func passingLeafAssertGo() string {
	return `func Assert(t *testing.T, req *Request, resp *Response, err error) {}`
}

func expandDisplayPath(p string) string {
	if strings.HasPrefix(p, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return p
		}
		return filepath.Join(home, p[2:])
	}
	return p
}

func parseGoTestRunDir(stderr string) string {
	for _, line := range strings.Split(stderr, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "cd ") || !strings.Contains(line, " && go ") {
			continue
		}
		rest := strings.TrimPrefix(line, "cd ")
		idx := strings.Index(rest, " && go ")
		if idx >= 0 {
			return strings.TrimSpace(rest[:idx])
		}
	}
	return ""
}

func doRun(t *testing.T, bin string, args []string) *Response {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Env = os.Environ()

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	resp := &Response{}
	err := cmd.Run()
	resp.Stdout = stdoutBuf.String()
	resp.Stderr = stderrBuf.String()
	if err == nil {
		return resp
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		resp.ExitCode = exitErr.ExitCode()
		return resp
	}
	resp.Err = err
	return resp
}
```