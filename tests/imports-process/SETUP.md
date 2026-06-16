# Scenario

**Feature**: the doctest module root is two levels above this test tree (`DOCTEST_ROOT/../..`)

```
# Go import processing during code generation
doctest build -> parse imports -> remove unused -> report syntax errors
```

## Preconditions
- The doctest module root is two levels above this test tree (`DOCTEST_ROOT/../..`).
- Each leaf creates a temporary doctest project and runs `doctest test` on it.
- The doctest binary is built fresh for each test.

## Steps
1. Build the doctest binary from the module root.
2. Create a temp Go project with a doctest tree containing specific imports.
3. Run `doctest test <test-dir>` and capture output.

## Context
- These tests verify the `imports.Process` behavior in `WriteGeneratedCase`.
- The generated test code must compile and pass for the unused-import case.
- The syntax error case must produce a clear error message.

```go
import (
    "bytes"
    "context"
    "errors"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
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

var bt = "\x60\x60\x60"

func doctestGoBlock(code string) string {
    return "\n## Test\n\n" + bt + "go\n" + code + "\n" + bt + "\n"
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

func createDoctestRoot(dir, setupContent string) error {
    if err := os.WriteFile(filepath.Join(dir, "DOCTEST.md"), []byte("# Test\n"), 0644); err != nil {
        return err
    }
    if err := os.WriteFile(filepath.Join(dir, "SETUP.md"), []byte(setupContent), 0644); err != nil {
        return err
    }
    return nil
}

func createDoctestLeaf(dir, setupContent string) error {
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }
    if err := os.WriteFile(filepath.Join(dir, "SETUP.md"), []byte(setupContent), 0644); err != nil {
        return err
    }
    if err := os.WriteFile(filepath.Join(dir, "ASSERT.md"), []byte(leafAssertContent()), 0644); err != nil {
        return err
    }
    return nil
}

```
