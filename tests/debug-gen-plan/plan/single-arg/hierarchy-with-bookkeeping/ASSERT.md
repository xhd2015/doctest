---
explanation: nested product binary + generate plan phase
---

## Expected

- CLI completes prepare path (bypass-go-test; exit may be 0 with BYPASS summary).
- Stderr contains gen-plan plan markers:
  - `gen-plan: arg[1/1]` (or `arg[1/1]` with gen-plan prefix)
  - bookkeeping under that plan tree: `go.mod` and `doctest.gen-manifest`
    (go.sum / doctest.tidy-done when written — soft if absent on some paths)
- Stderr does **not** require a separate `gen-plan: merged` for single-arg.
- Plan phase (before `gen-plan: result`) has **no** status annotations
  (`# new` / `# modified` / `# unchanged` / `# deleted`).
- Stdout does not contain `gen-plan:` lines (plan stays on stderr).

## Side Effects

- GenDir receives generated content (at least after successful generate).

## Errors

- Fail if gen-plan markers missing (Classic TDD RED pre-feature: unknown key
  or no plan emit).

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run error: %v\nstderr:\n%s", err, resp.Stderr)
	}
	// Pre-feature: DOCTEST_DEBUG=gen-plan=1 fails closed → non-zero / error text.
	combined := resp.Stderr + "\n" + resp.Stdout
	if strings.Contains(combined, `unknown key "gen-plan"`) || strings.Contains(combined, "unknown key \"gen-plan\"") {
		t.Fatalf("gen-plan not accepted yet (Classic RED until parse lands):\n%s", combined)
	}
	if !strings.Contains(resp.Stderr, "gen-plan:") {
		t.Fatalf("expected gen-plan markers on stderr, got:\n%s", resp.Stderr)
	}
	if !strings.Contains(resp.Stderr, "arg[1/1]") {
		t.Fatalf("expected gen-plan: arg[1/1] on stderr, got:\n%s", resp.Stderr)
	}
	// Bookkeeping in single-arg plan tree
	for _, name := range []string{"go.mod", genManifestName} {
		if !strings.Contains(resp.Stderr, name) {
			t.Fatalf("single-arg plan hierarchy missing %q:\n%s", name, resp.Stderr)
		}
	}
	// Plan phase (before gen-plan: result) must not carry status annotations.
	planPart := resp.Stderr
	if i := strings.Index(planPart, "gen-plan: result"); i >= 0 {
		planPart = planPart[:i]
	}
	for _, tag := range []string{"# new", "# modified", "# unchanged", "# deleted"} {
		if strings.Contains(planPart, tag) {
			t.Fatalf("plan phase must not include annotation %q:\n%s", tag, planPart)
		}
	}
	// Plan must not leak into stdout as gen-plan: lines
	for _, line := range strings.Split(resp.Stdout, "\n") {
		if strings.Contains(line, "gen-plan:") {
			t.Fatalf("gen-plan must stay on stderr; found on stdout: %q", line)
		}
	}
}
```
