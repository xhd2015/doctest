# Scenario

**Feature**: tests for mod b

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Name = "mod_b"
	return nil
}
```
