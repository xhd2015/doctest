# Scenario

**Feature**: quiet and `-v` prepare/run the same case count on the mega parent fixture

```
# dual CLI
quiet:   doctest test parent/     -> exit 0; planned Pq
verbose: doctest test -v parent/  -> exit 0; planned Pv
assert:  Pq == Pv == 1
```

## Preconditions

- Same mega fixture for both runs (one module, one parent path).
- **RED** until verbose prepare succeeds with the same planned count as quiet.

## Steps

1. Create mega parent nested fixture once.
2. Set `Op=dual_cli` with quiet and verbose arg lists targeting the same parent.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	_, parentDir := createMegaParentNestedFixture(t)
	req.Op = "dual_cli"
	req.QuietArgs = []string{"test", "--no-color", parentDir}
	req.VerboseArgs = []string{"test", "-v", "--no-color", parentDir}
	return nil
}
```
