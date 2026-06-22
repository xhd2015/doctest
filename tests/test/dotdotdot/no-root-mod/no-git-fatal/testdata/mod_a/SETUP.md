# Scenario

**Feature**: tests for mod a

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Name = "mod_a"
	return nil
}
```
