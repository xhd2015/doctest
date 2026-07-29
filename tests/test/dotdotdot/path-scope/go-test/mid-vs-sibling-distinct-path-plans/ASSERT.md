---
label: heavy
---

## Expected

- Mid plan ≠ sibling plan.
- Mid plan mentions `mid`; sibling plan mentions `sibling`.
- Neither collapses to an identical root-only suite plan for both scopes.

```go
import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func pathScopeRunVerbose(t *testing.T, req *Request, pattern string) string {
	t.Helper()
	genDir := filepath.Join(req.WorkDir, ".gen")
	cmd := exec.Command(req.Bin, "test", "-v", "--gen-dir", genDir, "-count=1", pattern)
	cmd.Dir = req.WorkDir
	cmd.Env = append(os.Environ(), "GOWORK=off", "DOCTEST_CACHE_HOME="+t.TempDir())
	b, err := cmd.CombinedOutput()
	out := string(b)
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok && ee.ExitCode() != 0 {
			// still inspect plans on failure
			return out
		}
		t.Fatalf("run %s: %v\n%s", pattern, err, out)
	}
	return out
}

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	midOut := pathScopeRunVerbose(t, req, "./tree/mid/...")
	sibOut := pathScopeRunVerbose(t, req, "./tree/sibling/...")

	midPlans := pathScopeGoTestPlanLines(midOut)
	sibPlans := pathScopeGoTestPlanLines(sibOut)
	if len(midPlans) == 0 {
		t.Fatalf("no go test plan for mid:\n%s", midOut)
	}
	if len(sibPlans) == 0 {
		t.Fatalf("no go test plan for sibling:\n%s", sibOut)
	}
	if midPlans[0] == sibPlans[0] {
		t.Fatalf("mid and sibling share the same go test plan (path scope lost):\n  mid: %s\n  sib: %s", midPlans[0], sibPlans[0])
	}
	// Expect ./tree/mid/... and ./tree/sibling/... (not */suite).
	if !strings.Contains(midPlans[0], "mid") || !strings.Contains(midPlans[0], "/...") {
		t.Fatalf("mid plan want ... under mid:\n  %s\n%s", midPlans[0], midOut)
	}
	if !strings.Contains(sibPlans[0], "sibling") || !strings.Contains(sibPlans[0], "/...") {
		t.Fatalf("sibling plan want ... under sibling:\n  %s\n%s", sibPlans[0], sibOut)
	}
	if strings.Contains(midPlans[0], "/suite") && !strings.Contains(midPlans[0], "/...") {
		t.Fatalf("mid plan must not hard-code suite package: %s", midPlans[0])
	}
}
```
