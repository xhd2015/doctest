## Preconditions
- A first run has completed and created the hash-based gen directory.
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
