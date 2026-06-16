# Scenario

**Feature**: a doctest tree may have no runnable leaves

```
# maps doc-style leaves to Go test packages
doctest build/test -> mapping-gen -> leaf <-> Go package

# gen-dir modes
auto gen-dir -> one package per leaf | explicit gen-dir -> user-specified layout
```

## Preconditions
- A doctest tree may have no runnable leaves.

## Steps
1. Create a project with a doctest root that has no ASSERT.md leaves.

```go
import (
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	t.Log("error-cases group")
	return nil
}
```
