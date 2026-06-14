## Preconditions
- A doc-style test tree with two groups (group-a with leaf-1, leaf-2; group-b with leaf-3).
- The parent cache-isolation SETUP.md provides multi-run Run.

## Steps
1. Set scenario to "full_after_subdir".
2. First run executes doctest on group-a sub-directory (caches only group-a tests).
3. Second run executes doctest on the full tree root (should run all tests).

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    cisoCfg.Scenario = "full_after_subdir"
    return nil
}
```
