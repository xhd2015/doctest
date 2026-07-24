# Scenario

**Feature**: a --gen-dir is specified to control where output goes

```
# maps doc-style leaves to Go test packages
doctest build/test -> mapping-gen -> leaf <-> Go package

# gen-dir modes
auto gen-dir -> one package per leaf | explicit gen-dir -> user-specified layout
```

## Preconditions
- A --gen-dir is specified to control where output goes.

## Steps
1. Add --gen-dir to the test arguments so the generated output is written to a known location.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	// Request-local gen dir: Parallel-safe (no package-level var genDir).
	req.GenDir = filepath.Join(t.TempDir(), "generated")
	if err := os.MkdirAll(req.GenDir, 0755); err != nil {
		t.Fatalf("mkdir gen dir: %v", err)
	}
	req.Args = append(req.Args, "--gen-dir", req.GenDir)
	return nil
}
```
