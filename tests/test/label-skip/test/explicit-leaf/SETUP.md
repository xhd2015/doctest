# Scenario

**Feature**: explicit leaf path runs labeled tests

```
# concrete leaf directory
doctest test <leaf-dir> -> execute labeled leaf
```

## Steps

1. Build a labeled temp tree.
2. Run `doctest test <leaf-dir>` (not tree root).

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