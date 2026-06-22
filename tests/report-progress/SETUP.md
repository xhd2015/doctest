# Scenario

**Feature**: tests in this group run the `report-progress` binary directly

```
# sub-agents report progress to a file
sub-agent --writes--> progress file (env var DOCTEST_PROGRESS_FILE)

# multiple entries append
each step -> structured JSON entry -> append to file
```

## Preconditions
- Tests in this group run the `report-progress` binary directly.
- The binary is dispatched via the doctest binary copied as `report-progress`.

## Steps
1. Copy the built doctest binary to a temp dir as `report-progress`.
2. Run the binary with the leaf's args and env.

```go
import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
	libdocbuild "github.com/xhd2015/doctest/libdoc/build"
)

func Setup(t *testing.T, req *Request) error {
	req.Timeout = 60 * time.Second
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
	req.Env = append(req.Env, "TEST_GROUP=report-progress")
	return nil
}
```
