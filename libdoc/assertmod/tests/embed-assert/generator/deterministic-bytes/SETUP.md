# Scenario

**Feature**: two embed script runs produce identical output bytes

```
# stable filename sort
run 1 MD5 == run 2 MD5
```

## Steps

1. Run embed script twice with `SecondRun` enabled.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.SecondRun = true
	return nil
}
```