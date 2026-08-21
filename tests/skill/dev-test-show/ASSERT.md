## Expected

- The command succeeds.
- stdout contains the concise, self-contained dev-test workflow.

## Exit Code

- Exit code 0.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d, stderr:\n%s", resp.ExitCode, resp.Stderr)
	}
	for _, want := range []string{
		"doctest-dev-test",
		"without sub-agent delegation",
		"Develop the requested change first, then add tests",
		"Significance-first",
		"Mutually exclusive and collectively exhaustive",
		"L2 in-process doctest",
		"70–85%",
		"L3 e2e doctest",
		"label: e2e",
	} {
		if !strings.Contains(resp.Stdout, want) {
			t.Fatalf("stdout missing %q:\n%s", want, resp.Stdout)
		}
	}
	for _, absent := range []string{
		"design-principle",
		"You are the **orchestrator**",
		"designer → RED",
		"tests are sealed — do not modify",
		"__DOCTEST_SPEC__",
	} {
		if strings.Contains(resp.Stdout, absent) {
			t.Fatalf("stdout must not contain %q:\n%s", absent, resp.Stdout)
		}
	}
}
```
