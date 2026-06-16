# Scenario

**Feature**: color disabled — no ANSI on dot progress or summary

```
# color mode
ColorAuto -> TTY check on stdout | ColorNever -> plain

# colored regions (non-verbose only)
fail dot -> red | pass dot -> plain | summary Pass/Fail/Cached -> green/red/gray
```

## Preconditions
- Color is disabled via `ColorAuto` on a pipe or `ColorNever`.

## Steps
1. Configure a single-package sub-tree unless a leaf overrides counts.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.PassCount = 1
	req.FailCount = 0
	return nil
}
```