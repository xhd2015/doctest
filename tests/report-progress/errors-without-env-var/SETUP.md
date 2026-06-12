## Preconditions
- `PROGRESS_FILE` environment variable is NOT set.

## Steps
1. Invoke `report-progress` without `PROGRESS_FILE` env var.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"some progress"}
    return nil
}
```
