## Preconditions
- A temporary test project exists with a valid doctest tree.
- The auto-gen-dir parent provides multi-run Run.

## Steps
1. Use parent defaults: two runs, no modifications, no count override.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    cfg.TestDir = createTempTestProject(t, "mytest")
    return nil
}
```
