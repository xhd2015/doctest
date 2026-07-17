---
label: heavy
---

## Expected

- Exit code 0.
- Exactly one `*_test.go` file is generated (changed leaf only).

## Exit Code

0

```go
import (
	"os"
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	genDir := ""
	for _, e := range req.Env {
		if strings.HasPrefix(e, "CHANGED_GEN_DIR=") {
			genDir = strings.TrimPrefix(e, "CHANGED_GEN_DIR=")
			break
		}
	}
	if genDir == "" {
		t.Fatal("CHANGED_GEN_DIR not set in req.Env")
	}
	count := countGeneratedTestGoFiles(t, genDir)
	if count != 1 {
		t.Fatalf("expected 1 generated *_test.go file, got %d in %s", count, genDir)
	}
	if _, err := os.Stat(genDir); err != nil {
		t.Fatalf("gen dir missing: %v", err)
	}
}
```