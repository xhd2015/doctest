## Preconditions
- Two consecutive runs on the same test tree produce identical test results.
- The auto-gen-dir parent provides multi-run Run.

## Steps
1. Use parent defaults: two runs, no modifications.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    cfg.TestDir = createTempTestProject(t, "mytest")
    return nil
}
```
