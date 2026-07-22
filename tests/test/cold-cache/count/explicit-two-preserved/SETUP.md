# Scenario

**Feature**: explicit `-count=2` is preserved under `--cold-cache`

```
# count already set
doctest test --cold-cache -count=2 <tree>
  -> stderr: cd … && go test … -count=2 …
  -> not forced down to -count=1
```

## Preconditions

- Parent set cache sandbox + tiny test dir.

## Steps

1. Run with both `--cold-cache` and `-count=2`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"test", "--cold-cache", "-count=2", req.CCTestDir}
	return nil
}
```
