# Scenario

**Feature**: nested testcase go.mod replace for session package

```
# session import -> replace directive
replace github.com/xhd2015/doctest/session => <session-mod cache>
```

## Preconditions

- Nested-module compile path (no parent `internal/` import forcing internal-compile only).
- Public temp module with session-importing leaf.
- Cache tests serialize via flock so parallel leaves do not race on wipe/create
  of the shared session-mod directory.

## Steps

1. Acquire cache lock (same as cache/* leaves).
2. Create module with session import.
3. Run doctest with `--gen-dir` or inspect generated nested go.mod under work dir.
4. Assert replace line or consumer Once success.

```go
import (
	"path/filepath"
	"testing"
)

var genDir string

func Setup(t *testing.T, req *Request) error {
	lockCacheTests(t)
	genDir = t.TempDir()
	return nil
}
```
