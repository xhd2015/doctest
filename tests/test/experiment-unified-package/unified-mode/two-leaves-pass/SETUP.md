# Scenario

**Feature**: 2-leaf fixture passes under default unified generation

```
doctest test --gen-dir tmp fixture/{a,b}
  -> exit 0; both leaves via suite
```

## Preconditions

- Default hierarchical unified path.
- Fixture leaves `a` and `b` are trivial pass asserts.

## Steps

1. `Op=run_gen` with empty Dir (auto fixture).
2. Assert `RunErr` empty.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Op = "run_gen"
	return nil
}
```
