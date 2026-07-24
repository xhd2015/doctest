# Scenario

**Feature**: tests for simple

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    t.Logf("setup: %s", req.Name)
    return nil
}
```
