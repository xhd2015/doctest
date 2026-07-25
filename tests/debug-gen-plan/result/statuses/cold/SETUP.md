# Scenario

**Feature**: cold first generate reports gen-plan: result with modified files

```
fresh GenDir
  -> gen-plan: result
  -> summary: modified>=1 (created files) unchanged=? deleted=0|K
```

## Preconditions

- Parent single-tree Args; Mode=cli once.

## Steps

1. Assert result marker and summary on cold run.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	_ = req
	return nil
}
```
