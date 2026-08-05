# Scenario

**Feature**: full vet accepts a tree at the L3 share ceiling (≤10% e2e, ≥10 leaves)

```
# 10 leaves, 1 labeled e2e → L3 share 10.0% ≤ max 10%
writeShareFixture(10 leaves, 1 e2e)
  -> runner.VetArgs(["vet", dir])
  -> exit 0
```

## Preconditions

- Fixture: 10 leaves; first has `label: e2e`; remaining unlabeled.
- Share is exactly 10% → within budget (fail only when **>** 10%).

## Steps

1. Write fixture under `t.TempDir()`.
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
	writeShareFixture(t, dir, shareSpecs(10, 1))
	req.Args = []string{"vet", dir}
	return nil
}
```
