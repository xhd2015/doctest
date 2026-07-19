# Scenario

**Feature**: successful single-package runs with profile/cover flags write output files

```
doctest test -cpuprofile <session-path>/cpu.out <single-package-dir>
  -> go test writes profile file -> file exists (non-empty preferred)
```

## Preconditions
- Use the single-package `basic-request-runner` fixture (no multi-package profile preflight).
- Paths are under session temp dirs to avoid races and leftover pollution.

## Steps
1. Leaves choose absolute profile/coverprofile paths under a session temp directory.
2. After a successful run, assert the file exists.

```go
import (
"github.com/xhd2015/doctest/session"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.Timeout < 120*time.Second {
		req.Timeout = 120 * time.Second
	}
	exampleDir := filepath.Join(d.DOCTEST_ROOT, "testdata", "basic-request-runner")
	if _, err := os.Stat(exampleDir); err != nil {
		t.Fatalf("fixture %s: %v", exampleDir, err)
	}
	return nil
}
```

