## Preconditions
- The doc-style test code specification exists under `agents/doctest/doc`.

## Steps
1. Run `doctest skill code-spec show`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"skill", "code-spec", "show"}
    return nil
}
```

