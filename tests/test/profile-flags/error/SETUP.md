# Scenario

**Feature**: missing arguments for string profile flags are rejected at parse time

```
doctest test -cpuprofile   # missing value
  -> non-zero exit -> stderr mentions -cpuprofile
```

## Preconditions
- Parse fails before building/running the fixture (directory unused).

## Steps
1. Leaves pass incomplete profile flags and assert non-zero exit.

```go
import (
	"testing"
	"time"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	// Parse failures should be fast; keep a short bound.
	if req.Timeout <= 0 || req.Timeout > 30*time.Second {
		req.Timeout = 30 * time.Second
	}
	return nil
}
```

