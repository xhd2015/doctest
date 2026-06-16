# Scenario

**Feature**: `doctest test` color flags are parsed at the CLI boundary

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- The doctest binary is built by the root `tests/SETUP.md` harness.
- Color flags are mutually exclusive when both `--color` and `--no-color` are passed.

## Steps
1. Configure `req.Args` for the color-flag scenario under test.

```go
import (
	"testing"
	"time"
)

func Setup(t *testing.T, req *Request) error {
	req.Timeout = 20 * time.Second
	return nil
}
```