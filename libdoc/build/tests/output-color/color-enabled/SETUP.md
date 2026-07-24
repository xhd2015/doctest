# Scenario

**Feature**: color enabled via `ColorAlways` for deterministic ANSI assertions

```
# color mode
ColorAlways -> force ANSI

# colored regions (non-verbose only)
fail dot -> red | pass dot -> plain | summary Pass/Fail/Cached -> green/red/gray
```

## Preconditions
- `ColorAlways` forces ANSI even when stdout is a pipe.

## Steps
1. Set `Color` to `ColorAlways` unless a leaf overrides other fields.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Color = core.ColorAlways
	return nil
}
```