# Scenario

**Feature**: rule R2: Every non-root SETUP.md must have func Setup

```
# enforce SETUP.md rules at build time
root missing Run -> build error
child redefines Run -> build error
non-root missing Setup -> build error
```

## Preconditions
- Rule R2: Every non-root SETUP.md must have func Setup.
- Run is reserved for root; non-root SETUP.md without at least Setup is invalid.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    _ = req
    return nil
}
```
