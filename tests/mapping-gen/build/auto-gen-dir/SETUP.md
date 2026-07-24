# Scenario

**Feature**: no --gen-dir is specified; build uses a temp directory

```
# maps doc-style leaves to Go test packages
doctest build/test -> mapping-gen -> leaf <-> Go package

# gen-dir modes
auto gen-dir -> one package per leaf | explicit gen-dir -> user-specified layout
```

## Preconditions
- No --gen-dir is specified; build uses a temp directory.

## Steps
1. Create a project with 1 leaf.
2. Run `doctest build <test-dir> -v`.

```go
import (
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Log("using auto gen dir (mapping-gen cache, no --gen-dir)")
	return nil
}
```
