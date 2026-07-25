# Scenario

**Feature**: DOCTEST_DEBUG parser accepts or rejects gen-plan keys (library only)

```
debug.Parse(DebugEnv)
  -> Settings / error
```

## Preconditions

- L2 in-process: `Mode=parse`; no product CLI, no fixtures.
- Unlabeled (default discovery).

## Steps

1. Leaves set `req.DebugEnv` only.
2. Run calls `debug.Parse`; Assert checks error and flags.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.Mode = "parse"
	return nil
}
```
