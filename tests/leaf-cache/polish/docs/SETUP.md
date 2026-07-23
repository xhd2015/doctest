# Scenario

**Feature**: product CLI help mentions leaf-cache flags (**L2 in-process**)

```
doctest test --help
  -> usage includes -a and --no-leaf-cache
```

## Preconditions

- Short-path help: no product binary (`runtime_once` with empty Bin →
  `cli.RunWithWriter`).
- Unlabeled so default discovery runs it.

## Steps

1. Leave Bin empty; isolate env not required for help.
2. Child sets Op=`runtime_once` with `test --help`.

```go
import (
	"testing"
	"time"
)

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Timeout = 30 * time.Second
	// L2: Bin stays empty → Run uses in-process CLI for runtime_once.
	return nil
}
```
