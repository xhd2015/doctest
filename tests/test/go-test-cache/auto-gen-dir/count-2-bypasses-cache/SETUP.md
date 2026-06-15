## Preconditions
- A first run has completed successfully.
- The auto-gen-dir parent provides multi-run Run.

## Steps
1. Set cfg.UseCountTwo = true so second run uses `-count=2`.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    cfg.TestDir = createTempTestProject(t, "mytest")
    cfg.UseCountTwo = true
    doMultiRun(t, req)
    return nil
}
```
