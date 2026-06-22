# Scenario

**Feature**: build the doctest binary and invoke it as a subprocess to test CLI behavior

```
# build the doctest binary from module source
go build ./cmd/doctest -> doctest binary

# invoke binary as subprocess, capture everything
doctest <args> -> {stdout, stderr, exit code}
```

## Preconditions
- The doctest module root is the parent of this test tree (`DOCTEST_ROOT/..`).
- The tests are executed by the doc-style test runner from this test tree.
- Each leaf sets the doctest arguments it wants to execute.

## Steps
1. Build the doctest binary from the module root.
2. Execute the binary given by `req.Bin`.
3. Capture stdout, stderr, exit code, and the raw execution error.

## Context
- These are real integration tests, not mocked unit tests.
- Agent tests expect `fake-codex` to be in PATH.

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

func Setup(t *testing.T, req *Request) error {
	req.Timeout = 30 * time.Second

	tmp := t.TempDir()
	doctestBin := filepath.Join(tmp, "doctest")
	buildDir := filepath.Join(DOCTEST_ROOT, "..")
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
```
