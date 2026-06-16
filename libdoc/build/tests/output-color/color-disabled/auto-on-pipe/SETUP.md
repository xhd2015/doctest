# Scenario

**Feature**: `ColorAuto` with stdout redirected to a pipe disables color

```
# color mode
ColorAuto -> TTY check on stdout | ColorAlways -> force ANSI | ColorNever -> plain
```

## Preconditions
- `ColorAuto` is the default; stdout is a pipe (not a char device).

## Steps
1. Set `Color` to `ColorAuto`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
)

func Setup(t *testing.T, req *Request) error {
	req.Color = core.ColorAuto
	return nil
}
```