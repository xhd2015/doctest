# Scenario

**Feature**: tests for mod a

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Name = "mod_a"
	return nil
}
```
