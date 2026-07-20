# Scenario

**Feature**: preserve prior fix — ineffective assert flags must not churn go.mod

```
# multi-tree ./... prepare calls WriteGoMod with different per-tree flags
modPath == github.com/xhd2015/doctest
  withAssertReplace false -> true (+ cache paths)
  -> no extra replace lines in go.mod
  -> go.mod mtime stable
  -> tidy-done retained
```

## Preconditions

- Parent module path is the doctest self-module.
- Assert/session replace flags are ineffective for this modPath.

## Steps

1. Leaves configure doctest modPath and dual-flag WriteGoMod sequence.

```go
import (
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	req.ModPath = "github.com/xhd2015/doctest"
	req.HasMod = true
	return nil
}
```
