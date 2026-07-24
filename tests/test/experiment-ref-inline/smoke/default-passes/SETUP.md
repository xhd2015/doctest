# Scenario

**Feature**: tiny tree still passes under default generation

```
RunTest(1-pass tree) -> RunErr empty
```

## Preconditions

- Default hierarchical unified path.

## Steps

1. `Op=mini_run`.
2. Assert `RunErr` empty.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Op = "mini_run"
	return nil
}
```
