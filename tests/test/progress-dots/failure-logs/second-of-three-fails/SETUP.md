# Scenario

**Feature**: only the failing package leaks detailed failure output

```
# output
progress dots -> summary -> FAIL lines -> detailed failure text (only from fail package)
```

## Preconditions
- Three leaves in order: `pass_1`, `fail_2`, `pass_3`; only `fail_2` fails with `SECOND_FAIL_LOG_MARKER`.

## Steps
1. Create ordered three-leaf tree and run non-verbose `doctest test`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	testDir := createSecondOfThreeFailsTree(t)
	req.Args = []string{"test", testDir}
	return nil
}
```