# Scenario

**Feature**: when nested go test times out, doctest surfaces a clear timeout message

```
# user passes short timeout; suite sleeps past it
doctest test --timeout=2s <sleep-tree> -> go test -timeout=2s -> panic: test timed out after 2s

# fail path must show an obvious Error (not only a buried goroutine dump)
doctest -> Error: go test timed out after 2s  (stdout/stderr)

# fast pass must not emit that Error
doctest test <pass-tree> -> exit 0, no "Error: go test timed out"
```

## Preconditions

- The doctest binary is resolved via `testbin.Ensure` from the module root.
- Temp fixture trees are created per leaf (sleep or fast-pass).
- Outer harness timeout is generous (compile + nested 2s timeout).

## Steps

1. Build or reuse the shared doctest binary.
2. Set a generous outer `req.Timeout` (nested runs compile + may wait on timeout).
3. Provide helpers to write sleep and fast-pass fixture trees.

## Context

- Nested CLI self-tests are labeled `heavy`.
- Default timeout policy is unchanged: only messaging when a timeout actually fires.
- Prefer message form `Error: go test timed out after 2s`; accept equivalent
  `test timed out after` / `timed out after` on the fail path.

```go
import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/testbin"
	"github.com/xhd2015/doctest/libdoc/testtree"
	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Timeout = 120 * time.Second
	req.Bin = testbin.Ensure(t, filepath.Join(d.DOCTEST_ROOT, "..", "..", ".."))
	return nil
}

// createSleepTree builds a 1-leaf fixture whose Run sleeps for sleepSec seconds
// so nested go test -timeout can fire.
func createSleepTree(t *testing.T, sleepSec int) string {
	t.Helper()
	tmp := t.TempDir()
	testtree.WriteMinimalRunnableTree(t, tmp, []testtree.LeafSpec{
		{
			Name:     "sleep",
			Steps:    "Run sleeps past go test -timeout",
			Expected: "would pass if not timed out",
		},
	})
	runGo := fmt.Sprintf(`import (
	"testing"
	"time"
)

type Request struct{}
type Response struct{}

func Run(t *testing.T, req *Request) (*Response, error) {
	time.Sleep(%d * time.Second)
	return &Response{}, nil
}`, sleepSec)
	// Overwrite DOCTEST.md Run so the nested suite actually sleeps.
	testtree.WriteFile(t, tmp, "DOCTEST.md", testtree.MinimalDOCTEST(runGo))
	return tmp
}

// createFastPassTree builds a 1-pass fixture that finishes well under any normal timeout.
func createFastPassTree(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	testtree.WritePassFailTree(t, tmp, 1, 0)
	return tmp
}

// combinedOutput joins stdout and stderr for fail-path assertions.
func combinedOutput(resp *Response) string {
	if resp == nil {
		return ""
	}
	return resp.Stdout + "\n" + resp.Stderr
}

// hasTimeoutSignal reports whether combined output contains a user-facing
// timeout message (preferred Error line or go-test panic phrasing).
func hasTimeoutSignal(combined string) bool {
	if strings.Contains(combined, "Error: go test timed out after") {
		return true
	}
	if strings.Contains(combined, "test timed out after") {
		return true
	}
	if strings.Contains(combined, "timed out after") {
		return true
	}
	return false
}
```
