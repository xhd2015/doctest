# Scenario

**Feature**: `doctest test` forwards Go profiling and cover flags to the underlying `go test` command

```
# build and run test binary, report results
doctest test <dir> [profile/cover flags] -> build -> run binary
  -> go test -cpuprofile=... -memprofile=... -coverprofile=... | rates | -cover | -trace | -outputdir
  -> pass/fail per leaf -> exit code

# path resolution
relative profile/outputdir/coverprofile paths -> abs-resolve against process cwd at parse time
  -> absolute paths on go command line

# rates
unset -> omit from go args | set (incl. 0) -> forward exact value
```

## Preconditions
- The fixture tree exists at `DOCTEST_ROOT/testdata/basic-request-runner`.
- Profile/cover flags are not yet implemented (Classic TDD): leaves must fail until forwarding lands.

## Steps
1. Verify the basic-request-runner fixture is present.
2. Descendant leaves set profile/cover args and assert go-test command lines, files, or parse errors.

```go
import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func Setup(t *testing.T, req *Request) error {
	// Nested go test + compile; keep headroom like go-test-cache leaves.
	req.Timeout = 120 * time.Second

	exampleDir := filepath.Join(DOCTEST_ROOT, "testdata", "basic-request-runner")
	if info, err := os.Stat(exampleDir); err != nil {
		t.Fatalf("testdata dir %s not found: %v", exampleDir, err)
	} else if !info.IsDir() {
		t.Fatalf("testdata dir %s is not a directory", exampleDir)
	}
	return nil
}
```
