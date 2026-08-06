# Scenario

**Feature**: `doctest build` still rejects `-overlay` (test-only flag)

```
cli.RunWithWriter(["build", "-overlay", "x.json", "<dir>"])
  -> non-zero / unrecognized or rejected flag
```

## Preconditions

- In-process CLI; no materialize helpers.
- Confirms scope: only `doctest test` gains overlay (build stays reject).

## Steps

1. Set build-reject mode and argv with `-overlay`.

```go
import (
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Mode = modeBuildReject
	// Dir need not be a real doctest tree: parse/flag reject should win first.
	req.BuildArgs = []string{"build", "-overlay", "user-overlay.json", filepath.Join(d.DOCTEST_ROOT, "testdata-missing")}
	return nil
}
```
