# Scenario

**Feature**: the generate command receives an idea and a target output directory

```
# agent reads requirement, invokes Fake Codex, writes output
doctest agent <cmd> --requirement req.md -> Fake Codex -> generated code

# session state tracked in event files
doctest <- Fake Codex (session id, events, progress)
```

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
	req.UseCLI = true // true e2e agent generate
	requireFakeCodex(t, req)
	outDir := filepath.Join(t.TempDir(), "generated-doc-tests")
	req.Args = []string{"agent", "generate", "a cli that prints invoices", "--agent-runner", "fake-codex", "--dir", outDir}
	return nil
}
```

