## Expected
- The command succeeds.
- stdout includes the TDD lite specification content.

## Exit Code
- Exit code 0.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed: %v", err)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("exit code = %d, stderr:\n%s", resp.ExitCode, resp.Stderr)
    }
    for _, want := range []string{
        "doc-style-test-based-tdd-lite",
        "single-agent doctest TDD",
        "Phase 1 — Requirements",
        "Phase 6 — Verify",
        "design tests → RED → seal → implement → GREEN",
        "<DOCTEST_DESIGN_SPEC>",
        "DSN (Domain Specific Notion)",
    } {
        if !strings.Contains(resp.Stdout, want) {
            t.Fatalf("stdout missing %q:\n%s", want, resp.Stdout)
        }
    }
    for _, absent := range []string{
        "__DOCTEST_DESIGN_SPEC__",
        "__DOCTEST_SPEC__",
        "__DOCTEST_VERSION__",
        "You are the **orchestrator**",
        "You NEVER touch source files",
        "All code changes happen exclusively through two sub-agents",
        "adversarial multi-agent TDD",
    } {
        if strings.Contains(resp.Stdout, absent) {
            t.Fatalf("stdout must not contain %q:\n%s", absent, resp.Stdout)
        }
    }
}
```