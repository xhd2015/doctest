---
label: heavy
explanation: multi-arg nested product generate plan + merged
---

## Expected

- Stderr plan shows `arg[1/2]` and `arg[2/2]`.
- Stderr includes `gen-plan: merged`.
- Merged section (or overall plan stderr) includes bookkeeping `go.mod` and
  `doctest.gen-manifest`.
- Per-arg sections are package-oriented: requirement is that go.mod is not
  repeated as the only content of each arg tree — assert merged mentions go.mod
  and both tree names appear in plan output.
- If hub/workspace is written, `__workspace` may appear under merged (soft:
  required only when present in gen root after run).
- arg + merged (before `gen-plan: result`) have **no** status annotations.

## Errors

- RED until multi-arg plan emit lands.

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run error: %v\nstderr:\n%s", err, resp.Stderr)
	}
	combined := resp.Stderr + "\n" + resp.Stdout
	if strings.Contains(combined, `unknown key "gen-plan"`) {
		t.Fatalf("gen-plan not accepted yet:\n%s", combined)
	}
	se := resp.Stderr
	if !strings.Contains(se, "arg[1/2]") || !strings.Contains(se, "arg[2/2]") {
		t.Fatalf("expected multi-arg plan arg[1/2] and arg[2/2], got:\n%s", se)
	}
	if !strings.Contains(se, "gen-plan: merged") && !strings.Contains(se, "merged") {
		t.Fatalf("expected gen-plan: merged for multi-arg, got:\n%s", se)
	}
	// Merged / overall must show bookkeeping once (go.mod + manifest)
	if !strings.Contains(se, "go.mod") {
		t.Fatalf("expected go.mod in multi-arg plan/merged:\n%s", se)
	}
	if !strings.Contains(se, genManifestName) {
		t.Fatalf("expected %s in multi-arg plan/merged:\n%s", genManifestName, se)
	}
	// Both tree identities appear somewhere in plan output
	if !strings.Contains(se, "tree-a") && !strings.Contains(se, "a1") {
		t.Fatalf("expected tree-a content in plan stderr:\n%s", se)
	}
	if !strings.Contains(se, "tree-b") && !strings.Contains(se, "b1") {
		t.Fatalf("expected tree-b content in plan stderr:\n%s", se)
	}
	// Soft: if gen root has __workspace, plan/merged should mention it
	ws := filepath.Join(req.GenDir, "__workspace")
	if st, err := os.Stat(ws); err == nil && st.IsDir() {
		if !strings.Contains(se, "__workspace") {
			t.Fatalf("gen root has __workspace but plan/merged omitted it:\n%s", se)
		}
	}
	// arg + merged (before result) must not carry status annotations.
	planPart := se
	if i := strings.Index(planPart, "gen-plan: result"); i >= 0 {
		planPart = planPart[:i]
	}
	for _, tag := range []string{"# new", "# modified", "# unchanged", "# deleted"} {
		if strings.Contains(planPart, tag) {
			t.Fatalf("plan/merged must not include annotation %q:\n%s", tag, planPart)
		}
	}
}
```
