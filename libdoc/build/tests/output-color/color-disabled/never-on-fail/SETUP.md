# Scenario

**Feature**: `ColorNever` leaves fail output uncolored

```
# color mode
ColorNever -> plain

# colored regions (non-verbose only)
fail dot -> red | pass dot -> plain | summary Pass/Fail/Cached -> green/red/gray
```

## Preconditions
- `ColorNever` forces plain output even when packages fail.

## Steps
1. Set `Color` to `ColorNever` with one failing leaf.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
)

func Setup(t *testing.T, req *Request) error {
	req.Color = core.ColorNever
	req.PassCount = 0
	req.FailCount = 1
	return nil
}
```