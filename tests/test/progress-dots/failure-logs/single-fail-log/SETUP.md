# Scenario

**Feature**: one failing leaf emits its failure marker on stdout after the summary

```
# output
progress dots -> summary -> FAIL lines -> detailed failure text
```

## Preconditions
- One failing leaf with `SINGLE_FAIL_LOG_MARKER` in Assert.

## Steps
1. Create single-fail tree and run non-verbose `doctest test`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	testDir := createSingleFailLogTree(t)
	req.Args = []string{"test", testDir}
	return nil
}
```