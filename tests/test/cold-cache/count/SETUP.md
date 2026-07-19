# Scenario

**Feature**: `--cold-cache` forces `-count=1` when count is unset

```
# count policy under cold-cache
--cold-cache (no -count) -> go test line includes -count=1
--cold-cache -count=2    -> go test line includes -count=2 (not forced down)
```

## Preconditions

- Parent helpers create a tiny fixture and cache sandbox.
- Observation is the stderr `cd … && go test …` preview line (always printed).

## Steps

1. Leaves run `--cold-cache` with or without an explicit `-count`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	withCacheSandbox(t, req)
	st.TestDir = createTempTestProject(t)
	return nil
}
```
