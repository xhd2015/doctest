# Scenario

**Feature**: all-pass run colors summary metrics but not pass dots

```
# colored regions (non-verbose only)
pass dot -> plain | summary Pass -> green when N>0 | summary Cached -> gray always
```

## Preconditions
- Two passing packages.

## Steps
1. Set `PassCount` to 2.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.PassCount = 2
	req.FailCount = 0
	return nil
}
```