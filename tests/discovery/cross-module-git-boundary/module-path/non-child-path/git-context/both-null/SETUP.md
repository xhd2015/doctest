# Scenario

**Feature**: neither ancestor nor nested module is inside git

## Preconditions
- `gitRoot` returns `""` for both sides.
- Two nulls are equal → discovery continues.

## Steps
1. Create non-child module project without git.
2. Run `doctest test -v ./...` from root.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    t.Logf("both-null git group")
    return nil
}
```
