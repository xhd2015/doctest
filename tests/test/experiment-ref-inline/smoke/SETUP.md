# Scenario

**Feature**: mini RunTest under default hierarchical generation

```
RunTest(1-pass tree) -> pass
```

## Preconditions

- Default generation (no experiment flags).

## Steps

1. Set `Op=mini_run`.
2. Assert success.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Op = "mini_run"
	return nil
}
```
