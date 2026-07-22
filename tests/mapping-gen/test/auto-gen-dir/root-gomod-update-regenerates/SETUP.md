# Scenario

**Feature**: mapping-gen cache regenerates go.mod when the project root go.mod changes

```
# maps doc-style leaves to Go test packages
doctest build/test -> mapping-gen -> leaf <-> Go package

# gen-dir modes
auto gen-dir -> one package per leaf | explicit gen-dir -> user-specified layout
```

## Preconditions
- A temp project has a local replace directive in go.mod that is removed after the first doctest run warms the mapping-gen cache.
- The replaced local dependency directory is deleted so a stale cached replace would break go test.

## Steps
1. Create a project with `go.mod` containing `require localdep` and `replace localdep => ./dep`.
2. Run `doctest test` once to populate the mapping-gen cache.
3. Update root `go.mod` to remove the replace and require.
4. Delete the `dep/` directory.
5. Run `doctest test` again via the outer `Run` call.

```go
import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

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

func Setup(t *testing.T, req *Request) error {
	// Request-local multi-phase state (no package vars).
	req.MRFirst = nil
	req.MRGenDir = ""

	if req.Bin == "" {
		t.Fatalf("req.Bin is not set")
	}

	proj := t.TempDir()

	depDir := filepath.Join(proj, "dep")
	if err := os.MkdirAll(depDir, 0755); err != nil {
		t.Fatalf("mkdir dep: %v", err)
	}
	if err := os.WriteFile(filepath.Join(depDir, "go.mod"), []byte("module localdep\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatalf("write dep go.mod: %v", err)
	}

	goModV1 := "module testproj\n\ngo 1.21\n\nrequire localdep v0.0.0\n\nreplace localdep => ./dep\n"
	if err := os.WriteFile(filepath.Join(proj, "go.mod"), []byte(goModV1), 0644); err != nil {
		t.Fatalf("write go.mod v1: %v", err)
	}

	testDir := filepath.Join(proj, "tests")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("mkdir test dir: %v", err)
	}
	runCode := "func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }"
	if err := createDoctestRoot(testDir, runCode); err != nil {
		t.Fatalf("create doctest root: %v", err)
	}
	if err := createDoctestLeaf(filepath.Join(testDir, "simple")); err != nil {
		t.Fatalf("create leaf: %v", err)
	}

	warmArgs := []string{"test", "-v", testDir}
	warmResp := doRun(t, req.Bin, warmArgs)
	genRoot := parseGenDir(warmResp.Stderr)
	if genRoot == "" {
		t.Fatalf("could not parse gen root from warm run stderr:\n%s", warmResp.Stderr)
	}
	if warmResp.ExitCode != 0 {
		t.Fatalf("warm run expected exit 0, got %d\nstderr:\n%s", warmResp.ExitCode, warmResp.Stderr)
	}

	goModV2 := "module testproj\n\ngo 1.21\n"
	if err := os.WriteFile(filepath.Join(proj, "go.mod"), []byte(goModV2), 0644); err != nil {
		t.Fatalf("write go.mod v2: %v", err)
	}
	if err := os.RemoveAll(depDir); err != nil {
		t.Fatalf("remove dep dir: %v", err)
	}

	req.MRFirst = warmResp
	req.MRGenDir = genRoot
	req.Args = append(req.Args, "-count=1", "-v", testDir)
	return nil
}
```
