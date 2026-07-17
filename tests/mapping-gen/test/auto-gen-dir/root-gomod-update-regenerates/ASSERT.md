---
label: heavy
---

## Expected
- Warm-up run succeeds (exit 0) and populates the mapping-gen cache.
- After root `go.mod` removes the local replace and the `dep/` directory is deleted, the second run still succeeds (exit 0).
- The cached mapping-gen `go.mod` no longer contains a `replace localdep` directive.

## Exit Code
- Exit code 0.

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v\nstderr:\n%s", err, resp.Stderr)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("second run expected exit 0 after go.mod update, got %d\nstderr:\n%s\nstdout:\n%s",
			resp.ExitCode, resp.Stderr, resp.Stdout)
	}

	if rootGoModState.WarmResp == nil {
		t.Fatal("warm run state not set by Setup")
	}
	if rootGoModState.WarmResp.ExitCode != 0 {
		t.Fatalf("warm run expected exit 0, got %d\nstderr:\n%s", rootGoModState.WarmResp.ExitCode, rootGoModState.WarmResp.Stderr)
	}

	genRoot := rootGoModState.GenRoot
	if strings.HasPrefix(genRoot, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			t.Fatalf("home dir: %v", err)
		}
		genRoot = filepath.Join(home, genRoot[2:])
	}
	cachedGoMod, err := os.ReadFile(filepath.Join(genRoot, "go.mod"))
	if err != nil {
		t.Fatalf("read cached go.mod at %s: %v", genRoot, err)
	}
	if strings.Contains(string(cachedGoMod), "replace localdep") {
		t.Fatalf("cached go.mod still has stale replace localdep after root go.mod update:\n%s", string(cachedGoMod))
	}
}
```