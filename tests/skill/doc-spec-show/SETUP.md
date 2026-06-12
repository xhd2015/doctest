## Preconditions
- The doc-style test specification exists under `agents/doctest/doc`.

## Steps
1. Run `doctest skill doc-spec show`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"skill", "doc-spec", "show"}
    return nil
}
```

