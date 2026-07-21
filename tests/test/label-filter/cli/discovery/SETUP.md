# Scenario

**Feature**: discovery mode applies label filter across the whole mod

```
doctest test <mod-root> --label EXPR -> run matching labeled leaves only
```

## Steps

1. Use standard five-leaf fixture mod.

```go
import (
	"testing"
	"time"
)

func Setup(t *testing.T, req *Request) error {
	if req.Timeout == 0 {
		req.Timeout = 120 * time.Second
	}
	_ = t
	return nil
}
```