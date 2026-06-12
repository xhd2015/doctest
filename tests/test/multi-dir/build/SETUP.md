## Preconditions
- Build tests have a longer timeout due to compilation.

## Steps
1. Extend timeout for build operations.

```go
import (
    "testing"
    "time"
)

func Setup(t *testing.T, req *Request) error {
    req.Timeout = 120 * time.Second
    return nil
}
```
