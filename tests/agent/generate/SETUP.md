## Preconditions
- The generate command receives an idea and a target output directory.

## Steps
1. Create a temporary target directory path.
2. Run `doctest agent generate`.

```go
import (
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    outDir := filepath.Join(t.TempDir(), "generated-doc-tests")
    req.Args = []string{"agent", "generate", "a cli that prints invoices", "--agent-runner", "fake-codex", "--dir", outDir}
    return nil
}
```

