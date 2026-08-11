# Scenario

**Feature**: multi-leaf parent-internal subject tests PASS under suite-only go test

```
RunTest(multi-leaf parent internal, GenDir)
  -> both subject leaves PASS
  -> go test package args: single suite (not ./leaf-a ./leaf-b)
```

## Preconditions

- Fixture from parent multi-leaf Setup (ModuleRoot/TestDir/GenDir already set).
- Packaging focus: no coverprofile.

## Steps

1. Ensure cover flags off.
2. Run subject tree with GenDir; assert suite-only success.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	_ = d
	// Explicit packaging-only run (fixture already on req from multi-leaf Setup).
	req.WithCover = false
	req.CoverPath = ""
	return nil
}
```
