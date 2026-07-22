# Scenario

**Feature**: unset `-count` becomes `-count=1` under `--cold-cache`

```
# count unset
doctest test --cold-cache <tree>   # no -count flag
  -> stderr: cd … && go test … -count=1 …
```

## Preconditions

- Parent set cache sandbox + tiny test dir.
- No `-count` on the command line.

## Steps

1. Run `doctest test --cold-cache <testDir>` without `-count`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"test", "--cold-cache", req.CCTestDir}
	return nil
}
```
