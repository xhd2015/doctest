## Preconditions
- The vet command detects anti-patterns in test file content.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	_ = req
	return nil
}
```
