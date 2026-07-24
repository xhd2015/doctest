# Scenario

**Feature**: non-verbose `doctest test` forwards detailed go test failure logs to stdout

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# output
progress dots -> summary -> FAIL lines -> detailed failure text (stderr/stdout from go test)
```

## Preconditions
- Tests run `doctest test` without `-v`.
- Temp trees are built by parent helpers with unique failure markers in Assert.

## Steps
1. Build a temp tree for the failure-log scenario under test.
2. Run `doctest test <dir>` (non-verbose).

```go
import (
	"testing"
	"time"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.Timeout < 120*time.Second {
		req.Timeout = 120 * time.Second
	}
	return nil
}
```