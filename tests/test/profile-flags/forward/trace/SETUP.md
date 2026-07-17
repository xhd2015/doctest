# Scenario

**Feature**: `-trace` is forwarded to the go test command line (path abs-resolved when relative)

```
doctest test -v -trace traces/run.out <dir>
  -> stderr contains -trace=<abs>/traces/run.out
```

## Preconditions
- WorkDir is session-scoped so relative trace path resolves predictably.

## Steps
1. Create WorkDir; run with relative `-trace traces/run.out`.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	wd := filepath.Join(os.TempDir(), "doctest-profile-flags-"+DOCTEST_SESSION_ID, "trace")
	if err := os.MkdirAll(wd, 0o755); err != nil {
		t.Fatalf("mkdir workdir: %v", err)
	}
	req.WorkDir = wd

	exampleDir := filepath.Join(DOCTEST_ROOT, "testdata", "basic-request-runner")
	req.Args = []string{
		"test", "-v",
		"-trace", "traces/run.out",
		exampleDir,
	}
	return nil
}
```
