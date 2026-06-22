## Expected

- SIGINT is delivered during verbose `writeCases` (trigger: `leaf15_test.go` in stderr).
- `.doctest_run_*` is removed under `moduleRoot` after SIGINT (interrupt-safe cleanup).
- Stderr references `.doctest_run_` temp compile path during the run.
- `go test` has not started yet (no `cd ... go test` line in stderr).

## Exit Code

- Non-zero or interrupted (doctest does not complete normally).

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("SIGINT cleanup test failed before assert: %v", err)
	}

	assertNoDoctestRunDirs(t, moduleRoot)
	assertStderrUsesTempCompile(t, resp)
	if strings.Contains(resp.Stderr, "cd ") && strings.Contains(resp.Stderr, "go test") {
		t.Fatalf("expected SIGINT before go test started, but stderr contains go test invocation:\n%s", resp.Stderr)
	}
}
```