# Scenario

**Feature**: `doctest test --help` lists `-a` hard force (no `--no-leaf-cache`)

```
doctest test --help -> usage includes -a wipe/force; omits --no-leaf-cache
```

## Preconditions

- Parent set Op and Args.

## Steps

1. Run help once.
2. Assert combined output mentions `-a` and omits `--no-leaf-cache`.

## Context

- Backfill / docs leaf; expected GREEN when CLI usage string is up to date.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.Op = "runtime_once"
	req.Args = []string{"test", "--help"}
	return nil
}
```
