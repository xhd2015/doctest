# Scenario

**Feature**: `ColorAlways` colors the Pass summary metric on a pipe

```
# color mode
ColorAlways -> force ANSI

# colored regions (non-verbose only)
summary Pass/Fail/Cached -> green/red/gray
```

## Preconditions
- One passing package; stdout is a pipe.

## Steps
1. Set one passing leaf.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.PassCount = 1
	req.FailCount = 0
	return nil
}
```