# Root Missing Run Fixture

## Preconditions
- This fixture tree has a root SETUP.md that defines Request, Response, and Setup,
  but deliberately has no Run function.
- The missing Run should trigger a validation error during tree discovery.

## Steps
1. Build the tree. Discovery should fail because the root has no Run.

```go
import "fmt"

type Request struct {
    Input string
}

type Response struct {
    Output string
}

func Setup(t *testing.T, req *Request) error {
    if req.Input == "" {
        return fmt.Errorf("Input is required")
    }
    return nil
}
```
