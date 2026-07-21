# Scenario

**Feature**: `doctest build` compiles labeled leaves

```
# build ignores test-only skip
doctest build <labeled-tree> -> exit 0
```

## Steps

1. Configure `req.Args` with `doctest build` and a labeled temp tree.

```go
import (
	"testing"
	"time"
)

func Setup(t *testing.T, req *Request) error {
	if req.Timeout == 0 {
		req.Timeout = 120 * time.Second
	}
	return nil
}
```