# Scenario

**Feature**: build the doctest binary for `--changed` integration tests

```
# build doctest from module source
go build ./cmd/doctest -> doctest binary

# invoke with --changed against ephemeral git repos
doctest <subcmd> --changed <fixture-tree> -> filter leaves by git status
```

## Preconditions

- The doctest module root is two levels above this tree (`DOCTEST_ROOT/../..`).
- Each leaf configures subcommand args and git fixture layout.

## Steps

1. Build the doctest binary from the module root.
2. Store the binary path in `req.Bin` for `Run` to execute.

## Context

- Fixture test trees are created inside ephemeral git repos by descendant `Setup` functions.
- `req.WorkDir` is set to the git repo root when running against fixtures.

```go
import (
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	libdocbuild "github.com/xhd2015/doctest/libdoc/build"
)

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
```