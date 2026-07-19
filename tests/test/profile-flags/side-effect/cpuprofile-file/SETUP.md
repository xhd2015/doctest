# Scenario

**Feature**: `-cpuprofile` causes a CPU profile file to be written after a successful run

```
doctest test -cpuprofile <session>/cpu.out <fixture>
  -> exit 0 -> <session>/cpu.out exists
```

## Preconditions
- Absolute profile path under session temp; fixture is single-package.

## Steps
1. Create session profile path; run without requiring `-v` (file side effect only).

```go
import (
"github.com/xhd2015/doctest/session"
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := filepath.Join(os.TempDir(), "doctest-profile-flags-"+d.DOCTEST_SESSION_ID, "side-cpu")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	cpuPath := filepath.Join(dir, "cpu.out")
	// Remove any leftover so Assert only passes if this run writes the file.
	_ = os.Remove(cpuPath)

	exampleDir := filepath.Join(d.DOCTEST_ROOT, "testdata", "basic-request-runner")
	req.Args = []string{
		"test",
		"-cpuprofile", cpuPath,
		exampleDir,
	}
	return nil
}
```
