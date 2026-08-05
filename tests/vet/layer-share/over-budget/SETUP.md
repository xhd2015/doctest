# Scenario

**Feature**: full vet hard-fails when L3 (e2e) share exceeds 10% on a large enough tree

```
# 10 leaves, 2 labeled e2e → L3 share 20% > max 10%
writeShareFixture(10 leaves, 2 e2e)
  -> runner.VetArgs(["vet", dir])
  -> non-zero + L3 share message
```

## Preconditions

- Fixture: 10 leaves; first two have `label: e2e`; remaining unlabeled.
- Classic TDD: this leaf is **RED** until implementer lands the share check.

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
	writeShareFixture(t, dir, shareSpecs(10, 2))
	req.Args = []string{"vet", dir}
	return nil
}
```
