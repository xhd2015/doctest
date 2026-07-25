# Scenario

**Feature**: result summary counts modified / unchanged / deleted

```
gen-plan: result
  …
  summary: modified=N unchanged=M deleted=K
```

## Preconditions

- Single-tree fixture; --no-color so counts are plain text.

## Steps

1. prepareSingleTree with --no-color.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	prepareSingleTree(t, req)
	req.Args = baseTestArgs(req, "--no-color", req.TreeRoot)
	return nil
}
```
