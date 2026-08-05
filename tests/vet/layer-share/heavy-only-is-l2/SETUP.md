# Scenario

**Feature**: `label: heavy` without `e2e` counts as L2 — all-heavy large tree stays within budget

```
# 10 leaves all label: heavy (no e2e) → L3=0 → exit 0
writeShareFixture(10 leaves, all heavy)
  -> runner.VetArgs(["vet", dir])
  -> exit 0
```

## Preconditions

- Align with `doctest list` inventory: L3 is **`e2e` only**.
- Cost label `heavy` alone must not inflate L3 share.

## Steps

1. Write fixture: 10 leaves each with `label: heavy`.
2. Run `vet <dir>` in-process.

```go
import (
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	dir := filepath.Join(t.TempDir(), "tree")
	writeShareFixture(t, dir, heavySpecs(10))
	req.Args = []string{"vet", dir}
	return nil
}
```
