# Scenario

**Feature**: non-`e2e` labels (e.g. `slow`) count as L2 — large tree with only `slow` stays within L3 budget

```
# 10 leaves all label: slow (no e2e) → L3=0 → exit 0
writeShareFixture(10 leaves, all slow)
  -> runner.VetArgs(["vet", dir])
  -> exit 0
```

## Preconditions

- Align with `doctest list` inventory: L3 is **`e2e` only**.
- Other labels must not inflate L3 share.

## Steps

1. Write fixture: 10 leaves each with `label: slow`.
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
	writeShareFixture(t, dir, slowOnlySpecs(10))
	req.Args = []string{"vet", dir}
	return nil
}
```
