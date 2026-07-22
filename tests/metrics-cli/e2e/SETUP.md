# Scenario

**Feature**: sparse L3 binary smokes for `doctest metrics` CLI wiring

```
# build once per session
testbin.Ensure -> req.Bin
req.UseCLI = true
  -> doctest metrics --help | unknown subcommand
```

## Preconditions

- These leaves intentionally spawn the product binary (process boundary is SUT).
- Labeled `heavy` so default discovery skips them.
- Analyze combinatorics stay in-process under path/last/top/summary/show/prune.

## Steps

1. Set `UseCLI` and ensure `req.Bin` via session-shared `testbin.Ensure`.
2. Leaf sets `Args` for help or unknown subcommand.

## Context

- Module root: `DOCTEST_ROOT/../..`.
- No MetricsRoot fixtures required for help/unknown.

```go
import (
	"path/filepath"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/testbin"
	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.UseCLI = true
	if req.Timeout == 0 {
		req.Timeout = 45 * time.Second
	}
	if req.Bin == "" {
		req.Bin = testbin.Ensure(t, filepath.Join(d.DOCTEST_ROOT, "..", ".."))
	}
	return nil
}
```
