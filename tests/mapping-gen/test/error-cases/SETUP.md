## Preconditions
- A doctest tree may have no runnable leaves.

## Steps
1. Create a project with a doctest root that has no ASSERT.md leaves.

```go
import (
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	t.Log("error-cases group")
	return nil
}
```
