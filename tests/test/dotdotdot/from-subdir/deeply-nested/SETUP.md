# Scenario

**Feature**: a temporary project with a deeply nested DOCTEST.md tree and sibling trees elsewhere

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- A temporary project with a deeply nested DOCTEST.md tree and sibling trees elsewhere.

## Steps
1. Create temp project with DOCTEST.md at `group/subgroup/deep_tests/` and another at `other/`.
2. Run `doctest test ./...` from the `group/subgroup/deep_tests/` directory.

```go
import (
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    projDir := createTempProject(t, req)

    // create deep doctest tree: group/subgroup/deep_tests
    deepParent := filepath.Join(projDir, "group", "subgroup")
    if err := createTestTree(deepParent, "deep_tests"); err != nil {
        t.Fatalf("create deep_tests: %v", err)
    }

    // create sibling doctest tree: other
    if err := createTestTree(projDir, "other"); err != nil {
        t.Fatalf("create other: %v", err)
    }

    req.WorkDir = filepath.Join(deepParent, "deep_tests")
    req.Args = []string{"test", "-v", "./..."}
    return nil
}
```
