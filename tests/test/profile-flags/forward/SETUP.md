# Scenario

**Feature**: with `-v`, the printed `go test` command line includes forwarded profile/cover flags

```
# verbose run prints the go test invocation on stderr
doctest test -v <dir> [flags] -> stderr: "cd ... && go test ... -cpuprofile=... ..."
```

## Preconditions
- Fixture `basic-request-runner` is available (parent setup).
- Leaves use a single-package/leaf tree so multi-package profile rules do not block the run.

## Steps
1. Run `doctest test -v` against the fixture with one or more profile/cover flags.
2. Assert the stderr command line contains the expected flag forms (absolute paths where required).

```go
import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func Setup(t *testing.T, req *Request) error {
	// Forward leaves run a nested go test; keep generous timeout.
	if req.Timeout < 120*time.Second {
		req.Timeout = 120 * time.Second
	}
	exampleDir := filepath.Join(DOCTEST_ROOT, "testdata", "basic-request-runner")
	if _, err := os.Stat(exampleDir); err != nil {
		t.Fatalf("fixture %s: %v", exampleDir, err)
	}
	return nil
}
```

