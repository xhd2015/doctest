# Scenario

**Feature**: L3 share check is skipped when leaf count is below MinLeaves=10

```
# 3 leaves, 1 e2e → 33% would exceed 10%, but leaves < 10 → skip share check
writeShareFixture(3 leaves, 1 e2e)
  -> runner.VetArgs(["vet", dir])
  -> exit 0
```

## Preconditions

- Fixture: 3 leaves; one `label: e2e`.
- High e2e % alone must not fail when `leaves < 10`.

## Steps

1. Write small fixture under `t.TempDir()`.
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
	writeShareFixture(t, dir, shareSpecs(3, 1))
	req.Args = []string{"vet", dir}
	return nil
}
```
