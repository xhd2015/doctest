# Scenario

**Feature**: `-run` is rejected at parse time (name-based selector)

```
runner.ParseTestOptions([-run, TestFoo, somedir])
  -> error: not supported; use path or --label
```

## Preconditions
- Nested L2 root: parse only.

## Steps
1. Pass `-run` with a dummy pattern and directory.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"-run", "TestFoo", "somedir"}
	return nil
}
```
