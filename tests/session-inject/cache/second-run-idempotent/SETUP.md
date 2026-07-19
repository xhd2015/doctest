# Scenario

**Feature**: second doctest run with session import does not rewrite session-mod cache

```
# warm cache
first materialize -> second doctest test -> same bytes / still complete layout
```

## Preconditions

- First run populates cache; second run reuses it.

## Steps

1. Ensure cache exists via first run (or leave if warm).
2. Record content of go.mod.
3. Run doctest test again.
4. Assert go.mod unchanged and layout still valid.

```go
import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

var goModBefore []byte

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	// Ensure warm: materialize via first subprocess if needed.
	createPublicModuleProject(t, "", defaultSessionAssertGo(), true)
	setupModuleEnv(t, req)
	req.Args = []string{"test", testDir, "-v"}
	// Warm-up run
	warm := *req
	if _, err := Run(t, d, &warm); err != nil {
		// continue; first-run may still create cache before failing
		_ = err
	}
	cacheDir := expectedSessionCacheDir(t)
	var err error
	goModBefore, err = os.ReadFile(filepath.Join(cacheDir, "go.mod"))
	if err != nil {
		t.Fatalf("warm cache missing go.mod: %v", err)
	}
	// Second run is the measured Run.
	return nil
}
```
