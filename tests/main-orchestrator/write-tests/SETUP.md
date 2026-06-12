## Preconditions
- The main-orchestrator SETUP.md provides shared infrastructure.

## Steps
1. Build binaries and provide helpers via parent SETUP.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env, "TEST_GROUP=write-tests")
    return nil
}
```
