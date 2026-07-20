# Scenario

**Feature**: `doctest test --help` lists `-a` and `--no-leaf-cache`

```
doctest test --help -> usage text includes force and no-leaf-cache options
```

## Preconditions

- Parent set Op and Args.

## Steps

1. Run help once.
2. Assert combined output mentions the flags.

## Context

- Backfill / docs leaf; expected GREEN when CLI usage string is up to date.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "runtime_once"
	req.Args = []string{"test", "--help"}
	return nil
}
```
