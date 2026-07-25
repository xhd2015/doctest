# Scenario

**Feature**: single-arg plan hierarchy lists bookkeeping files under arg[1/1]

```
stderr plan:
  gen-plan: arg[1/1] …
    go.mod
    doctest.gen-manifest  (and go.sum / tidy-done when written)
    <tree packages>
  # no gen-plan: merged for single-tree
```

## Preconditions

- Parent prepared single-tree fixture and Args.

## Steps

1. No extra setup — assert plan markers and bookkeeping names on stderr.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	_ = req
	return nil
}
```
