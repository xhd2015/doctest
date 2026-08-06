# Scenario

**Feature**: materialize `testdata/realmod` into an isolated workdir per leaf

```
testdata/realmod -> t.TempDir()/realmod + overlay.json (abs paths)
leaf Setup sets Scenario -> Run exec go
```

## Preconditions

- `go` on PATH.
- Fixture at `DOCTEST_ROOT/testdata/realmod` (discovery skips `testdata`).

## Steps

1. Resolve fixture path from `d.DOCTEST_ROOT`.
2. Copy into `t.TempDir()/realmod`.
3. Write `overlay.json` with absolute Replace keys.
4. Set `req.WorkDir`; leaves set `req.Scenario`.

## Context

- Parallel-safe: each leaf gets its own copy.
- Overlay bodies live under `overlay-src/`; virtual paths never exist on disk
  until `-overlay` maps them.

```go
import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	if d == nil || d.DOCTEST_ROOT == "" {
		t.Fatal("DOCTEST_ROOT required")
	}
	fixture := filepath.Join(d.DOCTEST_ROOT, "testdata", "realmod")
	if st, err := os.Stat(fixture); err != nil || !st.IsDir() {
		t.Fatalf("fixture missing: %s (%v)", fixture, err)
	}
	dest := filepath.Join(t.TempDir(), "realmod")
	materializeFixture(t, fixture, dest)
	req.WorkDir = dest
	req.Scenario = ""
	return nil
}
```
